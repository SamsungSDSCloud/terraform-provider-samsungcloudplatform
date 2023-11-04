package sqlserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/image"
	objectstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/object-storage"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver2"
	"github.com/antihax/optional"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_sqlserver", ResourceSqlServer())
}

func ResourceSqlServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlServerCreate,
		ReadContext:   resourceSqlServerRead,
		UpdateContext: resourceSqlServerUpdate,
		DeleteContext: resourceSqlServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SQL Server standard image id.",
			},
			"virtual_server_name_prefix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Prefix of virtual server. (3 to 13 alpha-numerics with dash)",
				ValidateDiagFunc: common.ValidateName3to13AlphaNumberDash,
			},
			"server_group_name": {
				// cluster
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
			},
			"db_service_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of SQL server database service. (Starts with a capital letter, 1 to 15 alphabet only)",
				ValidateDiagFunc: common.ValidateName1to15AlphaOnlyStartsWithCapitalLetter,
			},
			"db_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database.",
				ValidateDiagFunc: nil, // TODO
			},
			"db_user_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "User account id. (2 to 20 alpha-numerics)",
				ValidateDiagFunc: common.ValidateName2to20AlphaNumeric,
			},
			"db_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "User account password",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"db_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port number of database.",
			},
			"cpu_count": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "CPU core count (2, 4, 8,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"memory_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Memory size in gigabytes(4, 8, 16,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"contract_discount": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Contract : None, 1-year, 3-year",
			},
			"data_disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data storage disk type. (SSD, HDD)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					val := v.(string)
					if val == "SSD" {
						return diags
					} else if val == "HDD" {
						return diags
					}

					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       fmt.Sprintf("Must be either SSD or HDD"),
						AttributePath: path,
					})
					return diags
				},
			},
			"data_block_storage_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Data Block Storage size in gigabytes.",
				ValidateDiagFunc: common.ValidateBlockStorageSize,
			},
			"encrypt_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"additional_block_storages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Additional block storages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage product name. (only SSD)",
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								var diags diag.Diagnostics
								val := v.(string)
								if val == "SSD" {
									return diags
								}

								diags = append(diags, diag.Diagnostic{
									Severity:      diag.Error,
									Summary:       fmt.Sprintf("Must be SSD"),
									AttributePath: path,
								})
								return diags
							},
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage usage. (DATA, ARCHIVE)",
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								var diags diag.Diagnostics
								val := v.(string)
								if val == "DATA" {
									return diags
								} else if val == "ARCHIVE" {
									return diags
								}

								diags = append(diags, diag.Diagnostic{
									Severity:      diag.Error,
									Summary:       fmt.Sprintf("Must be either DATA or ARCHIVE"),
									AttributePath: path,
								})
								return diags
							},
						},
						"storage_size_gb": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Default data storage size in gigabytes. (10~7,168 GB)",
							ValidateDiagFunc: common.ValidateBlockStorageSize,
						},
					},
				},
			},
			"backup": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceSqlServerBackup(),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID.",
			},
			"license_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "License key.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Timezone setting of this database.",
			},
			//"nat_enabled":{},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this sql server group. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"db_collation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Commands that specify how to sort and compare data",
			},
			"additional_db": {
				Type:     schema.TypeList,
				Required: true, // why?
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Names of additional database.",
			},
			"high_availability": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem:     resourceSqlServerHighAvailability(),
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual IP.",
			},
			"external_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External virtual IP.",
			},
			"cluster_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster virtual IP.",
			},
		},
		Description: "Provide Microsoft SQL Server resource.",
	}
}

func resourceSqlServerBackup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_storage_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage ID where backup files will be stored.",
			},
			"backup_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup Method (s3api|cdp) ",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					val := v.(string)
					if val == "s3api" || val == "cdp" {
						return diags
					}

					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       fmt.Sprintf("Must be either s3api or cdp"),
						AttributePath: path,
					})
					return diags
				},
			},
			"backup_retention_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Backup File Retention Day.(7 <= day <= 35) ",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(int)
					if value < 7 || value > 35 {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("Backup retion day's value must be between 7 and 35 days"),
							AttributePath: path,
						})
					}
					return diags
				},
			},
			"backup_start_hour": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time at which the backup starts. (from 0 to 23)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(int)
					if value < 0 || value > 23 {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("Backup start hour can be set from 0 to 23."),
							AttributePath: path,
						})
					}
					return diags
				},
			},
		},
	}
}

func resourceSqlServerHighAvailability() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_vip_nat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "use Virtual IP NAT.",
			},
			"reserved_nat_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of Reserved Virtual NAT IP.",
			},
			"virtual_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Virtual IP for database cluster access.",
				ValidateDiagFunc: common.ValidateIpv4,
			},
			"active_server_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Static IP to assign to the ACTIVE server.",
				ValidateDiagFunc: common.ValidateIpv4,
			},
			"standby_server_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Static IP to assign to the STANDBy server.",
				ValidateDiagFunc: common.ValidateIpv4,
			},
			"active_availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Active Availability Zone Name",
			},
			"standby_availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Standby Availability Zone Name",
			},
		},
	}
}

func resourceSqlServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	vpcId := rd.Get("vpc_id").(string)
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}
	//blockId := vpcInfo.BlockId	// vpcInfo.BlockId is bogus

	subnetId := rd.Get("subnet_id").(string)
	subnetInfo, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return
	}

	serviceZoneId := vpcInfo.ServiceZoneId

	// block id, AZ
	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return
	}
	var blockId string
	var isMultiAvailabilityZone *bool
	for _, zoneInfo := range projectInfo.ServiceZones {
		if zoneInfo.ServiceZoneId == serviceZoneId {
			blockId = zoneInfo.BlockId
			isMultiAvailabilityZone = zoneInfo.IsMultiAvailabilityZone
			break
		}
	}

	serverNamePrefix := rd.Get("virtual_server_name_prefix").(string)
	networks := expandNetworkSetting(serverNamePrefix, rd.Get("high_availability").(*schema.Set), isMultiAvailabilityZone)
	ha := expandHASetting(rd.Get("high_availability").(*schema.Set))

	backup, err := createBackupSettings(ctx, inst.Client, serviceZoneId, rd.Get("backup").(*schema.Set))
	if err != nil {
		return
	}

	dbPort := rd.Get("db_port").(int)
	dbUserPassword := base64.StdEncoding.EncodeToString([]byte(rd.Get("db_user_password").(string)))
	licenseKey := base64.StdEncoding.EncodeToString([]byte(rd.Get("license_key").(string)))
	//licenseKey := rd.Get("license_key").(string)

	encryptEnabled := rd.Get("encrypt_enabled").(bool)

	dataDiskType := rd.Get("data_disk_type").(string)

	// block storage
	var blockStorages []sqlserver2.DatabaseBlockStorage
	additionalStorageList := rd.Get("additional_block_storages").(common.HclListObject)
	additionalStorageInfoList := common.ConvertAdditionalStorageList(additionalStorageList)
	for _, additionalStorageInfo := range additionalStorageInfoList {
		blockStorages = append(blockStorages, sqlserver2.DatabaseBlockStorage{
			BlockStorageType: additionalStorageInfo.StorageUsage,
			BlockStorageSize: int32(additionalStorageInfo.StorageSize),
			DiskType:         additionalStorageInfo.ProductName,
		})
	}
	blockStorageSize := rd.Get("data_block_storage_size_gb").(int)

	// additional db
	dbSWOption := make(map[string]string)
	dbSWOption["sqlserver_collation"] = rd.Get("db_collation").(string)
	iAdditionalDbs := rd.Get("additional_db").([]interface{})
	additionalDbs := make([]string, len(iAdditionalDbs))
	for i, aDbName := range iAdditionalDbs {
		additionalDbs[i] = aDbName.(string)
	}

	// product group
	imageId := rd.Get("image_id").(string)
	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, serviceZoneId, image.ActiveState, common.ServicedGroupDatabase, common.ServicedForSqlServer)
	if err != nil {
		return
	}

	var targetProductGroupId string
	for _, c := range standardImages.Contents {
		if c.ImageId == imageId {
			targetProductGroupId = c.ProductGroupId
		}
	}

	productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
	if err != nil {
		return
	}

	// scale
	scaleProductId, err := client.FindScaleProduct(ctx, inst.Client, targetProductGroupId, rd.Get("cpu_count").(int), rd.Get("memory_size_gb").(int))

	// contract
	contractDiscount := rd.Get("contract_discount").(string)
	contractId, err := common.FindProductId(common.ContractProductType, contractDiscount, &productGroup)
	if err != nil {
		return
	}

	valueFalse := false
	dbName := rd.Get("db_name").(string)
	serverGroupName := rd.Get("server_group_name").(string)

	_, _, err = inst.Client.SqlServer.CreateSqlServer(ctx, sqlserver2.CreateSqlServerRequest{
		ServiceZoneId:           serviceZoneId,
		BlockId:                 blockId,
		ImageId:                 imageId,
		ProductGroupId:          targetProductGroupId,
		VirtualServerNamePrefix: serverNamePrefix,
		ServerGroupName:         serverGroupName,
		DbServiceName:           rd.Get("db_service_name").(string),
		DbName:                  dbName,
		DbUserId:                rd.Get("db_user_id").(string),
		DbUserPassword:          dbUserPassword,
		DbPort:                  int32(dbPort),
		DeploymentEnvType:       "DEV",
		VirtualServer: &sqlserver2.InstanceSpec{
			ScaleProductId:            scaleProductId,
			ContractDiscountProductId: contractId,
			DataBlockStorageSize:      int32(blockStorageSize),
			EncryptEnabled:            &encryptEnabled,
			AdditionalBlockStorages:   blockStorages,
			DataDiskType:              dataDiskType,
		},
		HighAvailability: ha,
		Replica:          nil,
		Network: &sqlserver2.DatabaseNetwork{
			NetworkEnvType: strings.ToUpper(subnetInfo.SubnetType),
			VpcId:          vpcId,
			SubnetId:       subnetId,
			UseNat:         &valueFalse, //todo
			ServerNetworks: networks,
		},
		SecurityGroupIds:  getSecurityGroupIds(rd),
		Maintenance:       nil,
		Backup:            backup,
		UseDbLoggingAudit: &valueFalse,
		LicenseKey:        licenseKey,
		Timezone:          rd.Get("timezone").(string),
		DbSoftwareOptions: dbSWOption,
		AdditionalDb:      additionalDbs,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty

	resultList, _, err := inst.Client.SqlServer.ListSqlServer(ctx, &sqlserver2.MsSqlConfigurationControllerApiListSqlserverOpts{
		DbName:            optional.NewString(dbName),
		ServerGroupName:   optional.NewString(serverGroupName),
		VirtualServerName: optional.NewString(serverNamePrefix),
	})
	if err != nil {
		return
	}
	if len(resultList.Contents) == 0 {
		diagnostics = diag.Errorf("no pending create found")
		return
	}

	dbServerGroupId := resultList.Contents[0].ServerGroupId

	if len(dbServerGroupId) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	err = waitForSqlServer(ctx, inst.Client, dbServerGroupId, common.DatabaseProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return
	}

	rd.SetId(dbServerGroupId)

	return resourceSqlServerRead(ctx, rd, meta)
}

func resourceSqlServerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	sqlServer, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if len(sqlServer.VirtualServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	port64, err := strconv.ParseInt(sqlServer.DbPort, 10, 32)
	if err != nil {
		return
	}

	vsInfo := sqlServer.VirtualServers[0]
	serverNamePrefix := vsInfo.VirtualServerName[:len(vsInfo.VirtualServerName)-2]
	if len(serverNamePrefix) == 0 {
		diagnostics = diag.Errorf("server name prefix is invalid")
		return
	}

	contractDiscountName := vsInfo.ContractDiscountName
	if len(contractDiscountName) == 0 {
		diagnostics = diag.Errorf("contract discount information not found")
		return
	}

	cpuCount, memorySize, err := client.FindScaleInfo(ctx, inst.Client, sqlServer.ProductGroupId, vsInfo.ScaleProductId)
	if err != nil {
		return
	}

	// additional db 가 db_name에 합쳐져서 나오고 있기 때문에 parsing해야함 -> API가 수정되는게 맞을 것 같지만...
	var dbName string
	var additionalDb []string
	if strings.Contains(sqlServer.DbName, ",") {
		dbNames := strings.Split(sqlServer.DbName, ",")
		for i, dn := range dbNames {
			if i == 0 {
				dbName = dn
			} else {
				additionalDb = append(additionalDb, dn)
			}
		}
	} else {
		dbName = sqlServer.DbName
	}
	//rd.Set("service_zone_id", sqlServer.ServiceZoneId)
	rd.Set("image_id", sqlServer.ImageId)
	rd.Set("virtual_server_name_prefix", serverNamePrefix)
	rd.Set("server_group_name", sqlServer.ServerGroupName)
	rd.Set("db_service_name", sqlServer.SqlServerServiceName)
	rd.Set("db_name", dbName)
	rd.Set("db_user_id", sqlServer.DbUserId)
	//rd.Set("db_user_password",...)
	rd.Set("db_port", int(port64))
	rd.Set("cpu_count", cpuCount)
	rd.Set("memory_size_gb", memorySize)
	rd.Set("contract_discount", vsInfo.ContractDiscountName)
	//rd.Set("encrypt_enabled", ...)
	rd.Set("vpc_id", vsInfo.VpcId)
	rd.Set("subnet_id", vsInfo.Network.NetworkId)
	//rd.Set("license_key", ...)
	rd.Set("timezone", sqlServer.Timezone)
	rd.Set("vip", sqlServer.Vip)
	rd.Set("external_vip", sqlServer.ExternalVip)
	rd.Set("cluster_vip", sqlServer.ClusterVip)
	rd.Set("additional_db", additionalDb)

	// security group
	var securityGroupIds []string
	for _, sg := range sqlServer.SecurityGroups {
		securityGroupIds = append(securityGroupIds, sg.SecurityGroupId)
	}
	rd.Set("security_group_id", securityGroupIds)

	// additional block storages
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, sqlServer.ProductGroupId)
	if err != nil {
		return
	}
	var storageProductName string
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState == common.ProductAvailableState {
				storageProductName = productInfo.ProductName
				break
			}
		}
	}

	additionalStorages := common.HclListObject{}
	for i, bs := range vsInfo.BlockStorages {
		// Skip OS Storage (i == 0)
		if i == 0 {
			continue
		}
		// First data storage is default storage
		if i == 1 {
			rd.Set("data_block_storage_size_gb", bs.BlockStorageSize)
			continue
		}

		// Additional storages
		storageInfo := common.HclKeyValueObject{}
		storageInfo["id"] = bs.BlockStorageId
		storageInfo["product_name"] = storageProductName
		storageInfo["storage_size_gb"] = bs.BlockStorageSize
		storageInfo["storage_usage"] = bs.BlockStorageType

		additionalStorages = append(additionalStorages, storageInfo)
	}
	rd.Set("additional_storage", additionalStorages)

	if sqlServer.Backup != nil {
		backup := map[string]interface{}{
			"backup_method":        sqlServer.Backup.BackupMethod,
			"object_storage_id":    sqlServer.Backup.ObjectStorageId,
			"backup_retention_day": int(sqlServer.Backup.BackupRetentionDay),
			"backup_start_hour":    int(sqlServer.Backup.BackupStartHour),
		}

		backupSchema := resourceSqlServerBackup()
		rd.Set("backup", schema.NewSet(schema.HashResource(backupSchema), []interface{}{backup}))
	}

	return nil
}

type UpdateSqlServerParam struct {
	Ctx       context.Context
	Rd        *schema.ResourceData
	Inst      *client.Instance
	SqlServer *sqlserver2.DetailDatabaseResponse
}

func resourceSqlServerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	sqlServer, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
	if err != nil {
		return
	}
	if len(sqlServer.VirtualServers) == 0 {
		diagnostics = diag.Errorf("database id not found")
		return
	}

	param := UpdateSqlServerParam{
		Ctx:       ctx,
		Rd:        rd,
		Inst:      inst,
		SqlServer: &sqlServer,
	}

	var updateFuncs []func(serverParam UpdateSqlServerParam) error

	if rd.HasChanges("cpu_count", "memory_size_gb") {
		updateFuncs = append(updateFuncs, updateScale)
	}
	if rd.HasChanges("data_block_storage_size_gb") {
		updateFuncs = append(updateFuncs, updateDataStorageSize)
	}
	if rd.HasChanges("additional_block_storages") {
		updateFuncs = append(updateFuncs, updateAdditionalStorage)
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

	return resourceSqlServerRead(ctx, rd, meta)
}

func updateScale(param UpdateSqlServerParam) error {
	scaleProductId, err := client.FindScaleProduct(param.Ctx, param.Inst.Client, param.SqlServer.ProductGroupId, param.Rd.Get("cpu_count").(int), param.Rd.Get("memory_size_gb").(int))
	if err != nil {
		return err
	}

	if len(scaleProductId) == 0 {
		return fmt.Errorf("no server type found")
	}

	// update all
	for _, vsInfo := range param.SqlServer.VirtualServers {
		_, _, err = param.Inst.Client.SqlServer.UpdateSqlServerScale(param.Ctx, param.Rd.Id(), vsInfo.VirtualServerId, scaleProductId)

		err = waitForSqlServer(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateBackup(param UpdateSqlServerParam) error {
	o, n := param.Rd.GetChange("backup")

	oldBackups := expandBackupSettings(o.(*schema.Set))
	newBackups := expandBackupSettings(n.(*schema.Set))

	var backup *sqlserver2.DatabaseBackup
	var useBackup bool
	if len(newBackups) > 0 {
		useBackup = true
		backup = newBackups[0]
		if len(backup.ObjectStorageId) < 1 {
			objectStorageId, err := getObjectStorageId(param.Ctx, param.Inst.Client, param.SqlServer.ServiceZoneId)
			if err != nil {
				return err
			}
			backup.ObjectStorageId = objectStorageId
		}
	} else {
		useBackup = false
		backup = oldBackups[0]
	}
	updateBackupSettingRequest := sqlserver2.UpdateBackupSettingRequest{
		UseBackup: &useBackup,
		Backup:    backup,
	}
	if _, _, err := param.Inst.Client.SqlServer.UpdateBackupSetting(param.Ctx, param.Rd.Id(), updateBackupSettingRequest); err != nil {
		return err
	}
	if err := waitForSqlServer(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true); err != nil {
		return err
	}

	return nil
}

// Increase only
func updateAdditionalStorage(param UpdateSqlServerParam) error {
	o, n := param.Rd.GetChange("additional_block_storages")
	oldVal := o.(common.HclListObject)
	newVal := n.(common.HclListObject)

	oldList := common.ConvertAdditionalStorageList(oldVal)
	newList := common.ConvertAdditionalStorageList(newVal)

	if len(oldList) > len(newList) {
		return fmt.Errorf("removing additional storage is not allowed")
	}

	incIndices := make(map[int]int)
	for i := 0; i < len(oldList); i++ {
		if oldList[i].StorageUsage != newList[i].StorageUsage {
			return fmt.Errorf("changing storage usage is not allowed")
		}
		if oldList[i].ProductName != newList[i].ProductName {
			return fmt.Errorf("changing product name is not allowed")
		}
		if oldList[i].StorageSize > newList[i].StorageSize {
			return fmt.Errorf("decreasing size is not allowed")
		}
		incIndices[i] = newList[i].StorageSize
	}

	// Update all
	for _, vsInfo := range param.SqlServer.VirtualServers {
		for i, blockInfo := range vsInfo.BlockStorages {
			// Skip
			if i < 2 {
				continue
			}
			if newSize, ok := incIndices[i-2]; ok {
				if newSize != int(blockInfo.BlockStorageSize) {
					_, _, err := param.Inst.Client.SqlServer.UpdateSqlServerBlockSize(param.Ctx, param.Rd.Id(), vsInfo.VirtualServerId, blockInfo.BlockStorageId, newSize)
					if err != nil {
						return err
					}

					err = waitForSqlServer(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
					if err != nil {
						return err
					}
				} else {
					// No size change
				}
			}
		}
	}
	return nil
}

// Increase only
func updateDataStorageSize(param UpdateSqlServerParam) error {
	o, n := param.Rd.GetChange("data_block_storage_size_gb")
	oldSize := o.(int)
	newSize := n.(int)
	if oldSize >= newSize {
		return fmt.Errorf("storage size can only be increased. check size value : %d -> %d", oldSize, newSize)
	}

	// Update all
	for _, vs := range param.SqlServer.VirtualServers {
		var dataBlockId string
		for i, bs := range vs.BlockStorages {
			if i == 1 && bs.BlockStorageType == common.BlockStorageTypeData {
				dataBlockId = bs.BlockStorageId
				break
			}
		}
		if len(dataBlockId) < 1 {
			return fmt.Errorf("default data storage not found")
		}

		_, _, err := param.Inst.Client.SqlServer.UpdateSqlServerBlockSize(param.Ctx, param.Rd.Id(), vs.VirtualServerId, dataBlockId, newSize)
		if err != nil {
			return err
		}

		err = waitForSqlServer(param.Ctx, param.Inst.Client, param.Rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSqlServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.SqlServer.DeleteSqlServer(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForSqlServer(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSqlServer(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	baseInfo, baseStatusCode, baseErr := scpClient.SqlServer.GetSqlServer(ctx, id)
	if baseErr != nil {
		if baseStatusCode == 404 && !errorOnNotFound {
			return nil
		}
		if baseStatusCode == 403 && !errorOnNotFound {
			return nil
		}
		if baseStatusCode >= 500 && !errorOnNotFound {
			return nil
		}
		return baseErr
	}

	if len(baseInfo.VirtualServers) == 0 {
		return fmt.Errorf("no virtual server found")
	}

	// Check all virtual server software status
	for i := 0; i < len(baseInfo.VirtualServers); i++ {
		baseErr = client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
			info, c, err := scpClient.SqlServer.GetSqlServer(ctx, id)
			if err != nil {
				//virtual server not found
				if c == 400 && !errorOnNotFound {
					return "", common.DeletedState, nil
				}

				if c == 404 && !errorOnNotFound {
					return "", common.DeletedState, nil
				}
				/*
					if c == 403 && !errorOnNotFound {
						return "", common.DeletedState, nil
					}
					if c >= 500 && !errorOnNotFound {
						return "", common.DeletedState, nil
					}
				*/
				return nil, "", err
			}
			if i >= len(info.VirtualServers) {
				return nil, "", fmt.Errorf("invalid number of virtual servers")
			}

			vsSoftware := info.VirtualServers[i].Software
			if vsSoftware == nil {
				return nil, "", fmt.Errorf("virtual server software status not found")
			}
			return info, vsSoftware.SoftwareServiceState, nil
		})
		if baseErr != nil {
			return baseErr
		}
	}

	return nil
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getSecurityGroupIds(rd *schema.ResourceData) []string {
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	sgIds := make([]string, len(securityGroupIds))
	for i, valueIpv4 := range securityGroupIds {
		sgIds[i] = valueIpv4.(string)
	}
	return sgIds
}

func expandNetworkSetting(serverNamePrefix string, haSchemaSet *schema.Set, isMultiAvailabilityZone *bool) []sqlserver2.DatabaseServerNetwork {
	var networks []sqlserver2.DatabaseServerNetwork

	haSettings := haSchemaSet.List()

	if len(haSettings) > 0 {

		setting := haSettings[0].(map[string]interface{})

		network := sqlserver2.DatabaseServerNetwork{
			ServerName: serverNamePrefix + "01",
			NodeType:   "ACTIVE",
		}
		if v, ok := setting["active_server_ip"].(string); ok && v != "" {
			network.ServiceIp = v
		}

		if *isMultiAvailabilityZone {
			if v, ok := setting["active_availability_zone_name"].(string); ok && v != "" {
				network.AvailabilityZoneName = v
			}
		}

		networks = append(networks, network)

		network = sqlserver2.DatabaseServerNetwork{
			ServerName: serverNamePrefix + "02",
			NodeType:   "STANDBY",
		}
		if v, ok := setting["standby_server_ip"].(string); ok && v != "" {
			network.ServiceIp = v
		}

		if *isMultiAvailabilityZone {
			if v, ok := setting["standby_availability_zone_name"].(string); ok && v != "" {
				network.AvailabilityZoneName = v
			}
		}

		networks = append(networks, network)

	} else {
		if *isMultiAvailabilityZone {
			networks = append(networks, sqlserver2.DatabaseServerNetwork{
				ServerName:           serverNamePrefix + "01",
				NodeType:             "ACTIVE",
				ServiceIp:            "",
				AvailabilityZoneName: "AZ1",
				NatIpId:              "",
			})
		} else {
			networks = append(networks, sqlserver2.DatabaseServerNetwork{
				ServerName:           serverNamePrefix + "01",
				NodeType:             "ACTIVE",
				ServiceIp:            "",
				AvailabilityZoneName: "",
				NatIpId:              "",
			})
		}
	}

	return networks
}

func expandHASetting(haSchemaSet *schema.Set) *sqlserver2.HaNormal {
	var haSettings []*sqlserver2.HaNormal

	for _, haSet := range haSchemaSet.List() {

		setting := haSet.(map[string]interface{})

		haSetting := &sqlserver2.HaNormal{}

		if v, ok := setting["use_vip_nat"].(bool); ok {
			haSetting.UseVipNat = &v
		}

		if v, ok := setting["reserved_nat_ip_id"].(string); ok && v != "" {
			haSetting.ReservedNatIpId = v
		}

		if v, ok := setting["virtual_ip"].(string); ok && v != "" {
			haSetting.VirtualIp = v
		}
		haSettings = append(haSettings, haSetting)
	}

	var ha *sqlserver2.HaNormal
	if len(haSettings) > 0 {
		ha = haSettings[0]
	} else {
		ha = nil
	}

	return ha
}

func expandBackupSettings(vAdvancedBackupSettings *schema.Set) []*sqlserver2.DatabaseBackup {
	backupSettings := []*sqlserver2.DatabaseBackup{}

	for _, vAdvancedBackupSetting := range vAdvancedBackupSettings.List() {
		backupSetting := &sqlserver2.DatabaseBackup{}

		mAdvancedBackupSetting := vAdvancedBackupSetting.(map[string]interface{})

		if v, ok := mAdvancedBackupSetting["objectstorage_id"].(string); ok && v != "" {
			backupSetting.ObjectStorageId = v
		}

		if v, ok := mAdvancedBackupSetting["backup_method"].(string); ok && v != "" {
			backupSetting.BackupMethod = v
		}

		if v, ok := mAdvancedBackupSetting["backup_retention_day"].(int); ok {
			backupSetting.BackupRetentionDay = int32(v)
		}

		if v, ok := mAdvancedBackupSetting["backup_start_hour"].(int); ok {
			backupSetting.BackupStartHour = int32(v)
		}
		backupSettings = append(backupSettings, backupSetting)
	}

	return backupSettings
}

func getObjectStorageId(ctx context.Context, scpClient *client.SCPClient, serviceZoneId string) (string, error) {

	response, err := scpClient.ObjectStorage.ReadObjectStorageList(ctx, serviceZoneId, objectstorage.ObjectStorageV4ControllerApiListObjectStorage6Opts{})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	if response.TotalCount < 1 {
		return "", fmt.Errorf("No object storage list.")
	}
	objectStorageId := response.Contents[0].ObjectStorageId // idx 0 is OK?

	return objectStorageId, nil
}

func createBackupSettings(ctx context.Context, scpClient *client.SCPClient, serviceZoneId string, backupSettingsSet *schema.Set) (*sqlserver2.DatabaseBackup, error) {
	backupSetting := &sqlserver2.DatabaseBackup{}

	backupList := backupSettingsSet.List()
	if len(backupList) < 1 {
		return nil, nil
	}

	objectStorageId, err := getObjectStorageId(ctx, scpClient, serviceZoneId)
	if err != nil {
		return nil, err
	}

	backup := backupList[0].(map[string]interface{})

	backupSetting.ObjectStorageId = objectStorageId

	if v, ok := backup["backup_method"].(string); ok {
		backupSetting.BackupMethod = v
	}
	if v, ok := backup["backup_retention_day"].(int); ok {
		backupSetting.BackupRetentionDay = int32(v)
	}
	if v, ok := backup["backup_start_hour"].(int); ok {
		backupSetting.BackupStartHour = int32(v)
	}

	backupSetting.DbBackupArchMin = 60
	return backupSetting, nil
}
