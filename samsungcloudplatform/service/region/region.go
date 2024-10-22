package region

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_region", DatasourceRegion())
}

func DatasourceRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRegionRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of this region",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of this region",
			},
			"block_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block id of this region",
			},
		},
		Description: "Provides region details",
	}
}

func convertRegionListToHclSet(regions []project.ZoneResponseV3) (common.HclSetObject, []string) {
	var setRegions common.HclSetObject
	var ids []string
	// Convert to HclSet
	for _, z := range regions {
		if len(z.ServiceZoneId) == 0 {
			continue
		}
		ids = append(ids, z.ServiceZoneId)
		kv := common.HclKeyValueObject{
			"id":       z.ServiceZoneId,
			"name":     z.ServiceZoneName,
			"location": z.ServiceZoneLocation,
			"block_id": z.BlockId,
		}
		setRegions = append(setRegions, kv)

	}
	return setRegions, ids
}

func datasourceRegionRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
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

	setRegions, _ := convertRegionListToHclSet(projectInfo.ServiceZones)

	if f, ok := rd.GetOk("filter"); ok {
		setRegions = common.ApplyFilter(DatasourceRegions().Schema, f.(*schema.Set), setRegions)
	}

	if len(setRegions) == 0 {
		err = fmt.Errorf("no matching region data found for project")
		return
	}

	//set default region
	defaultRegion := setRegions[0]

	//set region from project info
	if len(setRegions) > 1 && projectInfo.DefaultZoneId != "" {
		for _, value := range setRegions {
			for k, v := range value {
				if k == "id" && v == projectInfo.DefaultZoneId {
					defaultRegion = value
					break
				}
			}
		}
	}

	//for k, v := range setRegions[0] {
	for k, v := range defaultRegion {
		if k == "id" {
			rd.SetId(v.(string))
		}
		rd.Set(k, v)
	}

	return nil
}
