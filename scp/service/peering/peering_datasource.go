package peering

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/peering"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_vpc_peerings", DataSourceVpcPeeringList())
	scp.RegisterDataSource("scp_vpc_peering_detail", DataSourceVpcPeeringDetail())
}

func DataSourceVpcPeeringList() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcPeeringListRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ApproverVpcId"):  {Type: schema.TypeString, Optional: true, Description: "Approver VPC Id"},
			common.ToSnakeCase("RequesterVpcId"): {Type: schema.TypeString, Optional: true, Description: "Requester VPC Id"},
			common.ToSnakeCase("VpcPeeringName"): {Type: schema.TypeString, Optional: true, Description: "VPC Peering Name"},
			common.ToSnakeCase("CreatedBy"):      {Type: schema.TypeString, Optional: true, Description: "Created By"},
			common.ToSnakeCase("Size"):           {Type: schema.TypeInt, Optional: true, Default: 50, Description: "Size"},
			common.ToSnakeCase("Page"):           {Type: schema.TypeInt, Optional: true, Description: "Page Number"},
			"contents":                           {Type: schema.TypeList, Optional: true, Description: "VPC Peering list", Elem: PeeringElem()},
			"total_count":                        {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a VPC Peering resource.",
	}
}

func PeeringElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ApproverProjectId"):     {Type: schema.TypeString, Computed: true, Description: "Approver Project Id"},
			common.ToSnakeCase("ApproverVpcId"):         {Type: schema.TypeString, Computed: true, Description: "Approver Vpc Id"},
			common.ToSnakeCase("Automated"):             {Type: schema.TypeBool, Computed: true, Description: "Is Automated"},
			common.ToSnakeCase("CompletedDt"):           {Type: schema.TypeString, Computed: true, Description: "Complated Date"},
			common.ToSnakeCase("RequesterProjectId"):    {Type: schema.TypeString, Computed: true, Description: "Requester Project Id"},
			common.ToSnakeCase("RequesterVpcId"):        {Type: schema.TypeString, Computed: true, Description: "Requester Vpc Id"},
			common.ToSnakeCase("VpcPeeringId"):          {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Id"},
			common.ToSnakeCase("VpcPeeringName"):        {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Name"},
			common.ToSnakeCase("VpcPeeringState"):       {Type: schema.TypeString, Computed: true, Description: "Vpc Peering State"},
			common.ToSnakeCase("VpcPeeringType"):        {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Type"},
			common.ToSnakeCase("VpcPeeringDescription"): {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Description"},
			common.ToSnakeCase("CreatedBy"):             {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):             {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):            {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):            {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}

func resourceVpcPeeringListRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	request := peering.VpcPeeringListRequest{
		ApproverVpcId:  rd.Get(common.ToSnakeCase("VpcPeeringName")).(string),
		RequesterVpcId: rd.Get(common.ToSnakeCase("RequesterVpcId")).(string),
		VpcPeeringName: rd.Get(common.ToSnakeCase("VpcPeeringName")).(string),
		CreatedBy:      rd.Get(common.ToSnakeCase("CreatedBy")).(string),
		Page:           (int32)(rd.Get(common.ToSnakeCase("Page")).(int)),
		Size:           (int32)(rd.Get(common.ToSnakeCase("Size")).(int)),
	}

	responses, err := inst.Client.Peering.GetVpcPeeringList(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DataSourceVpcPeeringDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcPeeringDetailRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("VpcPeeringId"):             {Type: schema.TypeString, Required: true, Description: "Vpc Peering Id"},
			common.ToSnakeCase("ApproverProjectId"):        {Type: schema.TypeString, Computed: true, Description: "Approver Project Id"},
			common.ToSnakeCase("ApproverVpcId"):            {Type: schema.TypeString, Computed: true, Description: "Approver Vpc Id"},
			common.ToSnakeCase("ApproverFirewallEnabled"):  {Type: schema.TypeBool, Computed: true, Description: "Approver Firewall Enabled"},
			common.ToSnakeCase("CompletedDt"):              {Type: schema.TypeString, Computed: true, Description: "Complated Date"},
			common.ToSnakeCase("RequesterProjectId"):       {Type: schema.TypeString, Computed: true, Description: "Requester Project Id"},
			common.ToSnakeCase("RequesterVpcId"):           {Type: schema.TypeString, Computed: true, Description: "Requester Vpc Id"},
			common.ToSnakeCase("RequesterFirewallEnabled"): {Type: schema.TypeBool, Computed: true, Description: "Requester Firewall Enabled"},
			common.ToSnakeCase("VpcPeeringName"):           {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Name"},
			common.ToSnakeCase("VpcPeeringState"):          {Type: schema.TypeString, Computed: true, Description: "Vpc Peering State"},
			common.ToSnakeCase("VpcPeeringType"):           {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Type"},
			common.ToSnakeCase("VpcPeeringDescription"):    {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Description"},
			common.ToSnakeCase("CreatedBy"):                {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):               {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):               {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
			common.ToSnakeCase("ProjectId"):                {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			common.ToSnakeCase("ApprovedBy"):               {Type: schema.TypeString, Computed: true, Description: "Approved By"},
			common.ToSnakeCase("ApprovedDt"):               {Type: schema.TypeString, Computed: true, Description: "Approved Date"},
			common.ToSnakeCase("BlockId"):                  {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			common.ToSnakeCase("ProductGroupId"):           {Type: schema.TypeString, Computed: true, Description: "Product Group Id"},
			common.ToSnakeCase("RequestedBy"):              {Type: schema.TypeString, Computed: true, Description: "Requested By"},
			common.ToSnakeCase("RequestedDt"):              {Type: schema.TypeString, Computed: true, Description: "Requested Date"},
			common.ToSnakeCase("ServiceZoneId"):            {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
		},
		Description: "Provides a VPC Peering detail.",
	}
}

func resourceVpcPeeringDetailRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	ruleInfo, _, err := inst.Client.Peering.GetVpcPeeringDetail(ctx, rd.Get(common.ToSnakeCase("VpcPeeringId")).(string))
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set(common.ToSnakeCase("ProjectId"), ruleInfo.ProjectId)
	rd.Set(common.ToSnakeCase("ApprovedBy"), ruleInfo.ApprovedBy)
	rd.Set(common.ToSnakeCase("ApprovedDt"), ruleInfo.ApprovedDt)
	rd.Set(common.ToSnakeCase("ApproverProjectId"), ruleInfo.ApproverProjectId)
	rd.Set(common.ToSnakeCase("ApproverVpcFirewallEnabled"), ruleInfo.ApproverVpcFirewallEnabled)
	rd.Set(common.ToSnakeCase("ApproverVpcId"), ruleInfo.ApproverVpcId)
	rd.Set(common.ToSnakeCase("BlockId"), ruleInfo.BlockId)
	rd.Set(common.ToSnakeCase("CompletedDt"), ruleInfo.CompletedDt)
	rd.Set(common.ToSnakeCase("ProductGroupId"), ruleInfo.ProductGroupId)
	rd.Set(common.ToSnakeCase("RequestedBy"), ruleInfo.RequestedBy)
	rd.Set(common.ToSnakeCase("RequestedDt"), ruleInfo.RequestedDt)
	rd.Set(common.ToSnakeCase("RequesterProjectId"), ruleInfo.RequesterProjectId)
	rd.Set(common.ToSnakeCase("RequesterVpcFirewallEnabled"), ruleInfo.RequesterVpcFirewallEnabled)
	rd.Set(common.ToSnakeCase("RequesterVpcId"), ruleInfo.RequesterVpcId)
	rd.Set(common.ToSnakeCase("ServiceZoneId"), ruleInfo.ServiceZoneId)
	rd.Set(common.ToSnakeCase("VpcPeeringId"), ruleInfo.VpcPeeringId)
	rd.Set(common.ToSnakeCase("VpcPeeringName"), ruleInfo.VpcPeeringName)
	rd.Set(common.ToSnakeCase("VpcPeeringState"), ruleInfo.VpcPeeringState)
	rd.Set(common.ToSnakeCase("VpcPeeringType"), ruleInfo.VpcPeeringType)
	rd.Set(common.ToSnakeCase("VpcPeeringDescription"), ruleInfo.VpcPeeringDescription)
	rd.Set(common.ToSnakeCase("CreatedBy"), ruleInfo.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), ruleInfo.CreatedDt)
	rd.Set(common.ToSnakeCase("ModifiedBy"), ruleInfo.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), ruleInfo.ModifiedDt)

	return nil
}
