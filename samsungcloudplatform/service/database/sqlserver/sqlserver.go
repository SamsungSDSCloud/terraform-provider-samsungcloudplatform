package sqlserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"sort"
	"strings"

	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_sqlserver", ResourceSqlserver())
}

func ResourceSqlserver() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlserverCreate,
		ReadContext:   resourceSqlserverRead,
		UpdateContext: resourceSqlserverUpdate,
		DeleteContext: resourceSqlserverDelete,
		CustomizeDiff: resourceSqlserverDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"audit_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use database audit logging.",
			},
			"contract_period": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract : None, 1-year, 3-year",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
			},
			"next_contract_period": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				Description:      "Next contract : None, 1-year, 3-year",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL Server standard image id.",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use nat.",
			},
			"nat_public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public IP for NAT. If it is null, it is automatically allocated.",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this MS SQL Server DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service Zone Id",
			},
			"sqlserver_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"sqlserver_cluster_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "MS SQL Server cluster state",
				ValidateDiagFunc: database_common.ValidateStringInOptions("RUNNING", "STOPPED"),
			},
			"database_service_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "MS SQL Server Database Service name",
				ValidateDiagFunc: common.ValidateName1to15AlphaOnlyStartsWithUpperCase,
			},
			"database_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database Name List",
			},
			"database_user_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "User account id of database. (2 to 20 alpha-numerics)",
				ValidateDiagFunc: database_common.ValidateDbUserName2to20AlphaNumeric,
			},
			"database_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "User account password of database.",
				ValidateDiagFunc: database_common.ValidatePassword8to30WithSpecialsExceptQuotesAndDollar,
			},
			"license": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "License key.",
			},
			"database_port": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Port number of this database. (1024 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1024, 65535),
			},
			"database_collation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commands that specify how to sort and compare data",
			},
			"sqlserver_active_directory": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "MS SQL Server Active directory",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Active Directory Domain name",
						},
						"domain_net_bios_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Active Directory NetBios name",
						},
						"dns_server_ips": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Active Directory DNS Server IPs",
						},
						"ad_server_user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Active Directory Server User ID",
						},
						"ad_server_user_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Active Directory Server User password",
						},
						"failover_cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Active Directory Failover Cluster name",
						},
					},
				},
			},
			"block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    10,
				Description: "block storage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_storage_size": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Block Storage Size (10 to 7168)",
							ValidateDiagFunc: database_common.ValidateIntegerInRange(10, 7168),
						},
						"block_storage_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage product name. (SSD|HDD)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("SSD", "HDD"),
						},
						"block_storage_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Block storage group id",
						},
					},
				},
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use storage encryption.",
			},
			"sqlserver_servers": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    2,
				Description: "MS SQL Server servers (HA configuration when entering two server specifications)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. (AZ1 | AZ2)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("AZ1", "AZ2", ""),
						},
						"sqlserver_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "MS SQL Server database server names. (3 to 15 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to15LowercaseNumberDashAndStartLowercase,
						},
						"server_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Server role type Enter 'ACTIVE' for a single server configuration. (ACTIVE | PRIMARY | SECONDARY)\",",
							ValidateDiagFunc: database_common.ValidateStringInOptions("ACTIVE", "PRIMARY", "SECONDARY"),
						},
					},
				},
			},
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to use storage encryption.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id of this database server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timezone setting of this database.",
			},
			"backup": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_storage_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Object storage ID where backup files will be stored.",
						},
						"archive_backup_schedule_frequency": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Backup File Schedule Frequency.(5M, 10M, 30M, 1H) ",
							ValidateDiagFunc: database_common.ValidateStringInOptions("5M", "10M", "30M", "1H"),
						},
						"backup_retention_period": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Backup File Retention Day.(7D <= day <= 35D) ",
							ValidateDiagFunc: database_common.ValidateBackupRetentionPeriod,
						},
						"backup_start_hour": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "The time at which the backup starts. (from 0 to 23)",
							ValidateDiagFunc: database_common.ValidateIntegerInRange(0, 23),
						},
						"full_backup_day_of_week": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Full backup schedule(Day). (MONDAY to SUNDAY)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"),
						},
					},
				},
			},
			"virtual_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "virtual ip address",
			},
			"nat_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "nat ip address",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vpc id",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provide Microsoft SQL Server resource.",
	}
}

func resourceSqlserverCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	auditEnabled := rd.Get("audit_enabled").(bool)
	contractPeriod := rd.Get("contract_period").(string)
	nextContractPeriod := rd.Get("next_contract_period").(string)
	imageId := rd.Get("image_id").(string)
	natEnabled := rd.Get("nat_enabled").(bool)
	natPublicIpId := rd.Get("nat_public_ip_id").(string)

	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	serviceZoneId := rd.Get("service_zone_id").(string)

	sqlserverClusterName := rd.Get("sqlserver_cluster_name").(string)
	sqlserverClusterState := rd.Get("sqlserver_cluster_state").(string)

	//sqlserver InitialConfig
	databaseCollation := rd.Get("database_collation").(string)
	databaseNames := rd.Get("database_names").([]interface{})
	databasePort := rd.Get("database_port").(int)
	databaseServiceName := rd.Get("database_service_name").(string)
	databaseUserName := rd.Get("database_user_name").(string)
	databaseUserPassword := rd.Get("database_user_password").(string)
	license := rd.Get("license").(string)

	// sqlserver active directory
	sqlserverActiveDirectoryRequest := &sqlserver.SqlserverActiveDirectory{}
	sqlserverActiveDirectory := rd.Get("sqlserver_active_directory").(*schema.Set).List()

	if len(sqlserverActiveDirectory) != 0 {
		adMap := sqlserverActiveDirectory[0].(map[string]interface{})
		err = database_common.MapToObjectWithCamel(adMap, sqlserverActiveDirectoryRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		sqlserverActiveDirectoryRequest = nil
	}

	//sqlserver ServerGroup
	blockStorages := rd.Get("block_storages").([]interface{})
	encryptionEnabled := rd.Get("encryption_enabled").(bool)
	sqlserverServers := rd.Get("sqlserver_servers").([]interface{})
	serverType := rd.Get("server_type").(string)

	subnetId := rd.Get("subnet_id").(string)
	timezone := rd.Get("timezone").(string)
	backup := rd.Get("backup").(*schema.Set).List()

	// block storage
	var SqlserverBlockStorageGroupCreateRequestList []sqlserver.SqlserverBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertObjectSliceToStructSlice(blockStorages)
	for _, blockStorage := range blockStoragesList {
		SqlserverBlockStorageGroupCreateRequestList = append(SqlserverBlockStorageGroupCreateRequestList, sqlserver.SqlserverBlockStorageGroupCreateRequest{
			BlockStorageSize: int32(blockStorage.BlockStorageSize),
			BlockStorageType: blockStorage.BlockStorageType,
		})
	}

	// sqlserver server (HclListObject to Slice)
	var SqlserverServerCreateRequestList []sqlserver.SqlserverServerCreateRequest
	sqlserverServerList := database_common.ConvertObjectSliceToStructSlice(sqlserverServers)
	for _, sqlserverServer := range sqlserverServerList {
		SqlserverServerCreateRequestList = append(SqlserverServerCreateRequestList, sqlserver.SqlserverServerCreateRequest{
			AvailabilityZoneName: sqlserverServer.AvailabilityZoneName,
			SqlserverServerName:  sqlserverServer.SqlserverServerName,
			ServerRoleType:       sqlserverServer.ServerRoleType,
		})
	}

	securityGroupIdList := database_common.ConvertToList(securityGroupIds)
	databaseNameList := database_common.ConvertToList(databaseNames)

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		diagnostics = diag.FromErr(err)
		return
	}
	var blockId string
	for _, zoneInfo := range projectInfo.ServiceZones {
		if zoneInfo.ServiceZoneId == serviceZoneId {
			blockId = zoneInfo.BlockId
			break
		}
	}
	if len(blockId) == 0 {
		return diag.Errorf("current service block not found")
	}

	_, _, err = inst.Client.Sqlserver.CreateSqlserverCluster(ctx, sqlserver.SqlserverClusterCreateRequest{
		AuditEnabled:         &auditEnabled,
		ContractPeriod:       contractPeriod,
		ImageId:              imageId,
		NatEnabled:           &natEnabled,
		NatPublicIpId:        natPublicIpId,
		SqlserverClusterName: sqlserverClusterName,
		SqlserverInitialConfig: &sqlserver.SqlserverInitialConfigCreateRequest{
			DatabaseCollation:    databaseCollation,
			DatabaseNames:        databaseNameList,
			DatabasePort:         int32(databasePort),
			DatabaseServiceName:  databaseServiceName,
			DatabaseUserName:     databaseUserName,
			DatabaseUserPassword: databaseUserPassword,
			License:              license,
			ActiveDirectory:      sqlserverActiveDirectoryRequest,
		},
		SqlserverServerGroup: &sqlserver.SqlserverServerGroupCreateRequest{
			BlockStorages:     SqlserverBlockStorageGroupCreateRequestList,
			EncryptionEnabled: &encryptionEnabled,
			SqlserverServers:  SqlserverServerCreateRequestList,
			ServerType:        serverType,
		},
		SecurityGroupIds: securityGroupIdList,
		ServiceZoneId:    serviceZoneId,
		SubnetId:         subnetId,
		Timezone:         timezone,
	}, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty
	resultList, _, err := inst.Client.Sqlserver.ListSqlserverClusters(ctx, &sqlserver.SqlserverSearchApiListSqlserverClustersOpts{
		SqlserverClusterName: optional.NewString(sqlserverClusterName),
		Page:                 optional.NewInt32(0),
		Size:                 optional.NewInt32(1000),
		Sort:                 optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	sqlserverClusterId := resultList.Contents[0].SqlserverClusterId

	if len(sqlserverClusterId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForSqlserver(ctx, inst.Client, sqlserverClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(sqlserverClusterId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifySqlserverClusterNextContract(UpdateSqlserverParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, nextContractPeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(backup) != 0 {
		backupObject := &sqlserver.SqlserverCreateFullBackupConfigRequest{}
		backupMap := backup[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}

		err = createSqlserverClusterFullBackupConfig(UpdateSqlserverParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if sqlserverClusterState == common.StoppedState {
		err := stopSqlserverCluster(UpdateSqlserverParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSqlserverRead(ctx, rd, meta)
}

func resourceSqlserverRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	const BlockStorageDataRoleTypeIndex = 1

	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	sqlserverClusterDetail, _, err := inst.Client.Sqlserver.DetailSqlserverCluster(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(sqlserverClusterDetail.SqlserverServerGroup.SqlserverServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	sqlserverServerGroup := sqlserverClusterDetail.SqlserverServerGroup
	sqlserverBlockStorages := database_common.HclListObject{}
	blockStorageInfo := database_common.HclKeyValueObject{}
	blockStorageInfo["block_storage_size"] = sqlserverServerGroup.BlockStorages[BlockStorageDataRoleTypeIndex].BlockStorageSize
	blockStorageInfo["block_storage_type"] = sqlserverServerGroup.BlockStorages[BlockStorageDataRoleTypeIndex].BlockStorageType
	blockStorageInfo["block_storage_group_id"] = sqlserverServerGroup.BlockStorages[BlockStorageDataRoleTypeIndex].BlockStorageGroupId
	sqlserverBlockStorages = append(sqlserverBlockStorages, blockStorageInfo)

	sqlserverServers := database_common.HclListObject{}
	for _, server := range sqlserverServerGroup.SqlserverServers {
		sqlserverServersInfo := database_common.HclKeyValueObject{}
		sqlserverServersInfo["availability_zone_name"] = server.AvailabilityZoneName
		sqlserverServersInfo["sqlserver_server_name"] = server.SqlserverServerName
		sqlserverServersInfo["server_role_type"] = server.ServerRoleType

		sqlserverServers = append(sqlserverServers, sqlserverServersInfo)
	}

	backup := database_common.HclListObject{}
	if sqlserverClusterDetail.BackupConfig != nil {
		backupInfo := database_common.HclKeyValueObject{}
		backupList := rd.Get("backup").(*schema.Set).List()
		if len(backupList) == 0 {
			backupInfo["object_storage_id"] = nil
		} else {
			backupInfo["object_storage_id"] = backupList[0].(map[string]interface{})["object_storage_id"]
		}
		backupInfo["archive_backup_schedule_frequency"] = sqlserverClusterDetail.BackupConfig.FullBackupConfig.ArchiveBackupScheduleFrequency
		backupInfo["backup_retention_period"] = sqlserverClusterDetail.BackupConfig.FullBackupConfig.BackupRetentionPeriod
		backupInfo["backup_start_hour"] = sqlserverClusterDetail.BackupConfig.FullBackupConfig.BackupStartHour
		backupInfo["full_backup_day_of_week"] = sqlserverClusterDetail.BackupConfig.FullBackupConfig.FullBackupDayOfWeek

		backup = append(backup, backupInfo)
	}

	sort.SliceStable(sqlserverServers, func(i, j int) bool {
		return sqlserverServers[i].(map[string]interface{})["server_role_type"].(string) < sqlserverServers[j].(map[string]interface{})["server_role_type"].(string)
	})

	err = rd.Set("audit_enabled", sqlserverClusterDetail.AuditEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("contract_period", sqlserverClusterDetail.Contract.ContractPeriod)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("image_id", sqlserverClusterDetail.ImageId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_public_ip_id", rd.Get("nat_public_ip_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("sqlserver_cluster_name", sqlserverClusterDetail.SqlserverClusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("sqlserver_cluster_state", sqlserverClusterDetail.SqlserverClusterState)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("database_collation", sqlserverClusterDetail.SqlserverInitialConfig.DatabaseCollation)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_names", sqlserverClusterDetail.SqlserverInitialConfig.DatabaseNames)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", sqlserverClusterDetail.SqlserverInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_service_name", sqlserverClusterDetail.SqlserverInitialConfig.DatabaseServiceName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_user_name", sqlserverClusterDetail.SqlserverInitialConfig.DatabaseUserName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("block_storages", sqlserverBlockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", sqlserverClusterDetail.SqlserverServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("sqlserver_servers", sqlserverServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", sqlserverClusterDetail.SqlserverServerGroup.ServerType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("security_group_ids", sqlserverClusterDetail.SecurityGroupIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("service_zone_id", sqlserverClusterDetail.ServiceZoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("subnet_id", sqlserverClusterDetail.SubnetId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("timezone", sqlserverClusterDetail.Timezone)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("backup", backup)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("virtual_ip_address", sqlserverClusterDetail.SqlserverServerGroup.VirtualIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_ip_address", sqlserverClusterDetail.NatIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("vpc_id", sqlserverClusterDetail.VpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tfTags.SetTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type UpdateSqlserverParam struct {
	Ctx       context.Context
	Rd        *schema.ResourceData
	Inst      *client.Instance
	Sqlserver *sqlserver.SqlserverClusterDetailResponse
}

func resourceSqlserverUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	sqlserverClusterDetail, _, err := inst.Client.Sqlserver.DetailSqlserverCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(sqlserverClusterDetail.SqlserverServerGroup.SqlserverServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateSqlserverParam{
		Ctx:       ctx,
		Rd:        rd,
		Inst:      inst,
		Sqlserver: &sqlserverClusterDetail,
	}

	var updateFuncs []func(serverParam UpdateSqlserverParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizeSqlserverClusterVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updateSqlserverClusterBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updateSqlserverClusterSecurityGroupIds)
	}
	if rd.HasChanges("contract_period") {
		updateFuncs = append(updateFuncs, updateContractPeriod)
	}
	if rd.HasChanges("next_contract_period") {
		updateFuncs = append(updateFuncs, updateNextContractPeriod)
	}
	if rd.HasChanges("backup") {
		updateFuncs = append(updateFuncs, updateBackup)
	}

	for _, f := range updateFuncs {
		err = f(param)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("sqlserver_cluster_state") {
		err = updateSqlserverClusterServerState(param)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSqlserverRead(ctx, rd, meta)
}

func resizeSqlserverClusterVirtualServers(param UpdateSqlserverParam) error {
	_, _, err := param.Inst.Client.Sqlserver.ResizeSqlserverClusterVirtualServers(param.Ctx, param.Rd.Id(), sqlserver.SqlserverClusterResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}
func updateSqlserverClusterBlockStorages(param UpdateSqlserverParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertObjectSliceToStructSlice(oldValue)
	newList := database_common.ConvertObjectSliceToStructSlice(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizeSqlserverClusterBlockStorages(param, oldList, newList)
	if err != nil {
		return err
	}

	return nil
}

func validateBlockStorageInput(oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	if len(oldList) > len(newList) {
		return fmt.Errorf("removing additional storage is not allowed")
	}

	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageType != newList[i].BlockStorageType {
			return fmt.Errorf("changing block storage type is not allowed")
		}
		if oldList[i].BlockStorageSize > newList[i].BlockStorageSize {
			return fmt.Errorf("decreasing size is not allowed")
		}
	}
	return nil
}

func resizeSqlserverClusterBlockStorages(param UpdateSqlserverParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.Sqlserver.ResizeSqlserverClusterBlockStorages(param.Ctx, param.Rd.Id(), sqlserver.SqlserverClusterResizeBlockStoragesRequest{
				BlockStorageGroupId: param.Sqlserver.SqlserverServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateSqlserverClusterSecurityGroupIds(param UpdateSqlserverParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertSecurityGroupIdList(oldValue)
	newList := database_common.ConvertSecurityGroupIdList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachSqlserverClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachSqlserverClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachSqlserverClusterSecurityGroup(param UpdateSqlserverParam, v string) error {
	_, _, err := param.Inst.Client.Sqlserver.AttachSqlserverClusterSecurityGroup(param.Ctx, param.Rd.Id(), sqlserver.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachSqlserverClusterSecurityGroup(param UpdateSqlserverParam, v string) error {
	_, _, err := param.Inst.Client.Sqlserver.DetachSqlserverClusterSecurityGroup(param.Ctx, param.Rd.Id(), sqlserver.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateSqlserverClusterServerState(param UpdateSqlserverParam) error {
	_, n := param.Rd.GetChange("sqlserver_cluster_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startSqlserverCluster(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopSqlserverCluster(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("MS SQL Server status update failed. ")
	}

	return nil
}

func startSqlserverCluster(param UpdateSqlserverParam) error {
	_, _, err := param.Inst.Client.Sqlserver.StartSqlserverCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopSqlserverCluster(param UpdateSqlserverParam) error {
	_, _, err := param.Inst.Client.Sqlserver.StopSqlserverCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdateSqlserverParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifySqlserverClusterContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifySqlserverClusterContract(param UpdateSqlserverParam, newValue string) error {
	_, _, err := param.Inst.Client.Sqlserver.ModifySqlserverClusterContract(param.Ctx, param.Rd.Id(), sqlserver.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdateSqlserverParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifySqlserverClusterNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifySqlserverClusterNextContract(param UpdateSqlserverParam, newValue string) error {
	_, _, err := param.Inst.Client.Sqlserver.ModifySqlserverClusterNextContract(param.Ctx, param.Rd.Id(), sqlserver.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdateSqlserverParam) error {
	o, n := param.Rd.GetChange("backup")

	oldValue := o.(*schema.Set)
	newValue := n.(*schema.Set)

	if oldValue.Len() == 0 {
		backupObject := &sqlserver.SqlserverCreateFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = createSqlserverClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deleteSqlserverClusterFullBackupConfig(param)
		if err != nil {
			return err
		}
	} else {
		backupObject := &sqlserver.SqlserverModifyFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = updateSqlserverClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSqlserverClusterFullBackupConfig(param UpdateSqlserverParam, value *sqlserver.SqlserverCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Sqlserver.CreateSqlserverClusterFullBackupConfig(param.Ctx, param.Rd.Id(), sqlserver.SqlserverCreateFullBackupConfigRequest{
		ObjectStorageId:                value.ObjectStorageId,
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
		FullBackupDayOfWeek:            value.FullBackupDayOfWeek,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateSqlserverClusterFullBackupConfig(param UpdateSqlserverParam, value *sqlserver.SqlserverModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Sqlserver.ModifySqlserverClusterFullBackupConfig(param.Ctx, param.Rd.Id(), sqlserver.SqlserverModifyFullBackupConfigRequest{
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
		FullBackupDayOfWeek:            value.FullBackupDayOfWeek,
	})
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deleteSqlserverClusterFullBackupConfig(param UpdateSqlserverParam) error {
	_, _, err := param.Inst.Client.Sqlserver.DeleteSqlserverClusterFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForSqlserver(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func resourceSqlserverDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Sqlserver.DeleteSqlserverCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForSqlserver(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSqlserver(ctx context.Context, scpClient *client.SCPClient, sqlserverClusterId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info sqlserver.SqlserverClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Sqlserver.DetailSqlserverCluster(ctx, sqlserverClusterId)
			if err != nil && statusCode >= 500 && statusCode < 600 {
				log.Println("API temporarily unavailable. Status code: ", statusCode)
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}

		if err != nil {
			if statusCode == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if statusCode == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}

		servers := info.SqlserverServerGroup.SqlserverServers
		log.Println("1. len(servers) : ", len(servers))

		if len(servers) == 0 {
			return nil, "", fmt.Errorf("no virtual server found")
		} else {
			for i := 0; i < len(servers); i++ {
				log.Printf("servers[%s]", i+1)
				log.Println(".sqlServerState : ", servers[i].SqlserverServerState)

				if servers[i].SqlserverServerState != common.RunningState {
					return info, servers[i].SqlserverServerState, nil
				}
			}
		}
		return info, servers[0].SqlserverServerState, nil
	})
}

func resourceSqlserverDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}

	var errorMessages []string
	mutableFields := []string{
		"server_type",
		"block_storages",
		"security_group_ids",
		"sqlserver_cluster_state",
		"contract_period",
		"next_contract_period",
		"backup",
		"tags",
	}
	resourceSqlserver := ResourceSqlserver().Schema

	for key, _ := range resourceSqlserver {
		if rd.HasChanges(key) && !database_common.Contains(mutableFields, key) {
			o, n := rd.GetChange(key)
			errorMessage := fmt.Sprintf("value ['%v'] change not allowed (old: '%v', new: '%v')", key, o, n)
			errorMessages = append(errorMessages, errorMessage)
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("CustomizeDiff Validation Failed: \n%v", strings.Join(errorMessages, "\n"))
	}

	return nil
}
