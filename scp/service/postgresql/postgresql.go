package postgresql

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/service/image"
	objectstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/object-storage"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/postgresql2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_postgresql", ResourcePostgresql())
}

/*const (
	BlockStorageTypeOS      string = "OS"
	BlockStorageTypeData    string = "DATA"
	BlockStorageTypeArchive string = "ARCHIVE"
)*/

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
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Postgresql virtual server image id.",
			},
			"server_name_prefix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Prefix of database server names. (3 to 13 alpha-numerics with dash)",
				ValidateDiagFunc: common.ValidateName3to13AlphaNumberDash,
			},
			"cluster_name": {
				// Server-Group name
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database cluster. (3 to 20 characters only)",
				ValidateDiagFunc: common.ValidateName3to20AlphaOnly,
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
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id of this database server.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id of this database server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this postgresql DB. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"db_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database. (3 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName3to20DashInMiddle,
			},
			"db_user_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "User account id of database. (2 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName2to20LowerAlphaOnly,
			},
			"db_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ForceNew:         true,
				Description:      "User account password of database.",
				ValidateDiagFunc: common.ValidatePassword8to30WithSpecialsExceptQuotes,
			},
			"db_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port number of this database.",
			},
			"pg_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "UTF8",
				ForceNew:    true,
				Description: "Postgresql encoding. (Only 'UTF8' for now)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if value != "UTF8" {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("value must be UTF8"),
							AttributePath: path,
						})
					}
					return diags
				},
			},
			"pg_locale": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "C",
				ForceNew:    true,
				Description: "Postgresql locale. (Only 'C' for now)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if value != "C" {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("value must be UTF8"),
							AttributePath: path,
						})
					}
					return diags
				},
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Timezone setting of this database.",
			},
			"data_storage_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Default data storage size in gigabytes. (At least 10 GB required)",
				ValidateDiagFunc: common.ValidateBlockStorageSize,
			},
			"additional_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "External block storage. (Only adding is allowed)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//"id": {
						//	Type:        schema.TypeString,
						//	Computed:    true,
						//	Description: "Block storage Id",
						//},
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
							Description:      "Default data storage size in gigabytes. (At least 10 GB required)",
							ValidateDiagFunc: common.ValidateBlockStorageSize,
						},
					},
				},
			},
			"backup": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     resourcePostgresqlBackup(),
			},

			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "private database endpoint",
			},
			"external_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public database endpoint",
			},
			"high_availability": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem:     resourcePostgresqlHighAvailability(),
			},
		},
		Description: "Provides a PostgreSQL Database resource.",
	}
}
func resourcePostgresqlBackup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"objectstorage_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage ID where backup files will be stored",
			},
			"retention_day": {
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
			"start_hour": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time at which the backup starts. (must be between 0 and 23)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(int)
					if value < 0 || value > 23 {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("Start hour's value must be between 0 and 23"),
							AttributePath: path,
						})
					}
					return diags
				},
			},
		},
	}
}

func resourcePostgresqlHighAvailability() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"single_auto_restart": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Storage product name. (only SSD)",
			},
			"use_vip_nat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "use Virtual IP NAT",
			},
			"reserved_natip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of Reserved Virtual NAT IP",
			},
			"virtual_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual IP for database cluster access",
			},
			"active_server_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Static IP to assign to the ACTIVE server",
			},
			"standby_server_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Static IP to assign to the STANDBy server",
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

	imageId := rd.Get("image_id").(string)

	serverNamePrefix := rd.Get("server_name_prefix").(string)
	clusterName := rd.Get("cluster_name").(string)

	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	contractDiscount := rd.Get("contract_discount").(string)

	vpcId := rd.Get("vpc_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	//sgId := rd.Get("security_group_id").([]interface{})

	dbName := rd.Get("db_name").(string)
	dbUserId := rd.Get("db_user_id").(string)
	dbUserPassword := base64.StdEncoding.EncodeToString([]byte(rd.Get("db_user_password").(string)))
	dbPort := rd.Get("db_port").(int)

	useLoggingAudit := false // := rd.Get("use_logging_audit").(bool)

	timezone := rd.Get("timezone").(string)

	dataStorageSize := rd.Get("data_storage_size_gb").(int)
	additionalStorageList := rd.Get("additional_storage").(common.HclListObject)

	// Check cluster name and server name duplication (server name, cluster name 중복 체크 open api 없음)
	diagnostics = checkNameDuplication(ctx, meta, dbName, clusterName, serverNamePrefix)
	if diagnostics != nil {
		return
	}

	// Get vpc info
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}

	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, vpcInfo.ServiceZoneId, image.ActiveState, common.ServicedGroupDatabase, common.ServicedForPostgresql)
	if err != nil {
		return
	}

	var targetProductGroupId string
	for _, c := range standardImages.Contents {
		if c.ImageId == imageId {
			targetProductGroupId = c.ProductGroupId
		}
	}

	scaleId, err := client.FindScaleProduct(ctx, inst.Client, targetProductGroupId, cpuCount, memorySizeGB)
	if err != nil {
		return
	}

	if len(scaleId) == 0 {
		return diag.Errorf("no server type found")
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
	if err != nil {
		return
	}

	// Find contract
	contractId, err := common.FindProductId(common.ContractProductType, contractDiscount, &productGroup)
	if err != nil {
		return
	}

	// Additional disks
	var blockStorages []postgresql2.DatabaseBlockStorage
	additionalStorageInfoList := common.ConvertAdditionalStorageList(additionalStorageList)
	for _, additionalStorageInfo := range additionalStorageInfoList {
		//if len(additionalStorageInfo.Id) != 0 {
		//	diagnostics = diag.Errorf("additional storage 'id' must be empty at creation.")
		//	return
		//}
		blockStorages = append(blockStorages, postgresql2.DatabaseBlockStorage{
			BlockStorageType: additionalStorageInfo.StorageUsage,
			BlockStorageSize: int32(additionalStorageInfo.StorageSize),
		})
	}

	pgSWOption := make(map[string]string)
	pgSWOption["pg_encoding"] = rd.Get("pg_encoding").(string)
	pgSWOption["pg_locale"] = rd.Get("pg_locale").(string)

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		diagnostics = diag.FromErr(err)
		return
	}
	var blockId string
	for _, zoneInfo := range projectInfo.ServiceZones {
		if zoneInfo.ServiceZoneId == vpcInfo.ServiceZoneId {
			blockId = zoneInfo.BlockId
			break
		}
	}
	if len(blockId) == 0 {
		diagnostics = diag.Errorf("current service block not found")
		return
	}

	subnetInfo, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return
	}

	// KJ HA Test
	networks := expandNetworkSetting(serverNamePrefix, rd.Get("high_availability").(*schema.Set))
	ha := expandHASetting(rd.Get("high_availability").(*schema.Set))

	// backup : 23.02.21 kj
	backup, err := createBackupSettings(ctx, inst.Client, vpcInfo.ServiceZoneId, rd.Get("backup").(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	falseVal := false
	_, _, err = inst.Client.Postgresql.CreatePostgresql(ctx, postgresql2.CreateRdbRequest{
		ServiceZoneId:           vpcInfo.ServiceZoneId,
		BlockId:                 blockId,
		ImageId:                 imageId,
		ProductGroupId:          targetProductGroupId,
		VirtualServerNamePrefix: serverNamePrefix,
		ServerGroupName:         clusterName,
		DbName:                  dbName,
		DbUserId:                dbUserId,
		DbUserPassword:          dbUserPassword,
		DbPort:                  int32(dbPort),
		DeploymentEnvType:       "DEV", // TODO: Check if changed
		VirtualServer: &postgresql2.InstanceSpecWithAddtionalDisks{
			ScaleProductId:            scaleId,
			ContractDiscountProductId: contractId,
			DataBlockStorageSize:      int32(dataStorageSize),
			EncryptEnabled:            &falseVal,
			AdditionalBlockStorages:   blockStorages,
		},
		HighAvailability: ha,
		Replica:          nil,
		Network: &postgresql2.DatabaseNetwork{
			NetworkEnvType: strings.ToUpper(subnetInfo.SubnetType),
			VpcId:          vpcId,
			SubnetId:       subnetId,
			//AutoServiceIp:  false,	// deprecated
			UseNat: &falseVal,
			//NatIp:          "",		// deprecated
			//NatIpId:        "",		// deprecated
			ServerNetworks: networks,
		},
		SecurityGroupIds:  getSecurityGroupIds(rd),
		Maintenance:       nil,
		Backup:            backup,
		UseDbLoggingAudit: &useLoggingAudit,
		Timezone:          timezone,
		DbSoftwareOptions: pgSWOption,
	})
	if err != nil {
		return
	}

	time.Sleep(50 * time.Second)

	// NOTE : response.ResourceId is empty
	resultList, err := inst.Client.Postgresql.ListPostgresql(ctx, dbName, clusterName, serverNamePrefix)
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

	if err := waitForPostgresql2(ctx, inst.Client, rd, dbServerGroupId, common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(dbServerGroupId)

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

	dbInfo, _, err := inst.Client.Postgresql.GetPostgresql(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	if len(dbInfo.VirtualServers) == 0 {
		diagnostics = diag.Errorf("no server found")
		return
	}

	port64, err := strconv.ParseInt(dbInfo.DbPort, 10, 32)
	if err != nil {
		return
	}

	// First server as representative
	vsInfo := dbInfo.VirtualServers[0]

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

	rd.Set("image_id", dbInfo.ImageId)
	rd.Set("server_name_prefix", serverNamePrefix)
	rd.Set("cluster_name", dbInfo.ServerGroupName)

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, dbInfo.ProductGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	serverTypeId := vsInfo.ScaleProductId
	// Set cpu / memory
	scale, err := client.FindProductById(ctx, inst.Client, dbInfo.ProductGroupId, serverTypeId)
	if err != nil {
		return
	}

	cpuFound := false
	memoryFound := false
	for _, item := range scale.Item {
		if item.ItemType == "cpu" {
			var cpuCount int
			cpuCount, err = strconv.Atoi(item.ItemValue)

			if err != nil {
				continue
			}

			rd.Set("cpu_count", cpuCount)
			cpuFound = true
		} else if item.ItemType == "memory" {
			var memorySize int
			memorySize, err = strconv.Atoi(item.ItemValue)

			if err != nil {
				continue
			}

			rd.Set("memory_size_gb", memorySize)
			memoryFound = true
		}
	}

	if !cpuFound || !memoryFound {
		return
	}

	rd.Set("contract_discount", contractDiscountName)

	rd.Set("vpc_id", vsInfo.VpcId)
	rd.Set("subnet_id", vsInfo.Network.NetworkId)

	securityGroupIds := make([]string, 0)
	for _, sg := range dbInfo.SecurityGroups {
		securityGroupIds = append(securityGroupIds, sg.SecurityGroupId)
	}

	rd.Set("security_group_ids", securityGroupIds)

	rd.Set("db_name", dbInfo.DbName)
	rd.Set("db_user_id", dbInfo.DbUserId)
	//rd.Set("db_user_password", "") // Remove Sensitive
	rd.Set("db_port", int(port64))

	//rd.Set("use_logging_audit")

	rd.Set("timezone", dbInfo.Timezone)
	if len(dbInfo.VirtualServers) > 1 {
		rd.Set("vip", dbInfo.Vip)
		rd.Set("external_vip", dbInfo.ExternalVip)
	} else {
		rd.Set("vip", dbInfo.VirtualServers[0].Software.SoftwareProperties["db.serviceIp"])
		rd.Set("external_vip", "")
	}

	if encoding, ok := dbInfo.DbConfigs["pg_encoding"]; ok {
		rd.Set("pg_encoding", encoding)
	} else {
		diagnostics = diag.Errorf("pg_encoding information not found")
		return
	}
	if locale, ok := dbInfo.DbConfigs["pg_locale"]; ok {
		rd.Set("pg_locale", locale)
	} else {
		diagnostics = diag.Errorf("pg_local information not found")
		return
	}

	// Storages
	var storageProductName string
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState == common.ProductAvailableState {
				storageProductName = productInfo.ProductName
				break
			}
		}
	}
	//var osStorage DatabaseStroagesResponse
	additionalStorages := common.HclListObject{}
	for i, bs := range vsInfo.BlockStorages {
		// Skip OS Storage (i == 0)
		if i == 0 {
			//if bs.BlockStorageType == BlockStorageTypeOS {
			//	//osStorage = bs
			//	continue
			//}
			continue
		}
		// First data storage is default storage
		if i == 1 {
			rd.Set("data_storage_size_gb", bs.BlockStorageSize)
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

	// 2023.02.21 KJ
	if dbInfo.Backup != nil {
		if err := rd.Set("backup", flattenBackupSettings(dbInfo.Backup)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourcePostgresqlUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	dbInfo, _, err := inst.Client.Postgresql.GetPostgresql(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dbInfo.VirtualServers) == 0 {
		return diag.Errorf("no server found")
	}

	// Only increase allowed
	if rd.HasChanges("data_storage_size_gb") {
		o, n := rd.GetChange("data_storage_size_gb")
		oldVal := o.(int)
		newVal := n.(int)
		if oldVal >= newVal {
			return diag.Errorf("storage size can only be increased. check size value : %d -> %d", oldVal, newVal)
		}

		// Update all
		for _, vs := range dbInfo.VirtualServers {
			var dataBlockId string
			for i, bs := range vs.BlockStorages {
				if i == 1 && bs.BlockStorageType == common.BlockStorageTypeData {
					dataBlockId = bs.BlockStorageId
					break
				}
			}
			if len(dataBlockId) == 0 {
				return diag.Errorf("default data storage not found")
			}

			_, _, err = inst.Client.Postgresql.UpdatePostgresqlBlockSize(ctx, rd.Id(), vs.VirtualServerId, dataBlockId, newVal)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := waitForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Only Addition & Increase allowed
	if rd.HasChanges("additional_storage") {
		o, n := rd.GetChange("additional_storage")
		oldVal := o.(common.HclListObject)
		newVal := n.(common.HclListObject)

		oldList := common.ConvertAdditionalStorageList(oldVal)
		newList := common.ConvertAdditionalStorageList(newVal)

		if len(oldList) > len(newList) {
			return diag.Errorf("removing additional storage is not allowed")
		}

		incIndices := make(map[int]int)
		for i := 0; i < len(oldList); i++ {
			if oldList[i].StorageUsage != newList[i].StorageUsage {
				return diag.Errorf("changing storage usage is not allowed")
			}
			if oldList[i].ProductName != newList[i].ProductName {
				return diag.Errorf("changing product name is not allowed")
			}
			if oldList[i].StorageSize > newList[i].StorageSize {
				return diag.Errorf("decreasing size is not allowed")
			}
			incIndices[i] = newList[i].StorageSize
		}

		// Update all
		for _, vsInfo := range dbInfo.VirtualServers {
			for i, blockInfo := range vsInfo.BlockStorages {
				// Skip
				if i < 2 {
					continue
				}
				if newSize, ok := incIndices[i-2]; ok {
					if newSize != int(blockInfo.BlockStorageSize) {
						_, _, err = inst.Client.Postgresql.UpdatePostgresqlBlockSize(ctx, rd.Id(), vsInfo.VirtualServerId, blockInfo.BlockStorageId, newSize)
						if err != nil {
							return diag.FromErr(err)
						}

						if err := waitForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
							return diag.FromErr(err)
						}
					} else {
						// No size change
					}
				}
			}
			// Newly added
			for i := len(oldList); i < len(newList); i++ {
				as := newList[i]

				_, _, err = inst.Client.Postgresql.AddPostgresqlBlock(ctx, rd.Id(), vsInfo.VirtualServerId, as.StorageUsage, as.StorageSize)
				if err != nil {
					return diag.FromErr(err)
				}

				if err := waitForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if rd.HasChanges("cpu_count", "memory_size_gb") {

		cpuCount := rd.Get("cpu_count").(int)
		memorySizeGB := rd.Get("memory_size_gb").(int)

		scaleId, err := client.FindScaleProduct(ctx, inst.Client, dbInfo.ProductGroupId, cpuCount, memorySizeGB)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(scaleId) == 0 {
			return diag.Errorf("no server type found")
		}

		// Update all
		for _, vsInfo := range dbInfo.VirtualServers {
			_, _, err = inst.Client.Postgresql.UpdatePostgresqlScale(ctx, rd.Id(), vsInfo.VirtualServerId, scaleId)

			if err := waitForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if rd.HasChanges("backup") {
		o, n := rd.GetChange("backup")

		oldbackups := expandBackupSettings(o.(*schema.Set))
		newbackups := expandBackupSettings(n.(*schema.Set))

		var backup *postgresql2.DatabaseBackup
		var useBackup bool
		if len(newbackups) > 0 {
			useBackup = true
			backup = newbackups[0]
			if len(backup.ObjectStorageId) < 1 {
				objectStorageId, err := getObjectStorageId(ctx, inst.Client, dbInfo.ServiceZoneId)
				if err != nil {
					return diag.FromErr(err)
				}
				backup.ObjectStorageId = objectStorageId
			}

		} else {
			useBackup = false
			backup = oldbackups[0]
		}
		updateBackupSettingRequest := postgresql2.UpdateBackupSettingRequest{
			UseBackup: &useBackup,
			Backup:    backup,
		}
		if _, _, err := inst.Client.Postgresql.UpdateBackupSetting(ctx, rd.Id(), updateBackupSettingRequest); err != nil {
			return diag.FromErr(err)
		}
		// if err := waitForPostgresql(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true); err != nil {
		// 	return diag.FromErr(err)
		// }
		if err := waitForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}); err != nil {
			return diag.FromErr(err)
		}

	}

	return resourcePostgresqlRead(ctx, rd, meta)
}
func resourcePostgresqlDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Postgresql.DeletePostgresql(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	if err := waitDeletedForPostgresql2(ctx, inst.Client, rd, rd.Id(), common.DatabaseProcessingStates(), []string{common.DeletedState}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForPostgresql2(ctx context.Context, scpClient *client.SCPClient, rd *schema.ResourceData, id string, pendingStates []string, targetStates []string) error {

	createStateConf := &resource.StateChangeConf{
		Pending: append(pendingStates, common.ActiveState),
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			info, statusCode, err := scpClient.Postgresql.GetPostgresql(ctx, id)

			// tflog.Debug(ctx, fmt.Sprintf("Wait -- %d\n", statusCode))

			if err != nil {
				if statusCode == 404 {
					return nil, "", nil
				}
				return "", "", err
			}
			running := 0
			for _, server := range info.VirtualServers {

				// tflog.Debug(ctx, fmt.Sprintf("Wait Start -- %d\n", statusCode))
				if server.VirtualServerState == common.RunningState && server.Software.SoftwareServiceState == common.RunningState {
					running = running + 1
				} else if server.VirtualServerState != common.RunningState && !contains(common.DatabaseProcessingStates(), server.VirtualServerState) {
					return "", "", fmt.Errorf("Error: %s", server.VirtualServerName)
				}
			}
			if running == len(info.VirtualServers) {
				return info, common.RunningState, nil
			}

			// tflog.Debug(ctx, fmt.Sprintf("Wait END -- %s\n", info.ServerGroupState))
			return info, info.ServerGroupState, nil
		},
		Timeout:                   rd.Timeout(schema.TimeoutCreate),
		Delay:                     20 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            5,
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for databse (%s) to be created: %s", id, err)
	}
	return nil
}

func waitDeletedForPostgresql2(ctx context.Context, scpClient *client.SCPClient, rd *schema.ResourceData, id string, pendingStates []string, targetStates []string) error {

	tflog.Debug(ctx, "waitForPostgresql")

	createStateConf := &resource.StateChangeConf{
		Pending: append(pendingStates, common.ActiveState),
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			info, statusCode, err := scpClient.Postgresql.GetPostgresql(ctx, id)

			// tflog.Debug(ctx, fmt.Sprintf("Wait -- %d\n", statusCode))

			if err != nil {
				if statusCode == 404 {
					return "", common.DeletedState, nil
				} else if statusCode == 400 {
					return nil, "", nil
				}

				return "", "", err
			}
			// tflog.Debug(ctx, fmt.Sprintf("Wait END -- %s\n", info.ServerGroupState))
			return info, info.ServerGroupState, nil
		},
		Timeout:                   rd.Timeout(schema.TimeoutCreate),
		Delay:                     20 * time.Second,
		MinTimeout:                20 * time.Second,
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            8,
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for databse (%s) to be created: %s", id, err)
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

func waitForPostgresql(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	baseInfo, baseStatusCode, baseErr := scpClient.Postgresql.GetPostgresql(ctx, id)
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
			info, c, err := scpClient.Postgresql.GetPostgresql(ctx, id)
			if err != nil {
				//virtual server not found
				if c == 400 && !errorOnNotFound {
					return "", common.DeletedState, nil
				}
				/*
					if c == 404 && !errorOnNotFound {
						return "", common.DeletedState, nil
					}
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

func getSecurityGroupIds(rd *schema.ResourceData) []string {
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	sgIds := make([]string, len(securityGroupIds))
	for i, valueIpv4 := range securityGroupIds {
		sgIds[i] = valueIpv4.(string)
	}
	return sgIds
}

func checkNameDuplication(ctx context.Context, meta interface{}, dbName string, clusterName string, serverNamePrefix string) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)
	resultList, err := inst.Client.Postgresql.ListPostgresql(ctx, dbName, clusterName, serverNamePrefix)
	if err != nil || len(resultList.Contents) == 0 {
		return nil
	}

	for _, content := range resultList.Contents {
		if clusterName == content.ServerGroupName {
			return diag.Errorf("Error: Input cluster name is invalid (maybe duplicated) : " + clusterName)
		}
	}

	return nil
}

func getObjectStorageId(ctx context.Context, scpClient *client.SCPClient, serviceZoneId string) (string, error) {
	response, err := scpClient.ObjectStorage.ReadObjectStorageList(ctx, serviceZoneId, objectstorage.ObjectStorageV3ControllerApiListObjectStorage3Opts{})
	if err != nil {
		return "", fmt.Errorf("failed while querying object storage list")
	}
	if response.TotalCount < 1 {
		return "", fmt.Errorf("no object storage list")
	}
	objectStorageId := response.Contents[0].ObsId

	tflog.Debug(ctx, fmt.Sprintf("Backup.objectStorageId : %s\n", objectStorageId))
	return objectStorageId, nil
}

func createBackupSettings(ctx context.Context, scpClient *client.SCPClient, serviceZoneId string, vAdvancedBackupSettings *schema.Set) (*postgresql2.DatabaseBackup, error) {
	backupSetting := &postgresql2.DatabaseBackup{}

	backUpList := vAdvancedBackupSettings.List()

	if len(backUpList) < 1 {
		return nil, nil
	}

	objectStorageId, err := getObjectStorageId(ctx, scpClient, serviceZoneId)
	if err != nil {
		return nil, err
	}
	backup := backUpList[0].(map[string]interface{})

	backupSetting.ObjectStorageId = objectStorageId

	if v, ok := backup["retention_day"].(int); ok {
		backupSetting.BackupRetentionDay = int32(v)
	}
	if v, ok := backup["start_hour"].(int); ok {
		backupSetting.BackupStartHour = int32(v)
	}
	return backupSetting, nil
}

func expandBackupSettings(vAdvancedBackupSettings *schema.Set) []*postgresql2.DatabaseBackup {
	backupSettings := []*postgresql2.DatabaseBackup{}

	for _, vAdvancedBackupSetting := range vAdvancedBackupSettings.List() {
		backupSetting := &postgresql2.DatabaseBackup{}

		mAdvancedBackupSetting := vAdvancedBackupSetting.(map[string]interface{})

		if v, ok := mAdvancedBackupSetting["objectstorage_id"].(string); ok && v != "" {
			backupSetting.ObjectStorageId = v
		}

		if v, ok := mAdvancedBackupSetting["retention_day"].(int); ok {
			backupSetting.BackupRetentionDay = int32(v)
		}

		if v, ok := mAdvancedBackupSetting["start_hour"].(int); ok {
			backupSetting.BackupStartHour = int32(v)
		}
		backupSettings = append(backupSettings, backupSetting)
	}

	return backupSettings
}

func flattenBackupSettings(backupSetting *postgresql2.DatabaseBackup) *schema.Set {

	mAdvancedBackupSetting := map[string]interface{}{
		"objectstorage_id": backupSetting.ObjectStorageId,
		"retention_day":    int(backupSetting.BackupRetentionDay),
		"start_hour":       int(backupSetting.BackupStartHour),
	}

	backupSchema := resourcePostgresqlBackup()
	return schema.NewSet(schema.HashResource(backupSchema), []interface{}{mAdvancedBackupSetting})
}

func expandHASetting(haSchemaSet *schema.Set) *postgresql2.HaSingleRestart {
	haSettings := []*postgresql2.HaSingleRestart{}

	for _, haSet := range haSchemaSet.List() {

		setting := haSet.(map[string]interface{})

		haSetting := &postgresql2.HaSingleRestart{}

		//if v, ok := setting["single_auto_restart"].(bool); ok {
		//	haSetting.SingleAutoRestart = v
		//}

		if v, ok := setting["use_vip_nat"].(bool); ok {
			haSetting.UseVipNat = &v
		}

		if v, ok := setting["reserved_natip_id"].(string); ok && v != "" {
			haSetting.ReservedNatIpId = v
		}

		if v, ok := setting["virtual_ip"].(string); ok && v != "" {
			haSetting.VirtualIp = v
		}
		haSettings = append(haSettings, haSetting)
	}

	var ha *postgresql2.HaSingleRestart
	if len(haSettings) > 0 {
		ha = haSettings[0]
	} else {
		ha = nil
	}

	return ha
}

func expandNetworkSetting(serverNamePrefix string, haSchemaSet *schema.Set) []postgresql2.DatabaseServerNetwork {
	networks := []postgresql2.DatabaseServerNetwork{}

	haSettings := haSchemaSet.List()

	if len(haSettings) > 0 {

		setting := haSettings[0].(map[string]interface{})

		network := postgresql2.DatabaseServerNetwork{
			ServerName: serverNamePrefix + "01",
			NodeType:   "ACTIVE",
		}
		if v, ok := setting["active_server_ip"].(string); ok && v != "" {
			network.ServiceIp = v
		}

		networks = append(networks, network)

		network = postgresql2.DatabaseServerNetwork{
			ServerName: serverNamePrefix + "02",
			NodeType:   "STANDBY",
		}
		if v, ok := setting["standby_server_ip"].(string); ok && v != "" {
			network.ServiceIp = v
		}
		networks = append(networks, network)

	} else {
		networks = append(networks, postgresql2.DatabaseServerNetwork{
			ServerName: serverNamePrefix + "01",
			NodeType:   "ACTIVE",
			ServiceIp:  "",
			NatIpId:    "",
		})
	}
	return networks
}
