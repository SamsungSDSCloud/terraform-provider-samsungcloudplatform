package redis

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/redis"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
	"strings"
	"time"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_redis", ResourceRedis())
}

func ResourceRedis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedisCreate,
		ReadContext:   resourceRedisRead,
		UpdateContext: resourceRedisUpdate,
		DeleteContext: resourceRedisDelete,
		CustomizeDiff: resourceRedisDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"tags": tfTags.TagsSchema(),
			"redis_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service Zone Id",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Redis virtual server image id.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timezone setting of this database.",
			},
			"contract_period": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract (None|1 Year|3 Year)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id of this database server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this redis DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use nat.",
			},
			"database_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "User account password of database.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"database_port": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Port number of database. (1024 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1024, 65535),
			},
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server type",
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use storage encryption.",
			},
			"block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "block storage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_storage_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage product name. (SSD|HDD)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("SSD", "HDD"),
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
			"redis_servers": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    2,
				Description: "redis servers (HA configuration when entering two server specifications)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redis_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Redis database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"nat_public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public IP for NAT. If it is null, it is automatically allocated.",
						},
						"server_role_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Server role type Enter 'ACTIVE' for a single server configuration. (MASTER | REPLICA)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("MASTER", "REPLICA"),
						},
						"nat_public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "nat ip address",
						},
					},
				},
			},
			"redis_sentinel_server": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "redis sentinel servers",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sentinel_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Redis database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"sentinel_nat_public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ublic IP for NAT. If it is null, it is automatically allocated.",
						},
						"sentinel_port": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Server role type Enter 'ACTIVE' for a single server configuration. (ACTIVE | STANDBY)",
							ValidateDiagFunc: database_common.ValidateIntegerInRange(1024, 65535),
						},
						"nat_public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "nat ip address",
						},
					},
				},
			},
			"redis_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "redis state (RUNNING|STOPPED)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("RUNNING", "STOPPED"),
			},
			"next_contract_period": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				Description:      "Next contract (None|1 Year|3 Year)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("None", database_common.OneYear, database_common.ThreeYear),
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
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vpc id",
			},
		},
		Description: "Provides a Redis Database resource.",
	}
}

func resourceRedisCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	redisName := rd.Get("redis_name").(string)
	serviceZoneId := rd.Get("service_zone_id").(string)
	imageId := rd.Get("image_id").(string)
	timezone := rd.Get("timezone").(string)
	contractPeriod := rd.Get("contract_period").(string)
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	subnetId := rd.Get("subnet_id").(string)
	natEnabled := rd.Get("nat_enabled").(bool)

	//redisInitialConfig
	databaseUserPassword := rd.Get("database_user_password").(string)
	databasePort := rd.Get("database_port").(int)

	//redisServerGroup
	serverType := rd.Get("server_type").(string)
	encryptionEnabled := rd.Get("encryption_enabled").(bool)

	//update value
	nextContractPeriod := rd.Get("next_contract_period").(string)
	redisState := rd.Get("redis_state").(string)

	//redisServerGroup
	blockStorages := rd.Get("block_storages").([]interface{})
	redisServers := rd.Get("redis_servers").([]interface{})
	redisSentinelServer := rd.Get("redis_sentinel_server").(*schema.Set).List()
	backup := rd.Get("backup").(*schema.Set).List()

	// block storage
	var RedisBlockStorageGroupCreateRequestList []redis.RedisBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertObjectSliceToStructSlice(blockStorages)
	for _, blockStorage := range blockStoragesList {
		RedisBlockStorageGroupCreateRequestList = append(RedisBlockStorageGroupCreateRequestList, redis.RedisBlockStorageGroupCreateRequest{
			BlockStorageSize: int32(blockStorage.BlockStorageSize),
			BlockStorageType: blockStorage.BlockStorageType,
		})
	}

	// redis server
	var RedisServerCreateRequestList []redis.RedisServerCreateRequest
	redisServerList := database_common.ConvertObjectSliceToStructSlice(redisServers)
	for _, redisServer := range redisServerList {
		RedisServerCreateRequestList = append(RedisServerCreateRequestList, redis.RedisServerCreateRequest{
			NatPublicIpId:   redisServer.NatPublicIpId,
			RedisServerName: redisServer.RedisServerName,
			ServerRoleType:  redisServer.ServerRoleType,
		})
	}

	sentinelObject := &redis.RedisSentinelServerCreateRequest{}
	if len(redisSentinelServer) != 0 {
		sentinelMap := redisSentinelServer[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(sentinelMap, sentinelObject)
		if err != nil {
			return diag.FromErr(err)
		}
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

	if len(redisSentinelServer) != 0 {
		_, _, err = inst.Client.Redis.CreateRedis(ctx, redis.RedisCreateRequest{
			RedisName:        redisName,
			ServiceZoneId:    serviceZoneId,
			ImageId:          imageId,
			Timezone:         timezone,
			ContractPeriod:   contractPeriod,
			SecurityGroupIds: securityGroupIdList,
			SubnetId:         subnetId,
			NatEnabled:       &natEnabled,
			RedisInitialConfig: &redis.RedisInitialConfigCreateRequest{
				DatabasePort:         int32(databasePort),
				DatabaseUserPassword: databaseUserPassword,
			},
			RedisServerGroup: &redis.RedisServerGroupCreateRequest{
				ServerType:        serverType,
				EncryptionEnabled: &encryptionEnabled,
				RedisServers:      RedisServerCreateRequestList,
				BlockStorages:     RedisBlockStorageGroupCreateRequestList,
			},
			RedisSentinelServer: sentinelObject,
		}, rd.Get("tags").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		_, _, err = inst.Client.Redis.CreateRedis(ctx, redis.RedisCreateRequest{
			RedisName:        redisName,
			ServiceZoneId:    serviceZoneId,
			ImageId:          imageId,
			Timezone:         timezone,
			ContractPeriod:   contractPeriod,
			SecurityGroupIds: securityGroupIdList,
			SubnetId:         subnetId,
			NatEnabled:       &natEnabled,
			RedisInitialConfig: &redis.RedisInitialConfigCreateRequest{
				DatabasePort:         int32(databasePort),
				DatabaseUserPassword: databaseUserPassword,
			},
			RedisServerGroup: &redis.RedisServerGroupCreateRequest{
				ServerType:        serverType,
				EncryptionEnabled: &encryptionEnabled,
				RedisServers:      RedisServerCreateRequestList,
				BlockStorages:     RedisBlockStorageGroupCreateRequestList,
			},
		}, rd.Get("tags").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty
	resultList, _, err := inst.Client.Redis.ListRedis(ctx, &redis.RedisSearchApiListRedisOpts{
		RedisName: optional.NewString(redisName),
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	redisId := resultList.Contents[0].RedisId

	if len(redisId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForRedis(ctx, inst.Client, redisId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(redisId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifyRedisNextContract(UpdateRedisParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, nextContractPeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(backup) != 0 {
		backupObject := &redis.RedisCreateFullBackupConfigRequest{}
		backupMap := backup[0].(map[string]interface{})
		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}

		err = createRedisFullBackupConfig(UpdateRedisParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if redisState == common.StoppedState {
		err := stopRedis(UpdateRedisParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRedisRead(ctx, rd, meta)
}

func resourceRedisRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Redis.DetailRedis(ctx, rd.Id())
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

	blockStorages := database_common.HclListObject{}
	for i, bs := range dbInfo.RedisServerGroup.BlockStorages {
		// Skip OS Storage
		if i == 0 {
			continue
		}
		blockStorageInfo := database_common.HclKeyValueObject{}
		blockStorageInfo["block_storage_size"] = bs.BlockStorageSize
		blockStorageInfo["block_storage_type"] = bs.BlockStorageType
		blockStorageInfo["block_storage_group_id"] = bs.BlockStorageGroupId

		blockStorages = append(blockStorages, blockStorageInfo)
	}

	redisServers := database_common.HclListObject{}
	for _, server := range dbInfo.RedisServerGroup.RedisServers {
		redisServersInfo := database_common.HclKeyValueObject{}
		redisServersInfo["redis_server_name"] = server.RedisServerName
		redisServersInfo["server_role_type"] = server.ServerRoleType
		redisServersInfo["nat_public_ip_address"] = server.NatPublicIpAddress

		redisServers = append(redisServers, redisServersInfo)
	}

	backup := database_common.HclListObject{}
	if dbInfo.BackupConfig != nil {
		backupInfo := database_common.HclKeyValueObject{}
		backupList := rd.Get("backup").(*schema.Set).List()
		if len(backupList) == 0 {
			backupInfo["object_storage_id"] = nil
		} else {
			backupInfo["object_storage_id"] = backupList[0].(map[string]interface{})["object_storage_id"]
		}
		backupInfo["backup_retention_period"] = dbInfo.BackupConfig.FullBackupConfig.BackupRetentionPeriod
		backupInfo["backup_start_hour"] = dbInfo.BackupConfig.FullBackupConfig.BackupStartHour

		backup = append(backup, backupInfo)
	}

	redisSentinelServer := database_common.HclListObject{}
	if dbInfo.SentinelServer != nil {
		redisSentinelServerInfo := database_common.HclKeyValueObject{}
		redisSentinelServerInfo["sentinel_server_name"] = dbInfo.SentinelServer.SentinelServerName
		redisSentinelServerInfo["nat_public_ip_address"] = dbInfo.SentinelServer.NatPublicIpAddress
		redisSentinelServerInfo["sentinel_port"] = rd.Get("redis_sentinel_server").(*schema.Set).List()[0].(map[string]interface{})["sentinel_port"]

		redisSentinelServer = append(redisSentinelServer, redisSentinelServerInfo)
	}

	sort.SliceStable(redisServers, func(i, j int) bool {
		return redisServers[i].(map[string]interface{})["server_role_type"].(string) < redisServers[j].(map[string]interface{})["server_role_type"].(string)
	})

	err = rd.Set("redis_name", dbInfo.RedisName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("service_zone_id", dbInfo.ServiceZoneId)
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
	err = rd.Set("contract_period", dbInfo.Contract.ContractPeriod)
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
	err = rd.Set("nat_enabled", rd.Get("nat_enabled").(bool))
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", dbInfo.RedisInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", dbInfo.RedisServerGroup.ServerType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", dbInfo.RedisServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("block_storages", blockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_servers", redisServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_sentinel_server", redisSentinelServer)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_state", dbInfo.RedisState)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("backup", backup)
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

type UpdateRedisParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *redis.RedisDetailResponse
}

func resourceRedisUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Redis.DetailRedis(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.RedisServerGroup.RedisServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateRedisParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdateRedisParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizeRedisVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updateRedisBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updateRedisSecurityGroupIds)
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

	if rd.HasChanges("redis_state") {
		err = updateRedisServerState(param)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceRedisRead(ctx, rd, meta)
}

func resourceRedisDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Redis.DeleteRedis(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForRedis(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRedisDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}

	var errorMessages []string
	mutableFields := []string{
		"server_type",
		"block_storages",
		"security_group_ids",
		"redis_state",
		"contract_period",
		"next_contract_period",
		"backup",
		"redis_sentinel_server",
		"tags",
	}
	resourceRedis := ResourceRedis().Schema

	for key := range resourceRedis {
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

func resizeRedisVirtualServers(param UpdateRedisParam) error {
	_, _, err := param.Inst.Client.Redis.ResizeRedisVirtualServers(param.Ctx, param.Rd.Id(), redis.RedisResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateRedisBlockStorages(param UpdateRedisParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertObjectSliceToStructSlice(oldValue)
	newList := database_common.ConvertObjectSliceToStructSlice(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizeRedisBlockStorages(param, oldList, newList)
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

func resizeRedisBlockStorages(param UpdateRedisParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.Redis.ResizeRedisBlockStorages(param.Ctx, param.Rd.Id(), redis.RedisResizeBlockStoragesRequest{
				BlockStorageGroupId: param.DbInfo.RedisServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateRedisSecurityGroupIds(param UpdateRedisParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertSecurityGroupIdList(oldValue)
	newList := database_common.ConvertSecurityGroupIdList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachRedisSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachRedisSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachRedisSecurityGroup(param UpdateRedisParam, v string) error {
	_, _, err := param.Inst.Client.Redis.AttachRedisSecurityGroup(param.Ctx, param.Rd.Id(), redis.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachRedisSecurityGroup(param UpdateRedisParam, v string) error {
	_, _, err := param.Inst.Client.Redis.DetachRedisSecurityGroup(param.Ctx, param.Rd.Id(), redis.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateRedisServerState(param UpdateRedisParam) error {
	_, n := param.Rd.GetChange("redis_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startRedis(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopRedis(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Redis status update failed. ")
	}

	return nil

}

func startRedis(param UpdateRedisParam) error {
	_, _, err := param.Inst.Client.Redis.StartRedis(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopRedis(param UpdateRedisParam) error {
	_, _, err := param.Inst.Client.Redis.StopRedis(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdateRedisParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifyRedisContract(param, newValue)
	if err != nil {
		return err
	}

	return nil

}

func modifyRedisContract(param UpdateRedisParam, newValue string) error {
	_, _, err := param.Inst.Client.Redis.ModifyRedisContract(param.Ctx, param.Rd.Id(), redis.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdateRedisParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifyRedisNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifyRedisNextContract(param UpdateRedisParam, newValue string) error {
	_, _, err := param.Inst.Client.Redis.ModifyRedisNextContract(param.Ctx, param.Rd.Id(), redis.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdateRedisParam) error {
	o, n := param.Rd.GetChange("backup")

	oldValue := o.(*schema.Set)
	newValue := n.(*schema.Set)

	if oldValue.Len() == 0 {
		backupObject := &redis.RedisCreateFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = createRedisFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deleteRedisFullBackupConfig(param)
		if err != nil {
			return err
		}
	} else {
		backupObject := &redis.RedisModifyFullBackupConfigRequest{}
		backupMap := newValue.List()[0].(map[string]interface{})

		err := database_common.MapToObjectWithCamel(backupMap, backupObject)
		if err != nil {
			return err
		}

		err = updateRedisFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createRedisFullBackupConfig(param UpdateRedisParam, value *redis.RedisCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Redis.CreateRedisFullBackupConfig(param.Ctx, param.Rd.Id(), redis.RedisCreateFullBackupConfigRequest{
		ObjectStorageId:       value.ObjectStorageId,
		BackupRetentionPeriod: value.BackupRetentionPeriod,
		BackupStartHour:       value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func updateRedisFullBackupConfig(param UpdateRedisParam, value *redis.RedisModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.Redis.ModifyRedisFullBackupConfig(param.Ctx, param.Rd.Id(), redis.RedisModifyFullBackupConfigRequest{
		BackupRetentionPeriod: value.BackupRetentionPeriod,
		BackupStartHour:       value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deleteRedisFullBackupConfig(param UpdateRedisParam) error {
	_, _, err := param.Inst.Client.Redis.DeleteRedisFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedis(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}

func waitForRedis(ctx context.Context, scpClient *client.SCPClient, redisId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info redis.RedisDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Redis.DetailRedis(ctx, redisId)
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

		state := info.RedisState
		log.Println("RedisState : ", state)

		return info, state, nil
	})
}
