package rediscluster

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
	samsungcloudplatform.RegisterResource("samsungcloudplatform_redis_cluster", ResourceRedisCluster())
}

func ResourceRedisCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedisClusterCreate,
		ReadContext:   resourceRedisClusterRead,
		UpdateContext: resourceRedisClusterUpdate,
		DeleteContext: resourceRedisClusterDelete,
		CustomizeDiff: resourceRedisClusterDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
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
				Description: "Redis Cluster virtual server image id.",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use nat.",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this redisCluster DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
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
			"redis_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"redis_cluster_state": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Redis Cluster state (RUNNING|STOPPED)",
				ValidateDiagFunc: database_common.ValidateStringInOptions("RUNNING", "STOPPED"),
			},
			//initialconfig
			"database_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "User account password of database.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"database_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Port number of this database. (1024 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1024, 65535),
			},
			//servergroup
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server type.",
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use storage encryption.",
			},
			"redis_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "RedisCluster servers",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redis_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "RedisCluster database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
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
							Description:      "Server role type Enter 'MASTER' for a single server configuration. (MASTER | REPLICA)\",",
							ValidateDiagFunc: database_common.ValidateStringInOptions("MASTER", "REPLICA"),
						},
						"nat_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "nat ip address",
						},
					},
				},
			},
			"block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "block storage. (It can't be deleted.)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_storage_role_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Storage usage. (Only DATA)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("DATA"),
						},
						"block_storage_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Storage product name. (SSD|HDD)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("SSD", "HDD"),
						},
						"block_storage_size": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Block Storage Size (50 to 5120)",
							ValidateDiagFunc: database_common.ValidateIntegerInRange(50, 5120),
						},
						"block_storage_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Block storage group id",
						},
					},
				},
			},
			"shards_count": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Number of Masters.",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(3, 40),
			},
			"shards_replica_count": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Number of Replicas created per Master.",
				ValidateDiagFunc: database_common.ValidateIntegerLessEqualThan(3),
			},
			"backup": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceRedisClusterBackup(),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vpc id",
			},
			"tags": tfTags.TagsSchema(),
		},
	}

}

func resourceRedisClusterBackup() *schema.Resource {
	return &schema.Resource{
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
	}
}

func resourceRedisClusterCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
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
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	serviceZoneId := rd.Get("service_zone_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	timezone := rd.Get("timezone").(string)
	redisClusterName := rd.Get("redis_cluster_name").(string)
	redisClusterState := rd.Get("redis_cluster_state").(string)

	//redisClusterInitialConfig
	databasePort := rd.Get("database_port").(int)
	databaseUserPassword := rd.Get("database_user_password").(string)

	//redisClusterServerGroup
	serverType := rd.Get("server_type").(string)
	encryptionEnabled := rd.Get("encryption_enabled").(bool)
	redisServers := rd.Get("redis_servers").([]interface{})
	blockStorages := rd.Get("block_storages").([]interface{})

	shardsCount := rd.Get("shards_count").(int)
	shardsReplicaCount := rd.Get("shards_replica_count").(int)

	backup := rd.Get("backup").(*schema.Set).List()

	// block storage (HclListObject to Slice)
	var RedisClusterBlockStorageGroupCreateRequestList []redis.RedisBlockStorageGroupCreateRequest
	blockStoragesList := database_common.ConvertObjectSliceToStructSlice(blockStorages)
	for _, blockStorage := range blockStoragesList {
		RedisClusterBlockStorageGroupCreateRequestList = append(RedisClusterBlockStorageGroupCreateRequestList, redis.RedisBlockStorageGroupCreateRequest{
			BlockStorageSize: int32(blockStorage.BlockStorageSize),
			BlockStorageType: blockStorage.BlockStorageType,
		})
	}

	// redis cluster server (HclListObject to Slice)
	var RedisClusterServerCreateRequestList []redis.RedisServerCreateRequest
	redisClusterServerList := database_common.ConvertObjectSliceToStructSlice(redisServers)
	for _, redisServer := range redisClusterServerList {
		RedisClusterServerCreateRequestList = append(RedisClusterServerCreateRequestList, redis.RedisServerCreateRequest{
			NatPublicIpId:   redisServer.NatPublicIpId,
			RedisServerName: redisServer.RedisServerName,
			ServerRoleType:  redisServer.ServerRoleType,
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

	_, _, err = inst.Client.RedisCluster.CreateRedisCluster(ctx, redis.RedisClusterCreateRequest{
		ContractPeriod: contractPeriod,
		ImageId:        imageId,
		NatEnabled:     &natEnabled,
		RedisName:      redisClusterName,
		RedisInitialConfig: &redis.RedisInitialConfigCreateRequest{
			DatabasePort:         int32(databasePort),
			DatabaseUserPassword: databaseUserPassword,
		},
		RedisServerGroup: &redis.RedisServerGroupCreateRequest{
			BlockStorages:     RedisClusterBlockStorageGroupCreateRequestList,
			EncryptionEnabled: &encryptionEnabled,
			RedisServers:      RedisClusterServerCreateRequestList,
			ServerType:        serverType,
		},
		ShardsCount:        int32(shardsCount),
		ShardsReplicaCount: int32(shardsReplicaCount),
		SecurityGroupIds:   securityGroupIdList,
		ServiceZoneId:      serviceZoneId,
		SubnetId:           subnetId,
		Timezone:           timezone,
	}, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty
	resultList, _, err := inst.Client.RedisCluster.ListRedisCluster(ctx, &redis.RedisClusterSearchApiListRedisClusterOpts{
		RedisName: optional.NewString(redisClusterName),
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

	redisClusterId := resultList.Contents[0].RedisId

	if len(redisClusterId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	log.Printf("redisClusterId : %s", redisClusterId)
	err = waitForRedisCluster(ctx, inst.Client, redisClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(redisClusterId)

	if nextContractPeriod == database_common.OneYear || nextContractPeriod == database_common.ThreeYear {
		err := modifyRedisClusterNextContract(UpdateRedisClusterParam{
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

		err = createRedisClusterFullBackupConfig(UpdateRedisClusterParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		}, backupObject)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if redisClusterState == common.StoppedState {
		err := stopRedisCluster(UpdateRedisClusterParam{
			Ctx:  ctx,
			Rd:   rd,
			Inst: inst,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRedisClusterRead(ctx, rd, meta)
}

func resourceRedisClusterRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.RedisCluster.DetailRedisCluster(ctx, rd.Id())
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

	//TODO BlockStorageGroupId, BlockStorageName 받아올 지 확인
	blockStorages := database_common.HclListObject{}
	for i, bs := range dbInfo.RedisServerGroup.BlockStorages {
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

	redisClusterServers := database_common.HclListObject{}
	for _, server := range dbInfo.RedisServerGroup.RedisServers {
		redisClusterServersInfo := database_common.HclKeyValueObject{}
		redisClusterServersInfo["redis_server_name"] = server.RedisServerName
		redisClusterServersInfo["nat_public_ip_id"] = server.NatPublicIpAddress
		redisClusterServersInfo["server_role_type"] = server.ServerRoleType
		redisClusterServersInfo["nat_ip_address"] = server.NatPublicIpAddress
		redisClusterServersInfo["created_dt"] = server.CreatedDt

		redisClusterServers = append(redisClusterServers, redisClusterServersInfo)
	}

	log.Printf("redisClusterServers : %s", redisClusterServers)

	// created_dt 로 sort
	sort.Slice(redisClusterServers, func(i, j int) bool {
		return redisClusterServers[i].(map[string]interface{})["created_dt"].(time.Time).Before(redisClusterServers[j].(map[string]interface{})["created_dt"].(time.Time))
	})

	log.Printf("redisClusterServers(After Sorting) : %s", redisClusterServers)

	// created_dt 제거
	redisClusterServersSortedCreatedDt := database_common.ConvertObjectSliceToStructSlice(redisClusterServers)

	redisClusterServersExcludeCreatedDt := database_common.HclListObject{}
	for _, server := range redisClusterServersSortedCreatedDt {
		redisClusterServersInfo := database_common.HclKeyValueObject{}
		redisClusterServersInfo["redis_server_name"] = server.RedisServerName
		redisClusterServersInfo["nat_public_ip_id"] = server.NatPublicIpId
		redisClusterServersInfo["server_role_type"] = server.ServerRoleType
		redisClusterServersInfo["nat_ip_address"] = server.NatPublicIpAddress

		redisClusterServersExcludeCreatedDt = append(redisClusterServersExcludeCreatedDt, redisClusterServersInfo)
	}
	log.Printf("redisClusterServers2(After deleting createdDt) : %s", redisClusterServersExcludeCreatedDt)

	backup := database_common.HclListObject{}
	if dbInfo.BackupConfig != nil {
		backupInfo := database_common.HclKeyValueObject{}
		backupInfo["object_storage_id"] = rd.Get("backup").(*schema.Set).List()[0].(map[string]interface{})["object_storage_id"]
		backupInfo["backup_retention_period"] = dbInfo.BackupConfig.FullBackupConfig.BackupRetentionPeriod
		backupInfo["backup_start_hour"] = dbInfo.BackupConfig.FullBackupConfig.BackupStartHour

		backup = append(backup, backupInfo)
	}

	err = rd.Set("contract_period", dbInfo.Contract.ContractPeriod)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("image_id", dbInfo.ImageId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_cluster_name", dbInfo.RedisName)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_cluster_state", dbInfo.RedisState)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("database_port", dbInfo.RedisInitialConfig.DatabasePort)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("block_storages", blockStorages)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("encryption_enabled", dbInfo.RedisServerGroup.EncryptionEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("redis_servers", redisClusterServersExcludeCreatedDt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = rd.Set("server_type", dbInfo.RedisServerGroup.ServerType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = rd.Set("backup", backup)
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

type UpdateRedisClusterParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *redis.RedisClusterDetailResponse
}

func resourceRedisClusterUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.RedisCluster.DetailRedisCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.RedisServerGroup.RedisServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateRedisClusterParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdateRedisClusterParam) error

	if rd.HasChanges("server_type") {
		updateFuncs = append(updateFuncs, resizeRedisClusterVirtualServers)
	}
	if rd.HasChanges("block_storages") {
		updateFuncs = append(updateFuncs, updateRedisClusterBlockStorages)
	}
	if rd.HasChanges("security_group_ids") {
		updateFuncs = append(updateFuncs, updateRedisClusterSecurityGroupIds)
	}
	if rd.HasChanges("redis_cluster_state") {
		updateFuncs = append(updateFuncs, updateRedisClusterServerState)
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
	return resourceRedisClusterRead(ctx, rd, meta)
}

func resourceRedisClusterDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.RedisCluster.DeleteRedisCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForRedisCluster(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRedisClusterDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}

	var errorMessages []string
	mutableFields := []string{
		"server_type",
		"block_storages",
		"security_group_ids",
		"redis_cluster_state",
		"contract_period",
		"next_contract_period",
		"backup",
		"tags",
	}
	resourceRedisCluster := ResourceRedisCluster().Schema

	for key, _ := range resourceRedisCluster {
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

func resizeRedisClusterVirtualServers(param UpdateRedisClusterParam) error {
	_, _, err := param.Inst.Client.RedisCluster.ResizeRedisClusterVirtualServers(param.Ctx, param.Rd.Id(), redis.RedisResizeVirtualServersRequest{
		ServerType: param.Rd.Get("server_type").(string),
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), database_common.DatabaseProcessingAndStoppedStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateRedisClusterBlockStorages(param UpdateRedisClusterParam) error {
	o, n := param.Rd.GetChange("block_storages")
	oldValue := o.([]interface{})
	newValue := n.([]interface{})

	oldList := database_common.ConvertObjectSliceToStructSlice(oldValue)
	newList := database_common.ConvertObjectSliceToStructSlice(newValue)

	err := validateBlockStorageInput(oldList, newList)
	if err != nil {
		return err
	}

	err = resizeRedisClusterBlockStorages(param, oldList, newList)
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

func resizeRedisClusterBlockStorages(param UpdateRedisClusterParam, oldList []database_common.ConvertedStruct, newList []database_common.ConvertedStruct) error {
	for i := 0; i < len(oldList); i++ {
		if oldList[i].BlockStorageSize < newList[i].BlockStorageSize {

			_, _, err := param.Inst.Client.RedisCluster.ResizeRedisClusterBlockStorages(param.Ctx, param.Rd.Id(), redis.RedisResizeBlockStoragesRequest{
				BlockStorageGroupId: param.DbInfo.RedisServerGroup.BlockStorages[i+1].BlockStorageGroupId,
				BlockStorageSize:    int32(newList[i].BlockStorageSize),
			})
			if err != nil {
				return err
			}

			err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateRedisClusterSecurityGroupIds(param UpdateRedisClusterParam) error {
	o, n := param.Rd.GetChange("security_group_ids")
	oldValue := o.(common.HclListObject)
	newValue := n.(common.HclListObject)

	oldList := database_common.ConvertToList(oldValue)
	newList := database_common.ConvertToList(newValue)

	for _, v := range newList {
		if !database_common.Contains(oldList, v) {
			if err := attachRedisClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	for _, v := range oldList {
		if !database_common.Contains(newList, v) {
			if err := detachRedisClusterSecurityGroup(param, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachRedisClusterSecurityGroup(param UpdateRedisClusterParam, v string) error {
	_, _, err := param.Inst.Client.RedisCluster.AttachRedisClusterSecurityGroup(param.Ctx, param.Rd.Id(), redis.DbClusterAttachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func detachRedisClusterSecurityGroup(param UpdateRedisClusterParam, v string) error {
	_, _, err := param.Inst.Client.RedisCluster.DetachRedisClusterSecurityGroup(param.Ctx, param.Rd.Id(), redis.DbClusterDetachSecurityGroupRequest{
		SecurityGroupId: v,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	return nil
}

func updateRedisClusterServerState(param UpdateRedisClusterParam) error {
	_, n := param.Rd.GetChange("redis_cluster_state")
	newVal := n.(string)

	if newVal == common.RunningState {
		err := startRedisCluster(param)
		if err != nil {
			return err
		}
	} else if newVal == common.StoppedState {
		err := stopRedisCluster(param)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Redis Clsuter status update failed. ")
	}

	return nil

}

func startRedisCluster(param UpdateRedisClusterParam) error {
	_, _, err := param.Inst.Client.RedisCluster.StartRedisCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func stopRedisCluster(param UpdateRedisClusterParam) error {
	_, _, err := param.Inst.Client.RedisCluster.StopRedisCluster(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.StoppedState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateContractPeriod(param UpdateRedisClusterParam) error {
	o, n := param.Rd.GetChange("contract_period")

	oldValue := o.(string)
	newValue := n.(string)

	if oldValue != database_common.None {
		return fmt.Errorf("changing contract period is not allowed")
	}

	err := modifyRedisClusterContract(param, newValue)
	if err != nil {
		return err
	}

	return nil

}

func modifyRedisClusterContract(param UpdateRedisClusterParam, newValue string) error {
	_, _, err := param.Inst.Client.RedisCluster.ModifyRedisClusterContract(param.Ctx, param.Rd.Id(), redis.DbClusterModifyContractRequest{
		ContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateNextContractPeriod(param UpdateRedisClusterParam) error {
	_, n := param.Rd.GetChange("next_contract_period")

	newValue := n.(string)

	err := modifyRedisClusterNextContract(param, newValue)
	if err != nil {
		return err
	}

	return nil
}

func modifyRedisClusterNextContract(param UpdateRedisClusterParam, newValue string) error {
	_, _, err := param.Inst.Client.RedisCluster.ModifyRedisClusterNextContract(param.Ctx, param.Rd.Id(), redis.DbClusterModifyNextContractRequest{
		NextContractPeriod: newValue,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateBackup(param UpdateRedisClusterParam) error {
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

		err = createRedisClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	} else if newValue.Len() == 0 {
		err := deleteRedisClusterFullBackupConfig(param)
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

		err = updateRedisClusterFullBackupConfig(param, backupObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func createRedisClusterFullBackupConfig(param UpdateRedisClusterParam, value *redis.RedisCreateFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.RedisCluster.CreateRedisClusterFullBackupConfig(param.Ctx, param.Rd.Id(), redis.RedisCreateFullBackupConfigRequest{
		ObjectStorageId:       value.ObjectStorageId,
		BackupRetentionPeriod: value.BackupRetentionPeriod,
		BackupStartHour:       value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func updateRedisClusterFullBackupConfig(param UpdateRedisClusterParam, value *redis.RedisModifyFullBackupConfigRequest) error {
	_, _, err := param.Inst.Client.RedisCluster.ModifyRedisClusterFullBackupConfig(param.Ctx, param.Rd.Id(), redis.RedisModifyFullBackupConfigRequest{
		BackupRetentionPeriod: value.BackupRetentionPeriod,
		BackupStartHour:       value.BackupStartHour,
	})
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil
}

func deleteRedisClusterFullBackupConfig(param UpdateRedisClusterParam) error {
	_, _, err := param.Inst.Client.RedisCluster.DeleteRedisClusterFullBackupConfig(param.Ctx, param.Rd.Id())
	if err != nil {
		return err
	}

	err = waitForRedisCluster(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}
	return nil

}
func waitForRedisCluster(ctx context.Context, scpClient *client.SCPClient, redisId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info redis.RedisClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.RedisCluster.DetailRedisCluster(ctx, redisId)
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
		log.Println("RedisClusterState : ", state)

		return info, state, nil

	})
}
