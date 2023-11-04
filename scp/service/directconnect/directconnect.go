package directconnect

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_direct_connect", ResourceDirectConnect())
}

func ResourceDirectConnect() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceDirectConnectCreate,
		ReadContext:   resourceDirectConnectRead,
		//UpdateContext: resourceDirectConnectUpdate,
		DeleteContext: resourceDirectConnectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DirectConnect name. (3 to 20 characters without specials)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters"),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "DirectConnect description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Bandwidth gbps. (1 or 10)",
			},
		},
		Description: "Provides a DirectConnect resource.",
	}
}

func resourceDirectConnectCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	dcName := rd.Get("name").(string)
	dcDescription := rd.Get("description").(string)
	dcLocation := rd.Get("region").(string)
	bandwidth := int32(rd.Get("bandwidth").(int))

	inst := meta.(*client.Instance)

	serviceZoneId, err := client.FindServiceZoneId(ctx, inst.Client, dcLocation)

	sbandwidth := fmt.Sprint(bandwidth)
	tflog.Debug(ctx, "Try create direct connect : "+dcName+","+dcDescription+","+sbandwidth)

	response, _, err := inst.Client.DirectConnect.CreateDirectConnect(ctx, bandwidth, dcName, serviceZoneId, dcDescription)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitDirectConnectCreating(ctx, inst.Client, response.ResourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	// Refresh
	return resourceDirectConnectRead(ctx, rd, meta)
}

func resourceDirectConnectRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	dcInfo, _, err := inst.Client.DirectConnect.GetDirectConnectInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", dcInfo.DirectConnectName)
	rd.Set("description", dcInfo.DirectConnectDescription)

	location, err := client.FindLocationName(ctx, inst.Client, dcInfo.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information")
	}

	rd.Set("region", location)

	return nil
}

func resourceDirectConnectDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.DirectConnect.DeleteDirectConnect(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitDirectConnectDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitDirectConnectCreating(ctx context.Context, scpClient *client.SCPClient, dcId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		dcInfo, _, err := scpClient.DirectConnect.GetDirectConnectInfo(ctx, dcId)
		if err != nil {
			return nil, "", err
		}
		return dcInfo, dcInfo.DirectConnectState, nil
	})
}

func waitDirectConnectDeleting(ctx context.Context, scpClient *client.SCPClient, dcId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		dcInfo, c, err := scpClient.DirectConnect.GetDirectConnectInfo(ctx, dcId)
		if err != nil {
			if c == 404 {
				return "", "DELETED", nil
			}
			// DirectConnect may return 403 for deleted resources
			if c == 403 {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return dcInfo, dcInfo.DirectConnectState, nil
	})
}
