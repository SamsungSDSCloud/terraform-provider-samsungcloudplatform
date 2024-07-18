package autoscaling

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_auto_scaling_group_virtual_servers", DataSourceAsgVirtualServers())
}

func DataSourceAsgVirtualServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAsgVirtualServerList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"asg_id":      {Type: schema.TypeString, Required: true, Description: "Auto-Scaling Group ID"},
			"page":        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":        {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Auto-Scaling Group policy list", Elem: dataSourceAsgVirtualServerElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Auto-Scaling Group Virtual Servers",
	}
}

func dataSourceAsgVirtualServerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter":               common.DatasourceFilter(),
			"virtual_server_id":    {Type: schema.TypeString, Computed: true, Description: "Virtual Server Id"},
			"virtual_server_name":  {Type: schema.TypeString, Computed: true, Description: "Virtual Server Name"},
			"virtual_server_state": {Type: schema.TypeString, Computed: true, Description: "Virtual Server State"},
			"ip":                   {Type: schema.TypeString, Computed: true, Description: "Ip"},
			"server_group_id":      {Type: schema.TypeString, Computed: true, Description: "Server Group Id"},
			"image_id":             {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			"serviced_for":         {Type: schema.TypeString, Computed: true, Description: "Serviced For"},
			"serviced_group_for":   {Type: schema.TypeString, Computed: true, Description: "Serviced Group For"},
			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
				MinItems: 0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Properties",
			},
			"file_storage_linked_state": {Type: schema.TypeString, Computed: true, Description: "File Storage Linked State"},
			"created_by":                {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":                {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":               {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":               {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
		Description: "Auto-Scaling Group Virtual Server Element",
	}
}

func dataSourceAsgVirtualServerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	response, _, err := inst.Client.AutoScaling.GetAutoScalingGroupVirtualServerList(ctx, rd.Get("asg_id").(string), &autoscaling2.AsgVirtualServerV2ApiGetAsgVirtualServerListV2Opts{
		Page: optional.NewInt32((int32)(rd.Get("page").(int))),
		Size: optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort: optional.NewInterface([]string{rd.Get("sort").(string)}),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(response.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DataSourceAsgVirtualServers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}
