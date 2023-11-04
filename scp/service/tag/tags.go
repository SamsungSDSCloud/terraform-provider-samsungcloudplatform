package tag

import (
	"context"
	"encoding/json"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/tag"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_resource_tags", DatasourceResourceTags())
	//scp.RegisterDataSource("scp_tag_resources", DatasourceTagResources())
}

func DatasourceResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceResourceTagsList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"resource_id": {Type: schema.TypeString, Required: true, Description: "Resource Id"},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total count"},
			"contents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "List of tags",
			},
		},
		Description: "Provides a list of tags for a resource",
	}
}

func datasourceResourceTagsList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	resourceId := rd.Get("resource_id").(string)
	response, _, err := inst.Client.Tag.ListResourceTags(ctx, resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	results := common.ConvertStructToMaps(response.Contents)
	if f, ok := rd.GetOk("filter"); ok {
		results = common.ApplyFilter(DatasourceResourceTags().Schema, f.(*schema.Set), results)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", results)
	rd.Set("total_count", len(results))

	return nil
}

func DatasourceTagResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceTagResourcesList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource IDs",
			},
			"resource_type_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource type filters",
			},
			"tag_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key",
						},
						"tag_values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Tag values",
						},
					},
				},
			},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total count"},
			"contents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type",
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value",
									},
								},
							},
							Description: "Tags list",
						},
					},
				},
			},
		},

		Description: "Provides a list of resources for tag",
	}
}

func datasourceTagResourcesList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceIds := common.ToStringList(rd.Get("resource_ids").([]interface{}))
	resourceTypeFilters := common.ToStringList(rd.Get("resource_type_filters").([]interface{}))
	inTagFilters := rd.Get("tag_filters").([]interface{})

	myTagFilters := make([]tag.Filter, 0)
	for _, itf := range inTagFilters {
		var inInterface map[string]interface{}
		inrec, _ := json.Marshal(itf)
		json.Unmarshal(inrec, &inInterface)

		tagKey := inInterface["tag_key"].(string)
		interfaceTagValues := inInterface["tag_values"].([]interface{})
		stringTagValues := common.ToStringList(interfaceTagValues)

		myTagFilters = append(myTagFilters, tag.Filter{
			TagKey:    tagKey,
			TagValues: stringTagValues,
		})
	}

	inst := meta.(*client.Instance)
	response, _, err := inst.Client.Tag.ListResources(ctx, resourceIds, resourceTypeFilters, myTagFilters)
	if err != nil {
		return diag.FromErr(err)
	}

	// set
	results := common.ConvertStructToMaps(response.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", results)
	rd.Set("total_count", len(results))

	return nil
}
