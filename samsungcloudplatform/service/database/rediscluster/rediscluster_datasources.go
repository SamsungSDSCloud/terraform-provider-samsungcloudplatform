package rediscluster

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
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_redis_clusters", DatasourceRedisClusters())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_redis_cluster", DatasourceRedisCluster())
}

func DatasourceRedisClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedisClusterList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"redis_name":  {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"page":        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort":        {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "Redis Cluster list", Elem: datasourceRedisClusterElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of Redis Cluster Servers.",
	}
}

func datasourceRedisClusterElem() interface{} {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":      {Type: schema.TypeString, Computed: true, Description: "Project ID."},
			"block_id":        {Type: schema.TypeString, Computed: true, Description: "Block ID."},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"redis_id":        {Type: schema.TypeString, Computed: true, Description: "Redis Cluster ID"},
			"redis_name":      {Type: schema.TypeString, Computed: true, Description: "Redis Cluster Name"},
			"redis_state":     {Type: schema.TypeString, Computed: true, Description: "Redis Cluster State"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceRedisClusterList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.RedisCluster.ListRedisCluster(ctx, &redis.RedisClusterSearchApiListRedisClusterOpts{
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

func DatasourceRedisCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedisClusterSingle,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"redis_cluster_id":   {Type: schema.TypeString, Required: true, Description: "Redis Cluster Id"},
			"redis_name":         {Type: schema.TypeString, Computed: true, Description: "Redis Cluster Name"},
			"redis_state":        {Type: schema.TypeString, Computed: true, Description: "Redis Cluster State"},
			"image_id":           {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			"timezone":           {Type: schema.TypeString, Computed: true, Description: "Timezone"},
			"vpc_id":             {Type: schema.TypeString, Computed: true, Description: "VPC Id"},
			"subnet_id":          {Type: schema.TypeString, Computed: true, Description: "Subnet Id"},
			"security_group_ids": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Security group ids"},
			"contract": {Type: schema.TypeList, Computed: true, Description: "contract",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_period":        {Type: schema.TypeString, Computed: true, Description: "Contract period"},
						"contract_start_date":    {Type: schema.TypeString, Computed: true, Description: "Contract start date"},
						"contract_end_date":      {Type: schema.TypeString, Computed: true, Description: "Contract end date"},
						"next_contract_period":   {Type: schema.TypeString, Computed: true, Description: "Next contract period"},
						"next_contract_end_date": {Type: schema.TypeString, Computed: true, Description: "Next contract end date"},
					},
				},
			},
			"redis_initial_config": {Type: schema.TypeList, Computed: true, Description: "Redis Cluster initial config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_port": {Type: schema.TypeInt, Computed: true, Description: "database port"},
					},
				},
			},
			"redis_server_group": {Type: schema.TypeList, Computed: true, Description: "Redis Cluster server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_enabled":     {Type: schema.TypeBool, Computed: true, Description: "encryption enabled"},
						"server_group_role_type": {Type: schema.TypeString, Computed: true, Description: "server group role type"},
						"server_type":            {Type: schema.TypeString, Computed: true, Description: "server type"},
						"redis_servers": {Type: schema.TypeList, Computed: true, Description: "Redis Cluster servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redis_server_id":       {Type: schema.TypeString, Computed: true, Description: "Redis Cluster server id"},
									"redis_server_name":     {Type: schema.TypeString, Computed: true, Description: "Redis Cluster server name"},
									"redis_server_state":    {Type: schema.TypeString, Computed: true, Description: "Redis Cluster server state"},
									"server_role_type":      {Type: schema.TypeString, Computed: true, Description: "server role type"},
									"subnet_ip_address":     {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
									"nat_public_ip_address": {Type: schema.TypeString, Computed: true, Description: "nat ip address"},
									"created_by":            {Type: schema.TypeString, Computed: true, Description: "created by"},
									"created_dt":            {Type: schema.TypeString, Computed: true, Description: "created dt"},
									"modified_by":           {Type: schema.TypeString, Computed: true, Description: "modified by"},
									"modified_dt":           {Type: schema.TypeString, Computed: true, Description: "modified dt"},
								},
							},
						},
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
									"object_storage_bucket_id": {Type: schema.TypeString, Computed: true, Description: "object storage bucket id"},
									"backup_retention_period":  {Type: schema.TypeString, Computed: true, Description: "backup retention period"},
									"backup_start_hour":        {Type: schema.TypeInt, Computed: true, Description: "backup start hour"},
								},
							},
						},
					},
				},
			},
			"project_id":      {Type: schema.TypeString, Computed: true, Description: "project id"},
			"block_id":        {Type: schema.TypeString, Computed: true, Description: "block id"},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "service zone id"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "created by"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "created dt"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "modified by"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "modified dt"},
		},
		Description: "Search single redis cluster database.",
	}
}

func dataSourceRedisClusterSingle(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.RedisCluster.DetailRedisCluster(ctx, rd.Get("redis_cluster_id").(string))
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

	redisClusterInitialConfig := database_common.HclListObject{}
	if dbInfo.RedisInitialConfig != nil {
		redisClusterInitialConfigInfo := database_common.HclKeyValueObject{}
		redisClusterInitialConfigInfo["database_port"] = dbInfo.RedisInitialConfig.DatabasePort

		redisClusterInitialConfig = append(redisClusterInitialConfig, redisClusterInitialConfigInfo)
	}

	err = rd.Set("redis_initial_config", redisClusterInitialConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	redisCluterServerGroup := database_common.HclListObject{}
	if dbInfo.RedisServerGroup != nil {
		redisClusterServerGroupInfo := database_common.HclKeyValueObject{}
		redisClusterServerGroupInfo["server_group_role_type"] = dbInfo.RedisServerGroup.ServerGroupRoleType
		redisClusterServerGroupInfo["server_type"] = dbInfo.RedisServerGroup.ServerType
		redisClusterServerGroupInfo["encryption_enabled"] = dbInfo.RedisServerGroup.EncryptionEnabled

		redisClusterServers := database_common.HclListObject{}
		for _, value := range dbInfo.RedisServerGroup.RedisServers {
			redisClusterServersInfo := database_common.HclKeyValueObject{}
			redisClusterServersInfo["redis_server_id"] = value.RedisServerId
			redisClusterServersInfo["redis_server_name"] = value.RedisServerName
			redisClusterServersInfo["redis_server_state"] = value.RedisServerState
			redisClusterServersInfo["server_role_type"] = value.ServerRoleType
			redisClusterServersInfo["subnet_ip_address"] = value.SubnetIpAddress
			redisClusterServersInfo["nat_public_ip_address"] = value.NatPublicIpAddress
			redisClusterServersInfo["created_by"] = value.CreatedBy
			redisClusterServersInfo["created_dt"] = value.CreatedDt.String()
			redisClusterServersInfo["modified_by"] = value.ModifiedBy
			redisClusterServersInfo["modified_dt"] = value.ModifiedDt.String()

			redisClusterServers = append(redisClusterServers, redisClusterServersInfo)
		}
		redisClusterServerGroupInfo["redis_servers"] = redisClusterServers

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
		redisClusterServerGroupInfo["block_storages"] = blockStorages

		redisCluterServerGroup = append(redisCluterServerGroup, redisClusterServerGroupInfo)
	}

	err = rd.Set("redis_server_group", redisCluterServerGroup)
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

	err = rd.Set("project_id", dbInfo.ProjectId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("service_zone_id", dbInfo.ServiceZoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("block_id", dbInfo.BlockId)
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
