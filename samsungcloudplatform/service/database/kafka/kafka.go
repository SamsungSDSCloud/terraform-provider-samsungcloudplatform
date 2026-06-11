package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kafka"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	samsungcloudplatform.RegisterResource("Kafka", "samsungcloudplatform_kafka", ResourceKafka())
}

func ResourceKafka() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKafkaCreate,
		ReadContext:   resourceKafkaRead,
		UpdateContext: resourceKafkaUpdate,
		DeleteContext: resourceKafkaDelete,
		CustomizeDiff: resourceKafkaDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"kafka_cluster_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"kafka_cluster_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kafka cluster state",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service Zone Id",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kafka virtual server image id.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vpc id",
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
				Description: "Security-Group ids of this Kafka. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use nat.",
			},
			"broker_sasl_account": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "SASL account of broker. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"broker_sasl_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "SASL account password of broker.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"broker_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Port number of broker. (1200 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1200, 65535),
			},
			"zookeeper_sasl_account": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "SASL account of zookeeper. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"zookeeper_sasl_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "SASL account password of zookeeper.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"zookeeper_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Port number of zookeeper. (1200 to 65535)",
				ValidateDiagFunc: database_common.ValidateIntegerInRange(1200, 65535),
			},
			"akhq_account": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Account of AKHQ. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"akhq_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port of AKHQ",
			},
			"akhq_password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Description:      "Password of AKHQ.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"broker_server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Broker Server type",
			},
			"broker_nodes": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    10,
				Description: "Broker nodes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_node_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Broker node names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"nat_public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public IP for NAT. If it is null, it is automatically allocated.",
						},
						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability zone name",
						},
					},
				},
			},
			"broker_block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Broker block storage.",
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
			"zookeeper_server_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zookeeper server type",
			},
			"zookeeper_nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    3,
				MaxItems:    3,
				Description: "Zookeeper nodes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zookeeper_node_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Zookeeper node names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"nat_public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public IP for NAT. If it is null, it is automatically allocated.",
						},
						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability zone name",
						},
					},
				},
			},
			"zookeeper_block_storages": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Zookeeper block storage.",
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
			"akhq_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to use AKHQ.",
			},
			"akhq_node": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "AKHQ node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"akhq_node_name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "AKHQ node names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)",
							ValidateDiagFunc: database_common.Validate3to20LowercaseNumberDashAndStartLowercase,
						},
						"nat_public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public IP for NAT. If it is null, it is automatically allocated.",
						},
						"akhq_availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. The single server does not input anything. (AZ1|AZ2|AZ3)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("AZ1", "AZ2", "AZ3", ""),
						},
						"akhq_server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AKHQ Server type",
						},
					},
				},
			},
			"availability_zone_config": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Availability Zone Config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_deployment_type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Availability Zone Deployment Type. UNIFORM is only available for 3AZ. (DESIGNATED|UNIFORM)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("DESIGNATED", "UNIFORM"),
						},
						"availability_zone_name": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Availability Zone Name. The single server does not input anything. (AZ1|AZ2|AZ3)",
							ValidateDiagFunc: database_common.ValidateStringInOptions("AZ1", "AZ2", "AZ3", ""),
						},
					},
				},
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timezone setting of this database.",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Kafka Database resource.",
	}
}

func resourceKafkaCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	kafkaClusterName := rd.Get("kafka_cluster_name").(string)

	serviceZoneId := rd.Get("service_zone_id").(string)
	imageId := rd.Get("image_id").(string)
	timezone := rd.Get("timezone").(string)
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	subnetId := rd.Get("subnet_id").(string)
	natEnabled := rd.Get("nat_enabled").(bool)

	brokerSaslAccount := rd.Get("broker_sasl_account").(string)
	brokerSaslPassword := rd.Get("broker_sasl_password").(string)
	brokerPort := rd.Get("broker_port").(int)
	zookeeperSaslAccount := rd.Get("zookeeper_sasl_account").(string)
	zookeeperSaslPassword := rd.Get("zookeeper_sasl_password").(string)
	zookeeperPort := rd.Get("zookeeper_port").(int)

	akhqEnabled := rd.Get("akhq_enabled").(bool)
	akhqAccount := rd.Get("akhq_account").(string)
	akhqPassword := rd.Get("akhq_password").(string)
	akhqNode := rd.Get("akhq_node").(*schema.Set).List()

	////////////////////////////////////////////////////////////////////////////////
	////////// BrokerNode 설정 시작 //////////
	brokerServerType := rd.Get("broker_server_type").(string)

	brokerNodes := rd.Get("broker_nodes").([]interface{})
	var brokerNodeCreateRequestList []kafka.BrokerNodeCreateRequest
	brokerNodeList := database_common.ConvertObjectSliceToStructSlice(brokerNodes)
	for _, brokerNode := range brokerNodeList {
		brokerNodeCreateRequestList = append(brokerNodeCreateRequestList, kafka.BrokerNodeCreateRequest{
			BrokerNodeName: brokerNode.BrokerNodeName,
			NatPublicIpId:  brokerNode.NatPublicIpId,
		})
	}

	brokerBlockStorages := rd.Get("broker_block_storages").([]interface{})
	brokerBlockStorageMap := brokerBlockStorages[0].(map[string]interface{})
	brokerBlockStorage := &kafka.KafkaNodeGroupBlockStorageGroupCreateRequest{
		BlockStorageType: brokerBlockStorageMap["block_storage_type"].(string),
		BlockStorageSize: int32(brokerBlockStorageMap["block_storage_size"].(int)),
	}
	////////// BrokerNode 설정 끝 //////////
	////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////
	////////// ZookeeperNode 설정 시작 //////////
	zookeeperServerType := rd.Get("zookeeper_server_type").(string)
	zookeeperNodes := rd.Get("zookeeper_nodes").([]interface{})
	zookeeperBlockStorages := rd.Get("zookeeper_block_storages").([]interface{})
	var zookeeperNodeGroup *kafka.ZookeeperNodeGroupCreateRequest = nil

	if len(zookeeperServerType) != 0 || len(zookeeperNodes) != 0 || len(zookeeperBlockStorages) != 0 {
		if len(zookeeperServerType) == 0 || len(zookeeperNodes) == 0 || len(zookeeperBlockStorages) == 0 {
			return diag.Errorf("zookeeper_server_type, zookeeper_nodes, zookeeper_block_storages are required to use the zookeeper node.")
		}

		var zookeeperNodeCreateRequestList []kafka.ZookeeperNodeCreateRequest
		zookeeperNodeList := database_common.ConvertObjectSliceToStructSlice(zookeeperNodes)
		for _, zookeeperNode := range zookeeperNodeList {
			zookeeperNodeCreateRequestList = append(zookeeperNodeCreateRequestList, kafka.ZookeeperNodeCreateRequest{
				ZookeeperNodeName: zookeeperNode.ZookeeperNodeName,
				NatPublicIpId:     zookeeperNode.NatPublicIpId,
			})
		}

		zookeeperBlockStorageMap := zookeeperBlockStorages[0].(map[string]interface{})
		zookeeperBlockStorage := &kafka.KafkaNodeGroupBlockStorageGroupCreateRequest{
			BlockStorageType: zookeeperBlockStorageMap["block_storage_type"].(string),
			BlockStorageSize: int32(zookeeperBlockStorageMap["block_storage_size"].(int)),
		}

		zookeeperNodeGroup = &kafka.ZookeeperNodeGroupCreateRequest{
			ServerType:     zookeeperServerType,
			ZookeeperNodes: zookeeperNodeCreateRequestList,
			BlockStorage:   zookeeperBlockStorage,
		}
	}
	////////// ZookeeperNode 설정 끝 //////////
	////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////
	////////// AkhqNode 설정 시작 //////////
	var akhqInitialConfig *kafka.AkhqInitialConfigCreateRequest = nil
	var akhqNodeGroup *kafka.AkhqNodeGroupCreateRequest = nil
	if akhqEnabled {
		if len(akhqAccount) == 0 || len(akhqPassword) == 0 || len(akhqNode) == 0 {
			return diag.Errorf("akhq_account, akhq_password, akhq_node are required to enable the AKHQ.")
		}

		akhqInitialConfig = &kafka.AkhqInitialConfigCreateRequest{
			AkhqAccount:  akhqAccount,
			AkhqPassword: akhqPassword,
		}

		akhqNodeMap := akhqNode[0].(map[string]interface{})
		akhqNodeName := akhqNodeMap["akhq_node_name"].(string)
		akhqNatPublicIpId := akhqNodeMap["nat_public_ip_id"].(string)
		akhqAvailabilityZoneName := akhqNodeMap["akhq_availability_zone_name"].(string)

		akhqNodeGroup = &kafka.AkhqNodeGroupCreateRequest{
			AkhqNodeName:             akhqNodeName,
			NatPublicIpId:            akhqNatPublicIpId,
			AkhqAvailabilityZoneName: akhqAvailabilityZoneName,
		}
	}
	////////// AkhqNode 설정 끝 //////////
	////////////////////////////////////////////////////////////////////////////////

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

	azConfig := rd.Get("availability_zone_config").(*schema.Set).List()
	var availabilityZoneConfig *kafka.KafkaClusterCreateAvailabilityZoneConfig = nil
	if len(azConfig) != 0 {
		azConfigMap := azConfig[0].(map[string]interface{})
		availabilityZoneConfig = &kafka.KafkaClusterCreateAvailabilityZoneConfig{
			AvailabilityZoneDeploymentType: azConfigMap["availability_zone_deployment_type"].(string),
			AvailabilityZoneName:           azConfigMap["availability_zone_name"].(string),
		}
	}

	_, _, err = inst.Client.Kafka.CreateKafkaCluster(ctx, kafka.KafkaClusterCreateRequest{
		KafkaClusterName: kafkaClusterName,
		ServiceZoneId:    serviceZoneId,
		ImageId:          imageId,
		Timezone:         timezone,
		SecurityGroupIds: securityGroupIdList,
		SubnetId:         subnetId,
		NatEnabled:       &natEnabled,
		KafkaInitialConfig: &kafka.KafkaInitialConfigCreateRequest{
			BrokerInitialConfig: &kafka.BrokerInitialConfigCreateRequest{
				BrokerSaslAccount:  brokerSaslAccount,
				BrokerSaslPassword: brokerSaslPassword,
				BrokerPort:         int32(brokerPort),
			},
			ZookeeperInitialConfig: &kafka.ZookeeperInitialConfigCreateRequest{
				ZookeeperSaslAccount:  zookeeperSaslAccount,
				ZookeeperSaslPassword: zookeeperSaslPassword,
				ZookeeperPort:         int32(zookeeperPort),
			},
			AkhqInitialConfig: akhqInitialConfig,
		},
		BrokerNodeGroup: &kafka.BrokerNodeGroupCreateRequest{
			ServerType:   brokerServerType,
			BrokerNodes:  brokerNodeCreateRequestList,
			BlockStorage: brokerBlockStorage,
		},
		ZookeeperNodeGroup:     zookeeperNodeGroup,
		AkhqEnabled:            &akhqEnabled,
		AkhqNodeGroup:          akhqNodeGroup,
		AvailabilityZoneConfig: availabilityZoneConfig,
	}, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		fmt.Printf("%s err...\n", err)
		return diag.FromErr(err)
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty
	resultList, _, err := inst.Client.Kafka.ListKafkaClusters(ctx, &kafka.KafkaSearchApiListKafkaClustersOpts{
		KafkaClusterName: optional.NewString(kafkaClusterName),
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

	kafkaClusterId := resultList.Contents[0].KafkaClusterId

	if len(kafkaClusterId) == 0 {
		diagnostics = diag.Errorf("Kafka_cluster_id not found")
		return
	}

	err = waitForKafka(ctx, inst.Client, kafkaClusterId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(kafkaClusterId)

	return resourceKafkaRead(ctx, rd, meta)
}

func resourceKafkaRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Kafka.DetailKafkaCluster(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	if len(dbInfo.BrokerNodeGroup.BrokerNodes) == 0 {
		diagnostics = diag.Errorf("Kafka is not found")
		return
	}
	fs := &common.FieldSetter{Rd: rd}
	fs.Set("kafka_cluster_name", dbInfo.KafkaClusterName)
	fs.Set("kafka_cluster_state", dbInfo.KafkaClusterState)
	fs.Set("service_zone_id", dbInfo.ServiceZoneId)
	fs.Set("image_id", dbInfo.ImageId)
	fs.Set("vpc_id", dbInfo.VpcId)
	fs.Set("subnet_id", dbInfo.SubnetId)
	fs.Set("security_group_ids", dbInfo.SecurityGroupIds)
	fs.Set("nat_enabled", rd.Get("nat_enabled").(bool))
	fs.Set("broker_sasl_account", dbInfo.KafkaInitialConfig.BrokerInitialConfig.BrokerSaslAccount)
	fs.Set("broker_port", dbInfo.KafkaInitialConfig.BrokerInitialConfig.BrokerPort)
	fs.Set("zookeeper_sasl_account", dbInfo.KafkaInitialConfig.ZookeeperInitialConfig.ZookeeperSaslAccount)
	fs.Set("zookeeper_port", dbInfo.KafkaInitialConfig.ZookeeperInitialConfig.ZookeeperPort)
	fs.Set("broker_server_type", dbInfo.BrokerNodeGroup.ServerType)
	brokerNodes := database_common.HclListObject{}
	for _, node := range dbInfo.BrokerNodeGroup.BrokerNodes {
		brokerNodeInfo := database_common.HclKeyValueObject{}
		brokerNodeInfo["broker_node_name"] = node.BrokerNodeName
		brokerNodeInfo["availability_zone_name"] = node.AvailabilityZoneName
		brokerNodes = append(brokerNodes, brokerNodeInfo)
	}

	brokerBlockStorages := database_common.HclListObject{}
	for _, bs := range dbInfo.BrokerNodeGroup.BlockStorages {
		if bs.BlockStorageRoleType == "OS" {
			// Skip OS Storage
			continue
		}

		blockStorageInfo := database_common.HclKeyValueObject{}
		blockStorageInfo["block_storage_size"] = bs.BlockStorageSize
		blockStorageInfo["block_storage_type"] = bs.BlockStorageType
		blockStorageInfo["block_storage_group_id"] = bs.BlockStorageGroupId
		brokerBlockStorages = append(brokerBlockStorages, blockStorageInfo)
	}
	fs.Set("broker_block_storages", brokerBlockStorages)

	if dbInfo.ZookeeperNodeGroup != nil {
		fs.Set("zookeeper_server_type", dbInfo.ZookeeperNodeGroup.ServerType)
		for _, node := range dbInfo.ZookeeperNodeGroup.ZookeeperNodes {
			brokerNodeInfo := database_common.HclKeyValueObject{}
			brokerNodeInfo["zookeeper_node_name"] = node.ZookeeperNodeName
			brokerNodeInfo["availability_zone_name"] = node.AvailabilityZoneName
			brokerNodes = append(brokerNodes, brokerNodeInfo)
		}

		zookeeperBlockStorages := database_common.HclListObject{}
		for _, bs := range dbInfo.ZookeeperNodeGroup.BlockStorages {
			if bs.BlockStorageRoleType == "OS" {
				// Skip OS Storage
				continue
			}

			blockStorageInfo := database_common.HclKeyValueObject{}
			blockStorageInfo["block_storage_size"] = bs.BlockStorageSize
			blockStorageInfo["block_storage_type"] = bs.BlockStorageType
			blockStorageInfo["block_storage_group_id"] = bs.BlockStorageGroupId
			zookeeperBlockStorages = append(zookeeperBlockStorages, blockStorageInfo)
		}
		fs.Set("zookeeper_block_storages", zookeeperBlockStorages)
	}

	fs.Set("akhq_enabled", rd.Get("akhq_enabled").(bool))

	if dbInfo.AkhqNodeGroup != nil {
		fs.Set("akhq_account", dbInfo.KafkaInitialConfig.AkhqInitialConfig.AkhqAccount)
		fs.Set("akhq_port", dbInfo.KafkaInitialConfig.AkhqInitialConfig.AkhqPort)
		akhqNode := database_common.HclListObject{}
		akhqNodeInfo := database_common.HclKeyValueObject{}
		akhqNodeInfo["akhq_node_name"] = dbInfo.AkhqNodeGroup.AkhqNodeName
		akhqNodeInfo["akhq_availability_zone_name"] = dbInfo.AkhqNodeGroup.AvailabilityZoneName
		akhqNodeInfo["akhq_server_type"] = dbInfo.AkhqNodeGroup.ServerType
		akhqNode = append(akhqNode, akhqNodeInfo)
		fs.Set("akhq_node", akhqNode)
	}

	availabilityZoneConfig := rd.Get("availability_zone_config").(*schema.Set).List()
	fs.Set("availability_zone_config", availabilityZoneConfig)
	fs.Set("timezone", dbInfo.Timezone)

	if diags := fs.Diags(); diags != nil {
		rd.SetId("")
		return diags
	}
	err = tfTags.SetTags(ctx, rd, meta, rd.Id())
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Kafka.DetailKafkaCluster(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.BrokerNodeGroup.BrokerNodes) == 0 {
		diagnostics = diag.Errorf("Kafka is not found")
		return
	}

	param := UpdateKafkaParam{
		Ctx:    ctx,
		Rd:     rd,
		Inst:   inst,
		DbInfo: &dbInfo,
	}

	var updateFuncs []func(serverParam UpdateKafkaParam) error

	for _, f := range updateFuncs {
		err = f(param)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKafkaRead(ctx, rd, meta)
}

func resourceKafkaDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Kafka.DeleteKafkaCluster(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitForKafka(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaDiff(ctx context.Context, rd *schema.ResourceDiff, meta interface{}) error {
	if rd.Id() == "" {
		return nil
	}
	return nil
}

func waitForKafka(ctx context.Context, scpClient *client.SCPClient, kafkaClusterId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		var info kafka.KafkaClusterDetailResponse
		var statusCode int
		var err error
		retryCount := 10

		for i := 0; i < retryCount; i++ {
			info, statusCode, err = scpClient.Kafka.DetailKafkaCluster(ctx, kafkaClusterId)
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

		state := info.KafkaClusterState

		return info, state, nil
	})
}

type UpdateKafkaParam struct {
	Ctx    context.Context
	Rd     *schema.ResourceData
	Inst   *client.Instance
	DbInfo *kafka.KafkaClusterDetailResponse
}
