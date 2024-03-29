package sqlserver

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/sqlserver2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_sqlservers", DatasourceSqlServers())
}

func DatasourceSqlServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlServerList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"db_name":             {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"region":              {Type: schema.TypeString, Optional: true, Description: "Region."},
			"server_group_name":   {Type: schema.TypeString, Optional: true, Description: "Server group name."},
			"virtual_server_name": {Type: schema.TypeString, Optional: true, Description: "Virtual server name."},
			"created_by":          {Type: schema.TypeString, Optional: true, Description: "Creator."},
			"page":                {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":                {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Sorting conditions",
			},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "Microsoft SQL server list", Elem: datasourceSqlServerElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of Microsoft SQL Servers.",
	}
}

func getKeyString(rd *schema.ResourceData, key string) optional.String {
	if len(rd.Get(key).(string)) > 0 {
		return optional.NewString(rd.Get(key).(string))
	} else {
		return optional.String{}
	}
}

func dataSourceSqlServerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.SqlServer.ListSqlServer(ctx, &sqlserver2.MsSqlConfigurationControllerApiListSqlserverOpts{
		DbName:            getKeyString(rd, "db_name"),
		Region:            getKeyString(rd, "region"),
		ServerGroupName:   getKeyString(rd, "server_group_name"),
		VirtualServerName: getKeyString(rd, "virtual_server_name"),
		CreatedBy:         getKeyString(rd, "created_by"),
		Page:              optional.NewInt32(0),
		Size:              optional.NewInt32(1000),
		Sort:              optional.Interface{},
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

func datasourceSqlServerVirtualServerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"virtual_server_id":       {Type: schema.TypeString, Computed: true, Description: "Virtual server ID."},
			"virtual_server_name":     {Type: schema.TypeString, Computed: true, Description: "Virtual server name."},
			"database_state":          {Type: schema.TypeString, Computed: true, Description: "Database state."},
			"data_block_storage_spec": {Type: schema.TypeString, Computed: true, Description: "Data block storage specifications."},
			"software": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Software information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"software_id":            {Type: schema.TypeString, Computed: true, Description: "Software id."},
						"software_name":          {Type: schema.TypeString, Computed: true, Description: "Software name."},
						"software_category":      {Type: schema.TypeString, Computed: true, Description: "Software category. (ex: DATABASE)"},
						"software_type":          {Type: schema.TypeString, Computed: true, Description: "Software type. (ex: MariaDB)"},
						"software_version":       {Type: schema.TypeString, Computed: true, Description: "Software version."},
						"software_state":         {Type: schema.TypeString, Computed: true, Description: "Software state."},
						"software_service_state": {Type: schema.TypeString, Computed: true, Description: "Software service state."},
						"software_properties":    {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "List of software properties."},
						"created_by":             {Type: schema.TypeString, Computed: true, Description: "The person who created the resource."},
						"created_dt":             {Type: schema.TypeString, Computed: true, Description: "Creation date"},
						"modified_by":            {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
						"modified_dt":            {Type: schema.TypeString, Computed: true, Description: "Modification date"},
					},
				},
			},
			"region":      {Type: schema.TypeString, Computed: true, Description: "Region."},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "The person who created the resource."},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func datasourceSqlServerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_group_id":    {Type: schema.TypeString, Computed: true, Description: "Server group id."},
			"server_group_name":  {Type: schema.TypeString, Computed: true, Description: "Server group name."},
			"db_name":            {Type: schema.TypeString, Computed: true, Description: "Database name."},
			"server_group_state": {Type: schema.TypeString, Computed: true, Description: "Server group tate."},
			"virtual_servers":    {Type: schema.TypeList, Computed: true, Description: "List of virtual servers.", Elem: datasourceSqlServerVirtualServerElem()},
			"region":             {Type: schema.TypeString, Computed: true, Description: "Region."},
			"created_by":         {Type: schema.TypeString, Computed: true, Description: "The person who created the resource."},
			"created_dt":         {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":        {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":        {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}
