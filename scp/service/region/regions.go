package region

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRegionsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Region information list",
				Elem:        common.GetDatasourceItemsSchema(DatasourceRegion()),
			},
		},
		Description: "Provides list of regions",
	}
}

func datasourceRegionsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return
	}

	setRegions, ids := convertRegionListToHclSet(projectInfo.ServiceZones)

	if f, ok := rd.GetOk("filter"); ok {
		setRegions = common.ApplyFilter(DatasourceRegions().Schema, f.(*schema.Set), setRegions)
	}

	rd.SetId(common.GenerateHash(ids))
	rd.Set("ids", ids)
	rd.Set("regions", setRegions)

	return nil

}
