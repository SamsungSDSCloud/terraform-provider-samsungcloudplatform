package directconnect

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	directconnect2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/direct-connect2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_dcon_vpc_connections", DatasourceDconVpcConnections())
}

func DatasourceDconVpcConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dconVpcConnectionList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			common.ToSnakeCase("ApproverVpcId"):               {Type: schema.TypeString, Optional: true, Description: "Vpc id of approver"},
			common.ToSnakeCase("DirectConnectConnectionName"): {Type: schema.TypeString, Optional: true, Description: "Direct connect connection name"},
			common.ToSnakeCase("RequesterDirectConnectId"):    {Type: schema.TypeString, Optional: true, Description: "Direct connect id of requester"},
			common.ToSnakeCase("CreatedBy"):                   {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			common.ToSnakeCase("Page"):                        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):                        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                                        {Type: schema.TypeList, Optional: true, Description: "Direct Connect list", Elem: dconVpcConnectionElem()},
			"total_count":                                     {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of direct connect vpc connections.",
	}
}

func dconVpcConnectionList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.DirectConnect.GetDconVpcConnectionList(ctx, &directconnect2.DirectConnectConnectionOpenApiControllerApiListDirectConnectConnectionsOpts{
		ApproverVpcId:               optional.NewString(rd.Get("approver_vpc_id").(string)),
		DirectConnectConnectionName: optional.NewString(rd.Get("direct_connect_connection_name").(string)),
		RequesterDirectConnectId:    optional.NewString(rd.Get("requester_direct_connect_id").(string)),
		CreatedBy:                   optional.NewString(rd.Get("created_by").(string)),
		Page:                        optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                        optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                        optional.Interface{},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dconVpcConnectionElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                          {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ApproverProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Project id of approver"},
			common.ToSnakeCase("ApproverVpcId"):                      {Type: schema.TypeString, Computed: true, Description: "Vpc id of approver"},
			common.ToSnakeCase("CompletedDt"):                        {Type: schema.TypeString, Computed: true, Description: "Complete date"},
			common.ToSnakeCase("DirectConnectConnectionId"):          {Type: schema.TypeString, Computed: true, Description: "DirectConnect connection id"},
			common.ToSnakeCase("DirectConnectConnectionName"):        {Type: schema.TypeString, Computed: true, Description: "DirectConnect connection name"},
			common.ToSnakeCase("DirectConnectConnectionState"):       {Type: schema.TypeString, Computed: true, Description: "DirectConnect connection state"},
			common.ToSnakeCase("DirectConnectConnectionType"):        {Type: schema.TypeString, Computed: true, Description: "DirectConnect connection type"},
			common.ToSnakeCase("RequesterDirectConnectId"):           {Type: schema.TypeString, Computed: true, Description: "DirectConnect id of requester"},
			common.ToSnakeCase("RequesterProjectId"):                 {Type: schema.TypeString, Computed: true, Description: "Project id of requester"},
			common.ToSnakeCase("DirectConnectConnectionDescription"): {Type: schema.TypeString, Computed: true, Description: "Uplink enabled"},
			common.ToSnakeCase("CreatedBy"):                          {Type: schema.TypeString, Computed: true, Description: "DirectConnect connection description"},
			common.ToSnakeCase("CreatedDt"):                          {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):                         {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):                         {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}
