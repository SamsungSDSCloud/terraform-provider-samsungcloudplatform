package publicip

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/public-ip2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_public_ip", ResourceVpcPublicIp())
	scp.RegisterDataSource("scp_public_ip", DatasourceVpcPublicIp())
}

func ResourceVpcPublicIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPublicIpCreate,
		ReadContext:   resourceVpcPublicIpRead,
		UpdateContext: resourceVpcPublicIpUpdate,
		DeleteContext: resourceVpcPublicIpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description of public IP",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of public IP",
			},
		},
		Description: "Provides a Public IP resource.",
	}
}

func resourceVpcPublicIpCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	description := rd.Get("description").(string)
	location := rd.Get("region").(string)

	inst := meta.(*client.Instance)

	serviceZoneId, productGroupId, err := client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, location, common.NetworkProductGroup, common.PublicIpProductName)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := inst.Client.PublicIp.CreatePublicIp(ctx, productGroupId, common.VpcPublicIpPurpose, serviceZoneId, common.VpcPublicIpUplinkType, description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForVpcPublicIpStatus(ctx, inst.Client, result.PublicIpAddressId, []string{}, []string{common.ReservedState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.PublicIpAddressId)

	return resourceVpcPublicIpRead(ctx, rd, meta)
}

func resourceVpcPublicIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	info, _, err := inst.Client.PublicIp.GetPublicIp(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("ipv4", info.IpAddress)
	rd.Set("description", info.PublicIpAddressDescription)

	location, err := client.FindLocationName(ctx, inst.Client, info.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information")
	}

	rd.Set("region", location)

	return nil
}

func resourceVpcPublicIpUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("description") {
		_, err := inst.Client.PublicIp.UpdatePublicIp(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpcPublicIpRead(ctx, rd, meta)
}

func resourceVpcPublicIpDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.PublicIp.DeletePublicIp(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForVpcPublicIpStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED", "FREE"}, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForVpcPublicIpStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.PublicIp.GetPublicIp(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.PublicIpState, nil
	})
}

func DatasourceVpcPublicIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVpcPublicIpRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public Ip Id",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of public IP",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region name",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of public IP",
			},
		},
		Description: "Provides list of public ips.",
	}
}

func convertPublicIpListToHclSet(publicIps []publicip2.DetailPublicIpResponse) (common.HclSetObject, []string) {
	var setPublicIps common.HclSetObject
	var ids []string
	for _, publicIpInfo := range publicIps {
		if len(publicIpInfo.PublicIpAddressId) == 0 {
			continue
		}
		ids = append(ids, publicIpInfo.PublicIpAddressId)
		kv := common.HclKeyValueObject{
			"id":                   publicIpInfo.PublicIpAddressId,
			"purpose":              publicIpInfo.PublicIpPurpose,
			"description":          publicIpInfo.PublicIpAddressDescription,
			"state":                publicIpInfo.PublicIpState,
			"ipv4":                 publicIpInfo.IpAddress,
			"ip_id":                publicIpInfo.IpAddressId,
			"zone_id":              publicIpInfo.ServiceZoneId,
			"project_id":           publicIpInfo.ProjectId,
			"attached_object_name": publicIpInfo.AttachedObjectName,
		}
		setPublicIps = append(setPublicIps, kv)
	}
	return setPublicIps, ids
}

func datasourceVpcPublicIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	location := rd.Get("region").(string)

	inst := meta.(*client.Instance)

	serviceZoneId, _, err := client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, location, common.NetworkProductGroup, common.PublicIpProductName)
	if err != nil {
		return diag.FromErr(err)
	}

	publicIpList, err := inst.Client.PublicIp.GetPublicIpList(ctx, serviceZoneId, &publicip2.PublicIpOpenApiControllerApiListPublicIpsV2Opts{
		IpAddress:       optional.String{},
		IsBillable:      optional.NewBool(true),
		IsViewable:      optional.NewBool(true),
		PublicIpPurpose: optional.NewString(common.VpcPublicIpPurpose),
		PublicIpState:   optional.String{},
		UplinkType:      optional.NewString(common.VpcPublicIpUplinkType),
		CreatedBy:       optional.String{},
		Page:            optional.NewInt32(0),
		Size:            optional.NewInt32(10000),
		Sort:            optional.NewInterface([]string{"createdDt:desc"}),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	setPublicIps, _ := convertPublicIpListToHclSet(publicIpList.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		setPublicIps = common.ApplyFilter(DatasourceVpcPublicIp().Schema, f.(*schema.Set), setPublicIps)
	}

	if len(setPublicIps) == 0 {
		return diag.Errorf("no matching public ip found")
	}

	kvPublicIp := setPublicIps[0]

	location, err = client.FindLocationName(ctx, inst.Client, kvPublicIp["zone_id"].(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(kvPublicIp["id"].(string))
	rd.Set("description", kvPublicIp["description"].(string))
	rd.Set("ipv4", kvPublicIp["ipv4"].(string))
	rd.Set("region", location)

	return nil
}
