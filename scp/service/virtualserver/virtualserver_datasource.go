package virtualserver

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	virtualserver2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/virtual-server2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func init() {
	scp.RegisterDataSource("scp_virtual_servers", DatasourceVirtualServer())
}

func DatasourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":                common.DatasourceFilter(),
			"virtual_server_id":     {Type: schema.TypeString, Optional: true, Description: "Virtual server id"},
			"virtual_server_name":   {Type: schema.TypeString, Optional: true, Description: "Virtual Server Name"},
			"auto_scaling_enabled":  {Type: schema.TypeBool, Optional: true, Description: "Auto Scaling Enabled"},
			"server_group_id":       {Type: schema.TypeString, Optional: true, Description: "Server Group Id"},
			"auto_scaling_group_id": {Type: schema.TypeString, Optional: true, Description: "Auto Scaling Group Id"},
			//"serviced_for_list":       {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Serviced For"}, Description: "Serviced For List"},
			//"serviced_group_for_list": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Serviced Group For"}, Description: "Serviced Group For List"},
			"page":        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":        {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Virtual Server list", Elem: datasourceElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Virtual Servers",
	}
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			//"filter":                      common.DatasourceFilter(),
			"project_id":                  {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			"autoscaling_enabled":         {Type: schema.TypeBool, Computed: true, Description: "Auto Scaling Enabled"},
			"availability_zone_name":      {Type: schema.TypeString, Computed: true, Description: "Availability Zone Name"},
			"block_id":                    {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			"block_storage_ids":           {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Block Id"}, Description: "Block Storage Ids"},
			"contract":                    {Type: schema.TypeString, Computed: true, Description: "Contract"},
			"contract_end_date":           {Type: schema.TypeString, Computed: true, Description: "Contract End Date"},
			"contract_id":                 {Type: schema.TypeString, Computed: true, Description: "Contract Id"},
			"contract_start_date":         {Type: schema.TypeString, Computed: true, Description: "Contract Start Date"},
			"deletion_protection_enabled": {Type: schema.TypeBool, Computed: true, Description: "Deletion Protection Enabled"},
			"dns_enabled":                 {Type: schema.TypeBool, Computed: true, Description: "Dns Enabled"},
			"encrypt_enabled":             {Type: schema.TypeBool, Computed: true, Description: "Encrypt Enabled"},
			"key_pair_id":                 {Type: schema.TypeString, Computed: true, Description: "Key Pair Id"},
			"image_id":                    {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			"initial_script_content":      {Type: schema.TypeString, Computed: true, Description: "Initial Script Content"},
			"ip":                          {Type: schema.TypeString, Computed: true, Description: "Ip"},
			"is_dr":                       {Type: schema.TypeBool, Computed: true, Description: "Is Dr"},
			"next_contract_end_date":      {Type: schema.TypeString, Computed: true, Description: "Next Contract End Date"},
			"next_contract_id":            {Type: schema.TypeString, Computed: true, Description: "Next Contract Id"},
			"nic_ids":                     {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Nic Id"}, Description: "Nic Id List"},
			"os_type":                     {Type: schema.TypeString, Computed: true, Description: "Os Type"},
			"placement_group_id":          {Type: schema.TypeString, Computed: true, Description: "Placement Group Id"},
			"product_group_id":            {Type: schema.TypeString, Computed: true, Description: "Product Group Id"},
			"security_group_ids": {Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id":           {Type: schema.TypeString, Computed: true, Description: "Security Group Id"},
						"security_group_member_state": {Type: schema.TypeString, Computed: true, Description: "Security Group Member State"},
					},
				},
				Description: "Security Group Ids",
				Computed:    true,
			},
			"server_group_id":      {Type: schema.TypeString, Computed: true, Description: "Server Group Id"},
			"server_type":          {Type: schema.TypeString, Computed: true, Description: "Server Type"},
			"server_type_id":       {Type: schema.TypeString, Computed: true, Description: "Server Type Id"},
			"service_zone_id":      {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			"serviced_for":         {Type: schema.TypeString, Computed: true, Description: "Serviced For"},
			"serviced_group_for":   {Type: schema.TypeString, Computed: true, Description: "Serviced Group For"},
			"virtual_server_dr_id": {Type: schema.TypeString, Computed: true, Description: "Virtual Server Dr Id"},
			"virtual_server_id":    {Type: schema.TypeString, Computed: true, Description: "Virtual Server Id"},
			"virtual_server_name":  {Type: schema.TypeString, Computed: true, Description: "Virtual Server Name"},
			"virtual_server_state": {Type: schema.TypeString, Computed: true, Description: "Virtual Server State"},
			"vpc_id":               {Type: schema.TypeString, Computed: true, Description: "Vpc Id"},
			"created_by":           {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":           {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":          {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":          {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
		Description: "Virtual Server Element",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	virtualServerId := rd.Get("virtual_server_id").(string)
	contents := make([]map[string]interface{}, 0)
	if strings.Compare(virtualServerId, "") == 0 {

		responses, err := inst.Client.VirtualServer.ListVirtualServers(ctx, *getListVirtualServersRequestParam(rd))
		if err != nil {
			return diag.FromErr(err)
		}

		var diagnostics diag.Diagnostics
		var done bool
		for _, response := range responses.Contents {
			contents, diagnostics, done = appendContentFromDetailVirtualServerApi(contents, ctx, inst, response.VirtualServerId)
			if done {
				return diagnostics
			}
		}
		if f, ok := rd.GetOk("filter"); ok {
			contents = common.ApplyFilter(datasourceElem().Schema, f.(*schema.Set), contents)
		}

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", len(contents))

	} else {
		// detail api
		contents, diagnostics, done := appendContentFromDetailVirtualServerApi(contents, ctx, inst, virtualServerId)
		if done {
			return diagnostics
		}

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", 1)

	}
	return nil
}

func getListVirtualServersRequestParam(rd *schema.ResourceData) *virtualserver.ListVirtualServersRequestParam {
	return &virtualserver.ListVirtualServersRequestParam{
		AutoscalingEnabled: getBoolPtrFromRd(rd, "auto_scaling_enabled"),
		ServerGroupId:      rd.Get("server_group_id").(string),
		VirtualServerName:  rd.Get("virtual_server_name").(string),
		AutoScalingGroupId: rd.Get("auto_scaling_group_id").(string),
		//ServicedForList:      convertToStringArray(rd.Get("serviced_for_list").([]interface{})),
		//ServicedGroupForList: convertToStringArray(rd.Get("serviced_group_for_list").([]interface{})),
		Page: int32(rd.Get("page").(int)),
		Size: int32(rd.Get("size").(int)),
		Sort: rd.Get("sort").(string),
	}
}

func getBoolPtrFromRd(rd *schema.ResourceData, key string) *bool {
	boolValueFromRd := rd.GetRawConfig().GetAttr(key)
	if !boolValueFromRd.IsNull() {
		boolVal := boolValueFromRd.True()
		return &boolVal
	}

	return nil
}

func convertToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, 0)
	for _, interfaceElem := range interfaceArray {
		stringArray = append(stringArray, interfaceElem.(string))
	}
	return stringArray
}

func appendContentFromDetailVirtualServerApi(contents []map[string]interface{}, ctx context.Context, inst *client.Instance, virtualServerId string) ([]map[string]interface{}, diag.Diagnostics, bool) {
	content, err := getContentFromDetailVirtualServerApi(ctx, inst, virtualServerId)
	if err != nil {
		return nil, diag.FromErr(err), true
	}
	contents = append(contents, content)
	return contents, nil, false
}

func getContentFromDetailVirtualServerApi(ctx context.Context, inst *client.Instance, virtualServerId string) (map[string]interface{}, error) {
	detailResponse, err := inst.Client.VirtualServer.DetailVirtualServer(ctx, virtualServerId)
	if err != nil {
		return nil, err
	}
	content := getContentMap(detailResponse)
	validContentMap := getContentMapMatchedWithSchemaAttr(content)
	return validContentMap, err
}

func getContentMapMatchedWithSchemaAttr(content map[string]interface{}) map[string]interface{} {
	resourceData := datasourceElem()
	validContentMap := map[string]interface{}{}
	for key, value := range content {
		if _, ok := resourceData.Schema[key]; ok {
			validContentMap[key] = value
		}
	}
	return validContentMap
}

func getContentMap(responses virtualserver2.DetailVirtualServerV3Response) map[string]interface{} {
	content := common.ToMap(responses)
	securityGroupIds := common.ConvertStructToMaps(responses.SecurityGroupIds)
	content["security_group_ids"] = securityGroupIds
	content["nic_ids"] = responses.NicIds
	content["block_storage_ids"] = responses.BlockStorageIds
	return content
}
