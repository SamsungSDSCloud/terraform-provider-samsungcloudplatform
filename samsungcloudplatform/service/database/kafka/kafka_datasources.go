package kafka

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/database/database_common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("Kafka", "samsungcloudplatform_kafka", DatasourceKafka())
}

func DatasourceKafka() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKafkaSingle,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"kafka_cluster_id":    {Type: schema.TypeString, Required: true, Description: "Kafka Cluster Id"},
			"kafka_cluster_name":  {Type: schema.TypeString, Computed: true, Description: "Kafka Cluster Name"},
			"kafka_cluster_state": {Type: schema.TypeString, Computed: true, Description: "Kafka Cluster State"},
			"image_id":            {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			"timezone":            {Type: schema.TypeString, Computed: true, Description: "Timezone"},
			"vpc_id":              {Type: schema.TypeString, Computed: true, Description: "VPC Id"},
			"subnet_id":           {Type: schema.TypeString, Computed: true, Description: "Subnet Id"},
			"security_group_ids":  {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Security group ids"},
			"database_version":    {Type: schema.TypeString, Computed: true, Description: "database version"},
			"nat_ip_address":      {Type: schema.TypeString, Computed: true, Description: "nat ip address"},
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
			"akhq_node_group": {Type: schema.TypeList, Computed: true, Description: "AKHQ Node group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"akhq_node_id":            {Type: schema.TypeString, Computed: true, Description: "AKHQ node ID"},
						"akhq_node_name":          {Type: schema.TypeString, Computed: true, Description: "AKHQ node name"},
						"akhq_node_state":         {Type: schema.TypeString, Computed: true, Description: "AKHQ node state"},
						"block_storage_group_id":  {Type: schema.TypeString, Computed: true, Description: "block storage group id"},
						"block_storage_name":      {Type: schema.TypeString, Computed: true, Description: "block storage name"},
						"block_storage_role_type": {Type: schema.TypeString, Computed: true, Description: "block storage role type"},
						"block_storage_size":      {Type: schema.TypeInt, Computed: true, Description: "block Storage size"},
						"block_storage_type":      {Type: schema.TypeString, Computed: true, Description: "block storage type"},
						"node_role_type":          {Type: schema.TypeString, Computed: true, Description: "node role type"},
						"server_type":             {Type: schema.TypeString, Computed: true, Description: "server type"},
						"subnet_ip_address":       {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
						"availability_zone_name":  {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
					},
				},
			},
			"broker_node_group": {Type: schema.TypeList, Computed: true, Description: "Broker Node group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_role_type": {Type: schema.TypeString, Computed: true, Description: "node role type"},
						"server_type":    {Type: schema.TypeString, Computed: true, Description: "server type"},
						"block_storages": {Type: schema.TypeList, Computed: true, Description: "block storages",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_storage_group_id":  {Type: schema.TypeString, Computed: true, Description: "block storage group id"},
									"block_storage_name":      {Type: schema.TypeString, Computed: true, Description: "block storage name"},
									"block_storage_role_type": {Type: schema.TypeString, Computed: true, Description: "block storage role type"},
									"block_storage_size":      {Type: schema.TypeInt, Computed: true, Description: "block Storage size"},
									"block_storage_type":      {Type: schema.TypeString, Computed: true, Description: "block storage type"},
								},
							},
						},
						"broker_nodes": {Type: schema.TypeList, Computed: true, Description: "Broker nodes",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"broker_node_id":         {Type: schema.TypeString, Computed: true, Description: "Broker node id"},
									"broker_node_name":       {Type: schema.TypeString, Computed: true, Description: "Broker node name"},
									"broker_node_state":      {Type: schema.TypeString, Computed: true, Description: "Broker node state"},
									"subnet_ip_address":      {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
									"created_by":             {Type: schema.TypeString, Computed: true, Description: "created by"},
									"created_dt":             {Type: schema.TypeString, Computed: true, Description: "created dt"},
									"modified_by":            {Type: schema.TypeString, Computed: true, Description: "modified by"},
									"modified_dt":            {Type: schema.TypeString, Computed: true, Description: "modified dt"},
								},
							},
						},
					},
				},
			},
			"kafka_initial_config": {Type: schema.TypeList, Computed: true, Description: "Kafka Cluster initial config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"akhq_initial_config": {Type: schema.TypeList, Computed: true, Description: "AKHQ initial config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"akhq_account": {Type: schema.TypeString, Computed: true, Description: "AKHQ Account"},
									"akhq_port":    {Type: schema.TypeInt, Computed: true, Description: "AKHQ Port"},
								},
							},
						},
						"broker_initial_config": {Type: schema.TypeList, Computed: true, Description: "Broker initial config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"broker_sasl_account": {Type: schema.TypeString, Computed: true, Description: "Broker SASL Account"},
									"broker_port":         {Type: schema.TypeInt, Computed: true, Description: "Broker Port"},
								},
							},
						},
						"zookeeper_initial_config": {Type: schema.TypeList, Computed: true, Description: "Zookeeper initial config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zookeeper_sasl_account": {Type: schema.TypeString, Computed: true, Description: "Zookeeper SASL Account"},
									"zookeeper_port":         {Type: schema.TypeInt, Computed: true, Description: "Zookeeper Port"},
								},
							},
						},
					},
				},
			},
			"zookeeper_node_group": {Type: schema.TypeList, Computed: true, Description: "Zookeeper Node group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_role_type": {Type: schema.TypeString, Computed: true, Description: "node role type"},
						"server_type":    {Type: schema.TypeString, Computed: true, Description: "server type"},
						"block_storages": {Type: schema.TypeList, Computed: true, Description: "block storages",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_storage_group_id":  {Type: schema.TypeString, Computed: true, Description: "block storage group id"},
									"block_storage_name":      {Type: schema.TypeString, Computed: true, Description: "block storage name"},
									"block_storage_role_type": {Type: schema.TypeString, Computed: true, Description: "block storage role type"},
									"block_storage_size":      {Type: schema.TypeInt, Computed: true, Description: "block Storage size"},
									"block_storage_type":      {Type: schema.TypeString, Computed: true, Description: "block storage type"},
								},
							},
						},
						"zookeeper_nodes": {Type: schema.TypeList, Computed: true, Description: "Zookeeper nodes",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zookeeper_node_id":      {Type: schema.TypeString, Computed: true, Description: "Zookeeper node id"},
									"zookeeper_node_name":    {Type: schema.TypeString, Computed: true, Description: "Zookeeper node name"},
									"zookeeper_node_state":   {Type: schema.TypeString, Computed: true, Description: "Zookeeper node state"},
									"subnet_ip_address":      {Type: schema.TypeString, Computed: true, Description: "subnet ip address"},
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
									"created_by":             {Type: schema.TypeString, Computed: true, Description: "created by"},
									"created_dt":             {Type: schema.TypeString, Computed: true, Description: "created dt"},
									"modified_by":            {Type: schema.TypeString, Computed: true, Description: "modified by"},
									"modified_dt":            {Type: schema.TypeString, Computed: true, Description: "modified dt"},
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
						"maintenance_period":            {Type: schema.TypeString, Computed: true, Description: "maintenance period"},
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
		Description: "Search kafka cluster database.",
	}
}

func dataSourceKafkaSingle(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Kafka.DetailKafkaCluster(ctx, rd.Get("kafka_cluster_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(dbInfo.BrokerNodeGroup.BrokerNodes) == 0 {
		return diag.Errorf("no server found")
	}

	rd.SetId(uuid.NewV4().String())
	fs := &common.FieldSetter{Rd: rd}

	fs.Set("kafka_cluster_name", dbInfo.KafkaClusterName)
	fs.Set("kafka_cluster_state", dbInfo.KafkaClusterState)
	fs.Set("image_id", dbInfo.ImageId)
	fs.Set("timezone", dbInfo.Timezone)
	fs.Set("vpc_id", dbInfo.VpcId)
	fs.Set("subnet_id", dbInfo.SubnetId)
	fs.Set("security_group_ids", dbInfo.SecurityGroupIds)
	fs.Set("database_version", dbInfo.DatabaseVersion)
	fs.Set("nat_ip_address", dbInfo.NatIpAddress)

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
	fs.Set("contract", contract)

	kafkaInitialConfig := database_common.HclListObject{}
	if dbInfo.KafkaInitialConfig != nil {
		kafkaInitialConfigInfo := database_common.HclKeyValueObject{}

		brokerInitialConfig := database_common.HclListObject{}
		brokerInitialConfigInfo := database_common.HclKeyValueObject{}
		brokerInitialConfigInfo["broker_sasl_account"] = dbInfo.KafkaInitialConfig.BrokerInitialConfig.BrokerSaslAccount
		brokerInitialConfigInfo["broker_port"] = dbInfo.KafkaInitialConfig.BrokerInitialConfig.BrokerPort
		brokerInitialConfig = append(brokerInitialConfig, brokerInitialConfigInfo)
		kafkaInitialConfigInfo["broker_initial_config"] = brokerInitialConfig

		zookeeperInitialConfig := database_common.HclListObject{}
		if dbInfo.KafkaInitialConfig.ZookeeperInitialConfig != nil {
			zookeeperInitialConfigInfo := database_common.HclKeyValueObject{}
			zookeeperInitialConfigInfo["zookeeper_sasl_account"] = dbInfo.KafkaInitialConfig.ZookeeperInitialConfig.ZookeeperSaslAccount
			zookeeperInitialConfigInfo["zookeeper_port"] = dbInfo.KafkaInitialConfig.ZookeeperInitialConfig.ZookeeperPort
			zookeeperInitialConfig = append(zookeeperInitialConfig, zookeeperInitialConfigInfo)
		}
		kafkaInitialConfigInfo["zookeeper_initial_config"] = zookeeperInitialConfig

		akhqInitialConfig := database_common.HclListObject{}
		if dbInfo.KafkaInitialConfig.AkhqInitialConfig != nil {
			akhqInitialConfigInfo := database_common.HclKeyValueObject{}
			akhqInitialConfigInfo["akhq_account"] = dbInfo.KafkaInitialConfig.AkhqInitialConfig.AkhqAccount
			akhqInitialConfigInfo["akhq_port"] = dbInfo.KafkaInitialConfig.AkhqInitialConfig.AkhqPort
			akhqInitialConfig = append(akhqInitialConfig, akhqInitialConfigInfo)
		}
		kafkaInitialConfigInfo["akhq_initial_config"] = akhqInitialConfig

		kafkaInitialConfig = append(kafkaInitialConfig, kafkaInitialConfigInfo)
	}
	fs.Set("kafka_initial_config", kafkaInitialConfig)

	// BrokerNodeGroup 세팅
	brokerNodeGroup := database_common.HclListObject{}
	if dbInfo.BrokerNodeGroup != nil {
		brokerNodeGroupInfo := database_common.HclKeyValueObject{}
		brokerNodeGroupInfo["server_type"] = dbInfo.BrokerNodeGroup.ServerType
		brokerNodeGroupInfo["node_role_type"] = dbInfo.BrokerNodeGroup.NodeRoleType

		brokerNodes := database_common.HclListObject{}
		for _, value := range dbInfo.BrokerNodeGroup.BrokerNodes {
			brokerNodeInfo := database_common.HclKeyValueObject{}
			brokerNodeInfo["broker_node_id"] = value.BrokerNodeId
			brokerNodeInfo["broker_node_name"] = value.BrokerNodeName
			brokerNodeInfo["broker_node_state"] = value.BrokerNodeState
			brokerNodeInfo["subnet_ip_address"] = value.SubnetIpAddress
			brokerNodeInfo["availability_zone_name"] = value.AvailabilityZoneName
			brokerNodeInfo["created_by"] = value.CreatedBy
			brokerNodeInfo["created_dt"] = value.CreatedDt.String()
			brokerNodeInfo["modified_by"] = value.ModifiedBy
			brokerNodeInfo["modified_dt"] = value.ModifiedDt.String()
			brokerNodes = append(brokerNodes, brokerNodeInfo)
		}
		brokerNodeGroupInfo["broker_nodes"] = brokerNodes

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.BrokerNodeGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize
			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		brokerNodeGroupInfo["block_storages"] = blockStorages

		brokerNodeGroup = append(brokerNodeGroup, brokerNodeGroupInfo)
	}
	fs.Set("broker_node_group", brokerNodeGroup)

	// ZookeeperNodeGroup 세팅
	zookeeperNodeGroup := database_common.HclListObject{}
	if dbInfo.ZookeeperNodeGroup != nil {
		zookeeperNodeGroupInfo := database_common.HclKeyValueObject{}
		zookeeperNodeGroupInfo["server_type"] = dbInfo.ZookeeperNodeGroup.ServerType
		zookeeperNodeGroupInfo["node_role_type"] = dbInfo.ZookeeperNodeGroup.NodeRoleType

		zookeeperNodes := database_common.HclListObject{}
		for _, value := range dbInfo.ZookeeperNodeGroup.ZookeeperNodes {
			zookeeperNodeInfo := database_common.HclKeyValueObject{}
			zookeeperNodeInfo["zookeeper_node_id"] = value.ZookeeperNodeId
			zookeeperNodeInfo["zookeeper_node_name"] = value.ZookeeperNodeName
			zookeeperNodeInfo["zookeeper_node_state"] = value.ZookeeperNodeState
			zookeeperNodeInfo["subnet_ip_address"] = value.SubnetIpAddress
			zookeeperNodeInfo["availability_zone_name"] = value.AvailabilityZoneName
			zookeeperNodeInfo["created_by"] = value.CreatedBy
			zookeeperNodeInfo["created_dt"] = value.CreatedDt.String()
			zookeeperNodeInfo["modified_by"] = value.ModifiedBy
			zookeeperNodeInfo["modified_dt"] = value.ModifiedDt.String()
			zookeeperNodes = append(zookeeperNodes, zookeeperNodeInfo)
		}
		zookeeperNodeGroupInfo["zookeeper_nodes"] = zookeeperNodes

		blockStorages := database_common.HclListObject{}
		for _, value := range dbInfo.ZookeeperNodeGroup.BlockStorages {
			blockStoragesInfo := database_common.HclKeyValueObject{}
			blockStoragesInfo["block_storage_group_id"] = value.BlockStorageGroupId
			blockStoragesInfo["block_storage_name"] = value.BlockStorageName
			blockStoragesInfo["block_storage_role_type"] = value.BlockStorageRoleType
			blockStoragesInfo["block_storage_type"] = value.BlockStorageType
			blockStoragesInfo["block_storage_size"] = value.BlockStorageSize
			blockStorages = append(blockStorages, blockStoragesInfo)
		}
		zookeeperNodeGroupInfo["block_storages"] = blockStorages

		zookeeperNodeGroup = append(zookeeperNodeGroup, zookeeperNodeGroupInfo)
	}
	fs.Set("zookeeper_node_group", zookeeperNodeGroup)

	// AkhqNodeGroup 세팅
	akhqNodeGroup := database_common.HclListObject{}
	if dbInfo.AkhqNodeGroup != nil {
		akhqNodeGroupInfo := database_common.HclKeyValueObject{}
		akhqNodeGroupInfo["akhq_node_id"] = dbInfo.AkhqNodeGroup.AkhqNodeId
		akhqNodeGroupInfo["akhq_node_name"] = dbInfo.AkhqNodeGroup.AkhqNodeName
		akhqNodeGroupInfo["akhq_node_state"] = dbInfo.AkhqNodeGroup.AkhqNodeState
		akhqNodeGroupInfo["block_storage_group_id"] = dbInfo.AkhqNodeGroup.BlockStorageGroupId
		akhqNodeGroupInfo["block_storage_name"] = dbInfo.AkhqNodeGroup.BlockStorageName
		akhqNodeGroupInfo["block_storage_role_type"] = dbInfo.AkhqNodeGroup.BlockStorageRoleType
		akhqNodeGroupInfo["block_storage_size"] = dbInfo.AkhqNodeGroup.BlockStorageSize
		akhqNodeGroupInfo["block_storage_type"] = dbInfo.AkhqNodeGroup.BlockStorageType
		akhqNodeGroupInfo["node_role_type"] = dbInfo.AkhqNodeGroup.NodeRoleType
		akhqNodeGroupInfo["server_type"] = dbInfo.AkhqNodeGroup.ServerType
		akhqNodeGroupInfo["subnet_ip_address"] = dbInfo.AkhqNodeGroup.SubnetIpAddress
		akhqNodeGroupInfo["availability_zone_name"] = dbInfo.AkhqNodeGroup.AvailabilityZoneName
		akhqNodeGroup = append(akhqNodeGroup, akhqNodeGroupInfo)
	}
	fs.Set("akhq_node_group", akhqNodeGroup)

	// Maintenance
	maintenance := database_common.HclListObject{}
	if dbInfo.Maintenance != nil {
		maintenanceInfo := database_common.HclKeyValueObject{}
		maintenanceInfo["maintenance_start_day_of_week"] = dbInfo.Maintenance.MaintenanceStartDayOfWeek
		maintenanceInfo["maintenance_start_time"] = dbInfo.Maintenance.MaintenanceStartTime
		maintenanceInfo["maintenance_period"] = dbInfo.Maintenance.MaintenancePeriod
		maintenance = append(maintenance, maintenanceInfo)
	}
	fs.Set("maintenance", maintenance)

	fs.Set("project_id", dbInfo.ProjectId)
	fs.Set("service_zone_id", dbInfo.ServiceZoneId)
	fs.Set("block_id", dbInfo.BlockId)
	fs.Set("created_by", dbInfo.CreatedBy)
	fs.Set("created_dt", dbInfo.CreatedDt.String())
	fs.Set("modified_by", dbInfo.ModifiedBy)
	fs.Set("modified_dt", dbInfo.ModifiedDt.String())

	return fs.Diags()
}
