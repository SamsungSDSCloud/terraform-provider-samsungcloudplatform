package sqlserver

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver2"
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

func dataSourceSqlServerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.SqlServer.ListSqlServer(ctx, &sqlserver2.MsSqlConfigurationControllerApiListSqlserverOpts{
		DbName:            common.GetKeyString(rd, "db_name"),
		Region:            common.GetKeyString(rd, "region"),
		ServerGroupName:   common.GetKeyString(rd, "server_group_name"),
		VirtualServerName: common.GetKeyString(rd, "virtual_server_name"),
		CreatedBy:         common.GetKeyString(rd, "created_by"),
		Page:              optional.NewInt32(0),
		Size:              optional.NewInt32(1000),
		Sort:              optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertSqlServerListToHclSet(responses)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func convertSqlServerListToHclSet(responses sqlserver2.ListResponseOfDbServerGroupsResponse) common.HclSetObject {
	var sqlserverList common.HclSetObject

	for _, sqlserver := range responses.Contents {
		if len(sqlserver.ServerGroupId) == 0 {
			continue
		}

		var virtualServersList common.HclListObject
		for _, virtualServer := range sqlserver.VirtualServers {
			var softwareList common.HclListObject
			softwareKv := common.HclKeyValueObject{
				"software_id":            virtualServer.Software.SoftwareId,
				"software_name":          virtualServer.Software.SoftwareName,
				"software_category":      virtualServer.Software.SoftwareCategory,
				"software_type":          virtualServer.Software.SoftwareType,
				"software_version":       virtualServer.Software.SoftwareVersion,
				"software_state":         virtualServer.Software.SoftwareState,
				"software_service_state": virtualServer.Software.SoftwareServiceState,
				"created_by":             virtualServer.Software.CreatedBy,
				"created_dt":             virtualServer.Software.CreatedDt.String(),
				"modified_by":            virtualServer.Software.ModifiedBy,
				"modified_dt":            virtualServer.Software.ModifiedDt.String(),
			}
			softwareList = append(softwareList, softwareKv)

			virtualServerKv := common.HclKeyValueObject{
				"virtual_server_id":       virtualServer.VirtualServerId,
				"virtual_server_name":     virtualServer.VirtualServerName,
				"database_state":          virtualServer.DatabaseState,
				"data_block_storage_spec": virtualServer.DataBlockStorageSpec,
				"region":                  virtualServer.Region,
				"software":                softwareList,
				"created_by":              virtualServer.CreatedBy,
				"created_dt":              virtualServer.CreatedDt.String(),
				"modified_by":             virtualServer.ModifiedBy,
				"modified_dt":             virtualServer.ModifiedDt.String(),
			}
			virtualServersList = append(virtualServersList, virtualServerKv)
		}

		kv := common.HclKeyValueObject{
			"server_group_id":    sqlserver.ServerGroupId,
			"server_group_name":  sqlserver.ServerGroupName,
			"db_name":            sqlserver.DbName,
			"server_group_state": sqlserver.ServerGroupState,
			"virtual_servers":    virtualServersList,
			"region":             sqlserver.Region,
			"created_by":         sqlserver.CreatedBy,
			"created_dt":         sqlserver.CreatedDt.String(),
			"modified_by":        sqlserver.ModifiedBy,
			"modified_dt":        sqlserver.ModifiedDt.String(),
		}
		sqlserverList = append(sqlserverList, kv)
	}
	return sqlserverList
}
