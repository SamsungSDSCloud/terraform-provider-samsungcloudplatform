package internetgateway

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func init() {
	scp.RegisterResource("scp_internet_gateway", ResourceInternetGateway())
}

func ResourceInternetGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInternetGatewayCreate,
		ReadContext:   resourceInternetGatewayRead,
		UpdateContext: resourceInternetGatewayUpdate,
		DeleteContext: resourceInternetGatewayDelete,
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
			"igw_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Internet-Gateway Type. One of SHARED , DEDICATED, SHARED_GROUP, SECURE",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SHARED", "DEDICATED", "SHARED_GROUP", "SECURE"}, false)),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Internet-Gateway description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Internet Gateway resource.",
	}
}

func resourceInternetGatewayCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()
	vpcId := rd.Get("vpc_id").(string)
	description := rd.Get("description").(string)
	igwType := rd.Get("igw_type").(string)
	tags := rd.Get("tags").(map[string]interface{})

	inst := meta.(*client.Instance)

	result, _, err := inst.Client.InternetGateway.CreateInternetGateway(ctx, vpcId, igwType, description, false, false, tags)
	if err != nil {
		return
	}

	err = waitForInternetGatewayStatus(ctx, inst.Client, result.ResourceId, []string{"ATTACHING"}, []string{"ATTACHED"}, true)
	if err != nil {
		return
	}

	rd.SetId(result.ResourceId)

	return resourceInternetGatewayRead(ctx, rd, meta)
}

func resourceInternetGatewayRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	info, _, err := inst.Client.InternetGateway.GetInternetGateway(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("vpc_id", info.VpcId)
	rd.Set("description", info.InternetGatewayDescription)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceInternetGatewayUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges("description") {
		_, _, err := inst.Client.InternetGateway.UpdateInternetGateway(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForInternetGatewayStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ATTACHED"}, true)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceInternetGatewayRead(ctx, rd, meta)
}

func resourceInternetGatewayDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.InternetGateway.DeleteInternetGateway(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	time.Sleep(10 * time.Second)

	err = waitForInternetGatewayStatus(ctx, inst.Client, rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForInternetGatewayStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.InternetGateway.GetInternetGateway(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.InternetGatewayState, nil
	})
}
