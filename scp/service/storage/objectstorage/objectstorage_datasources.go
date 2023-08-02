package objectstorage

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/antihax/optional"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	objectstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/object-storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_obs_storages", DatasourceObjectStorages())
	scp.RegisterDataSource("scp_obs_buckets", DatasourceObjectStorageBuckets())
	scp.RegisterDataSource("scp_obs_bucket", DatasourceObjectStorageBucketInfo())
}

func DatasourceObjectStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceObjectStoragesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"multi_az_yn":       {Type: schema.TypeString, Optional: true, Description: "Multi AZ Y/N"},
			"obs_name":          {Type: schema.TypeString, Optional: true, Description: "Object Storage Name"},
			"obs_rest_endpoint": {Type: schema.TypeString, Optional: true, Description: "Object Storage REST Endpoint"},
			"region":            {Type: schema.TypeString, Optional: true, Description: "Region"},
			"zone_id":           {Type: schema.TypeString, Required: true, Description: "Service Zone ID"},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":              {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Object Storage List", Elem: datasourceObjectStorageElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of Object Storages.",
	}
}

func datasourceObjectStoragesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var multiAzYn optional.String
	var obsName optional.String
	var obsRestEndpoint optional.String
	var region optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("multi_az_yn"); ok {
		multiAzYn = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_name"); ok {
		obsName = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_rest_endpoint"); ok {
		obsRestEndpoint = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("region"); ok {
		region = optional.NewString(v.(string))
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

	responses, err := inst.Client.ObjectStorage.ReadObjectStorageList(ctx, rd.Get("zone_id").(string), objectstorage.ObjectStorageV3ControllerApiListObjectStorage3Opts{
		MultiAzYn:       multiAzYn,
		ObsName:         obsName,
		ObsRestEndpoint: obsRestEndpoint,
		Region:          region,
		Page:            page,
		Size:            size,
		Sort:            sort,
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
			"endpoint_ip":                {Type: schema.TypeString, Computed: true, Description: "Endpoint IP"},
			"multi_az_yn":                {Type: schema.TypeString, Computed: true, Description: "Multi AZ Y/N"},
			"obs_id":                     {Type: schema.TypeString, Computed: true, Description: "Object Storage ID"},
			"obs_internal_rest_endpoint": {Type: schema.TypeString, Computed: true, Description: "Object Storage Internal Rest Endpoint"},
			"obs_name":                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			"obs_rest_endpoint":          {Type: schema.TypeString, Computed: true, Description: "Object Storage Rest Endpoint"},
			"obs_vendor":                 {Type: schema.TypeString, Computed: true, Description: "Object Storage Vendor"},
			"region":                     {Type: schema.TypeString, Computed: true, Description: "Region"},
			"system_url":                 {Type: schema.TypeString, Computed: true, Description: "System URL"},
			"obs_description":            {Type: schema.TypeString, Computed: true, Description: "Object Storage Description"},
			"created_by":                 {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":                 {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":                {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":                {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
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
			"is_obs_bucket_sync":           {Type: schema.TypeBool, Optional: true, Description: "Perform Object Storage Bucket sync (true | false)"},
			"is_obs_system_bucket_enabled": {Type: schema.TypeBool, Optional: true, Description: "Is Object Storage System Bucket enabled (true | false)"},
			"obs_storage_id":               {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Name"},
			"obs_bucket_name":              {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Name (Like)"},
			"obs_bucket_name_exact":        {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Name (Equal)"},
			"obs_bucket_query_end_dt":      {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Query End Date"},
			"obs_bucket_query_start_dt":    {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket Query Start Date"},
			"obs_bucket_state":             {Type: schema.TypeString, Optional: true, Description: "Object Storage Bucket State"},
			"obs_bucket_used_type":         {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket Used Type", Elem: &schema.Schema{Type: schema.TypeString}},
			"obs_quota_id":                 {Type: schema.TypeString, Optional: true, Description: "Object Storage Quota ID"},
			"pool_region":                  {Type: schema.TypeString, Optional: true, Description: "Region"},
			"obs_bucket_id_list":           {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket ID List", Elem: &schema.Schema{Type: schema.TypeString}},
			"obs_bucket_state_in":          {Type: schema.TypeList, Optional: true, Description: "Object Storage Bucket State List", Elem: &schema.Schema{Type: schema.TypeString}},
			"created_by":                   {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"page":                         {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                         {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":                     {Type: schema.TypeList, Computed: true, Description: "Object Storage Bucket List", Elem: datasourceObjectStorageBucketElem()},
			"total_count":                  {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of Object Storage Buckets.",
	}
}

func datasourceObjectStorageBucketsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var isObsBucketSync optional.Bool
	var isObsSystemBucketEnabled optional.Bool
	var objectStorageId optional.String
	var obsBucketName optional.String
	var obsBucketNameExact optional.String
	var obsBucketQueryEndDt optional.String
	var obsBucketQueryStartDt optional.String
	var obsBucketState optional.String
	var obsBucketUsedType optional.Interface
	var obsQuotaId optional.String
	var poolRegion optional.String
	var obsBucketIdList optional.Interface
	var obsBucketStateIn optional.Interface
	var createdBy optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("is_obs_bucket_sync"); ok {
		isObsBucketSync = optional.NewBool(v.(bool))
	}
	if v, ok := rd.GetOk("is_obs_system_bucket_enabled"); ok {
		isObsSystemBucketEnabled = optional.NewBool(v.(bool))
	}
	if v, ok := rd.GetOk("obs_storage_id"); ok {
		objectStorageId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_name"); ok {
		obsBucketName = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_name_exact"); ok {
		obsBucketNameExact = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_query_end_dt"); ok {
		obsBucketQueryEndDt = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_query_start_dt"); ok {
		obsBucketQueryStartDt = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_state"); ok {
		obsBucketState = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_used_type"); ok {
		obsBucketUsedType = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("obs_quota_id"); ok {
		obsQuotaId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("pool_region"); ok {
		poolRegion = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("obs_bucket_id_list"); ok {
		obsBucketIdList = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("obs_bucket_state_in"); ok {
		obsBucketStateIn = optional.NewInterface(v.([]interface{}))
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

	responses, err := inst.Client.ObjectStorage.ReadBucketList(ctx, objectstorage.ObjectStorageBucketV3ControllerApiListBucket3Opts{
		IsObsBucketSync:          isObsBucketSync,
		IsObsSystemBucketEnabled: isObsSystemBucketEnabled,
		ObjectStorageId:          objectStorageId,
		ObsBucketName:            obsBucketName,
		ObsBucketNameExact:       obsBucketNameExact,
		ObsBucketQueryEndDt:      obsBucketQueryEndDt,
		ObsBucketQueryStartDt:    obsBucketQueryStartDt,
		ObsBucketState:           obsBucketState,
		ObsBucketUsedType:        obsBucketUsedType,
		ObsQuotaId:               obsQuotaId,
		PoolRegion:               poolRegion,
		ObsBucketIdList:          obsBucketIdList,
		ObsBucketStateIn:         obsBucketStateIn,
		CreatedBy:                createdBy,
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

func datasourceObjectStorageBucketElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":                              {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"is_obs_bucket_dr_enabled":                {Type: schema.TypeBool, Computed: true, Description: "Is Object Storage Bucket DR enabled"},
			"is_obs_bucket_ip_address_filter_enabled": {Type: schema.TypeBool, Computed: true, Description: "Is Object Storage Bucket IP Address Filter enabled"},
			"is_obs_object_creation_enabled":          {Type: schema.TypeBool, Computed: true, Description: "Is Object Storage Object Creation enabled"},
			"is_obs_system_bucket_enabled":            {Type: schema.TypeBool, Computed: true, Description: "Is Object Storage System Bucket enabled"},
			"obs_bucket_dr_type":                      {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket DR Type"},
			"obs_bucket_id":                           {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket ID"},
			"obs_bucket_name":                         {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Name"},
			"obs_bucket_state":                        {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket State"},
			"obs_bucket_used_type":                    {Type: schema.TypeString, Computed: true, Description: "Object Storage Bucket Used Type"},
			"obs_id":                                  {Type: schema.TypeString, Computed: true, Description: "Object Storage ID"},
			"obs_name":                                {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			"obs_quota_id":                            {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota ID"},
			"obs_quota_name":                          {Type: schema.TypeString, Computed: true, Description: "Object Storage Quota Name"},
			"obs_tenant_name":                         {Type: schema.TypeString, Computed: true, Description: "Object Storage Tenant Name"},
			"pool_region":                             {Type: schema.TypeString, Computed: true, Description: "Region"},
			"system_id":                               {Type: schema.TypeString, Computed: true, Description: "System ID"},
			"system_name":                             {Type: schema.TypeString, Computed: true, Description: "System Name"},
			"zone_name":                               {Type: schema.TypeString, Computed: true, Description: "Service Zone Name"},
			"created_by":                              {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":                              {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":                             {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":                             {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
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
			common.ToSnakeCase("ObsBucketId"):                       {Type: schema.TypeString, Required: true, Description: "Obs Bucket Id"},
			common.ToSnakeCase("ProjectId"):                         {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			common.ToSnakeCase("IsObsBucketDrEnabled"):              {Type: schema.TypeBool, Computed: true, Description: "Dr Enabled"},
			common.ToSnakeCase("IsObsBucketIpAddressFilterEnabled"): {Type: schema.TypeBool, Computed: true, Description: "Ip Filter Enabled"},
			common.ToSnakeCase("IsObsObjectCreationEnabled"):        {Type: schema.TypeBool, Computed: true, Description: "Object Creation Enabled"},
			common.ToSnakeCase("IsObsSystemBucketEnabled"):          {Type: schema.TypeBool, Computed: true, Description: "System Bucket Enabled"},
			common.ToSnakeCase("IsReplicationInProgress"):           {Type: schema.TypeBool, Computed: true, Description: "Replication In Progress"},
			common.ToSnakeCase("MultiAzYn"):                         {Type: schema.TypeString, Computed: true, Description: "Multi Az Y/N"},
			common.ToSnakeCase("ObsBucketAccessIpAddressRanges"): {Type: schema.TypeList, Computed: true, Description: "Bucket Access Ip Ranges",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address_range": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
					},
				},
			},
			common.ToSnakeCase("ObsBucketAccessUrl"):               {Type: schema.TypeString, Computed: true, Description: "Bucket Access Url"},
			common.ToSnakeCase("ObsBucketDrType"):                  {Type: schema.TypeString, Computed: true, Description: "Dr Type"},
			common.ToSnakeCase("ObsBucketFileEncryptionAlgorithm"): {Type: schema.TypeString, Computed: true, Description: "Bucket Encryption Algorithm"},
			common.ToSnakeCase("ObsBucketFileEncryptionEnabled"):   {Type: schema.TypeBool, Computed: true, Description: "Is Encryption Enabled"},
			common.ToSnakeCase("ObsBucketFileEncryptionType"):      {Type: schema.TypeString, Computed: true, Description: "Bucket Encryption Type"},

			common.ToSnakeCase("ObsBucketName"):           {Type: schema.TypeString, Computed: true, Description: "Bucket Name"},
			common.ToSnakeCase("ObsBucketState"):          {Type: schema.TypeString, Computed: true, Description: "Bucket State"},
			common.ToSnakeCase("ObsBucketUsedSize"):       {Type: schema.TypeInt, Computed: true, Description: "Bucket Used Size"},
			common.ToSnakeCase("ObsBucketUsedType"):       {Type: schema.TypeString, Computed: true, Description: "Bucket Used Type"},
			common.ToSnakeCase("ObsBucketVersionEnabled"): {Type: schema.TypeBool, Computed: true, Description: "Versioning Enabled"},
			common.ToSnakeCase("ObsId"):                   {Type: schema.TypeString, Computed: true, Description: "Object Storage Id"},
			common.ToSnakeCase("ObsName"):                 {Type: schema.TypeString, Computed: true, Description: "Object Storage Name"},
			common.ToSnakeCase("ObsQuotaId"):              {Type: schema.TypeString, Computed: true, Description: "Obs Quota Id"},
			common.ToSnakeCase("ObsSyncBucketId"):         {Type: schema.TypeString, Computed: true, Description: "Obs Quota Name"},
			common.ToSnakeCase("ObsSyncBucketName"):       {Type: schema.TypeString, Computed: true, Description: "Obs Tenant Name"},
			common.ToSnakeCase("ObsSyncBucketObsName"):    {Type: schema.TypeString, Computed: true, Description: "Pool Region"},
			common.ToSnakeCase("ObsSyncBucketRegion"):     {Type: schema.TypeString, Computed: true, Description: "System Id"},
			common.ToSnakeCase("ObsSyncBucketZoneName"):   {Type: schema.TypeString, Computed: true, Description: "Obs Quota Name"},
			/*  TypeSet has problem on mapping (2023-02-06 seoro)
			common.ToSnakeCase("ObsUrls"): {Type: schema.TypeSet, Optional: true, Description: "Obs Tenant Name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						common.ToSnakeCase("ObsExternalServiceUrl"):        {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsExternalSystemUrl"):         {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsInternalCompanyNetworkUrl"): {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsInternalGroupNetworkUrl"):   {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsInternalInternetUrl"):       {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsVpcV2ExternalUrl"):          {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsVpcV2InternalDbBackupUrl"):  {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsVpcV2InternalNatUrl"):       {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsVpcV2InternalPbBackupUrl"):  {Type: schema.TypeString, Computed: true, Description: ""},
						common.ToSnakeCase("ObsVpcV2InternalUrl"):          {Type: schema.TypeString, Computed: true, Description: ""},
					},
				},
			},
			*/
			common.ToSnakeCase("ObsUrls"): {Type: schema.TypeMap, Optional: true, Description: "Obs Urls",
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
			common.ToSnakeCase("ProjectName"): {Type: schema.TypeString, Computed: true, Description: "Project Name"},
			common.ToSnakeCase("Region"):      {Type: schema.TypeString, Computed: true, Description: "Region"},
			common.ToSnakeCase("SystemId"):    {Type: schema.TypeString, Computed: true, Description: "System Id"},
			common.ToSnakeCase("SystemName"):  {Type: schema.TypeString, Computed: true, Description: "System Name"},
			common.ToSnakeCase("ZoneId"):      {Type: schema.TypeString, Computed: true, Description: "Zone Id"},
			common.ToSnakeCase("ZoneName"):    {Type: schema.TypeString, Computed: true, Description: "Zone Name"},
			common.ToSnakeCase("CreatedBy"):   {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):   {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):  {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):  {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
		Description: "Provides Object Bucket Info.",
	}
}

func datasourceObjectStorageBucketInfoRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	BucketId := rd.Get("obs_bucket_id").(string)

	info, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, BucketId)

	s := common.HclSetObject{}
	for _, svc := range info.ObsBucketAccessIpAddressRanges {
		s = append(s, common.HclKeyValueObject{
			"obs_bucket_access_ip_address_range": svc.ObsBucketAccessIpAddressRange,
			"type":                               svc.Type_,
		})
	}
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	m := common.HclKeyValueObject{}
	m[common.ToSnakeCase("ObsExternalServiceUrl")] = info.ObsUrls.ObsExternalServiceUrl
	m[common.ToSnakeCase("ObsExternalSystemUrl")] = info.ObsUrls.ObsExternalSystemUrl
	m[common.ToSnakeCase("ObsInternalCompanyNetworkUrl")] = info.ObsUrls.ObsInternalCompanyNetworkUrl
	m[common.ToSnakeCase("ObsInternalGroupNetworkUrl")] = info.ObsUrls.ObsInternalGroupNetworkUrl
	m[common.ToSnakeCase("ObsInternalInternetUrl")] = info.ObsUrls.ObsInternalInternetUrl
	m[common.ToSnakeCase("ObsVpcV2ExternalUrl")] = info.ObsUrls.ObsVpcV2ExternalUrl
	m[common.ToSnakeCase("ObsVpcV2InternalDbBackupUrl")] = info.ObsUrls.ObsVpcV2InternalDbBackupUrl
	m[common.ToSnakeCase("ObsVpcV2InternalNatUrl")] = info.ObsUrls.ObsVpcV2InternalNatUrl
	m[common.ToSnakeCase("ObsVpcV2InternalPbBackupUrl")] = info.ObsUrls.ObsVpcV2InternalPbBackupUrl
	m[common.ToSnakeCase("ObsVpcV2InternalUrl")] = info.ObsUrls.ObsVpcV2InternalUrl

	rd.Set(common.ToSnakeCase("ProjectId"), info.ProjectId)
	rd.Set(common.ToSnakeCase("IsObsBucketDrEnabled"), info.IsObsBucketDrEnabled)
	rd.Set(common.ToSnakeCase("IsObsBucketIpAddressFilterEnabled"), info.IsObsBucketIpAddressFilterEnabled)
	rd.Set(common.ToSnakeCase("IsObsObjectCreationEnabled"), info.IsObsObjectCreationEnabled)
	rd.Set(common.ToSnakeCase("IsObsSystemBucketEnabled"), info.IsObsSystemBucketEnabled)
	rd.Set(common.ToSnakeCase("IsReplicationInProgress"), info.IsReplicationInProgress)
	rd.Set(common.ToSnakeCase("MultiAzYn"), info.MultiAzYn)
	rd.Set(common.ToSnakeCase("ObsBucketAccessIpAddressRanges"), s)
	rd.Set(common.ToSnakeCase("ObsBucketAccessUrl"), info.ObsBucketAccessUrl)
	rd.Set(common.ToSnakeCase("ObsBucketDrType"), info.ObsBucketDrType)
	rd.Set(common.ToSnakeCase("ObsBucketFileEncryptionAlgorithm"), info.ObsBucketFileEncryptionAlgorithm)
	rd.Set(common.ToSnakeCase("ObsBucketFileEncryptionEnabled"), info.ObsBucketFileEncryptionEnabled)
	rd.Set(common.ToSnakeCase("ObsBucketFileEncryptionType"), info.ObsBucketFileEncryptionType)
	rd.Set(common.ToSnakeCase("ObsBucketId"), info.ObsBucketId)
	rd.Set(common.ToSnakeCase("ObsBucketName"), info.ObsBucketName)
	rd.Set(common.ToSnakeCase("ObsBucketState"), info.ObsBucketState)
	rd.Set(common.ToSnakeCase("ObsBucketUsedSize"), info.ObsBucketUsedSize)
	rd.Set(common.ToSnakeCase("ObsBucketUsedType"), info.ObsBucketUsedType)
	rd.Set(common.ToSnakeCase("ObsBucketVersionEnabled"), info.ObsBucketVersionEnabled)
	rd.Set(common.ToSnakeCase("ObsId"), info.ObsId)
	rd.Set(common.ToSnakeCase("ObsName"), info.ObsName)
	rd.Set(common.ToSnakeCase("ObsQuotaId"), info.ObsQuotaId)
	rd.Set(common.ToSnakeCase("ObsSyncBucketId"), info.ObsSyncBucketId)
	rd.Set(common.ToSnakeCase("ObsSyncBucketName"), info.ObsSyncBucketName)
	rd.Set(common.ToSnakeCase("ObsSyncBucketObsName"), info.ObsSyncBucketObsName)
	rd.Set(common.ToSnakeCase("ObsSyncBucketRegion"), info.ObsSyncBucketRegion)
	rd.Set(common.ToSnakeCase("ObsSyncBucketZoneName"), info.ObsSyncBucketZoneName)
	rd.Set(common.ToSnakeCase("ObsUrls"), m)
	rd.Set(common.ToSnakeCase("ProjectName"), info.ProjectName)
	rd.Set(common.ToSnakeCase("Region"), info.Region)
	rd.Set(common.ToSnakeCase("SystemId"), info.SystemId)
	rd.Set(common.ToSnakeCase("SystemName"), info.SystemName)
	rd.Set(common.ToSnakeCase("ZoneId"), info.ZoneId)
	rd.Set(common.ToSnakeCase("ZoneName"), info.ZoneName)
	rd.Set(common.ToSnakeCase("CreatedBy"), info.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), info.CreatedDt)
	rd.Set(common.ToSnakeCase("ModifiedBy"), info.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), info.ModifiedDt)

	rd.SetId(uuid.NewV4().String())

	return nil
}
