package mysql

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/mysql"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_mysql", ResourceMysql())
}

func ResourceMysql() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlCreate,
		ReadContext:   resourceMysqlRead,
		UpdateContext: resourceMysqlUpdate,
		DeleteContext: resourceMysqlDelete,
		CustomizeDiff: resourceMysqlDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"contract_period": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract (None|1 Year|3 Year)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
			},
			"next_contract_period": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				Description:      "Next contract (None|1 Year|3 Year)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mysql virtual server image id.",
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
			"mysql_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"mysql_cluster_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mysql cluster state (RUNNING|STOPPED)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("RUNNING", "STOPPED"),
			},
			"database_case_sensitivity": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Mysql database case sensitivity.",
			},
			"database_character_set": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mysql encoding. (utf8mb3|utf8mb4)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("utf8mb3", "utf8mb4"),
			},
			"database_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database. (only English alphabets or numbers between 3 and 20 characters)",
				ValidateDiagFunc: database_common.ValidateAlphaNumeric3to20,
			},
			"database_port": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Port number of database. (1024 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1024, 65535),
			},
			"database_user_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "User account id of database. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"database_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
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
							ValidateDiagFunc: database_common.ValidateStringInOptions("SSD", "HDD"),
						},
						"block_storage_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage usage. (DATA|ARCHIVE|TEMP|BACKUP)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("DATA", "ARCHIVE", "TEMP", "BACKUP"),
						},
						"block_storage_size": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Block Storage Size (10 to 5120)",
							ValidateDiagFunc: database_common.ValidateIntegerInRange(10, 5120),
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
			"mysql_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Mysql servers (HA configuration when entering two server specifications)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. The single server does not input anything. (AZ1|AZ2)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("AZ1", "AZ2"),
						},
						"mysql_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Mysql database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"server_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Server role type Enter 'ACTIVE' for a single server configuration. (ACTIVE | STANDBY)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("ACTIVE", "STANDBY"),
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
				Description: "Security-Group ids of this Mysql DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service Zone Id",
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
							Optional:    true,
							Description: "Object storage ID where backup files will be stored.",
						},
						"archive_backup_schedule_frequency": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Backup File Schedule Frequency.(5M|10M|30M|1H) ",
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
		Description: "Provides a Mysql Database resource.",
	}
}

func resourceMysqlCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	contractPeriod := rd.Get("contract_period").(string)
	nextContractPeriod := rd.Get("next_contract_period").(string)
	imageId := rd.Get("image_id").(string)
	natEnabled := rd.Get("nat_enabled").(bool)
	natPublicIpId := rd.Get("nat_public_ip_id").(string)
	mysqlClusterName := rd.Get("mysql_cluster_name").(string)
	mysqlClusterState := rd.Get("mysql_cluster_state").(string)

	//MysqlInitialConfig
	databaseCaseSensitivity := rd.Get("database_case_sensitivity").(bool)
	databaseCharacterSet := rd.Get("database_character_set").(string)
	databaseName := rd.Get("database_name").(string)
	databasePort := rd.Get("database_port").(int)
	databaseUserName := rd.Get("database_user_name").(string)
	databaseUserPassword := rd.Get("database_user_password").(string)

	//MysqlServerGroup
	blockStorages := rd.Get("block_storages").([]interface{})
	encryptionEnabled := rd.Get("encryption_enabled").(bool)
	mysqlServers := rd.Get("mysql_servers").([]interface{})
	serverType := rd.Get("server_type").(string)

	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	serviceZoneId := rd.Get("service_zone_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	timezone := rd.Get("timezone").(string)
	backup := rd.Get("backup").(*schema.Set).List()

	// block storage (HclListObject to Slice)
	var MysqlBlockStorageGroupCreateRequestList []mysql.MysqlBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertObjectSliceToStructSlice(blockStorages)
	for _, blockStorage := range blockStoragesList {
		MysqlBlockStorageGroupCreateRequestList = append(MysqlBlockStorageGroupCreateRequestList, mysql.MysqlBlockStorageGroupCreateRequest{
			BlockStorageRoleType: blockStorage.BlockStorageRoleType,
			BlockStorageSize:     int32(blockStorage.BlockStorageSize),
			BlockStorageType:     blockStorage.BlockStorageType,
		})
	}

	// Mysql server (HclListObject to Slice)
	var MysqlServerCreateRequestList []mysql.MysqlServerCreateRequest
	MysqlServerList := database_common.ConvertObjectSliceToStructSlice(mysqlServers)
	for _, MysqlServer := range MysqlServerList {
		MysqlServerCreateRequestList = append(MysqlServerCreateRequestList, mysql.MysqlServerCreateRequest{
			AvailabilityZoneName: MysqlServer.AvailabilityZoneName,
			MysqlServerName:      MysqlServer.MysqlServerName,
			ServerRoleType:       MysqlServer.ServerRoleType,
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

	_, _, err = inst.Client.Mysql.CreateMysqlCluster(ctx, mysql.MysqlClusterCreateRequest{
		ContractPeriod:   contractPeriod,
		ImageId:          imageId,
		NatEnabled:       &natEnabled,
		NatPublicIpId:    natPublicIpId,
		MysqlClusterName: mysqlClusterName,
		MysqlInitialConfig: &mysql.MysqlInitialConfigCreateRequest{
			DatabaseCaseSensitivity: &databaseCaseSensitivity,
			DatabaseCharacterSet:    databaseCharacterSet,
			DatabaseName:            databaseName,
			DatabasePort:            int32(databasePort),
			DatabaseUserName:        databaseUserName,
			DatabaseUserPassword:    databaseUserPassword,
		},
		MysqlServerGroup: &mysql.MysqlServerGroupCreateRequest{
			BlockStorages:     MysqlBlockStorageGroupCreateRequestList,
			EncryptionEnabled: &encryptionEnabled,
			MysqlServers:      MysqlServerCreateRequestList,
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
	resultList, _, err := inst.Client.Mysql.ListMysqlClusters(ctx, &mysql.MysqlSearchApiListMysqlClustersOpts{
		MysqlClusterName: optional.NewString(mysqlClusterName),
		Page:             optional.NewInt32(0),
		Size:             optional.NewInt32(1000),
		Sort:             optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	MysqlClusterId := resultList.Contents[0].MysqlClusterId

	if len(MysqlClusterId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForMysql(ctx, inst.Client, MysqlClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(MysqlClusterId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifyMysqlClusterNextContract(UpdateMysqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, nextContractPeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(backup) != 0 {
		backupObject := &mysql.DbClusterCreateFullBackupConfigRequest{}
		backupMap := backup[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}

		err = createMysqlClusterFullBackupConfig(UpdateMysqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if mysqlClusterState == common.StoppedState {
		err := stopMysqlCluster(UpdateMysqlParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMysqlRead(ctx, rd, meta)
}

func resourceMysqlRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Mysql.DetailMysqlCluster(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	if len(dbInfo.MysqlServerGroup.MysqlServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	blockStorages := database_common.HclListObject{}
	for i, bs := range dbInfo.MysqlServerGroup.BlockStorages {
		// Skip OS Storage
		if i == 0 {
			continue
		}
		blockStorageInfo := database_common.HclKeyValueObject{}
		blockStorageInfo["block_storage_role_type"] = bs.BlockStorageRoleType
		blockStorageInfo["block_storage_size"] = bs.BlockStorageSize
		blockStorageInfo["block_storage_type"] = bs.BlockStorageType
		blockStorageInfo["block_storage_group_id"] = bs.BlockStorageGroupId

		blockStorages = append(blockStorages, blockStorageInfo)
	}

	MysqlServers := database_common.HclListObject{}
	for _, server := range dbInfo.MysqlServerGroup.MysqlServers {
		MysqlServersInfo := database_common.HclKeyValueObject{}
		MysqlServersInfo["availability_zone_name"] = server.AvailabilityZoneName
		MysqlServersInfo["mysql_server_name"] = server.MysqlServerName
		MysqlServersInfo["server_role_type"] = server.ServerRoleType

		MysqlServers = append(MysqlServers, MysqlServersInfo)
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

	sort.SliceStable(MysqlServers, func(i, j int) bool {
		return MysqlServers[i].(map[string]interface{})["server_role_type"].(string) < MysqlServers[j].(map[string]interface{})["server_role_type"].(string)
	})

	err = rd.Set("contract_period", dbInfo.Contract.ContractPeriod)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("image_id", dbInfo.ImageId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_public_ip_id", rd.Get("nat_public_ip_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mysql_cluster_name", dbInfo.MysqlClusterName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mysql_cluster_state", dbInfo.MysqlClusterState)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("database_case_sensitivity", rd.Get("database_case_sensitivity").(bool))
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_character_set", dbInfo.MysqlInitialConfig.DatabaseCharacterSet)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_name", dbInfo.MysqlInitialConfig.DatabaseName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", dbInfo.MysqlInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_user_name", dbInfo.MysqlInitialConfig.DatabaseUserName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("block_storages", blockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", dbInfo.MysqlServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mysql_servers", MysqlServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", dbInfo.MysqlServerGroup.ServerType)
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
	err = rd.Set("virtual_ip_address", dbInfo.MysqlServerGroup.VirtualIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("nat_ip_address", dbInfo.NatIpAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("vpc_id", dbInfo.VpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tfTags.SetTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type UpdateMysqlParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *mysql.MysqlClusterDetailResponse
}

func resourceMysqlUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Mysql.DetailMysqlCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.MysqlServerGroup.MysqlServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateMysqlParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdateMysqlParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizeMysqlClusterVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updateMysqlClusterBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updateMysqlClusterSecurityGroupIds)
	}
	if rd.HasChanges("mysql_cluster_state") {
		updateFuncs = append(updateFuncs, updateMysqlClusterServerState)
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

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMysqlRead(ctx, rd, meta)
}

func resourceMysqlDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Mysql.DeleteMysqlCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForMysql(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMysqlDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}

	var errorMessages []string
	mutableFields := []string{
		"server_type",
		"block_storages",
		"security_group_ids",
		"mysql_cluster_state",
		"contract_period",
		"next_contract_period",
		"backup",
		"tags",
	}
	resourceMysql := ResourceMysql().Schema

	for key := range resourceMysql {
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

func resizeMysqlClusterVirtualServers(param UpdateMysqlParam) error {
	_, _, err := param.Inst.Client.Mysql.ResizeMysqlClusterVirtualServers(param.Ctx, param.Rd.Id(), mysql.MysqlClusterResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateMysqlClusterBlockStorages(param UpdateMysqlParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertObjectSliceToStructSlice(oldValue)
	newList := database_common.ConvertObjectSliceToStructSlice(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizeMysqlClusterBlockStorages(param, oldList, newList)
	if err != nil {
		return err
	}

	err = addMysqlClusterBlockStorages(param, oldList, newList)
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

func resizeMysqlClusterBlockStorages(param UpdateMysqlParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.Mysql.ResizeMysqlClusterBlockStorages(param.Ctx, param.Rd.Id(), mysql.MysqlClusterResizeBlockStoragesRequest{
				BlockStorageGroupId: param.DbInfo.MysqlServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addMysqlClusterBlockStorages(param UpdateMysqlParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(newList)-len(oldList); i++ {
		_, _, err := param.Inst.Client.Mysql.AddMysqlClusterBlockStorages(param.Ctx, param.Rd.Id(), mysql.MysqlClusterAddBlockStoragesRequest{
			BlockStorageRoleType: newList[len(oldList)+i].BlockStorageRoleType,
			BlockStorageType:     newList[len(oldList)+i].BlockStorageType,
			BlockStorageSize:     int32(newList[len(oldList)+i].BlockStorageSize),
		})

		if err != nil {
			return err
		}

		err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateMysqlClusterSecurityGroupIds(param UpdateMysqlParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertSecurityGroupIdList(oldValue)
	newList := database_common.ConvertSecurityGroupIdList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachMysqlClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachMysqlClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachMysqlClusterSecurityGroup(param UpdateMysqlParam, v string) error {
	_, _, err := param.Inst.Client.Mysql.AttachMysqlClusterSecurityGroup(param.Ctx, param.Rd.Id(), mysql.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachMysqlClusterSecurityGroup(param UpdateMysqlParam, v string) error {
	_, _, err := param.Inst.Client.Mysql.DetachMysqlClusterSecurityGroup(param.Ctx, param.Rd.Id(), mysql.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateMysqlClusterServerState(param UpdateMysqlParam) error {
	_, n := param.Rd.GetChange("mysql_cluster_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startMysqlCluster(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopMysqlCluster(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Mysql status update failed. ")
	}

	return nil

}

func startMysqlCluster(param UpdateMysqlParam) error {
	_, _, err := param.Inst.Client.Mysql.StartMysqlCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopMysqlCluster(param UpdateMysqlParam) error {
	_, _, err := param.Inst.Client.Mysql.StopMysqlCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdateMysqlParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifyMysqlClusterContract(param, newValue)
	if err != nil {
		return err
	}

	return nil

}

func modifyMysqlClusterContract(param UpdateMysqlParam, newValue string) error {
	_, _, err := param.Inst.Client.Mysql.ModifyMysqlClusterContract(param.Ctx, param.Rd.Id(), mysql.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdateMysqlParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifyMysqlClusterNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifyMysqlClusterNextContract(param UpdateMysqlParam, newValue string) error {
	_, _, err := param.Inst.Client.Mysql.ModifyMysqlClusterNextContract(param.Ctx, param.Rd.Id(), mysql.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdateMysqlParam) error {
	o, n := param.Rd.GetChange("backup")

	oldValue := o.(*schema.Set)
	newValue := n.(*schema.Set)

	if oldValue.Len() == 0 {
		backupObject := &mysql.DbClusterCreateFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = createMysqlClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deleteMysqlClusterFullBackupConfig(param)
		if err != nil {
			return err
		}
	} else {
		backupObject := &mysql.DbClusterModifyFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = updateMysqlClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMysqlClusterFullBackupConfig(param UpdateMysqlParam, value *mysql.DbClusterCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Mysql.CreateMysqlClusterFullBackupConfig(param.Ctx, param.Rd.Id(), mysql.DbClusterCreateFullBackupConfigRequest{
		ObjectStorageId:                value.ObjectStorageId,
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func updateMysqlClusterFullBackupConfig(param UpdateMysqlParam, value *mysql.DbClusterModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Mysql.ModifyMysqlClusterFullBackupConfig(param.Ctx, param.Rd.Id(), mysql.DbClusterModifyFullBackupConfigRequest{
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deleteMysqlClusterFullBackupConfig(param UpdateMysqlParam) error {
	_, _, err := param.Inst.Client.Mysql.DeleteMysqlClusterFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMysql(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func waitForMysql(ctx context.Context, scpClient *client.SCPClient, MysqlClusterId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info mysql.MysqlClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Mysql.DetailMysqlCluster(ctx, MysqlClusterId)
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

		servers := info.MysqlServerGroup.MysqlServers
		log.Println("1. len(servers) : ", len(servers))

		switch len(servers) {
		case 0:
			return nil, "", fmt.Errorf("no virtual server found")
		case 1:
			log.Println("2. servers[0].MysqlServerState : ", servers[0].MysqlServerState)
			return info, servers[0].MysqlServerState, nil
		case 2:
			log.Println("3. servers[0].MysqlServerState : ", servers[0].MysqlServerState, ", servers[1].MysqlServerState : ", servers[1].MysqlServerState)
			if servers[0].MysqlServerState == common.RunningState && servers[1].MysqlServerState == common.RunningState {
				return info, servers[0].MysqlServerState, nil
			} else {
				if servers[0].MysqlServerState != common.RunningState {
					return info, servers[0].MysqlServerState, nil
				} else {
					return info, servers[1].MysqlServerState, nil
				}
			}
		default:
			return nil, "", fmt.Errorf("invalid number of virtual servers")
		}
	})
}
