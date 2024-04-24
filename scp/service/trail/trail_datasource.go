package trail

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	loggingaudit2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/logging-audit"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_trails", DatasourceTrails())
	scp.RegisterDataSource("scp_trail", DatasourceTrail())
}

func DatasourceTrails() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceTrailsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"is_mine": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Search my trail",
			},
			"logging_target_regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target region list",
			},
			"logging_target_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target resource ID list",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State (ACTIVE | STOPPED)",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trail name",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count",
			},
			"contents": {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceTrailElem()},
		},
		Description: "Provides list of trails",
	}
}

func datasourceTrailsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var isMinePtr *bool
	isMine := false
	if rd.Get("is_mine") == nil {
		isMinePtr = nil
	} else {
		isMine = rd.Get("is_mine").(bool)
		isMinePtr = &isMine
	}
	state := rd.Get("state").(string)
	name := rd.Get("name").(string)
	regions := common.ToStringList(rd.Get("logging_target_regions").(*schema.Set).List())
	resourceIds := common.ToStringList(rd.Get("logging_target_resource_ids").(*schema.Set).List())

	response, err := inst.Client.Loggingaudit.ListTrails(ctx, isMinePtr, regions, resourceIds, state, name)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertTrailListToHclSet(response)
	//contents := common.ConvertStructToMaps(response.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceTrails().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func convertTrailListToHclSet(trail loggingaudit2.PageResponseV2TrailResponse) common.HclSetObject {
	var trailList common.HclSetObject
	for _, tr := range trail.Contents {
		if len(tr.TrailId) == 0 {
			continue
		}

		var targetUserList common.HclListObject
		for _, user := range tr.LoggingTargetUsers {
			userKv := common.HclKeyValueObject{
				"email":    user.Email,
				"id":       user.UserId,
				"login_id": user.UserLoginId,
				"name":     user.UserName,
			}
			targetUserList = append(targetUserList, userKv)
		}

		//target_logging_resource_list
		var targetResourceList common.HclListObject
		for _, resource := range tr.TargetLoggingResourceList {
			resKv := common.HclKeyValueObject{
				"logging_item":       resource.LoggingItem,
				"resource_type_name": resource.ResourceTypeName,
			}
			targetResourceList = append(targetResourceList, resKv)
		}

		kv := common.HclKeyValueObject{
			"project_id":                     tr.ProjectId,
			"project_name":                   tr.ProjectName,
			"region":                         tr.Region,
			"service_zone_id":                tr.ServiceZoneId,
			"is_logging_target_all_region":   tr.IsLoggingTargetAllRegion,
			"is_logging_target_all_resource": tr.IsLoggingTargetAllResource,
			"is_logging_target_all_user":     tr.IsLoggingTargetAllUser,
			"is_trail_deleted":               tr.IsTrailDeleted,
			"last_digest_file_create_dt":     tr.LastDigestFileCreateDt.String(),
			"logging_target_regions":         tr.LoggingTargetRegions,
			"logging_target_resource_ids":    tr.LoggingTargetResourceIds,
			"logging_target_users":           targetUserList,
			"object_storage_folder_name":     tr.ObjectStorageFolderName,
			"object_storage_name":            tr.ObjectStorageName,
			"obs_bucket_id":                  tr.ObsBucketId,
			"obs_bucket_name":                tr.ObsBucketName,
			"target_logging_resource_list":   targetResourceList,
			"trail_batch_end_dt":             tr.TrailBatchEndDt.String(),
			"trail_batch_first_start_dt":     tr.TrailBatchFirstStartDt.String(),
			"trail_batch_last_success_dt":    tr.TrailBatchLastSuccessDt.String(),
			"trail_batch_start_dt":           tr.TrailBatchStartDt.String(),
			"trail_batch_state":              tr.TrailBatchState,
			"trail_id":                       tr.TrailId,
			"trail_name":                     tr.TrailName,
			"trail_save_type":                tr.TrailSaveType,
			"trail_state":                    tr.TrailState,
			"validation_yn":                  tr.ValidationYn,
			"description":                    tr.TrailDescription,
			"created_by":                     tr.CreatedBy,
			"created_by_name":                tr.CreatedByName,
			"created_dt":                     tr.CreatedDt.String(),
			"modified_by":                    tr.ModifiedBy,
			"modified_by_name":               tr.ModifiedByName,
			"modified_dt":                    tr.ModifiedDt.String(),
		}
		trailList = append(trailList, kv)
	}
	return trailList
}

func datasourceTrailElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":                     {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"project_name":                   {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"region":                         {Type: schema.TypeString, Computed: true, Description: "Region"},
			"service_zone_id":                {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"is_logging_target_all_region":   {Type: schema.TypeString, Computed: true, Description: "Whether to target all regions"},
			"is_logging_target_all_resource": {Type: schema.TypeString, Computed: true, Description: "Whether to target all resources"},
			"is_logging_target_all_user":     {Type: schema.TypeString, Computed: true, Description: "Whether to target all users"},
			"is_trail_deleted":               {Type: schema.TypeString, Computed: true, Description: "Whether the trail is deleted"},
			"last_digest_file_create_dt":     {Type: schema.TypeString, Computed: true, Description: "Last digest file create date and time"},
			"logging_target_regions":         {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Logging target region list"},
			"logging_target_resource_ids":    {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Logging target resource ID list"},
			"logging_target_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email":    {Type: schema.TypeString, Computed: true, Description: "Email"},
						"id":       {Type: schema.TypeString, Computed: true, Description: "User ID"},
						"login_id": {Type: schema.TypeString, Computed: true, Description: "User login ID"},
						"name":     {Type: schema.TypeString, Computed: true, Description: "User name"},
					},
				},
				Description: "Logging target user list",
			},
			"target_logging_resource_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logging_item":       {Type: schema.TypeString, Computed: true, Description: "Logging item"},
						"resource_type_name": {Type: schema.TypeString, Computed: true, Description: "Resource type"},
					},
				},
				Description: "Target logging resource list",
			},
			"object_storage_folder_name":  {Type: schema.TypeString, Computed: true, Description: "Object storage folder name"},
			"object_storage_name":         {Type: schema.TypeString, Computed: true, Description: "Object storage name"},
			"obs_bucket_id":               {Type: schema.TypeString, Computed: true, Description: "Object storage bucket ID"},
			"obs_bucket_name":             {Type: schema.TypeString, Computed: true, Description: "Object storage bucket name"},
			"trail_batch_end_dt":          {Type: schema.TypeString, Computed: true, Description: "Batch processing end date and time"},
			"trail_batch_first_start_dt":  {Type: schema.TypeString, Computed: true, Description: "Batch processing first start date and time"},
			"trail_batch_last_success_dt": {Type: schema.TypeString, Computed: true, Description: "Batch processing last success date and time"},
			"trail_batch_start_dt":        {Type: schema.TypeString, Computed: true, Description: "Batch processing start date and time"},
			"trail_batch_state":           {Type: schema.TypeString, Computed: true, Description: "Batch processing status"},
			"trail_id":                    {Type: schema.TypeString, Computed: true, Description: "Trail ID"},
			"trail_name":                  {Type: schema.TypeString, Computed: true, Description: "Trail name"},
			"trail_save_type":             {Type: schema.TypeString, Computed: true, Description: "Trail save type"},
			"trail_state":                 {Type: schema.TypeString, Computed: true, Description: "Trail status"},
			"validation_yn":               {Type: schema.TypeString, Computed: true, Description: "Trail verification"},
			"description":                 {Type: schema.TypeString, Computed: true, Description: "Description"},
			"created_by":                  {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":             {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_dt":                  {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":                 {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":            {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_dt":                 {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}

func DatasourceTrail() *schema.Resource {
	var trailResource schema.Resource
	trailResource.ReadContext = datasourceTrailRead
	trailResource.Schema = datasourceTrailElem().Schema

	delete(trailResource.Schema, "trail_id")
	trailResource.Schema["trail_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Trail ID",
	}
	trailResource.Description = "Provides detailed trail information for a given trail id"

	return &trailResource
}

func datasourceTrailRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	trailId := rd.Get("trail_id").(string)
	info, _, err := inst.Client.Loggingaudit.ReadTrail(ctx, trailId)
	if err != nil {
		return diag.FromErr(err)
	}

	var targetUserList common.HclListObject
	for _, user := range info.LoggingTargetUsers {
		userKv := common.HclKeyValueObject{
			"email":    user.Email,
			"id":       user.UserId,
			"login_id": user.UserLoginId,
			"name":     user.UserName,
		}
		targetUserList = append(targetUserList, userKv)
	}

	//target_logging_resource_list
	var targetResourceList common.HclListObject
	for _, resource := range info.TargetLoggingResourceList {
		resKv := common.HclKeyValueObject{
			"logging_item":       resource.LoggingItem,
			"resource_type_name": resource.ResourceTypeName,
		}
		targetResourceList = append(targetResourceList, resKv)
	}

	rd.SetId(trailId)
	rd.Set("project_id", info.ProjectId)
	rd.Set("project_name", info.ProjectName)
	rd.Set("region", info.Region)
	rd.Set("service_zone_id", info.ServiceZoneId)
	rd.Set("is_logging_target_all_region", info.IsLoggingTargetAllRegion)
	rd.Set("is_logging_target_all_resource", info.IsLoggingTargetAllResource)
	rd.Set("is_logging_target_all_user", info.IsLoggingTargetAllUser)
	rd.Set("is_trail_deleted", info.IsTrailDeleted)
	rd.Set("last_digest_file_create_dt", info.LastDigestFileCreateDt.String())
	rd.Set("logging_target_regions", info.LoggingTargetRegions)
	rd.Set("logging_target_resource_ids", info.LoggingTargetResourceIds)
	rd.Set("logging_target_users", targetUserList)
	rd.Set("object_storage_folder_name", info.ObjectStorageFolderName)
	rd.Set("object_storage_name", info.ObjectStorageName)
	rd.Set("obs_bucket_id", info.ObsBucketId)
	rd.Set("obs_bucket_name", info.ObsBucketName)
	rd.Set("target_logging_resource_list", targetResourceList)
	rd.Set("trail_batch_end_dt", info.TrailBatchEndDt.String())
	rd.Set("trail_batch_first_start_dt", info.TrailBatchFirstStartDt.String())
	rd.Set("trail_batch_last_success_dt", info.TrailBatchLastSuccessDt.String())
	rd.Set("trail_batch_start_dt", info.TrailBatchStartDt.String())
	rd.Set("trail_batch_state", info.TrailBatchState)
	rd.Set("trail_name", info.TrailName)
	rd.Set("trail_save_type", info.TrailSaveType)
	rd.Set("trail_state", info.TrailState)
	rd.Set("validation_yn", info.ValidationYn)
	rd.Set("description", info.TrailDescription)
	rd.Set("created_by", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by", info.ModifiedBy)
	rd.Set("modified_by_name", info.ModifiedByName)
	rd.Set("modified_dt", info.ModifiedDt.String())

	return nil
}
