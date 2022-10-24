package sqlserver

import (
	"context"
	"fmt"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/ScpDevTerra/trf-provider/scp/service/image"
	"github.com/ScpDevTerra/trf-sdk/library/postgresql2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

const (
	BlockStorageTypeOS      string = "OS"
	BlockStorageTypeData    string = "DATA"
	BlockStorageTypeArchive string = "ARCHIVE"
)

func ResourceSqlServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlServerCreate,
		ReadContext:   resourceSqlServerRead,
		UpdateContext: resourceSqlServerUpdate,
		DeleteContext: resourceSqlServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SqlServer virtual server image id.",
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
				Description:      "CPU core count",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"memory_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Memory size in gigabytes",
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
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Name of database. (3 to 20 lowercase alphabets)",
				ValidateDiagFunc: common.ValidateName3to20LowerAlphaOnly,
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
				Description: "SqlServer encoding. (Only 'UTF8' for now)",
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
				Description: "SqlServer locale. (Only 'C' for now)",
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
				Description:      "Default data storage size in gigabytes. (At least 10 GB required and size must be multiple of 10 : 10 GB, 20 GB, 30GB, ... )",
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
							Description: "Storage product name. (SSD)",
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
							Description:      "Default data storage size in gigabytes. (At least 10 GB required and size must be multiple of 10 : 10 GB, 20 GB, 30GB, ... )",
							ValidateDiagFunc: common.ValidateBlockStorageSize,
						},
					},
				},
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

	imageId := rd.Get("image_id").(string)

	serverNamePrefix := rd.Get("server_name_prefix").(string)
	clusterName := rd.Get("cluster_name").(string)

	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	contractDiscount := rd.Get("contract_discount").(string)

	vpcId := rd.Get("vpc_id").(string)
	subnetId := rd.Get("subnet_id").(string)
	sgIds := rd.Get("security_group_id").([]string)

	dbName := rd.Get("db_name").(string)
	dbUserId := rd.Get("db_user_id").(string)
	dbUserPassword := rd.Get("db_user_password").(string)
	dbPort := rd.Get("db_port").(int)

	useLoggingAudit := false // := rd.Get("use_logging_audit").(bool)

	timezone := rd.Get("timezone").(string)

	dataStorageSize := rd.Get("data_storage_size_gb").(int)
	additionalStorageList := rd.Get("additional_storage").(common.HclListObject)

	// Get vpc info
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}

	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, vpcInfo.ServiceZoneId, image.ActiveState, common.ServicedGroupDatabase, common.ServicedForSqlServer)
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

	// TODO : Replica & HA
	var networks []postgresql2.DatabaseServerNetwork
	networks = append(networks, postgresql2.DatabaseServerNetwork{
		ServerName:  serverNamePrefix + "01",
		NodeType:    "ACTIVE",
		ServiceIp:   "",
		ServiceIpId: "",
		NatIp:       "",
		NatIpId:     "",
	})

	_, _, err = inst.Client.Postgresql.CreatePostgresql(ctx, postgresql2.CreatePostgreSqlRequest{
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
			EncryptEnabled:            false,
			AdditionalBlockStorages:   blockStorages,
		},
		HighAvailability: nil,
		Replica:          nil,
		Network: &postgresql2.DatabaseNetwork{
			NetworkEnvType: strings.ToUpper(subnetInfo.SubnetType),
			VpcId:          vpcId,
			SubnetId:       subnetId,
			AutoServiceIp:  false,
			UseNat:         false,
			NatIp:          "",
			NatIpId:        "",
			ServerNetworks: networks,
		},
		SecurityGroupIds:  sgIds,
		Maintenance:       nil,
		Backup:            nil,
		UseDbLoggingAudit: useLoggingAudit,
		Timezone:          timezone,
		DbSoftwareOptions: pgSWOption,
	})
	if err != nil {
		return
	}

	// NOTE : response.ResourceId is empty
	resultList, err := inst.Client.SqlServer.ListSqlServer(ctx, dbName, clusterName, serverNamePrefix)
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

	dbInfo, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
	if err != nil {
		return
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
	rd.Set("security_group_id", dbInfo.SecurityGroup.SecurityGroupId)

	rd.Set("db_name", dbInfo.DbName)
	rd.Set("db_user_id", dbInfo.DbUserId)
	//rd.Set("db_user_password", "") // Remove Sensitive
	rd.Set("db_port", int(port64))

	//rd.Set("use_logging_audit")

	rd.Set("timezone", dbInfo.Timezone)

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
			if bs.BlockStorageType != BlockStorageTypeData {
				// Default storage not found
				rd.Set("data_storage_id", "")
				continue
			}
			rd.Set("data_storage_id", bs.BlockStorageId)
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

	return nil
}

func resourceSqlServerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// Only increase allowed
	if rd.HasChanges("data_storage_size_gb") {
		o, n := rd.GetChange("data_storage_size_gb")
		oldVal := o.(int)
		newVal := n.(int)
		if oldVal >= newVal {
			return diag.Errorf("storage size can only be increased. check size value : %d -> %d", oldVal, newVal)
		}

		dbInfo, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if len(dbInfo.VirtualServers) == 0 {
			return diag.Errorf("no server found")
		}

		// Update all
		for _, vs := range dbInfo.VirtualServers {
			var dataBlockId string
			for i, bs := range vs.BlockStorages {
				if i == 1 && bs.BlockStorageType == BlockStorageTypeData {
					dataBlockId = bs.BlockStorageId
					break
				}
			}
			if len(dataBlockId) == 0 {
				return diag.Errorf("default data storage not found")
			}

			_, _, err = inst.Client.SqlServer.UpdateSqlServerBlockSize(ctx, rd.Id(), vs.VirtualServerId, dataBlockId, newVal)
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForSqlServer(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
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

		dbInfo, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if len(dbInfo.VirtualServers) == 0 {
			return diag.Errorf("no server found")
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
						_, _, err = inst.Client.SqlServer.UpdateSqlServerBlockSize(ctx, rd.Id(), vsInfo.VirtualServerId, blockInfo.BlockStorageId, newSize)
						if err != nil {
							return diag.FromErr(err)
						}

						err = waitForSqlServer(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
						if err != nil {
							return diag.FromErr(err)
						}
					} else {
						// No size change
					}
				}
			}
		}
	}

	if rd.HasChanges("cpu_count", "memory_size_gb") {

		cpuCount := rd.Get("cpu_count").(int)
		memorySizeGB := rd.Get("memory_size_gb").(int)

		dbInfo, _, err := inst.Client.SqlServer.GetSqlServer(ctx, rd.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if len(dbInfo.VirtualServers) == 0 {
			return diag.Errorf("no server found")
		}

		scaleId, err := client.FindScaleProduct(ctx, inst.Client, dbInfo.ProductGroupId, cpuCount, memorySizeGB)
		if err != nil {
			return diag.FromErr(err)
		}
		// Update all
		for _, vsInfo := range dbInfo.VirtualServers {
			_, _, err = inst.Client.SqlServer.UpdateSqlServerScale(ctx, rd.Id(), vsInfo.VirtualServerId, scaleId)

			err = waitForSqlServer(ctx, inst.Client, rd.Id(), common.DatabaseProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceSqlServerRead(ctx, rd, meta)
}
func resourceSqlServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.SqlServer.DeleteSqlServer(ctx, rd.Id())
	if err != nil {
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
				if c == 404 && !errorOnNotFound {
					return "", common.DeletedState, nil
				}
				if c == 403 && !errorOnNotFound {
					return "", common.DeletedState, nil
				}
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
