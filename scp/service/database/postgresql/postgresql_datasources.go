package postgresql

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/postgresql"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_postgresqls", DatasourcePostgresqls())
}

func DatasourcePostgresqls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePostgresqlList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"postgresql_cluster_name": {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"page":                    {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":                    {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort":                    {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":                {Type: schema.TypeList, Optional: true, Description: "PostgreSQL list", Elem: datasourcePostgresqlElem()},
			"total_count":             {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of Microsoft SQL Servers.",
	}
}

func datasourcePostgresqlElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":               {Type: schema.TypeString, Computed: true, Description: "Project ID."},
			"block_id":                 {Type: schema.TypeString, Computed: true, Description: "Block ID."},
			"service_zone_id":          {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"postgresql_cluster_id":    {Type: schema.TypeString, Computed: true, Description: "PostgreSQL Cluster ID"},
			"postgresql_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "PostgreSQL Cluster Name"},
			"postgresql_cluster_state": {Type: schema.TypeString, Computed: true, Description: "PostgreSQL Cluster State"},
			"created_by":               {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":               {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":              {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":              {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourcePostgresqlList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Postgresql.ListPostgresqlClusters(ctx, &postgresql.PostgresqlSearchApiListPostgresqlClustersOpts{
		PostgresqlClusterName: optional.NewString(rd.Get("postgresql_cluster_name").(string)),
		Page:                  optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                  optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                  optional.NewInterface(rd.Get("sort").(string)),
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
