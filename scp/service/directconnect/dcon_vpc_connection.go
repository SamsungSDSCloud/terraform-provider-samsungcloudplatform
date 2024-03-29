package directconnect

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_dcon_vpc_connection", ResourceDconVpcConnection())
}

func ResourceDconVpcConnection() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceDconVpcConnectionCreate,
		ReadContext:   resourceDconVpcConnectionRead,
		//UpdateContext: resourceDconVpcConnectionUpdate,
		DeleteContext: resourceDconVpcConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Vpc id of approver",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"firewall_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall enabled",
			},
			"direct_connect_id": {
				Type:         schema.TypeString,
				Required:     true, // optional
				ForceNew:     true,
				Description:  "Direct connect id of requester",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Dcon-Vpc connection description. (0 to 100 characters)",
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
		},
		Description: "Provides a Dcon-Vpc connection resource.",
	}
}

func resourceDconVpcConnectionCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	approverVpcId := rd.Get("vpc_id").(string)
	firewallEnabled := rd.Get("firewall_enabled").(bool)
	requesterDcId := rd.Get("direct_connect_id").(string)
	connectionDescription := rd.Get("description").(string)
	connectionType := "INTERNAL"

	inst := meta.(*client.Instance)

	//Search ->  approverProjectId
	approverVpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, approverVpcId)
	if err != nil {
		return diag.FromErr(err)
	}
	//Search ->  requesterProjectId
	requestVpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, approverVpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err := inst.Client.DirectConnect.CreateDconVpcConnection(ctx, approverVpcInfo.ProjectId, approverVpcId, connectionType, firewallEnabled, requesterDcId, requestVpcInfo.ProjectId, connectionDescription)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitDconVpcConnectionCreating(ctx, inst.Client, response.ResourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	// Refresh
	return resourceDconVpcConnectionRead(ctx, rd, meta)
}

func resourceDconVpcConnectionRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.DirectConnect.GetDconVpcConnectionInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", info.DirectConnectConnectionName)
	rd.Set("type", info.DirectConnectConnectionType)
	rd.Set("description", info.DirectConnectConnectionDescription)
	rd.Set("state", info.DirectConnectConnectionState)

	return nil
}

func resourceDconVpcConnectionDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.DirectConnect.DeleteDconVpcConnection(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitDconVpcConnectionDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitDconVpcConnectionCreating(ctx context.Context, scpClient *client.SCPClient, dconVpcConId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		dconVpcConInfo, _, err := scpClient.DirectConnect.GetDconVpcConnectionInfo(ctx, dconVpcConId)
		if err != nil {
			return nil, "", err
		}
		return dconVpcConInfo, dconVpcConInfo.DirectConnectConnectionState, nil
	})
}

func waitDconVpcConnectionDeleting(ctx context.Context, scpClient *client.SCPClient, dconVpcConId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		dconVpcConInfo, c, err := scpClient.DirectConnect.GetDconVpcConnectionInfo(ctx, dconVpcConId)
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
		return dconVpcConInfo, dconVpcConInfo.DirectConnectConnectionState, nil
	})
}
