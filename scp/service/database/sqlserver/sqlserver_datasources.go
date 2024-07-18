package sqlserver

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/database/database_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_sqlservers", DatasourceSqlServers())
	scp.RegisterDataSource("scp_sqlserver", DatasourceSqlserver())
}

func DatasourceSqlServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlServerList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"sqlserver_cluster_name": {Type: schema.TypeString, Optional: true, Description: "Database name."},
			"page":                   {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list."},
			"size":                   {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list."},
			"sort":                   {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":               {Type: schema.TypeList, Optional: true, Description: "MS SQL Server list", Elem: datasourceSqlserverElem()},
			"total_count":            {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of Microsoft SQL Servers.",
	}
}

func datasourceSqlserverElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":              {Type: schema.TypeString, Computed: true, Description: "Project ID."},
			"block_id":                {Type: schema.TypeString, Computed: true, Description: "Block ID."},
			"service_zone_id":         {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"sqlserver_cluster_id":    {Type: schema.TypeString, Computed: true, Description: "MS SQL Server Cluster ID"},
			"sqlserver_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "MS SQL Server Cluster Name"},
			"sqlserver_cluster_state": {Type: schema.TypeString, Computed: true, Description: "MS SQL Server Cluster State"},
			"created_by":              {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":              {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":             {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":             {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceSqlServerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Sqlserver.ListSqlserverClusters(ctx, &sqlserver.SqlserverSearchApiListSqlserverClustersOpts{
		SqlserverClusterName: optional.NewString(rd.Get("sqlserver_cluster_name").(string)),
		Page:                 optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                 optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                 optional.NewInterface(rd.Get("sort").(string)),
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

func DatasourceSqlserver() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlserverDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"sqlserver_cluster_id":    {Type: schema.TypeString, Required: true, Description: "MS SQL Server Cluster Id"},
			"sqlserver_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "MS SQL Server Cluster Name"},
			"sqlserver_cluster_state": {Type: schema.TypeString, Computed: true, Description: "MS SQL Server Cluster State"},
			"image_id":                {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			"timezone":                {Type: schema.TypeString, Computed: true, Description: "Timezone"},
			"vpc_id":                  {Type: schema.TypeString, Computed: true, Description: "VPC Id"},
			"subnet_id":               {Type: schema.TypeString, Computed: true, Description: "Subnet Id"},
			"security_group_ids":      {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Security group ids"},
			"nat_ip_address":          {Type: schema.TypeString, Computed: true, Description: "nat ip address"},
			"audit_enabled":           {Type: schema.TypeBool, Computed: true, Description: "audit enabled"},
			"database_version":        {Type: schema.TypeString, Computed: true, Description: "MS SQL Server version"},
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
			"sqlserver_initial_config": {Type: schema.TypeList, Computed: true, Description: "MS SQL Server initial config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_names":        {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Database Name List"},
						"database_user_name":    {Type: schema.TypeString, Computed: true, Description: "database user name"},
						"database_service_name": {Type: schema.TypeString, Required: true, Description: "MS SQL Server Database Service name"},
						"database_port":         {Type: schema.TypeInt, Computed: true, Description: "database port"},
						"database_collation":    {Type: schema.TypeString, Computed: true, Description: "Commands that specify how to sort and compare data"},
						"sqlserver_active_directory": {Type: schema.TypeSet, Computed: true, Description: "MS SQL Server Active directory",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_name":             {Type: schema.TypeString, Computed: true, Description: "Active Directory Domain name"},
									"domain_net_bios_name":    {Type: schema.TypeString, Computed: true, Description: "Active Directory NetBios name"},
									"dns_server_ips":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Active Directory DNS Server IPs"},
									"ad_server_user_id":       {Type: schema.TypeString, Computed: true, Description: "Active Directory Server User ID"},
									"ad_server_user_password": {Type: schema.TypeString, Computed: true, Sensitive: true, Description: "Active Directory Server User password"},
									"failover_cluster_name":   {Type: schema.TypeString, Computed: true, Description: "Active Directory Failover Cluster name"},
								},
							},
						},
					},
				},
			},
			"sqlserver_server_group": {Type: schema.TypeList, Computed: true, Description: "MS SQL Server server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_role_type": {Type: schema.TypeString, Computed: true, Description: "server group role type"},
						"server_type":            {Type: schema.TypeString, Computed: true, Description: "server type"},
						"sqlserver_servers": {Type: schema.TypeList, Computed: true, Description: "MS SQL Server servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqlserver_server_id":    {Type: schema.TypeString, Computed: true, Description: "MS SQL Server server id"},
									"sqlserver_server_name":  {Type: schema.TypeString, Computed: true, Description: "MS SQL Server server name"},
									"sqlserver_server_state": {Type: schema.TypeString, Computed: true, Description: "MS SQL Server server state"},
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "availability zone name"},
									"server_role_type":       {Type: schema.TypeString, Computed: true, Description: "server role type"},
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
			"quorum_server_group": {Type: schema.TypeList, Computed: true, Description: "MS SQL Server quorum server group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_role_type": {Type: schema.TypeString, Computed: true, Description: "server group role type"},
						"server_type":            {Type: schema.TypeString, Computed: true, Description: "server type"},
						"sqlserver_servers": {Type: schema.TypeList, Computed: true, Description: "MS SQL Server quorum servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqlserver_server_id":    {Type: schema.TypeString, Computed: true, Description: "MS SQL Server quorum server id"},
									"sqlserver_server_name":  {Type: schema.TypeString, Computed: true, Description: "MS SQL Server quorum server name"},
									"sqlserver_server_state": {Type: schema.TypeString, Computed: true, Description: "MS SQL Server quorum server state"},
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "availability zone name"},
									"server_role_type":       {Type: schema.TypeString, Computed: true, Description: "server role type"},
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
			"sqlserver_master_cluster_id":    {Type: schema.TypeString, Computed: true, Description: "MS SQL Server master cluster id"},
			"sqlserver_secondary_cluster_id": {Type: schema.TypeString, Computed: true, Description: "MS SQL Server secondary cluster id"},
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
									"full_backup_day_of_week":           {Type: schema.TypeString, Computed: true, Description: "Full backup schedule(Day)."},
								},
							},
						},
					},
				},
			},
			"project_id":      {Type: schema.TypeString, Computed: true, Description: "project id"},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "service zone id"},
			"block_id":        {Type: schema.TypeString, Computed: true, Description: "Block id"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "created by"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "created dt"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "modified by"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "modified dt"},
		},
		Description: "Search Detail MS SQL Server database.",
	}
}

func dataSourceSqlserverDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Sqlserver.DetailSqlserverCluster(ctx, rd.Get("sqlserver_cluster_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(dbInfo.SqlserverServerGroup.SqlserverServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	rd.SetId(uuid.NewV4().String())
	err = rd.Set("sqlserver_cluster_name", dbInfo.SqlserverClusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("sqlserver_cluster_state", dbInfo.SqlserverClusterState)
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
	err = rd.Set("database_version", dbInfo.DatabaseVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("security_group_ids", dbInfo.SecurityGroupIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_ip_address", dbInfo.NatIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("audit_enabled", dbInfo.AuditEnabled)
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

	sqlserverInitialConfig := database_common.HclListObject{}
	if dbInfo.SqlserverInitialConfig != nil {
		sqlserverInitialConfigInfo := database_common.HclKeyValueObject{}
		sqlserverInitialConfigInfo["database_names"] = dbInfo.SqlserverInitialConfig.DatabaseNames
		sqlserverInitialConfigInfo["database_user_name"] = dbInfo.SqlserverInitialConfig.DatabaseUserName
		sqlserverInitialConfigInfo["database_service_name"] = dbInfo.SqlserverInitialConfig.DatabaseServiceName
		sqlserverInitialConfigInfo["database_port"] = dbInfo.SqlserverInitialConfig.DatabasePort
		sqlserverInitialConfigInfo["database_collation"] = dbInfo.SqlserverInitialConfig.DatabaseCollation

		activeDirectory := database_common.HclListObject{}
		if dbInfo.SqlserverInitialConfig.ActiveDirectory != nil {
			activeDirectoryInfo := database_common.HclKeyValueObject{}
			activeDirectoryInfo["ad_server_user_id"] = dbInfo.SqlserverInitialConfig.ActiveDirectory.AdServerUserId
			activeDirectoryInfo["dns_server_ips"] = dbInfo.SqlserverInitialConfig.ActiveDirectory.DnsServerIps
			activeDirectoryInfo["domain_name"] = dbInfo.SqlserverInitialConfig.ActiveDirectory.DomainName
			activeDirectoryInfo["domain_net_bios_name"] = dbInfo.SqlserverInitialConfig.ActiveDirectory.DomainNetBiosName
			activeDirectoryInfo["failover_cluster_name"] = dbInfo.SqlserverInitialConfig.ActiveDirectory.FailoverClusterName

			activeDirectory = append(activeDirectory, activeDirectoryInfo)
		}
		sqlserverInitialConfigInfo["sqlserver_active_directory"] = activeDirectory

		sqlserverInitialConfig = append(sqlserverInitialConfig, sqlserverInitialConfigInfo)
	}
	err = rd.Set("sqlserver_initial_config", sqlserverInitialConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	sqlserverServerGroup := database_common.HclListObject{}
	if dbInfo.SqlserverServerGroup != nil {
		sqlserverServerGroupInfo := database_common.HclKeyValueObject{}
		sqlserverServerGroupInfo["server_group_role_type"] = dbInfo.SqlserverServerGroup.ServerGroupRoleType
		sqlserverServerGroupInfo["server_type"] = dbInfo.SqlserverServerGroup.ServerType

		sqlserverServers := database_common.HclListObject{}
		for _, value := range dbInfo.SqlserverServerGroup.SqlserverServers {
			sqlserverServersInfo := database_common.HclKeyValueObject{}
			sqlserverServersInfo["sqlserver_server_id"] = value.SqlserverServerId
			sqlserverServersInfo["sqlserver_server_name"] = value.SqlserverServerName
			sqlserverServersInfo["sqlserver_server_state"] = value.SqlserverServerState
			sqlserverServersInfo["availability_zone_name"] = value.AvailabilityZoneName
			sqlserverServersInfo["server_role_type"] = value.ServerRoleType
			sqlserverServersInfo["subnet_ip_address"] = value.SubnetIpAddress
			sqlserverServersInfo["created_by"] = value.CreatedBy
			sqlserverServersInfo["created_dt"] = value.CreatedDt.String()
			sqlserverServersInfo["modified_by"] = value.ModifiedBy
			sqlserverServersInfo["modified_dt"] = value.ModifiedDt.String()

			sqlserverServers = append(sqlserverServers, sqlserverServersInfo)
		}
		sqlserverServerGroupInfo["sqlserver_servers"] = sqlserverServers
		sqlserverServerGroupInfo["encryption_enabled"] = dbInfo.SqlserverServerGroup.EncryptionEnabled

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.SqlserverServerGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize

			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		sqlserverServerGroupInfo["block_storages"] = blockStorages
		sqlserverServerGroupInfo["virtual_ip_address"] = dbInfo.SqlserverServerGroup.VirtualIpAddress

		sqlserverServerGroup = append(sqlserverServerGroup, sqlserverServerGroupInfo)
	}
	err = rd.Set("sqlserver_server_group", sqlserverServerGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	quorumServerGroup := database_common.HclListObject{}
	if dbInfo.QuorumServerGroup != nil {
		quorumServerGroupInfo := database_common.HclKeyValueObject{}
		quorumServerGroupInfo["server_group_role_type"] = dbInfo.QuorumServerGroup.ServerGroupRoleType
		quorumServerGroupInfo["server_type"] = dbInfo.QuorumServerGroup.ServerType

		quorumServerServers := database_common.HclListObject{}
		for _, value := range dbInfo.QuorumServerGroup.SqlserverServers {
			quorumServersInfo := database_common.HclKeyValueObject{}
			quorumServersInfo["sqlserver_server_id"] = value.SqlserverServerId
			quorumServersInfo["sqlserver_server_name"] = value.SqlserverServerName
			quorumServersInfo["sqlserver_server_state"] = value.SqlserverServerState
			quorumServersInfo["availability_zone_name"] = value.AvailabilityZoneName
			quorumServersInfo["server_role_type"] = value.ServerRoleType
			quorumServersInfo["subnet_ip_address"] = value.SubnetIpAddress
			quorumServersInfo["created_by"] = value.CreatedBy
			quorumServersInfo["created_dt"] = value.CreatedDt.String()
			quorumServersInfo["modified_by"] = value.ModifiedBy
			quorumServersInfo["modified_dt"] = value.ModifiedDt.String()

			quorumServerServers = append(quorumServerServers, quorumServersInfo)
		}
		quorumServerGroupInfo["sqlserver_servers"] = quorumServerServers
		quorumServerGroupInfo["encryption_enabled"] = dbInfo.QuorumServerGroup.EncryptionEnabled

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.QuorumServerGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize

			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		quorumServerGroupInfo["block_storages"] = blockStorages
		quorumServerGroupInfo["virtual_ip_address"] = dbInfo.SqlserverServerGroup.VirtualIpAddress

		quorumServerGroup = append(quorumServerGroup, quorumServerGroupInfo)
	}
	err = rd.Set("quorum_server_group", quorumServerGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("sqlserver_master_cluster_id", dbInfo.SqlserverMasterClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("sqlserver_secondary_cluster_id", dbInfo.SqlserverSecondaryClusterId)
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
			fullBackupConfigInfo["full_backup_day_of_week"] = dbInfo.BackupConfig.FullBackupConfig.FullBackupDayOfWeek

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
