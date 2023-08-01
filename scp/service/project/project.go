package project

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project name",
			},
			"project_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project state",
			},
			"account_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account code",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account id",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account name",
			},
			"account_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account type",
			},
			"service_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Service zone id",
				},
				Description: "Service zone id list",
			},
		},
	}
}

func datasourceProjectRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(projectInfo.ServiceZones) == 0 {
		return diag.Errorf("This project does not have any valid service zones")
	}

	rd.SetId(common.GenerateHash([]string{projectInfo.ProjectId}))
	rd.Set("project_id", projectInfo.ProjectId)
	rd.Set("project_name", projectInfo.ProjectName)
	rd.Set("project_state", projectInfo.ProjectState)
	//rd.Set("account_code", projectInfo.AccountCode) 	// deprecated
	rd.Set("account_id", projectInfo.AccountId)
	rd.Set("account_name", projectInfo.AccountName)
	rd.Set("account_type", projectInfo.AccountType)
	var serviceZoneList []string
	for _, serviceZone := range projectInfo.ServiceZones {
		serviceZoneList = append(serviceZoneList, serviceZone.ServiceZoneId)
	}
	rd.Set("service_zones", serviceZoneList)

	return nil
}
