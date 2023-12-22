package postgresql

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/postgresql"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
	"time"
)

func init() {
	scp.RegisterResource("scp_postgresql", ResourcePostgresql())
}

func ResourcePostgresql() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostgresqlCreate,
		ReadContext:   resourcePostgresqlRead,
		UpdateContext: resourcePostgresqlUpdate,
		DeleteContext: resourcePostgresqlDelete,
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
				ForceNew:    true,
				Description: "Whether to use database audit logging.",
			},
			"contract_period": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract (None|1-year|3-year)",
				ValidateDiagFunc: database_common.ValidateContractPeriod,
			},
			"next_contract_period": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				Description:      "Next contract (None|1-year|3-year)",
				ValidateDiagFunc: database_common.ValidateContractPeriod,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Postgresql virtual server image id.",
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
			"postgresql_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"postgresql_cluster_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "postgresql cluster state (RUNNING|STOPPED)",
				ValidateDiagFunc: database_common.ValidServerState,
			},
			"database_encoding": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Postgresql encoding. (Only 'UTF8' for now)",
				ValidateDiagFunc: database_common.ValidateSameValue("UTF8"),
			},
			"database_locale": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Postgresql locale. (Only 'C' for now)",
				ValidateDiagFunc: database_common.ValidateSameValue("C"),
			},
			"database_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database. (only English alphabets or numbers between 3 and 20 characters)",
				ValidateDiagFunc: database_common.ValidateAlphaNumeric3to20,
			},
			"database_port": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				Description:      "Port number of database. (1024 to 65535)",
				ValidateDiagFunc: database_common.ValidatePortNumber,
			},
			"database_user_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "User account id of database. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"database_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ForceNew:         true,
				Description:      "User account password of database.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "block storage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_storage_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage product name. (SSD|HDD)",
							ValidateDiagFunc: database_common.ValidateBlockStorageType,
						},
						"block_storage_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage usage. (DATA|ARCHIVE|TEMP|BACKUP)",
							ValidateDiagFunc: database_common.ValidateBlockStorageRoleType,
						},
						"block_storage_size": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Block Storage Size (10 to 5120)",
							ValidateDiagFunc: database_common.ValidateBlockStorageSize,
						},
					},
				},
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use storage encryption.",
			},
			"postgresql_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "postgresql servers (HA configuration when entering two server specifications)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. The single server does not input anything. (AZ1|AZ2)",
							ValidateDiagFunc: database_common.ValidateAvailabilityZone,
						},
						"postgresql_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							Description:      "Postgresql database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"server_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							Description:      "Server role type Enter 'ACTIVE' for a single server configuration. (ACTIVE | STANDBY)",
							ValidateDiagFunc: database_common.ValidateServerRoleType,
						},
					},
				},
			},
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server type",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this postgresql DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone Id",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id of this database server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Timezone setting of this database.",
			},
			"backup": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     resourcePostgreSQLBackup(),
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a PostgreSQL Database resource.",
	}
}

func resourcePostgreSQLBackup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object storage ID where backup files will be stored.",
			},
			"archive_backup_schedule_frequency": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Backup File Schedule Frequency.(5M|10M|30M|1H) ",
				ValidateDiagFunc: database_common.ValidateBackupScheduleFrequency,
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
				ValidateDiagFunc: database_common.ValidateBackupStartHour,
			},
		},
	}
}

func resourcePostgresqlCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
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
	postgresqlClusterName := rd.Get("postgresql_cluster_name").(string)
	postgresqlClusterState := rd.Get("postgresql_cluster_state").(string)

	//postgresqlInitialConfig
	databaseEncoding := rd.Get("database_encoding").(string)
	databaseLocale := rd.Get("database_locale").(string)
	databaseName := rd.Get("database_name").(string)
	databasePort := rd.Get("database_port").(int)
	databaseUserName := rd.Get("database_user_name").(string)
	databaseUserPassword := rd.Get("database_user_password").(string)

	//postgresqlServerGroup
	blockStorages := rd.Get("block_storages").([]interface{})
	encryptionEnabled := rd.Get("encryption_enabled").(bool)
	postgresqlServers := rd.Get("postgresql_servers").([]interface{})
	serverType := rd.Get("server_type").(string)

	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	serviceZoneId := rd.Get("service_zone_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	timezone := rd.Get("timezone").(string)
	backup := rd.Get("backup").(*schema.Set).List()

	// block storage (HclListObject to Slice)
	var PostgresqlBlockStorageGroupCreateRequestList []postgresql.PostgresqlBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertBlockStorageList(blockStorages)
	for _, blockStorage := range blockStoragesList {
		PostgresqlBlockStorageGroupCreateRequestList = append(PostgresqlBlockStorageGroupCreateRequestList, postgresql.PostgresqlBlockStorageGroupCreateRequest{
			BlockStorageRoleType: blockStorage.BlockStorageRoleType,
			BlockStorageSize:     int32(blockStorage.BlockStorageSize),
			BlockStorageType:     blockStorage.BlockStorageType,
		})
	}

	// postgresql server (HclListObject to Slice)
	var PostgresqlServerCreateRequestList []postgresql.PostgresqlServerCreateRequest
	postgresqlServerList := database_common.ConvertServerList(postgresqlServers)
	for _, postgresqlServer := range postgresqlServerList {
		PostgresqlServerCreateRequestList = append(PostgresqlServerCreateRequestList, postgresql.PostgresqlServerCreateRequest{
			AvailabilityZoneName: postgresqlServer.AvailabilityZoneName,
			PostgresqlServerName: postgresqlServer.PostgresqlServerName,
			ServerRoleType:       postgresqlServer.ServerRoleType,
		})
	}

	securityGroupIdList := database_common.ConvertSecurityGroupIdList(securityGroupIds)

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

	_, _, err = inst.Client.Postgresql.CreatePostgresqlCluster(ctx, postgresql.PostgresqlClusterCreateRequest{
		AuditEnabled:          &auditEnabled,
		ContractPeriod:        contractPeriod,
		ImageId:               imageId,
		NatEnabled:            &natEnabled,
		NatPublicIpId:         natPublicIpId,
		PostgresqlClusterName: postgresqlClusterName,
		PostgresqlInitialConfig: &postgresql.PostgresqlInitialConfigCreateRequest{
			DatabaseEncoding:     databaseEncoding,
			DatabaseLocale:       databaseLocale,
			DatabaseName:         databaseName,
			DatabasePort:         int32(databasePort),
			DatabaseUserName:     databaseUserName,
			DatabaseUserPassword: databaseUserPassword,
		},
		PostgresqlServerGroup: &postgresql.PostgresqlServerGroupCreateRequest{
			BlockStorages:     PostgresqlBlockStorageGroupCreateRequestList,
			EncryptionEnabled: &encryptionEnabled,
			PostgresqlServers: PostgresqlServerCreateRequestList,
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
	resultList, _, err := inst.Client.Postgresql.ListPostgresqlClusters(ctx, &postgresql.PostgresqlSearchApiListPostgresqlClustersOpts{
		PostgresqlClusterName: optional.NewString(postgresqlClusterName),
		Page:                  optional.NewInt32(0),
		Size:                  optional.NewInt32(1000),
		Sort:                  optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	postgresqlClusterId := resultList.Contents[0].PostgresqlClusterId

	if len(postgresqlClusterId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForPostgresql(ctx, inst.Client, postgresqlClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(postgresqlClusterId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifyPostgresqlClusterNextContract(UpdatePostgresqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, nextContractPeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(backup) != 0 {
		backupObject := &postgresql.DbClusterCreateFullBackupConfigRequest{}
		backupMap := backup[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}

		err = createPostgresqlClusterFullBackupConfig(UpdatePostgresqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if postgresqlClusterState == common.StoppedState {
		err := stopPostgresqlCluster(UpdatePostgresqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePostgresqlRead(ctx, rd, meta)
}

func resourcePostgresqlRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Postgresql.DetailPostgresqlCluster(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	if len(dbInfo.PostgresqlServerGroup.PostgresqlServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	//TODO BlockStorageGroupId, BlockStorageName 받아올 지 확인
	blockStorages := database_common.HclListObject{}
	for i, bs := range dbInfo.PostgresqlServerGroup.BlockStorages {
		// Skip OS Storage
		if i == 0 {
			continue
		}
		blockStorageInfo := database_common.HclKeyValueObject{}
		blockStorageInfo["block_storage_role_type"] = bs.BlockStorageRoleType
		blockStorageInfo["block_storage_size"] = bs.BlockStorageSize
		blockStorageInfo["block_storage_type"] = bs.BlockStorageType

		blockStorages = append(blockStorages, blockStorageInfo)
	}

	postgresqlServers := database_common.HclListObject{}
	for _, server := range dbInfo.PostgresqlServerGroup.PostgresqlServers {
		postgresqlServersInfo := database_common.HclKeyValueObject{}
		postgresqlServersInfo["availability_zone_name"] = server.AvailabilityZoneName
		postgresqlServersInfo["postgresql_server_name"] = server.PostgresqlServerName
		postgresqlServersInfo["server_role_type"] = server.ServerRoleType

		postgresqlServers = append(postgresqlServers, postgresqlServersInfo)
	}

	backup := database_common.HclListObject{}
	if dbInfo.BackupConfig != nil {
		backupInfo := database_common.HclKeyValueObject{}
		backupInfo["object_storage_id"] = rd.Get("backup").(*schema.Set).List()[0].(map[string]interface{})["object_storage_id"]
		backupInfo["archive_backup_schedule_frequency"] = dbInfo.BackupConfig.FullBackupConfig.ArchiveBackupScheduleFrequency
		backupInfo["backup_retention_period"] = dbInfo.BackupConfig.FullBackupConfig.BackupRetentionPeriod
		backupInfo["backup_start_hour"] = dbInfo.BackupConfig.FullBackupConfig.BackupStartHour

		backup = append(backup, backupInfo)
	}

	sort.SliceStable(postgresqlServers, func(i, j int) bool {
		return postgresqlServers[i].(map[string]interface{})["server_role_type"].(string) < postgresqlServers[j].(map[string]interface{})["server_role_type"].(string)
	})

	err = rd.Set("audit_enabled", dbInfo.AuditEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("contract_period", dbInfo.Contract.ContractPeriod)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("image_id", dbInfo.ImageId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_public_ip_id", dbInfo.NatIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("postgresql_cluster_name", dbInfo.PostgresqlClusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("postgresql_cluster_state", dbInfo.PostgresqlClusterState)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("database_encoding", dbInfo.PostgresqlInitialConfig.DatabaseEncoding)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_locale", dbInfo.PostgresqlInitialConfig.DatabaseLocale)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_name", dbInfo.PostgresqlInitialConfig.DatabaseName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", dbInfo.PostgresqlInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_user_name", dbInfo.PostgresqlInitialConfig.DatabaseUserName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("block_storages", blockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", dbInfo.PostgresqlServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("postgresql_servers", postgresqlServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", dbInfo.PostgresqlServerGroup.ServerType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("security_group_ids", dbInfo.SecurityGroupIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("service_zone_id", dbInfo.ServiceZoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("subnet_id", dbInfo.SubnetId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("timezone", dbInfo.Timezone)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("backup", backup)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tfTags.SetTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type UpdatePostgresqlParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *postgresql.PostgresqlClusterDetailResponse
}

func resourcePostgresqlUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Postgresql.DetailPostgresqlCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.PostgresqlServerGroup.PostgresqlServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdatePostgresqlParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdatePostgresqlParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizePostgresqlClusterVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updatePostgresqlClusterBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updatePostgresqlClusterSecurityGroupIds)
	}
	if rd.HasChanges("postgresql_cluster_state") {
		updateFuncs = append(updateFuncs, updatePostgresqlClusterServerState)
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
			return
		}
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return
	}
	return resourcePostgresqlRead(ctx, rd, meta)
}

func resourcePostgresqlDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Postgresql.DeletePostgresqlCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForPostgresql(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resizePostgresqlClusterVirtualServers(param UpdatePostgresqlParam) error {
	_, _, err := param.Inst.Client.Postgresql.ResizePostgresqlClusterVirtualServers(param.Ctx, param.Rd.Id(), postgresql.PostgresqlClusterResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updatePostgresqlClusterBlockStorages(param UpdatePostgresqlParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertBlockStorageList(oldValue)
	newList := database_common.ConvertBlockStorageList(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizePostgresqlClusterBlockStorages(param, oldList, newList)
	if err != nil {
		return err
	}

	err = addPostgresqlClusterBlockStorages(param, oldList, newList)
	if err != nil {
		return err
	}

	return nil
}

func validateBlockStorageInput(oldList []database_common.BlockStorage, newList []database_common.BlockStorage) error {
	if len(oldList) > len(newList) {
		return fmt.Errorf("removing additional storage is not allowed")
	}

	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageRoleType != newList[i].BlockStorageRoleType {
			return fmt.Errorf("changing block storage role type is not allowed")
		}
		if oldList[i].BlockStorageType != newList[i].BlockStorageType {
			return fmt.Errorf("changing block storage type is not allowed")
		}
		if oldList[i].BlockStorageSize > newList[i].BlockStorageSize {
			return fmt.Errorf("decreasing size is not allowed")
		}
	}
	return nil
}

func resizePostgresqlClusterBlockStorages(param UpdatePostgresqlParam, oldList []database_common.BlockStorage, newList []database_common.BlockStorage) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.Postgresql.ResizePostgresqlClusterBlockStorages(param.Ctx, param.Rd.Id(), postgresql.PostgresqlClusterResizeBlockStoragesRequest{
				BlockStorageGroupId: param.DbInfo.PostgresqlServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addPostgresqlClusterBlockStorages(param UpdatePostgresqlParam, oldList []database_common.BlockStorage, newList []database_common.BlockStorage) error {
	for i := 0; i < len(newList)-len(oldList); i++ {
		_, _, err := param.Inst.Client.Postgresql.AddPostgresqlClusterBlockStorages(param.Ctx, param.Rd.Id(), postgresql.PostgresqlClusterAddBlockStoragesRequest{
			BlockStorageRoleType: newList[len(oldList)+i].BlockStorageRoleType,
			BlockStorageType:     newList[len(oldList)+i].BlockStorageType,
			BlockStorageSize:     int32(newList[len(oldList)+i].BlockStorageSize),
		})

		if err != nil {
			return err
		}

		err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func updatePostgresqlClusterSecurityGroupIds(param UpdatePostgresqlParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertSecurityGroupIdList(oldValue)
	newList := database_common.ConvertSecurityGroupIdList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachPostgresqlClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachPostgresqlClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachPostgresqlClusterSecurityGroup(param UpdatePostgresqlParam, v string) error {
	_, _, err := param.Inst.Client.Postgresql.AttachPostgresqlClusterSecurityGroup(param.Ctx, param.Rd.Id(), postgresql.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachPostgresqlClusterSecurityGroup(param UpdatePostgresqlParam, v string) error {
	_, _, err := param.Inst.Client.Postgresql.DetachPostgresqlClusterSecurityGroup(param.Ctx, param.Rd.Id(), postgresql.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updatePostgresqlClusterServerState(param UpdatePostgresqlParam) error {
	_, n := param.Rd.GetChange("postgresql_cluster_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startPostgresqlCluster(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopPostgresqlCluster(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Postgresql status update failed. ")
	}

	return nil

}

func startPostgresqlCluster(param UpdatePostgresqlParam) error {
	_, _, err := param.Inst.Client.Postgresql.StartPostgresqlCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopPostgresqlCluster(param UpdatePostgresqlParam) error {
	_, _, err := param.Inst.Client.Postgresql.StopPostgresqlCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdatePostgresqlParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifyPostgresqlClusterContract(param, newValue)
	if err != nil {
		return err
	}

	return nil

}

func modifyPostgresqlClusterContract(param UpdatePostgresqlParam, newValue string) error {
	_, _, err := param.Inst.Client.Postgresql.ModifyPostgresqlClusterContract(param.Ctx, param.Rd.Id(), postgresql.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdatePostgresqlParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifyPostgresqlClusterNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifyPostgresqlClusterNextContract(param UpdatePostgresqlParam, newValue string) error {
	_, _, err := param.Inst.Client.Postgresql.ModifyPostgresqlClusterNextContract(param.Ctx, param.Rd.Id(), postgresql.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdatePostgresqlParam) error {
	o, n := param.Rd.GetChange("backup")

	oldValue := o.(*schema.Set)
	newValue := n.(*schema.Set)

	if oldValue.Len() == 0 {
		backupObject := &postgresql.DbClusterCreateFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = createPostgresqlClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deletePostgresqlClusterFullBackupConfig(param)
		if err != nil {
			return err
		}
	} else {
		backupObject := &postgresql.DbClusterModifyFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = updatePostgresqlClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPostgresqlClusterFullBackupConfig(param UpdatePostgresqlParam, value *postgresql.DbClusterCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Postgresql.CreatePostgresqlClusterFullBackupConfig(param.Ctx, param.Rd.Id(), postgresql.DbClusterCreateFullBackupConfigRequest{
		ObjectStorageId:                value.ObjectStorageId,
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func updatePostgresqlClusterFullBackupConfig(param UpdatePostgresqlParam, value *postgresql.DbClusterModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Postgresql.ModifyPostgresqlClusterFullBackupConfig(param.Ctx, param.Rd.Id(), postgresql.DbClusterModifyFullBackupConfigRequest{
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deletePostgresqlClusterFullBackupConfig(param UpdatePostgresqlParam) error {
	_, _, err := param.Inst.Client.Postgresql.DeletePostgresqlClusterFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForPostgresql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func waitForPostgresql(ctx context.Context, scpClient *client.SCPClient, postgresqlClusterId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info postgresql.PostgresqlClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Postgresql.DetailPostgresqlCluster(ctx, postgresqlClusterId)
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

		servers := info.PostgresqlServerGroup.PostgresqlServers
		log.Println("1. len(servers) : ", len(servers))

		switch len(servers) {
		case 0:
			return nil, "", fmt.Errorf("no virtual server found")
		case 1:
			log.Println("2. servers[0].PostgresqlServerState : ", servers[0].PostgresqlServerState)
			return info, servers[0].PostgresqlServerState, nil
		case 2:
			log.Println("3. servers[0].PostgresqlServerState : ", servers[0].PostgresqlServerState, ", servers[1].PostgresqlServerState : ", servers[1].PostgresqlServerState)
			if servers[0].PostgresqlServerState == common.RunningState && servers[1].PostgresqlServerState == common.RunningState {
				return info, servers[0].PostgresqlServerState, nil
			} else {
				if servers[0].PostgresqlServerState != common.RunningState {
					return info, servers[0].PostgresqlServerState, nil
				} else {
					return info, servers[1].PostgresqlServerState, nil
				}
			}
		default:
			return nil, "", fmt.Errorf("invalid number of virtual servers")
		}
	})
}
