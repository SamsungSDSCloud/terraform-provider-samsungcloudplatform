package endpoint

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
)

func init() {
	scp.RegisterResource("scp_endpoint", ResourceEndpoint())
}

func ResourceEndpoint() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceEndpointCreate,
		ReadContext:   resourceEndpointRead,
		UpdateContext: resourceEndpointUpdate,
		DeleteContext: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC id",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true, //필수 작성
				ForceNew:    true, //해당 필드 수정시 자원 삭제 후 다시 생성됨
				Description: "Endpoint name. (3 to 20 characters without specials)",
				ValidateFunc: validation.All( //입력 값 Validation 체크
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters"),
				),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint type ('VPC_DNS'|'FS'|'OBS'|'CONTAINER_REGISTRY')",
				ValidateFunc: validation.StringInSlice([]string{
					"VPC_DNS",
					"FS",
					"OBS",
					"CONTAINER_REGISTRY",
				}, false),
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Endpoint IP address",
				ValidateFunc: validation.IsIPAddress,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target object id",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true, //선택 입력
				Description:  "Endpoint description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a VPC resource.",
	}
}

func resourceEndpointCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	endpointIpAddress := rd.Get("ip_address").(string)
	endpointName := rd.Get("name").(string)
	endpointType := strings.ToUpper(rd.Get("type").(string))
	objectId := rd.Get("object_id").(string)
	vpcId := rd.Get("vpc_id").(string)
	endpointDescription := rd.Get("description").(string)
	endpointLocation := rd.Get("region").(string)
	tags := rd.Get("tags").(map[string]interface{})

	inst := meta.(*client.Instance)

	serviceZoneId, err := client.FindServiceZoneId(ctx, inst.Client, endpointLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Try create endpoint : "+endpointIpAddress+", "+endpointName+", "+endpointType+", "+objectId+", "+vpcId+", "+endpointDescription+", "+serviceZoneId)

	response, err := inst.Client.Endpoint.CreateEndpoint(ctx, endpointIpAddress, endpointName, endpointType, objectId, vpcId, endpointDescription, serviceZoneId, tags)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitEndpointCreating(ctx, inst.Client, response.ResourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	// Refresh
	return resourceEndpointRead(ctx, rd, meta)
}

func resourceEndpointRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	endpointInfo, _, err := inst.Client.Endpoint.GetEndpoint(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	location, err := client.FindLocationName(ctx, inst.Client, endpointInfo.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information")
	}

	rd.Set("vpc_id", endpointInfo.VpcId)
	rd.Set("ip_address", endpointInfo.EndpointIpAddress)
	rd.Set("name", endpointInfo.EndpointName)
	rd.Set("type", endpointInfo.EndpointType)
	rd.Set("object_id", endpointInfo.ObjectId)
	rd.Set("description", endpointInfo.EndpointDescription)
	rd.Set("region", location)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceEndpointUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("description") {
		_, err := inst.Client.Endpoint.UpdateEndpoint(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEndpointRead(ctx, rd, meta)
}
func resourceEndpointDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	err := inst.Client.Endpoint.DeleteEndpoint(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitEndpointDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitEndpointCreating(ctx context.Context, scpClient *client.SCPClient, endpointId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		endpointInfo, _, err := scpClient.Endpoint.GetEndpoint(ctx, endpointId)
		if err != nil {
			return nil, "", err
		}
		return endpointInfo, endpointInfo.EndpointState, nil
	})
}

func waitEndpointDeleting(ctx context.Context, scpClient *client.SCPClient, endpointId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		endpointInfo, c, err := scpClient.Endpoint.GetEndpoint(ctx, endpointId)
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
		return endpointInfo, endpointInfo.EndpointState, nil
	})
}
