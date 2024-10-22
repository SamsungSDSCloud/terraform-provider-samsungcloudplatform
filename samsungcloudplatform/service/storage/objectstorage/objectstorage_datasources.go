package objectstorage

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	scp.RegisterDataSource("samsungcloudplatform_obs_storages", DatasourceObjectStorages())
	scp.RegisterDataSource("samsungcloudplatform_obs_buckets", DatasourceObjectStorageBuckets())
	scp.RegisterDataSource("samsungcloudplatform_obs_bucket", DatasourceObjectStorageBucketInfo())
}

func DatasourceObjectStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceObjectStoragesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"is_multi_availability_zone": {Type: schema.TypeBool, Optional: true, Description: "Is Multi Availability Zone"},
			"object_storage_name":        {Type: schema.TypeString, Optional: true, Description: "Object Storage Name"},
			"service_zone_id":            {Type: schema.TypeString, Required: true, Description: "Service Zone ID"},
			"page":                       {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                       {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                       {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":                   {Type: schema.TypeList, Computed: true, Description: "Object Storage List", Elem: datasourceObjectStorageElem()},
			"total_count":                {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of Object Storages.",
	}
}

func datasourceObjectStoragesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var isMultiAvailabilityZone optional.Bool
	var objectStorageName optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("is_multi_availability_zone"); ok {
		isMultiAvailabilityZone = optional.NewBool(v.(bool))
	}
	if v, ok := rd.GetOk("object_storage_name"); ok {
		objectStorageName = optional.NewString(v.(string))
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

	responses, err := inst.Client.ObjectStorage.ReadObjectStorageList(ctx, rd.Get("service_zone_id").(string), objectstorage.ReadObjectStorageListRequest{
		IsMultiAvailabilityZone: isMultiAvailabilityZone,
		ObjectStorageName:       objectStorageName,
		Page:                    page,
		Size:                    size,
		Sort:                    sort,
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

func datasourceObjectStorageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_id":                   {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"is_multi_availability_zone": {Type: schema.TypeBool, Computed: true, Description: "Is Multi Availability Zone"},
			"object_storage_id":          {Type: schema.TypeString, Computed: true, Description: "Object Storage ID"},
			"object_storage_name":        {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			"service_zone_id":            {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"object_storage_description": {Type: schema.TypeString, Computed: true, Description: "Object Storage Description"},
			"created_by":                 {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":                 {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":                {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":                {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
			"project_id":                 {Type: schema.TypeString, Computed: true, Description: "Project ID"},
		},
	}
}

func DatasourceObjectStorageBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceObjectStorageBucketsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"object_storage_system_bucket_enabled": {Type: schema.TypeBool, Optional: true, Description: "Object Storage System Bucket Enabled"},
			"object_storage_id":                    {Type: schema.TypeString, Optional: true, Description: "Object Storage ID"},
			"object_storage_bucket_name":           {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Name"},
			"end_modified_dt":                      {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Query End Date"},
			"start_modified_dt":                    {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Query Start Date"},
			"object_storage_bucket_state":          {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket State"},
			"object_storage_bucket_states":         {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket State List", Elem: &schema.Schema{Type: schema.TypeString}},
			"object_storage_bucket_purposes":       {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket Purpose Type List", Elem: &schema.Schema{Type: schema.TypeString}},
			"object_storage_bucket_user_purpose":   {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket User Purpose"},
			"object_storage_quota_id":              {Type: schema.TypeString, Optional: true, Description: "Object Storage Quota ID"},
			"object_storage_bucket_ids":            {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket ID List", Elem: &schema.Schema{Type: schema.TypeString}},
			"service_zone_id":                      {Type: schema.TypeString, Optional: true, Description: "Service Zone ID"},
			"created_by":                           {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"page":                                 {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                                 {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                                 {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":                             {Type: schema.TypeList, Computed: true, Description: "Object Storage Bucket List", Elem: datasourceObjectStorageBucketElem()},
			"total_count":                          {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of Object Storage Buckets.",
	}
}

func datasourceObjectStorageBucketsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var objectStorageSystemBucketEnabled optional.Bool
	var objectStorageId optional.String
	var objectStorageBucketName optional.String
	var endModifiedDt optional.Time
	var startModifiedDt optional.Time
	var objectStorageBucketState optional.String
	var objectStorageBucketStates optional.Interface
	var objectStorageBucketPurposes optional.Interface
	var objectStorageBucketUserPurpose optional.String
	var objectStorageQuotaId optional.String
	var serviceZoneId optional.String
	var objectStorageBucketIds optional.Interface
	var createdBy optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("object_storage_system_bucket_enabled"); ok {
		objectStorageSystemBucketEnabled = optional.NewBool(v.(bool))
	}
	if v, ok := rd.GetOk("object_storage_id"); ok {
		objectStorageId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("object_storage_bucket_name"); ok {
		objectStorageBucketName = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("end_modified_dt"); ok {
		endModifiedDt = optional.NewTime(v.(time.Time))
	}
	if v, ok := rd.GetOk("start_modified_dt"); ok {
		startModifiedDt = optional.NewTime(v.(time.Time))
	}
	if v, ok := rd.GetOk("object_storage_bucket_state"); ok {
		objectStorageBucketState = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("object_storage_bucket_states"); ok {
		objectStorageBucketStates = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("object_storage_bucket_purposes"); ok {
		objectStorageBucketPurposes = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("object_storage_bucket_user_purpose"); ok {
		objectStorageBucketUserPurpose = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("object_storage_quota_id"); ok {
		objectStorageQuotaId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("service_zone_id"); ok {
		serviceZoneId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("object_storage_bucket_ids"); ok {
		objectStorageBucketIds = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("created_by"); ok {
		createdBy = optional.NewString(v.(string))
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

	responses, err := inst.Client.ObjectStorage.ReadBucketList(ctx, objectstorage.ReadBucketListRequest{
		ObjectStorageSystemBucketEnabled: objectStorageSystemBucketEnabled,
		ObjectStorageId:                  objectStorageId,
		ObjectStorageBucketName:          objectStorageBucketName,
		EndModifiedDt:                    endModifiedDt,
		StartModifiedDt:                  startModifiedDt,
		ObjectStorageBucketState:         objectStorageBucketState,
		ObjectStorageBucketStates:        objectStorageBucketStates,
		ObjectStorageBucketPurposes:      objectStorageBucketPurposes,
		ObjectStorageBucketUserPurpose:   objectStorageBucketUserPurpose,
		ObjectStorageQuotaId:             objectStorageQuotaId,
		ObjectStorageBucketIds:           objectStorageBucketIds,
		ServiceZoneId:                    serviceZoneId,
		CreatedBy:                        createdBy,
		Page:                             page,
		Size:                             size,
		Sort:                             sort,
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

func datasourceObjectStorageBucketElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"object_storage_bucket_access_control_enabled": {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket Access Control Enabled"},
			"object_storage_bucket_dr_enabled":             {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket DR Enabled"},
			"object_storage_bucket_dr_type":                {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket DR Type"},
			"object_storage_bucket_id":                     {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket ID"},
			"object_storage_bucket_name":                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Name"},
			"object_storage_bucket_purpose":                {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Purpose"},
			"object_storage_bucket_user_purpose":           {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket User Purpose"},
			"object_storage_bucket_state":                  {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket State"},
			"object_storage_bucket_version_enabled":        {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket Version Enabled"},
			"object_storage_id":                            {Type: schema.TypeString, Computed: true, Description: "Object Storage ID"},
			"object_storage_name":                          {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			"object_storage_quota_id":                      {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota ID"},
			"object_storage_quota_name":                    {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota Name"},
			"object_storage_system_bucket_enabled":         {Type: schema.TypeBool, Computed: true, Description: "Object Storage System Bucket Enabled"},
			"object_storage_tenant_name":                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Tenant Name"},
			"service_zone_id":                              {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"created_by":                                   {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":                                   {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":                                  {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":                                  {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}

func DatasourceObjectStorageBucketInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceObjectStorageBucketInfoRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                               {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			common.ToSnakeCase("IsMultiAvailabilityZone"):                 {Type: schema.TypeBool, Computed: true, Description: "Is Multi Availability Zone"},
			common.ToSnakeCase("IsSyncInProgress"):                        {Type: schema.TypeBool, Computed: true, Description: "Is Sync In-progress"},
			common.ToSnakeCase("ObjectStorageBucketAccessControlEnabled"): {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket Access Control Enabled"},
			common.ToSnakeCase("ObjectStorageBucketAccessControlRules"): {Type: schema.TypeList, Computed: true, Description: "Object Storage Bucket Access Control Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access Control Rule Type",
						},
						"rule_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access Control Rule Value",
						},
					},
				},
			},
			common.ToSnakeCase("ObjectStorageBucketDrEnabled"):               {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket DR Enabled"},
			common.ToSnakeCase("ObjectStorageBucketDrType"):                  {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket DR Type"},
			common.ToSnakeCase("ObjectStorageBucketFileEncryptionAlgorithm"): {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket File Encryption Algorithm"},
			common.ToSnakeCase("ObjectStorageBucketFileEncryptionEnabled"):   {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket File Encryption Enabled"},
			common.ToSnakeCase("ObjectStorageBucketFileEncryptionType"):      {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket File Encryption Type"},
			common.ToSnakeCase("ObjectStorageBucketId"):                      {Type: schema.TypeString, Required: true, Description: "Object Storage Bucket ID"},
			common.ToSnakeCase("ObjectStorageBucketName"):                    {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Name"},
			common.ToSnakeCase("ObjectStorageBucketObjectUploadEnabled"):     {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket Object Upload Enabled"},
			common.ToSnakeCase("ObjectStorageBucketPrivateEndpointUrl"):      {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Private Endpoint URL"},
			common.ToSnakeCase("ObjectStorageBucketPublicEndpointUrl"):       {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Public Endpoint URL"},
			common.ToSnakeCase("ObjectStorageBucketPurpose"):                 {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Purpose"},
			common.ToSnakeCase("ObjectStorageBucketUserPurpose"):             {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket User Purpose"},
			common.ToSnakeCase("ObjectStorageBucketState"):                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket State"},
			common.ToSnakeCase("ObjectStorageBucketUsage"):                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Usage"},
			common.ToSnakeCase("ObjectStorageBucketVersionEnabled"):          {Type: schema.TypeBool, Computed: true, Description: "Object Storage Bucket Version Enabled"},
			common.ToSnakeCase("ObjectStorageDeviceUserId"):                  {Type: schema.TypeString, Computed: true, Description: "Object Storage Device User ID"},
			common.ToSnakeCase("ObjectStorageId"):                            {Type: schema.TypeString, Computed: true, Description: "Object Storage ID"},
			common.ToSnakeCase("ObjectStorageName"):                          {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			common.ToSnakeCase("ObjectStorageQuotaId"):                       {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota ID"},
			common.ToSnakeCase("ObjectStorageQuotaName"):                     {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota Name"},
			common.ToSnakeCase("ObjectStorageSystemBucketEnabled"):           {Type: schema.TypeBool, Computed: true, Description: "Object Storage System Bucket Enabled"},
			common.ToSnakeCase("ObjectStorageTenantName"):                    {Type: schema.TypeString, Computed: true, Description: "Object Storage Tenant Name"},
			common.ToSnakeCase("ServiceZoneId"):                              {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			common.ToSnakeCase("SyncObjectStorageBucketId"):                  {Type: schema.TypeString, Computed: true, Description: "Sync Object Storage Bucket ID"},
			common.ToSnakeCase("SyncObjectStorageBucketName"):                {Type: schema.TypeString, Computed: true, Description: "Sync Object Storage Bucket Name"},
			common.ToSnakeCase("SyncObjectStorageBucketServiceZoneId"):       {Type: schema.TypeString, Computed: true, Description: "Sync Object Storage Bucket Service Zone ID"},

			//common.ToSnakeCase("ObsUrls"): {Type: schema.TypeMap, Optional: true, Description: "Obs Urls",
			//	Elem: &schema.Schema{
			//		Type:     schema.TypeString,
			//		Computed: true,
			//	},
			//},

			common.ToSnakeCase("CreatedBy"):  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):  {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"): {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"): {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
		Description: "Provides information of Object Storage Bucket.",
	}
}

func datasourceObjectStorageBucketInfoRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	BucketId := rd.Get("object_storage_bucket_id").(string)

	info, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, BucketId)

	s := common.HclSetObject{}
	for _, svc := range info.ObjectStorageBucketAccessControlRules {
		s = append(s, common.HclKeyValueObject{
			"rule_value": svc.RuleValue,
			"rule_type":  svc.RuleType,
		})
	}
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}
	//
	//m := common.HclKeyValueObject{}
	//m[common.ToSnakeCase("ObsExternalServiceUrl")] = info.ObsUrls.ObsExternalServiceUrl
	//m[common.ToSnakeCase("ObsExternalSystemUrl")] = info.ObsUrls.ObsExternalSystemUrl
	//m[common.ToSnakeCase("ObsInternalCompanyNetworkUrl")] = info.ObsUrls.ObsInternalCompanyNetworkUrl
	//m[common.ToSnakeCase("ObsInternalGroupNetworkUrl")] = info.ObsUrls.ObsInternalGroupNetworkUrl
	//m[common.ToSnakeCase("ObsInternalInternetUrl")] = info.ObsUrls.ObsInternalInternetUrl
	//m[common.ToSnakeCase("ObsVpcV2ExternalUrl")] = info.ObsUrls.ObsVpcV2ExternalUrl
	//m[common.ToSnakeCase("ObsVpcV2InternalDbBackupUrl")] = info.ObsUrls.ObsVpcV2InternalDbBackupUrl
	//m[common.ToSnakeCase("ObsVpcV2InternalNatUrl")] = info.ObsUrls.ObsVpcV2InternalNatUrl
	//m[common.ToSnakeCase("ObsVpcV2InternalPbBackupUrl")] = info.ObsUrls.ObsVpcV2InternalPbBackupUrl
	//m[common.ToSnakeCase("ObsVpcV2InternalUrl")] = info.ObsUrls.ObsVpcV2InternalUrl

	rd.Set(common.ToSnakeCase("ProjectId"), info.ProjectId)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketObjectUploadEnabled"), info.ObjectStorageBucketObjectUploadEnabled)
	rd.Set(common.ToSnakeCase("ObjectStorageSystemBucketEnabled"), info.ObjectStorageSystemBucketEnabled)
	rd.Set(common.ToSnakeCase("IsSyncInProgress"), info.IsSyncInProgress)
	rd.Set(common.ToSnakeCase("IsMultiAvailabilityZone"), info.IsMultiAvailabilityZone)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketAccessControlEnabled"), info.ObjectStorageBucketAccessControlEnabled)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketAccessControlRules"), s)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketDrEnabled"), info.ObjectStorageBucketDrEnabled)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketDrType"), info.ObjectStorageBucketDrType)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketFileEncryptionAlgorithm"), info.ObjectStorageBucketFileEncryptionAlgorithm)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketFileEncryptionEnabled"), info.ObjectStorageBucketFileEncryptionEnabled)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketFileEncryptionType"), info.ObjectStorageBucketFileEncryptionType)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketId"), info.ObjectStorageBucketId)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketName"), info.ObjectStorageBucketName)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketState"), info.ObjectStorageBucketState)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketUsage"), info.ObjectStorageBucketUsage)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketPurpose"), info.ObjectStorageBucketPurpose)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketUserPurpose"), info.ObjectStorageBucketUserPurpose)
	rd.Set(common.ToSnakeCase("ObjectStorageBucketVersionEnabled"), info.ObjectStorageBucketVersionEnabled)
	rd.Set(common.ToSnakeCase("ObjectStorageId"), info.ObjectStorageId)
	rd.Set(common.ToSnakeCase("ObjectStorageName"), info.ObjectStorageName)
	rd.Set(common.ToSnakeCase("ObjectStorageQuotaId"), info.ObjectStorageQuotaId)
	rd.Set(common.ToSnakeCase("SyncObjectStorageBucketId"), info.SyncObjectStorageBucketId)
	rd.Set(common.ToSnakeCase("SyncObjectStorageBucketName"), info.SyncObjectStorageBucketName)
	rd.Set(common.ToSnakeCase("SyncObjectStorageBucketServiceZoneId"), info.SyncObjectStorageBucketServiceZoneId)
	rd.Set(common.ToSnakeCase("ServiceZoneId"), info.ServiceZoneId)
	rd.Set(common.ToSnakeCase("CreatedBy"), info.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), info.CreatedDt)
	rd.Set(common.ToSnakeCase("ModifiedBy"), info.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), info.ModifiedDt)

	rd.SetId(uuid.NewV4().String())

	return nil
}
