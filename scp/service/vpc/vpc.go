package vpc

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_vpc", ResourceVpc())
}

func ResourceVpc() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcCreate,
		ReadContext:   resourceVpcRead,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true, //필수 작성
				ForceNew:    true, //해당 필드 수정시 자원 삭제 후 다시 생성됨
				Description: "VPC name. (3 to 20 characters without specials)",
				ValidateFunc: validation.All( //입력 값 Validation 체크
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters"),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true, //선택 입력
				Description:  "VPC description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
		},
		Description: "Provides a VPC resource.",
	}
}

func resourceVpcCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	vpcName := rd.Get("name").(string)
	vpcDescription := rd.Get("description").(string)
	vpcLocation := rd.Get("region").(string)

	inst := meta.(*client.Instance)

	isVpcNameInvalid, err := inst.Client.Vpc.CheckVpcName(ctx, vpcName)
	if err != nil {
		return diag.FromErr(err)
	}
	if isVpcNameInvalid {
		return diag.Errorf("Input vpc name is invalid (maybe duplicated) : " + vpcName)
	}

	serviceZoneId, productGroupId, err := client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, vpcLocation, common.NetworkProductGroup, common.VpcProductName)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Try create vpc : "+vpcName+", "+vpcDescription+", "+serviceZoneId+", "+productGroupId)

	response, err := inst.Client.Vpc.CreateVpc(ctx, vpcName, vpcDescription, productGroupId, serviceZoneId)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitVpcCreating(ctx, inst.Client, response.ResourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	// Refresh
	return resourceVpcRead(ctx, rd, meta)
}

func resourceVpcRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", vpcInfo.VpcName)
	rd.Set("description", vpcInfo.VpcDescription)

	location, err := client.FindLocationName(ctx, inst.Client, vpcInfo.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information")
	}

	rd.Set("region", location)

	return nil
}
func resourceVpcUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("description") {
		_, err := inst.Client.Vpc.UpdateVpc(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpcRead(ctx, rd, meta)
}
func resourceVpcDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	err := inst.Client.Vpc.DeleteVpc(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitVpcDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitVpcCreating(ctx context.Context, scpClient *client.SCPClient, vpcId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		vpcInfo, _, err := scpClient.Vpc.GetVpcInfo(ctx, vpcId)
		if err != nil {
			return nil, "", err
		}
		return vpcInfo, vpcInfo.VpcState, nil
	})
}

func waitVpcDeleting(ctx context.Context, scpClient *client.SCPClient, vpcId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		vpcInfo, c, err := scpClient.Vpc.GetVpcInfo(ctx, vpcId)
		if err != nil {
			if c == 404 {
				return "", "DELETED", nil
			}
			// VPC may return 403 for deleted resources
			if c == 403 {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return vpcInfo, vpcInfo.VpcState, nil
	})
}
