package backup

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/backup2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_backups", DatasourceBackups())
}

func DatasourceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"backup_name":                 {Type: schema.TypeString, Optional: true, Description: "Backup Name"},
			"created_by":                  {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"backup_policy_type_category": {Type: schema.TypeString, Optional: true, Description: "Backup Policy Type Category"},
			"page":                        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                        {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":                    {Type: schema.TypeList, Computed: true, Description: "Backup List", Elem: datasourceElem()},
			"total_count":                 {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of backups",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var backupName optional.String
	var createdBy optional.String
	var backupPolicyTypeCategory optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("backup_name"); ok {
		backupName = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("created_by"); ok {
		createdBy = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("backup_policy_type_category"); ok {
		backupPolicyTypeCategory = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("page"); ok {
		page = optional.NewInt32(int32(v.(int)))
	}
	if v, ok := rd.GetOk("size"); ok {
		size = optional.NewInt32(int32(v.(int)))
	}
	if v, ok := rd.GetOk("sort"); ok {
		sort = optional.NewInterface(v.([]interface{}))
	}

	responses, err := inst.Client.Backup.ReadBackupList(ctx, backup2.BackupSearchOpenApiApiListBackupsOpts{
		BackupName:               backupName,
		CreatedBy:                createdBy,
		BackupPolicyTypeCategory: backupPolicyTypeCategory,
		Page:                     page,
		Size:                     size,
		Sort:                     sort,
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

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":                  {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"backup_dr_id":                {Type: schema.TypeString, Computed: true, Description: "Backup DR ID"},
			"backup_id":                   {Type: schema.TypeString, Computed: true, Description: "Backup ID"},
			"backup_name":                 {Type: schema.TypeString, Computed: true, Description: "Backup Name"},
			"backup_policy_type_category": {Type: schema.TypeString, Computed: true, Description: "Backup Policy Type Category"},
			"backup_state":                {Type: schema.TypeString, Computed: true, Description: "Backup State"},
			"block_id":                    {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"is_backup_dr_enabled":        {Type: schema.TypeString, Computed: true, Description: "Is Backup Dr enabled"},
			"is_backup_dr_origin":         {Type: schema.TypeString, Computed: true, Description: "Is Backup Dr Origin"},
			"object_id":                   {Type: schema.TypeString, Computed: true, Description: "Object ID"},
			"object_type":                 {Type: schema.TypeString, Computed: true, Description: "Object Type"},
			"policy_type":                 {Type: schema.TypeString, Computed: true, Description: "Policy Type"},
			"region":                      {Type: schema.TypeString, Computed: true, Description: "Region"},
			"retention_period":            {Type: schema.TypeString, Computed: true, Description: "Retention Period"},
			"service_zone_id":             {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"zone_name":                   {Type: schema.TypeString, Computed: true, Description: "Zone Name"},
			"created_by":                  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":                  {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":                 {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":                 {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}
