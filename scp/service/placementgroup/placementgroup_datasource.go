package placementgroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/placementgroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	placementgroup2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/placement-group"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func init() {
	scp.RegisterDataSource("scp_placement_groups", DatasourcePlacementGroup())
}

func DatasourcePlacementGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":               common.DatasourceFilter(),
			"placement_group_id":   {Type: schema.TypeString, Optional: true, Description: "Placement Group Id"},
			"placement_group_name": {Type: schema.TypeString, Optional: true, Description: "Placement Group Name"},
			"service_zone_id":      {Type: schema.TypeString, Optional: true, Description: "Service Zone Id"},
			"virtual_server_type":  {Type: schema.TypeString, Optional: true, Description: "Virtual Server Type"},
			"created_by":           {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"page":                 {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                 {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                 {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":             {Type: schema.TypeList, Computed: true, Description: "Virtual Server list", Elem: datasourceElem()},
			"total_count":          {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Placement Groups",
	}
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":                  {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			"availability_zone_name":      {Type: schema.TypeString, Computed: true, Description: "Availability Zone Name"},
			"placement_group_id":          {Type: schema.TypeString, Computed: true, Description: "Placement Group Id"},
			"placement_group_name":        {Type: schema.TypeString, Computed: true, Description: "Placement Group  Name"},
			"virtual_server_type":         {Type: schema.TypeString, Computed: true, Description: "Virtual Server Type"},
			"service_zone_id":             {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			"placement_group_description": {Type: schema.TypeString, Computed: true, Description: "Description"},
			"placement_group_state":       {Type: schema.TypeString, Computed: true, Description: "Placement Group State"},
			"virtual_server_id_list": {Type: schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Virtual Server Ids",
				Computed:    true,
			},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
		Description: "Placement Group Element",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	placementGroupId := rd.Get("placement_group_id").(string)
	contents := make([]map[string]interface{}, 0)
	if strings.Compare(placementGroupId, "") == 0 {

		responses, err := inst.Client.PlacementGroup.ListPlacementGroups(ctx, *getListPlacementGroupsRequestParam(rd))
		if err != nil {
			return diag.FromErr(err)
		}

		var diagnostics diag.Diagnostics
		var done bool
		for _, response := range responses.Contents {
			contents, diagnostics, done = appendContentFromDetailPlacementGroupApi(contents, ctx, inst, response.PlacementGroupId)
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
		contents, diagnostics, done := appendContentFromDetailPlacementGroupApi(contents, ctx, inst, placementGroupId)
		if done {
			return diagnostics
		}

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", 1)

	}
	return nil

	return nil
}

func getListPlacementGroupsRequestParam(rd *schema.ResourceData) *placementgroup.ListPlacementGroupsRequestParam {
	return &placementgroup.ListPlacementGroupsRequestParam{
		PlacementGroupName: rd.Get("placement_group_name").(string),
		ServiceZoneId:      rd.Get("service_zone_id").(string),
		VirtualServerType:  rd.Get("virtual_server_type").(string),
		CreatedBy:          rd.Get("created_by").(string),
		Page:               int32(rd.Get("page").(int)),
		Size:               int32(rd.Get("size").(int)),
		Sort:               rd.Get("sort").(string),
	}
}

func appendContentFromDetailPlacementGroupApi(contents []map[string]interface{}, ctx context.Context, inst *client.Instance, placementGroupId string) ([]map[string]interface{}, diag.Diagnostics, bool) {
	content, err := getContentFromDetailPlacementGroupApi(ctx, inst, placementGroupId)
	if err != nil {
		return nil, diag.FromErr(err), true
	}
	contents = append(contents, content)
	return contents, nil, false
}

func getContentFromDetailPlacementGroupApi(ctx context.Context, inst *client.Instance, placementGroupId string) (map[string]interface{}, error) {
	detailResponse, err := inst.Client.PlacementGroup.DetailPlacementGroup(ctx, placementGroupId)
	if err != nil {
		return nil, err
	}
	content := getContentMap(detailResponse)
	return content, err
}

func getContentMap(responses placementgroup2.PlacementGroupDetailResponse) map[string]interface{} {
	content := common.ToMap(responses)
	content["virtual_server_id_list"] = responses.VirtualServerIdList
	return content
}
