package mariadb

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/database/database_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/mariadb"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_mariadbs", DatasourceMariadbList())
	scp.RegisterDataSource("scp_mariadb", DatasourceMariadb())
}

func DatasourceMariadbList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMariadbList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"mariadb_cluster_name": {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"page":                 {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":                 {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort":                 {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":             {Type: schema.TypeList, Optional: true, Description: "Mariadb list", Elem: datasourceMariadbElem()},
			"total_count":          {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of mariadb databases.",
	}
}

func datasourceMariadbElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":            {Type: schema.TypeString, Computed: true, Description: "Project ID."},
			"block_id":              {Type: schema.TypeString, Computed: true, Description: "Block ID."},
			"service_zone_id":       {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"mariadb_cluster_id":    {Type: schema.TypeString, Computed: true, Description: "Mariadb Cluster ID"},
			"mariadb_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "Mariadb Cluster Name"},
			"mariadb_cluster_state": {Type: schema.TypeString, Computed: true, Description: "Mariadb Cluster State"},
			"created_by":            {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":            {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":           {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":           {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceMariadbList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Mariadb.ListMariadbClusters(ctx, &mariadb.MariadbSearchApiListMariadbClustersOpts{
		MariadbClusterName: optional.NewString(rd.Get("mariadb_cluster_name").(string)),
		Page:               optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:               optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:               optional.NewInterface(rd.Get("sort").(string)),
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

func DatasourceMariadb() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMariadbSingle,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id":            {Type: schema.TypeString, Computed: true, Description: "project id"},
			"block_id":              {Type: schema.TypeString, Computed: true, Description: "block id"},
			"service_zone_id":       {Type: schema.TypeString, Computed: true, Description: "service zone id"},
			"mariadb_cluster_id":    {Type: schema.TypeString, Required: true, Description: "mariadb Cluster Id"},
			"mariadb_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "mariadb Cluster Name"},
			"mariadb_cluster_state": {Type: schema.TypeString, Computed: true, Description: "mariadb Cluster State"},
			"image_id":              {Type: schema.TypeString, Computed: true, Description: "image Id"},
			"database_version":      {Type: schema.TypeString, Computed: true, Description: "database version"},
			"timezone":              {Type: schema.TypeString, Computed: true, Description: "timezone"},
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
			"vpc_id":             {Type: schema.TypeString, Computed: true, Description: "vPC Id"},
			"subnet_id":          {Type: schema.TypeString, Computed: true, Description: "subnet Id"},
			"security_group_ids": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "security group ids"},
			"audit_enabled":      {Type: schema.TypeBool, Computed: true, Description: "audit enabled"},
			"nat_ip_address":     {Type: schema.TypeString, Computed: true, Description: "nat ip address"},
			"mariadb_initial_config": {Type: schema.TypeList, Computed: true, Description: "mariadb initial config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name":          {Type: schema.TypeString, Computed: true, Description: "database name"},
						"database_user_name":     {Type: schema.TypeString, Computed: true, Description: "database user name"},
						"database_port":          {Type: schema.TypeInt, Computed: true, Description: "database port"},
						"database_character_set": {Type: schema.TypeString, Computed: true, Description: "database character set"},
					},
				},
			},
			"mariadb_server_group": {Type: schema.TypeList, Computed: true, Description: "mariadb server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_role_type": {Type: schema.TypeString, Computed: true, Description: "server group role type"},
						"server_type":            {Type: schema.TypeString, Computed: true, Description: "server type"},
						"mariadb_servers": {Type: schema.TypeList, Computed: true, Description: "mariadb servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mariadb_server_id":      {Type: schema.TypeString, Computed: true, Description: "mariadb server id"},
									"mariadb_server_name":    {Type: schema.TypeString, Computed: true, Description: "mariadb server name"},
									"mariadb_server_state":   {Type: schema.TypeString, Computed: true, Description: "mariadb server state"},
									"server_role_type":       {Type: schema.TypeString, Computed: true, Description: "server role type"},
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "availability zone name"},
									"subnet_ip_address":      {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
									"created_by":             {Type: schema.TypeString, Computed: true, Description: "created by"},
									"created_dt":             {Type: schema.TypeString, Computed: true, Description: "created dt"},
									"modified_by":            {Type: schema.TypeString, Computed: true, Description: "modified by"},
									"modified_dt":            {Type: schema.TypeString, Computed: true, Description: "modified dt"},
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
						"virtual_ip_address": {Type: schema.TypeString, Computed: true, Description: "virtual ip address"},
					},
				},
			},
			"mariadb_replica_cluster_ids": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "mariadb replica cluster ids"},
			"mariadb_master_cluster_id":   {Type: schema.TypeString, Computed: true, Description: "mariadb master cluster id"},
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
						"incremental_backup_config": {Type: schema.TypeList, Computed: true, Description: "incremental_backup_config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_storage_bucket_id":          {Type: schema.TypeString, Computed: true, Description: "object storage bucket id"},
									"archive_backup_schedule_frequency": {Type: schema.TypeString, Computed: true, Description: "archive backup schedule frequency"},
									"backup_retention_period":           {Type: schema.TypeString, Computed: true, Description: "backup retention period"},
									"backup_schedule_frequency":         {Type: schema.TypeString, Computed: true, Description: "backup schedule frequency"},
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
		Description: "Search single mariadb database.",
	}
}

func dataSourceMariadbSingle(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Mariadb.DetailMariadbCluster(ctx, rd.Get("mariadb_cluster_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(dbInfo.MariadbServerGroup.MariadbServers) == 0 {
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
	err = rd.Set("mariadb_cluster_id", dbInfo.MariadbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mariadb_cluster_name", dbInfo.MariadbClusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mariadb_cluster_state", dbInfo.MariadbClusterState)
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
	err = rd.Set("audit_enabled", dbInfo.AuditEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_ip_address", dbInfo.NatIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}

	mariadbInitialConfig := database_common.HclListObject{}
	if dbInfo.MariadbInitialConfig != nil {
		mariadbInitialConfigInfo := database_common.HclKeyValueObject{}
		mariadbInitialConfigInfo["database_name"] = dbInfo.MariadbInitialConfig.DatabaseName
		mariadbInitialConfigInfo["database_user_name"] = dbInfo.MariadbInitialConfig.DatabaseUserName
		mariadbInitialConfigInfo["database_port"] = dbInfo.MariadbInitialConfig.DatabasePort
		mariadbInitialConfigInfo["database_character_set"] = dbInfo.MariadbInitialConfig.DatabaseCharacterSet

		mariadbInitialConfig = append(mariadbInitialConfig, mariadbInitialConfigInfo)
	}
	err = rd.Set("mariadb_initial_config", mariadbInitialConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	mariadbServerGroup := database_common.HclListObject{}
	if dbInfo.MariadbServerGroup != nil {
		mariadbServerGroupInfo := database_common.HclKeyValueObject{}
		mariadbServerGroupInfo["server_group_role_type"] = dbInfo.MariadbServerGroup.ServerGroupRoleType
		mariadbServerGroupInfo["server_type"] = dbInfo.MariadbServerGroup.ServerType

		mariadbServers := database_common.HclListObject{}
		for _, value := range dbInfo.MariadbServerGroup.MariadbServers {
			mariadbServersInfo := database_common.HclKeyValueObject{}
			mariadbServersInfo["mariadb_server_id"] = value.MariadbServerId
			mariadbServersInfo["mariadb_server_name"] = value.MariadbServerName
			mariadbServersInfo["availability_zone_name"] = value.AvailabilityZoneName
			mariadbServersInfo["server_role_type"] = value.ServerRoleType
			mariadbServersInfo["subnet_ip_address"] = value.SubnetIpAddress
			mariadbServersInfo["created_by"] = value.CreatedBy
			mariadbServersInfo["created_dt"] = value.CreatedDt.String()
			mariadbServersInfo["modified_by"] = value.ModifiedBy
			mariadbServersInfo["modified_dt"] = value.ModifiedDt.String()

			mariadbServers = append(mariadbServers, mariadbServersInfo)
		}
		mariadbServerGroupInfo["mariadb_servers"] = mariadbServers

		mariadbServerGroupInfo["encryption_enabled"] = dbInfo.MariadbServerGroup.EncryptionEnabled

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.MariadbServerGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize

			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		mariadbServerGroupInfo["block_storages"] = blockStorages

		mariadbServerGroupInfo["virtual_ip_address"] = dbInfo.MariadbServerGroup.VirtualIpAddress

		mariadbServerGroup = append(mariadbServerGroup, mariadbServerGroupInfo)
	}

	err = rd.Set("mariadb_server_group", mariadbServerGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("mariadb_replica_cluster_ids", dbInfo.MariadbReplicaClusterIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mariadb_master_cluster_id", dbInfo.MariadbMasterClusterId)
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
			fullBackupConfigInfo["archive_backup_schedule_frequency"] = dbInfo.BackupConfig.FullBackupConfig.ArchiveBackupScheduleFrequency
			fullBackupConfigInfo["backup_retention_period"] = dbInfo.BackupConfig.FullBackupConfig.BackupRetentionPeriod
			fullBackupConfigInfo["backup_start_hour"] = dbInfo.BackupConfig.FullBackupConfig.BackupStartHour

			fullBackupConfig = append(fullBackupConfig, fullBackupConfigInfo)
		}
		backupConfigInfo["full_backup_config"] = fullBackupConfig

		incrementalBackupConfig := database_common.HclListObject{}
		if dbInfo.BackupConfig.IncrementalBackupConfig != nil {
			incrementalBackupConfigInfo := database_common.HclKeyValueObject{}
			incrementalBackupConfigInfo["object_storage_bucket_id"] = dbInfo.BackupConfig.IncrementalBackupConfig.ObjectStorageBucketId
			incrementalBackupConfigInfo["archive_backup_schedule_frequency"] = dbInfo.BackupConfig.IncrementalBackupConfig.ArchiveBackupScheduleFrequency
			incrementalBackupConfigInfo["backup_retention_period"] = dbInfo.BackupConfig.IncrementalBackupConfig.BackupRetentionPeriod
			incrementalBackupConfigInfo["backup_schedule_frequency"] = dbInfo.BackupConfig.IncrementalBackupConfig.BackupScheduleFrequency

			incrementalBackupConfig = append(incrementalBackupConfig, incrementalBackupConfigInfo)
		}
		backupConfigInfo["incremental_backup_config"] = incrementalBackupConfig
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
