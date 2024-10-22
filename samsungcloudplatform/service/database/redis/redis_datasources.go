package redis

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/redis"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_redis_list", DatasourceRedisList())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_redis", DatasourceRedis())
}

func DatasourceRedisList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedisList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"redis_name":  {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"page":        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort":        {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "Redis list", Elem: datasourceRedisElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of redis databases.",
	}
}

func datasourceRedisElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":      {Type: schema.TypeString, Computed: true, Description: "Project ID."},
			"block_id":        {Type: schema.TypeString, Computed: true, Description: "Block ID."},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"redis_id":        {Type: schema.TypeString, Computed: true, Description: "Redis ID"},
			"redis_name":      {Type: schema.TypeString, Computed: true, Description: "Redis Name"},
			"redis_state":     {Type: schema.TypeString, Computed: true, Description: "Redis State"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceRedisList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Redis.ListRedis(ctx, &redis.RedisSearchApiListRedisOpts{
		RedisName: optional.NewString(rd.Get("redis_name").(string)),
		Page:      optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:      optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:      optional.NewInterface(rd.Get("sort").(string)),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	err = rd.Set("contents", contents)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("total_count", responses.TotalCount)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DatasourceRedis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedisSingle,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id":         {Type: schema.TypeString, Computed: true, Description: "project id"},
			"block_id":           {Type: schema.TypeString, Computed: true, Description: "block id"},
			"service_zone_id":    {Type: schema.TypeString, Computed: true, Description: "service zone id"},
			"redis_id":           {Type: schema.TypeString, Required: true, Description: "redis  Id"},
			"redis_name":         {Type: schema.TypeString, Computed: true, Description: "redis  Name"},
			"redis_state":        {Type: schema.TypeString, Computed: true, Description: "redis  State"},
			"image_id":           {Type: schema.TypeString, Computed: true, Description: "image Id"},
			"database_version":   {Type: schema.TypeString, Computed: true, Description: "database version"},
			"timezone":           {Type: schema.TypeString, Computed: true, Description: "timezone"},
			"vpc_id":             {Type: schema.TypeString, Computed: true, Description: "vPC Id"},
			"subnet_id":          {Type: schema.TypeString, Computed: true, Description: "subnet Id"},
			"security_group_ids": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "security group ids"},
			"contract": {Type: schema.TypeList, Computed: true, Description: "contract",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_period":        {Type: schema.TypeString, Computed: true, Description: "contract period"},
						"contract_start_date":    {Type: schema.TypeString, Computed: true, Description: "contract start date"},
						"contract_end_date":      {Type: schema.TypeString, Computed: true, Description: "contract end date"},
						"next_contract_period":   {Type: schema.TypeString, Computed: true, Description: "next contract period"},
						"next_contract_end_date": {Type: schema.TypeString, Computed: true, Description: "next contract end date"},
					},
				},
			},
			"nat_ip_address": {Type: schema.TypeString, Computed: true, Description: "nat ip address"},
			"redis_initial_config": {Type: schema.TypeList, Computed: true, Description: "redis initial config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_port": {Type: schema.TypeInt, Computed: true, Description: "database port"},
						"sentinel_port": {Type: schema.TypeInt, Computed: true, Description: "sentinel port"},
					},
				},
			},
			"redis_server_group": {Type: schema.TypeList, Computed: true, Description: "redis server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_role_type": {Type: schema.TypeString, Computed: true, Description: "server group role type"},
						"server_type":            {Type: schema.TypeString, Computed: true, Description: "server type"},
						"redis_servers": {Type: schema.TypeList, Computed: true, Description: "redis servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redis_server_id":       {Type: schema.TypeString, Computed: true, Description: "redis server id"},
									"redis_server_name":     {Type: schema.TypeString, Computed: true, Description: "redis server name"},
									"server_role_type":      {Type: schema.TypeString, Computed: true, Description: "server role type"},
									"redis_server_state":    {Type: schema.TypeString, Computed: true, Description: "redis server state"},
									"subnet_ip_address":     {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
									"nat_public_ip_address": {Type: schema.TypeString, Computed: true, Description: "nat public ip address"},
									"created_by":            {Type: schema.TypeString, Computed: true, Description: "created by"},
									"created_dt":            {Type: schema.TypeString, Computed: true, Description: "created dt"},
									"modified_by":           {Type: schema.TypeString, Computed: true, Description: "modified by"},
									"modified_dt":           {Type: schema.TypeString, Computed: true, Description: "modified dt"},
								},
							},
						},
						"encryption_enabled": {Type: schema.TypeBool, Computed: true, Description: "encryption enabled"},
						"block_storages": {Type: schema.TypeList, Computed: true, Description: "block storages",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_storage_group_id":  {Type: schema.TypeString, Computed: true, Description: "block storage group id"},
									"block_storage_name":      {Type: schema.TypeString, Computed: true, Description: "block storage name"},
									"block_storage_role_type": {Type: schema.TypeString, Computed: true, Description: "block storage role type"},
									"block_storage_type":      {Type: schema.TypeString, Computed: true, Description: "block storage type"},
									"block_storage_size":      {Type: schema.TypeInt, Computed: true, Description: "block Storage size"},
								},
							},
						},
					},
				},
			},
			"sentinel_server": {Type: schema.TypeList, Computed: true, Description: "redis server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_type":           {Type: schema.TypeString, Computed: true, Description: "server type"},
						"sentinel_server_id":    {Type: schema.TypeString, Computed: true, Description: "sentinel server id"},
						"sentinel_server_name":  {Type: schema.TypeString, Computed: true, Description: "sentinel server name"},
						"sentinel_server_state": {Type: schema.TypeString, Computed: true, Description: "sentinel server state"},
						"subnet_ip_address":     {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
						"nat_public_ip_address": {Type: schema.TypeString, Computed: true, Description: "nat public ip address"},
						"encryption_enabled":    {Type: schema.TypeBool, Computed: true, Description: "encryption enabled"},
						"block_storages": {Type: schema.TypeList, Computed: true, Description: "block storages",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_storage_group_id":  {Type: schema.TypeString, Computed: true, Description: "block storage group id"},
									"block_storage_name":      {Type: schema.TypeString, Computed: true, Description: "block storage name"},
									"block_storage_role_type": {Type: schema.TypeString, Computed: true, Description: "block storage role type"},
									"block_storage_type":      {Type: schema.TypeString, Computed: true, Description: "block storage type"},
									"block_storage_size":      {Type: schema.TypeInt, Computed: true, Description: "block Storage size"},
								},
							},
						},
						"created_by":  {Type: schema.TypeString, Computed: true, Description: "created by"},
						"created_dt":  {Type: schema.TypeString, Computed: true, Description: "created dt"},
						"modified_by": {Type: schema.TypeString, Computed: true, Description: "modified by"},
						"modified_dt": {Type: schema.TypeString, Computed: true, Description: "modified dt"},
					},
				},
			},
			"maintenance": {Type: schema.TypeList, Computed: true, Description: "maintenance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maintenance_start_day_of_week": {Type: schema.TypeString, Computed: true, Description: "maintenance start day of week"},
						"maintenance_start_time":        {Type: schema.TypeString, Computed: true, Description: "maintenance start time"},
						"maintenance_period":            {Type: schema.TypeInt, Computed: true, Description: "maintenance period"},
					},
				},
			},
			"backup_config": {Type: schema.TypeList, Computed: true, Description: "backup config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_backup_config": {Type: schema.TypeList, Computed: true, Description: "full backup config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_storage_bucket_id":          {Type: schema.TypeString, Computed: true, Description: "object storage bucket id"},
									"archive_backup_schedule_frequency": {Type: schema.TypeString, Computed: true, Description: "archive backup schedule frequency"},
									"backup_retention_period":           {Type: schema.TypeString, Computed: true, Description: "backup retention period"},
									"backup_start_hour":                 {Type: schema.TypeInt, Computed: true, Description: "backup start hour"},
								},
							},
						},
					},
				},
			},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "created by"},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "created dt"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "modified by"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "modified dt"},
		},
		Description: "Search single redis database.",
	}
}

func dataSourceRedisSingle(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Redis.DetailRedis(ctx, rd.Get("redis_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(dbInfo.RedisServerGroup.RedisServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	rd.SetId(uuid.NewV4().String())
	err = rd.Set("project_id", dbInfo.ProjectId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("block_id", dbInfo.BlockId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("service_zone_id", dbInfo.ServiceZoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_id", dbInfo.RedisId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_name", dbInfo.RedisName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_state", dbInfo.RedisState)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("image_id", dbInfo.ImageId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_version", dbInfo.DatabaseVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("timezone", dbInfo.Timezone)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("vpc_id", dbInfo.VpcId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("subnet_id", dbInfo.SubnetId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("security_group_ids", dbInfo.SecurityGroupIds)
	if err != nil {
		return diag.FromErr(err)
	}

	contract := database_common.HclListObject{}
	if dbInfo.Contract != nil {
		contractInfo := database_common.HclKeyValueObject{}
		contractInfo["contract_period"] = dbInfo.Contract.ContractPeriod
		contractInfo["contract_start_date"] = dbInfo.Contract.ContractStartDate
		contractInfo["contract_end_date"] = dbInfo.Contract.ContractEndDate
		contractInfo["next_contract_period"] = dbInfo.Contract.NextContractPeriod
		contractInfo["next_contract_end_date"] = dbInfo.Contract.NextContractEndDate

		contract = append(contract, contractInfo)
	}
	err = rd.Set("contract", contract)
	if err != nil {
		return diag.FromErr(err)
	}

	redisInitialConfig := database_common.HclListObject{}
	if dbInfo.RedisInitialConfig != nil {
		redisInitialConfigInfo := database_common.HclKeyValueObject{}
		redisInitialConfigInfo["database_port"] = dbInfo.RedisInitialConfig.DatabasePort
		redisInitialConfigInfo["sentinel_port"] = dbInfo.RedisInitialConfig.SentinelPort

		redisInitialConfig = append(redisInitialConfig, redisInitialConfigInfo)
	}
	err = rd.Set("redis_initial_config", redisInitialConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	redisServerGroup := database_common.HclListObject{}
	if dbInfo.RedisServerGroup != nil {
		redisServerGroupInfo := database_common.HclKeyValueObject{}
		redisServerGroupInfo["server_group_role_type"] = dbInfo.RedisServerGroup.ServerGroupRoleType
		redisServerGroupInfo["server_type"] = dbInfo.RedisServerGroup.ServerType

		redisServers := database_common.HclListObject{}
		for _, value := range dbInfo.RedisServerGroup.RedisServers {
			redisServersInfo := database_common.HclKeyValueObject{}
			redisServersInfo["redis_server_id"] = value.RedisServerId
			redisServersInfo["redis_server_name"] = value.RedisServerName
			redisServersInfo["server_role_type"] = value.ServerRoleType
			redisServersInfo["redis_server_state"] = value.RedisServerState
			redisServersInfo["subnet_ip_address"] = value.SubnetIpAddress
			redisServersInfo["nat_public_ip_address"] = value.NatPublicIpAddress
			redisServersInfo["created_by"] = value.CreatedBy
			redisServersInfo["created_dt"] = value.CreatedDt.String()
			redisServersInfo["modified_by"] = value.ModifiedBy
			redisServersInfo["modified_dt"] = value.ModifiedDt.String()

			redisServers = append(redisServers, redisServersInfo)
		}
		redisServerGroupInfo["redis_servers"] = redisServers

		redisServerGroupInfo["encryption_enabled"] = dbInfo.RedisServerGroup.EncryptionEnabled

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.RedisServerGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize

			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		redisServerGroupInfo["block_storages"] = blockStorages

		redisServerGroup = append(redisServerGroup, redisServerGroupInfo)
	}
	err = rd.Set("redis_server_group", redisServerGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	sentinelServer := database_common.HclListObject{}
	if dbInfo.SentinelServer != nil {
		sentinelServerInfo := database_common.HclKeyValueObject{}
		sentinelServerInfo["server_type"] = dbInfo.SentinelServer.ServerType
		sentinelServerInfo["sentinel_server_id"] = dbInfo.SentinelServer.SentinelServerId
		sentinelServerInfo["sentinel_server_name"] = dbInfo.SentinelServer.SentinelServerName
		sentinelServerInfo["sentinel_server_state"] = dbInfo.SentinelServer.SentinelServerState
		sentinelServerInfo["subnet_ip_address"] = dbInfo.SentinelServer.SubnetIpAddress
		sentinelServerInfo["nat_public_ip_address"] = dbInfo.SentinelServer.NatPublicIpAddress
		sentinelServerInfo["encryption_enabled"] = dbInfo.SentinelServer.EncryptionEnabled

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.SentinelServer.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize

			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		sentinelServerInfo["block_storages"] = blockStorages
		sentinelServerInfo["created_by"] = dbInfo.SentinelServer.CreatedBy
		sentinelServerInfo["created_dt"] = dbInfo.SentinelServer.CreatedDt.String()
		sentinelServerInfo["modified_by"] = dbInfo.SentinelServer.ModifiedBy
		sentinelServerInfo["modified_dt"] = dbInfo.SentinelServer.ModifiedDt.String()

		sentinelServer = append(sentinelServer, sentinelServerInfo)
	}
	err = rd.Set("sentinel_server", sentinelServer)
	if err != nil {
		return diag.FromErr(err)
	}

	maintenance := database_common.HclListObject{}
	if dbInfo.Maintenance != nil {
		maintenanceInfo := database_common.HclKeyValueObject{}
		maintenanceInfo["maintenance_start_day_of_week"] = dbInfo.Maintenance.MaintenanceStartDayOfWeek
		maintenanceInfo["maintenance_start_time"] = dbInfo.Maintenance.MaintenanceStartTime
		maintenanceInfo["maintenance_period"] = dbInfo.Maintenance.MaintenancePeriod

		maintenance = append(maintenance, maintenanceInfo)
	}
	err = rd.Set("maintenance", maintenance)
	if err != nil {
		return diag.FromErr(err)
	}

	backupConfig := database_common.HclListObject{}
	if dbInfo.BackupConfig != nil {
		backupConfigInfo := database_common.HclKeyValueObject{}

		fullBackupConfig := database_common.HclListObject{}
		if dbInfo.BackupConfig.FullBackupConfig != nil {
			fullBackupConfigInfo := database_common.HclKeyValueObject{}
			fullBackupConfigInfo["object_storage_bucket_id"] = dbInfo.BackupConfig.FullBackupConfig.ObjectStorageBucketId
			fullBackupConfigInfo["backup_retention_period"] = dbInfo.BackupConfig.FullBackupConfig.BackupRetentionPeriod
			fullBackupConfigInfo["backup_start_hour"] = dbInfo.BackupConfig.FullBackupConfig.BackupStartHour

			fullBackupConfig = append(fullBackupConfig, fullBackupConfigInfo)
		}
		backupConfigInfo["full_backup_config"] = fullBackupConfig
		backupConfig = append(backupConfig, backupConfigInfo)
	}
	err = rd.Set("backup_config", backupConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("created_by", dbInfo.CreatedBy)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("created_dt", dbInfo.CreatedDt.String())
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("modified_by", dbInfo.ModifiedBy)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("modified_dt", dbInfo.ModifiedDt.String())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
