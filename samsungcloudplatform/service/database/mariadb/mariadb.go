package mariadb

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/mariadb"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
	"strings"
	"time"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_mariadb", ResourceMariadb())
}

func ResourceMariadb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMariadbCreate,
		ReadContext:   resourceMariadbRead,
		UpdateContext: resourceMariadbUpdate,
		DeleteContext: resourceMariadbDelete,
		CustomizeDiff: resourceMariadbDiff,
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
				Description: "Mariadb virtual server image id.",
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
			"mariadb_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"mariadb_cluster_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mariadb cluster state (RUNNING|STOPPED)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("RUNNING", "STOPPED"),
			},
			"database_character_set": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Mariadb encoding. (utf8|utf8mb4)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("utf8", "utf8mb4"),
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
			"mariadb_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Mariadb servers (HA configuration when entering two server specifications)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. The single server does not input anything. (AZ1|AZ2)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("AZ1", "AZ2"),
						},
						"mariadb_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "MariaDB database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
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
				Description: "Security-Group ids of this Mariadb DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
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
		Description: "Provides a Mariadb Database resource.",
	}
}

func resourceMariadbCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
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
	mariadbClusterName := rd.Get("mariadb_cluster_name").(string)
	mariadbClusterState := rd.Get("mariadb_cluster_state").(string)

	//MariadbInitialConfig
	databaseCharacterSet := rd.Get("database_character_set").(string)
	databaseName := rd.Get("database_name").(string)
	databasePort := rd.Get("database_port").(int)
	databaseUserName := rd.Get("database_user_name").(string)
	databaseUserPassword := rd.Get("database_user_password").(string)

	//MariadbServerGroup
	blockStorages := rd.Get("block_storages").([]interface{})
	encryptionEnabled := rd.Get("encryption_enabled").(bool)
	mariadbServers := rd.Get("mariadb_servers").([]interface{})
	serverType := rd.Get("server_type").(string)

	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	serviceZoneId := rd.Get("service_zone_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	timezone := rd.Get("timezone").(string)
	backup := rd.Get("backup").(*schema.Set).List()

	// block storage (HclListObject to Slice)
	var MariadbBlockStorageGroupCreateRequestList []mariadb.MariadbBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertObjectSliceToStructSlice(blockStorages)
	for _, blockStorage := range blockStoragesList {
		MariadbBlockStorageGroupCreateRequestList = append(MariadbBlockStorageGroupCreateRequestList, mariadb.MariadbBlockStorageGroupCreateRequest{
			BlockStorageRoleType: blockStorage.BlockStorageRoleType,
			BlockStorageSize:     int32(blockStorage.BlockStorageSize),
			BlockStorageType:     blockStorage.BlockStorageType,
		})
	}

	// Mariadb server (HclListObject to Slice)
	var MariadbServerCreateRequestList []mariadb.MariadbServerCreateRequest
	MariadbServerList := database_common.ConvertObjectSliceToStructSlice(mariadbServers)
	for _, MariadbServer := range MariadbServerList {
		MariadbServerCreateRequestList = append(MariadbServerCreateRequestList, mariadb.MariadbServerCreateRequest{
			AvailabilityZoneName: MariadbServer.AvailabilityZoneName,
			MariadbServerName:    MariadbServer.MariadbServerName,
			ServerRoleType:       MariadbServer.ServerRoleType,
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

	_, _, err = inst.Client.Mariadb.CreateMariadbCluster(ctx, mariadb.MariadbClusterCreateRequest{
		AuditEnabled:       &auditEnabled,
		ContractPeriod:     contractPeriod,
		ImageId:            imageId,
		NatEnabled:         &natEnabled,
		NatPublicIpId:      natPublicIpId,
		MariadbClusterName: mariadbClusterName,
		MariadbInitialConfig: &mariadb.MariadbInitialConfigCreateRequest{
			DatabaseCharacterSet: databaseCharacterSet,
			DatabaseName:         databaseName,
			DatabasePort:         int32(databasePort),
			DatabaseUserName:     databaseUserName,
			DatabaseUserPassword: databaseUserPassword,
		},
		MariadbServerGroup: &mariadb.MariadbServerGroupCreateRequest{
			BlockStorages:     MariadbBlockStorageGroupCreateRequestList,
			EncryptionEnabled: &encryptionEnabled,
			MariadbServers:    MariadbServerCreateRequestList,
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
	resultList, _, err := inst.Client.Mariadb.ListMariadbClusters(ctx, &mariadb.MariadbSearchApiListMariadbClustersOpts{
		MariadbClusterName: optional.NewString(mariadbClusterName),
		Page:               optional.NewInt32(0),
		Size:               optional.NewInt32(1000),
		Sort:               optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	MariadbClusterId := resultList.Contents[0].MariadbClusterId

	if len(MariadbClusterId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForMariadb(ctx, inst.Client, MariadbClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(MariadbClusterId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifyMariadbClusterNextContract(UpdateMariadbParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, nextContractPeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(backup) != 0 {
		backupObject := &mariadb.DbClusterCreateFullBackupConfigRequest{}
		backupMap := backup[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}

		err = createMariadbClusterFullBackupConfig(UpdateMariadbParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if mariadbClusterState == common.StoppedState {
		err := stopMariadbCluster(UpdateMariadbParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMariadbRead(ctx, rd, meta)
}

func resourceMariadbRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Mariadb.DetailMariadbCluster(ctx, rd.Id())
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

	//TODO BlockStorageGroupId, BlockStorageName 받아올 지 확인
	blockStorages := database_common.HclListObject{}
	for i, bs := range dbInfo.MariadbServerGroup.BlockStorages {
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

	MariadbServers := database_common.HclListObject{}
	for _, server := range dbInfo.MariadbServerGroup.MariadbServers {
		MariadbServersInfo := database_common.HclKeyValueObject{}
		MariadbServersInfo["availability_zone_name"] = server.AvailabilityZoneName
		MariadbServersInfo["mariadb_server_name"] = server.MariadbServerName
		MariadbServersInfo["server_role_type"] = server.ServerRoleType

		MariadbServers = append(MariadbServers, MariadbServersInfo)
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

	sort.SliceStable(MariadbServers, func(i, j int) bool {
		return MariadbServers[i].(map[string]interface{})["server_role_type"].(string) < MariadbServers[j].(map[string]interface{})["server_role_type"].(string)
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
	err = rd.Set("nat_public_ip_id", rd.Get("nat_public_ip_id").(string))
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

	err = rd.Set("database_character_set", dbInfo.MariadbInitialConfig.DatabaseCharacterSet)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_name", dbInfo.MariadbInitialConfig.DatabaseName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", dbInfo.MariadbInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_user_name", dbInfo.MariadbInitialConfig.DatabaseUserName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("block_storages", blockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", dbInfo.MariadbServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("mariadb_servers", MariadbServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", dbInfo.MariadbServerGroup.ServerType)
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
	err = rd.Set("virtual_ip_address", dbInfo.MariadbServerGroup.VirtualIpAddress)
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

type UpdateMariadbParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *mariadb.MariadbClusterDetailResponse
}

func resourceMariadbUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Mariadb.DetailMariadbCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.MariadbServerGroup.MariadbServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateMariadbParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdateMariadbParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizeMariadbClusterVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updateMariadbClusterBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updateMariadbClusterSecurityGroupIds)
	}
	if rd.HasChanges("mariadb_cluster_state") {
		updateFuncs = append(updateFuncs, updateMariadbClusterServerState)
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
	return resourceMariadbRead(ctx, rd, meta)
}

func resourceMariadbDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Mariadb.DeleteMariadbCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForMariadb(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMariadbDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}

	var errorMessages []string
	mutableFields := []string{
		"server_type",
		"block_storages",
		"security_group_ids",
		"mariadb_cluster_state",
		"contract_period",
		"next_contract_period",
		"backup",
		"tags",
	}
	resourceMariadb := ResourceMariadb().Schema

	for key := range resourceMariadb {
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

func resizeMariadbClusterVirtualServers(param UpdateMariadbParam) error {
	_, _, err := param.Inst.Client.Mariadb.ResizeMariadbClusterVirtualServers(param.Ctx, param.Rd.Id(), mariadb.MariadbClusterResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateMariadbClusterBlockStorages(param UpdateMariadbParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertObjectSliceToStructSlice(oldValue)
	newList := database_common.ConvertObjectSliceToStructSlice(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizeMariadbClusterBlockStorages(param, oldList, newList)
	if err != nil {
		return err
	}

	err = addMariadbClusterBlockStorages(param, oldList, newList)
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

func resizeMariadbClusterBlockStorages(param UpdateMariadbParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.Mariadb.ResizeMariadbClusterBlockStorages(param.Ctx, param.Rd.Id(), mariadb.MariadbClusterResizeBlockStoragesRequest{
				BlockStorageGroupId: param.DbInfo.MariadbServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addMariadbClusterBlockStorages(param UpdateMariadbParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(newList)-len(oldList); i++ {
		_, _, err := param.Inst.Client.Mariadb.AddMariadbClusterBlockStorages(param.Ctx, param.Rd.Id(), mariadb.MariadbClusterAddBlockStoragesRequest{
			BlockStorageRoleType: newList[len(oldList)+i].BlockStorageRoleType,
			BlockStorageType:     newList[len(oldList)+i].BlockStorageType,
			BlockStorageSize:     int32(newList[len(oldList)+i].BlockStorageSize),
		})

		if err != nil {
			return err
		}

		err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateMariadbClusterSecurityGroupIds(param UpdateMariadbParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertSecurityGroupIdList(oldValue)
	newList := database_common.ConvertSecurityGroupIdList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachMariadbClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachMariadbClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachMariadbClusterSecurityGroup(param UpdateMariadbParam, v string) error {
	_, _, err := param.Inst.Client.Mariadb.AttachMariadbClusterSecurityGroup(param.Ctx, param.Rd.Id(), mariadb.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachMariadbClusterSecurityGroup(param UpdateMariadbParam, v string) error {
	_, _, err := param.Inst.Client.Mariadb.DetachMariadbClusterSecurityGroup(param.Ctx, param.Rd.Id(), mariadb.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateMariadbClusterServerState(param UpdateMariadbParam) error {
	_, n := param.Rd.GetChange("mariadb_cluster_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startMariadbCluster(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopMariadbCluster(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Mariadb status update failed. ")
	}

	return nil

}

func startMariadbCluster(param UpdateMariadbParam) error {
	_, _, err := param.Inst.Client.Mariadb.StartMariadbCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopMariadbCluster(param UpdateMariadbParam) error {
	_, _, err := param.Inst.Client.Mariadb.StopMariadbCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdateMariadbParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifyMariadbClusterContract(param, newValue)
	if err != nil {
		return err
	}

	return nil

}

func modifyMariadbClusterContract(param UpdateMariadbParam, newValue string) error {
	_, _, err := param.Inst.Client.Mariadb.ModifyMariadbClusterContract(param.Ctx, param.Rd.Id(), mariadb.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdateMariadbParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifyMariadbClusterNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifyMariadbClusterNextContract(param UpdateMariadbParam, newValue string) error {
	_, _, err := param.Inst.Client.Mariadb.ModifyMariadbClusterNextContract(param.Ctx, param.Rd.Id(), mariadb.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdateMariadbParam) error {
	o, n := param.Rd.GetChange("backup")

	oldValue := o.(*schema.Set)
	newValue := n.(*schema.Set)

	if oldValue.Len() == 0 {
		backupObject := &mariadb.DbClusterCreateFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = createMariadbClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deleteMariadbClusterFullBackupConfig(param)
		if err != nil {
			return err
		}
	} else {
		backupObject := &mariadb.DbClusterModifyFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = updateMariadbClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMariadbClusterFullBackupConfig(param UpdateMariadbParam, value *mariadb.DbClusterCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Mariadb.CreateMariadbClusterFullBackupConfig(param.Ctx, param.Rd.Id(), mariadb.DbClusterCreateFullBackupConfigRequest{
		ObjectStorageId:                value.ObjectStorageId,
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func updateMariadbClusterFullBackupConfig(param UpdateMariadbParam, value *mariadb.DbClusterModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Mariadb.ModifyMariadbClusterFullBackupConfig(param.Ctx, param.Rd.Id(), mariadb.DbClusterModifyFullBackupConfigRequest{
		ArchiveBackupScheduleFrequency: value.ArchiveBackupScheduleFrequency,
		BackupRetentionPeriod:          value.BackupRetentionPeriod,
		BackupStartHour:                value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deleteMariadbClusterFullBackupConfig(param UpdateMariadbParam) error {
	_, _, err := param.Inst.Client.Mariadb.DeleteMariadbClusterFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForMariadb(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func waitForMariadb(ctx context.Context, scpClient *client.SCPClient, MariadbClusterId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info mariadb.MariadbClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Mariadb.DetailMariadbCluster(ctx, MariadbClusterId)
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

		servers := info.MariadbServerGroup.MariadbServers
		log.Println("1. len(servers) : ", len(servers))

		switch len(servers) {
		case 0:
			return nil, "", fmt.Errorf("no virtual server found")
		case 1:
			log.Println("2. servers[0].MariadbServerState : ", servers[0].MariadbServerState)
			return info, servers[0].MariadbServerState, nil
		case 2:
			log.Println("3. servers[0].MariadbServerState : ", servers[0].MariadbServerState, ", servers[1].MariadbServerState : ", servers[1].MariadbServerState)
			if servers[0].MariadbServerState == common.RunningState && servers[1].MariadbServerState == common.RunningState {
				return info, servers[0].MariadbServerState, nil
			} else {
				if servers[0].MariadbServerState != common.RunningState {
					return info, servers[0].MariadbServerState, nil
				} else {
					return info, servers[1].MariadbServerState, nil
				}
			}
		default:
			return nil, "", fmt.Errorf("invalid number of virtual servers")
		}
	})
}
