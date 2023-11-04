package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/image"
	blockstorage2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/block-storage2"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/product"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/public-ip2"
	virtualserver2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/virtual-server2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_virtual_server", ResourceVirtualServer())
}

func ResourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualServerCreate,
		ReadContext:   resourceVirtualServerRead,
		UpdateContext: resourceVirtualServerUpdate,
		DeleteContext: resourceVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"virtual_server_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateName3to28AlphaDashStartsWithLowerCase,
				Description:      "Virtual server name",
			},
			"state": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateThatVmStateOnlyHasRunningOrStopped,
				Description:      "Virtual Server State",
			},
			/*"name_prefix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "VirtualServer name prefix",
				ValidateDiagFunc: common.ValidateName3to20NoSpecials,
			},*/
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable delete protection for this virtual server",
			},
			"anti_affinity": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable anti-affinity feature for this virtual server",
			},
			"cpu_count": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "CPU core count(2, 4, 8,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"memory_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Memory size in gigabytes(4, 8, 16,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"os_storage_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "OS(Boot) storage name. 3 to 28 alpha-numeric characters with space and dash starting with alphabet",
				ValidateDiagFunc: common.ValidateName3to28AlphaNumericWithSpaceAndDashStartsWithAlpha,
			},
			"os_storage_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "OS(Boot) storage size in gigabytes. (At least 100 GB required and size must be multiple of 10)",
				ValidateDiagFunc: common.ValidateBlockStorageSizeForOS,
			},
			"os_storage_encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable encryption feature in OS(Boot) storage. (WARNING) This option can not be changed after creation.",
			},
			"contract_discount": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract : None, 1 Year, 3 Year",
				ValidateDiagFunc: common.ValidateContractDesc,
			},
			"next_contract_discount": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "None",
				Description:      "Next Contract : None, 1 Year, 3 Year",
				ValidateDiagFunc: common.ValidateContractDesc,
			},
			"external_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "External block storage.",
				Elem:        common.ExternalStorageResourceSchema(),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id of this virtual server",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image id of this virtual server",
			},
			"initial_script_content": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Initialization script",
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				//ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security-Group ids of this virtual server. Each security-group must be a valid security-group resource which is attached to the VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id of this virtual server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"internal_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address for internal IP assignment.",
			},
			"local_subnet": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Local subnet id of this virtual server. Local subnet must be a valid local subnet resource which is attached to the Subnet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network interface id",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet Id",
						},
						"ipv4": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Subnet ip address.",
						},
					},
				},
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable NAT IP feature.",
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:    true,
				Description: "Public IP id of this virtual server. Public-IP must be a valid public-ip resource which is attached to the VPC.",
			},
			"use_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DNS feature for this virtual server.",
			},
			"admin_account": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateName3to20DashUnderscore,
				Description:      "Admin account for this virtual server OS. For linux, this must be 'root'. For Windows, this must not be 'administrator'.",
			},
			"admin_password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: common.ValidatePassword8to20,
				Description:      "Admin account password for this virtual server OS.",
			},
			"key_pair_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key Pair Id",
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Placement Group Id",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of this virtual server",
			},
			"nat_ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NAT IP address of this virtual server",
			},
			"availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Availability Zone Name",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags",
			},
		},
		Description: "Provides a Virtual Server resource.",
	}
}

func getSecurityGroupIds(rd *schema.ResourceData) []string {
	securityGroupIds := rd.Get("security_group_ids").([]interface{})
	sgIds := make([]string, len(securityGroupIds))
	for i, valueIpv4 := range securityGroupIds {
		sgIds[i] = valueIpv4.(string)
	}
	return sgIds
}

type LocalSubnetInfo struct {
	NicId    string
	SubnetId string
	Ipv4     string
}

func convertLocalSubnet(list common.HclListObject) ([]LocalSubnetInfo, error) {
	var result []LocalSubnetInfo
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := LocalSubnetInfo{}
		if nicId, ok := itemObject["id"]; ok {
			info.NicId = nicId.(string)
		}
		if subnetId, ok := itemObject["subnet_id"]; ok {
			info.SubnetId = subnetId.(string)
		} else {
			return result, fmt.Errorf("Subnet id not found")
		}
		if ipv4, ok := itemObject["ipv4"]; ok {
			info.Ipv4 = ipv4.(string)
		}

		result = append(result, info)
	}
	return result, nil
}

func resourceVirtualServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	vsName := rd.Get("virtual_server_name").(string)
	isDeleteProtected := rd.Get("delete_protection").(bool)
	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	osStorageName := rd.Get("os_storage_name").(string)
	osStorageSize := rd.Get("os_storage_size_gb").(int)
	osStorageEncrypted := rd.Get("os_storage_encrypted").(bool)

	//serviceLevel := rd.Get("service_level").(string)
	contractDiscount := rd.Get("contract_discount").(string)
	externalStorageList := rd.Get("external_storage").(common.HclListObject)

	subnetId := rd.Get("subnet_id").(string)
	//localSubnetId := rd.Get("local_subnet_id").(string)
	localSubnetInfos, err := convertLocalSubnet(rd.Get("local_subnet").(common.HclListObject))
	if err != nil {
		return
	}
	publicIpId := rd.Get("public_ip_id").(string)

	adminAccount := rd.Get("admin_account").(string)
	adminPassword := rd.Get("admin_password").(string)

	keyPairId := rd.Get("key_pair_id").(string)
	placementGroupId := rd.Get("placement_group_id").(string)

	if adminPassword == "" && keyPairId == "" {
		return diag.Errorf("Either admin_password or key_pair_id must be specified.")
	}

	antiAffinity := rd.Get("anti_affinity").(bool)

	vpcId := rd.Get("vpc_id").(string)

	imageId := rd.Get("image_id").(string)

	useDNS := rd.Get("use_dns").(bool)

	internalIpAddress := rd.Get("internal_ip_address").(string)
	if len(internalIpAddress) != 0 {
		res, err := inst.Client.Subnet.CheckAvailableSubnetIp(ctx, subnetId, internalIpAddress)
		if err != nil {
			return diag.FromErr(err)
		}
		if *res.Result == false {
			return diag.Errorf("Not Available Internal Ip Address")
		}
	}
	// Get vpc info
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}
	isOsWindows, targetProductGroupId, err := getImageInfo(ctx, vpcInfo.ServiceZoneId, imageId, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if !isOsWindows && adminAccount != common.LinuxAdminAccount {
		adminAccount = common.LinuxAdminAccount
		log.Println("Linux admin account must be root")
	}

	if isOsWindows && (adminAccount == common.WindowsAdminAccount || len(adminAccount) < 5) {
		diagnostics = diag.Errorf("Windows admin account must be 5 to 20 alpha-numeric characters with special character and not be 'administrator'.")
		return
	}

	if len(targetProductGroupId) == 0 {
		diagnostics = diag.Errorf("Product group id not found from image")
		return
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
	//productGroup, err := inst.Client.Product.GetProducesList(ctx, vpcInfo.ServiceZoneId, targetProductGroupId, "")
	if err != nil {
		return diag.FromErr(err)
	}

	// Find OS disk product
	//osDiskProductId, err := common.FirstProductId(common.ProductDefaultDisk, &productGroup)
	osDiskProductInfo, err := common.FirstProductId(common.ProductDisk, &productGroup)
	if err != nil {
		return
	}
	osDiskProductName := osDiskProductInfo.ProductName

	serviceLevelToId := common.ProductToIdMap(common.ProductServiceLevel, &productGroup)
	if len(serviceLevelToId) == 0 {
		diagnostics = diag.Errorf("Failed to find service level")
		return
	}

	externalDiskProductNameToId := common.ProductToIdMap(common.ProductDisk, &productGroup)
	if len(externalDiskProductNameToId) == 0 {
		diagnostics = diag.Errorf("Failed to find external disk product")
		return
	}

	diskProductIdToNameMap := common.ProductIdToNameMap(common.ProductDisk, &productGroup)

	contractToId := common.ProductToIdMap(common.ProductContractDiscount, &productGroup)
	if len(contractToId) == 0 {
		diagnostics = diag.Errorf("Failed to find contract info")
		return
	}

	// Find VM scaling
	var scaleProductInfo *product.ProductForCalculatorResponse
	scaleProductInfo = getScaleProductInfoFromProductGroup(productGroup, cpuCount, memorySizeGB)
	if scaleProductInfo == nil {
		return
	}
	// Find service level
	//serviceLevelId, ok := serviceLevelToId["None"]
	//if !ok {
	//	diagnostics = diag.Errorf("Invalid service level")
	//	return
	//}

	// Find contract
	//contractId, ok := contractToId[contractDiscount]
	//if !ok {
	//	diagnostics = diag.Errorf("Invalid contract")
	//	return
	//}

	// Find ServerGroup
	serverGroupList, err := inst.Client.ServerGroup.GetServerGroup(ctx, "", []string{common.ServicedForVirtualServer})
	if err != nil {
		return
	}

	var serverGroupId string
	if antiAffinity {
		for _, sg := range serverGroupList.Contents {
			if sg.AffinityPolicyType != "NONE" {
				serverGroupId = sg.ServerGroupId
				break
			}
		}
	}

	// Find PublicIP
	useNAT := false
	if len(publicIpId) != 0 {
		var publicIpInfo publicip2.DetailPublicIpResponse
		publicIpInfo, _, err = inst.Client.PublicIp.GetPublicIp(ctx, publicIpId)
		if err != nil {
			return
		}
		if len(publicIpInfo.IpAddress) != 0 {
			useNAT = true
		}
	}

	natEnabled := rd.Get("nat_enabled").(bool)
	if natEnabled == true {
		useNAT = true
	}

	var extStorages []virtualserver.BlockStorageInfo
	externalStorageInfoList := common.ConvertExternalStorageList(externalStorageList)

	// settings by types
	imageType, err := inst.Client.Image.GetImageType(ctx, imageId)
	if imageType == "STANDARD" {
		for _, extStorageInfo := range externalStorageInfoList {
			extStorages = append(extStorages, virtualserver.BlockStorageInfo{
				BlockStorageName: extStorageInfo.Name,
				DiskSize:         int32(extStorageInfo.StorageSize),
				EncryptEnabled:   extStorageInfo.Encrypted,
				DiskType:         extStorageInfo.ProductName,
			})
		}
	}

	if imageType == "CUSTOM" {
		info, _, err := inst.Client.CustomImage.GetCustomImage(ctx, imageId)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, d := range info.Disks {
			if *d.BootEnabled {
				osStorageSize = int(d.DiskSize) // int(bs.BlockStorageSize)
				osStorageEncrypted = *d.EncryptEnabled
				if osDiskType, ok := diskProductIdToNameMap[d.ProductId]; ok {
					osDiskProductName = osDiskType // bs.ProductId
				} // bs.EncryptEnabled
			}
		}
	}

	initialScriptShell := "bash"
	if isOsWindows {
		initialScriptShell = "pwsh"
	}
	initialScript := virtualserver.InitialScriptInfo{
		EncodingType:         "plain",
		InitialScriptContent: rd.Get("initial_script_content").(string),
		InitialScriptShell:   initialScriptShell,
		InitialScriptType:    "text",
	}

	createRequest := virtualserver.CreateRequest{
		BlockStorage: virtualserver.BlockStorageInfo{
			BlockStorageName: osStorageName,
			DiskSize:         int32(osStorageSize),
			EncryptEnabled:   osStorageEncrypted,
			DiskType:         osDiskProductName,
		},
		ContractDiscount:          contractDiscount,
		DeletionProtectionEnabled: isDeleteProtected,
		DnsEnabled:                useDNS,
		ExtraBlockStorages:        extStorages,
		ImageId:                   imageId,
		InitialScript:             initialScript,
		LocalSubnet:               virtualserver.LocalSubnetInfo{},
		Nic: virtualserver.NicInfo{
			InternalIpAddress: internalIpAddress,
			NatEnabled:        useNAT,
			PublicIpAddressId: publicIpId,
			SubnetId:          subnetId,
		},
		OsAdmin: virtualserver.OsAdminInfo{
			OsUserId:       adminAccount,
			OsUserPassword: adminPassword,
		},
		SecurityGroupIds:     getSecurityGroupIds(rd),
		ServerGroupId:        serverGroupId,
		ServerType:           scaleProductInfo.ProductName,
		ServiceZoneId:        vpcInfo.ServiceZoneId,
		VirtualServerName:    vsName,
		AvailabilityZoneName: rd.Get("availability_zone_name").(string),
		Tags:                 getTagRequestArray(rd),
		KeyPairId:            keyPairId,
		PlacementGroupId:     placementGroupId,
	}
	var createResponse virtualserver2.AsyncResponse
	if keyPairId != "" {
		createResponse, err = inst.Client.VirtualServer.CreateVirtualServerV4(ctx, createRequest)
	} else {
		createResponse, err = inst.Client.VirtualServer.CreateVirtualServer(ctx, createRequest)
	}
	if err != nil {
		return
	}

	err = WaitForVirtualServerStatus(ctx, inst.Client, createResponse.ResourceId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return
	}

	if len(localSubnetInfos) > 0 {
		// Even in RUNNING status, some update API may throw exceptions.
		// We have to wait for a while to sync up with the internal request status.
		time.Sleep(30 * time.Second)

		// in case of openapi, need additional attaching..
		for _, localSubnetInfo := range localSubnetInfos {
			_, _, err = inst.Client.VirtualServer.AttachLocalSubnet(ctx, createResponse.ResourceId, localSubnetInfo.SubnetId, localSubnetInfo.Ipv4)
			if err != nil {
				return
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, createResponse.ResourceId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
	}

	//if strings.Compare(strings.ToUpper(rd.Get("state").(string)), "STOPPED") == 0 {
	//	_, err = inst.Client.VirtualServer.StopVirtualServer(ctx, createResponse.ResourceId)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//	err = WaitForVirtualServerStatus(ctx, inst.Client, createResponse.ResourceId, common.VirtualServerProcessingStates(), []string{common.StoppedState}, true)
	//	if err != nil {
	//		return
	//	}
	//}

	rd.SetId(createResponse.ResourceId)

	return resourceVirtualServerRead(ctx, rd, meta)
}

func getScaleProductInfoFromProductGroup(productGroup product.ProductGroupDetailResponse, cpuCount int, memorySizeGB int) *product.ProductForCalculatorResponse {
	var resultProduct *product.ProductForCalculatorResponse = nil
	if scaleProductInfoArray, ok := productGroup.Products[common.ProductScale]; ok {
		cpuString := strconv.Itoa(cpuCount)
		memoryString := strconv.Itoa(memorySizeGB)
		resultProduct = getScaleProductInfo(scaleProductInfoArray, cpuString, memoryString)
		if resultProduct != nil {
			return resultProduct
		}
	}
	return resultProduct
}

func getScaleProductInfo(scaleProductInfoArray []product.ProductForCalculatorResponse, cpuString string, memoryString string) *product.ProductForCalculatorResponse {
	const FOUND = 3
	for _, product := range scaleProductInfoArray {
		if product.ProductType != common.ProductScale {
			continue
		}
		if isMatchedCpuMemoryItemValueFound(product, cpuString, memoryString) == FOUND {
			return &product
		}
	}
	return nil
}

func isMatchedCpuMemoryItemValueFound(product product.ProductForCalculatorResponse, cpuString string, memoryString string) int {
	const CPU_FLAG = 1
	const MEMORY_FLAG = 2
	const NOT_FOUND = 0
	flag := NOT_FOUND

	for _, item := range product.Item {
		if hasMatchedCpuItemValue(item, cpuString) {
			flag |= CPU_FLAG
		} else if hasMatchedMemoryItemValue(item, memoryString) {
			flag |= MEMORY_FLAG
		}
	}
	return flag
}

func hasMatchedMemoryItemValue(item product.ItemForCalculatorResponse, memoryString string) bool {
	return item.ItemType == "memory" && item.ItemValue == memoryString
}

func hasMatchedCpuItemValue(item product.ItemForCalculatorResponse, cpuString string) bool {
	return item.ItemType == "cpu" && item.ItemValue == cpuString
}

func getImageInfo(ctx context.Context, serviceZoneId string, imageId string, meta interface{}) (bool, string, error) {
	var err error = nil
	inst := meta.(*client.Instance)

	isOsWindows := false
	var targetProductGroupId string
	imageType, err := inst.Client.Image.GetImageType(ctx, imageId)
	if err != nil {
		return false, "", err
	}

	if imageType == "STANDARD" {
		standardImages, err := inst.Client.Image.GetStandardImageList(ctx, serviceZoneId, image.ActiveState, common.ServicedGroupCompute, common.ServicedForVirtualServer)
		if err != nil {
			return false, "", err
		}

		for _, c := range standardImages.Contents {
			if c.ImageId == imageId {
				targetProductGroupId = c.ProductGroupId
				if c.OsType == common.OsTypeWindows {
					isOsWindows = true
				}
			}
		}
	} else if imageType == "CUSTOM" {
		customImages, err := inst.Client.CustomImage.GetCustomImageList(ctx, image2.CustomImageV2ApiListCustomImagesOpts{
			ImageState:       optional.NewString(image.ActiveState),
			ServicedGroupFor: optional.NewString(common.ServicedGroupCompute),
			ServicedFor:      optional.NewString(common.ServicedForVirtualServer),
			ServiceZoneId:    optional.NewString(serviceZoneId),
		})
		if err != nil {
			return false, "", err
		}

		for _, c := range customImages.Contents {
			if c.ImageId == imageId {
				targetProductGroupId = c.ProductGroupId
				if c.OsType == common.OsTypeWindows {
					isOsWindows = true
				}
			}
		}
	} else if imageType == "MIGRATION" {
		migrationImages, err := inst.Client.MigrationImage.GetMigrationImageList(ctx, image2.MigrationImageV2ApiListMigrationImagesOpts{
			ServiceZoneId:    optional.NewString(serviceZoneId),
			ImageState:       optional.NewString(image.ActiveState),
			ServicedFor:      optional.NewString(common.ServicedForVirtualServer),
			ServicedGroupFor: optional.NewString(common.ServicedGroupCompute),
			Page:             optional.NewInt32(0),
			Size:             optional.NewInt32(10000),
			Sort:             optional.NewInterface([]string{"imageName:asc"}),
		})
		if err != nil {
			return false, "", err
		}

		for _, c := range migrationImages.Contents {
			if c.ImageId == imageId {
				targetProductGroupId = c.ProductGroupId
				if c.OsType == common.OsTypeWindows {
					isOsWindows = true
				}
			}
		}
	}

	return isOsWindows, targetProductGroupId, err
}

func getTagRequestArray(rd *schema.ResourceData) []virtualserver.TagRequest {
	tags := rd.Get("tags").(map[string]interface{})
	tagsRequests := make([]virtualserver.TagRequest, 0)
	for key, value := range tags {
		tagsRequests = append(tagsRequests, virtualserver.TagRequest{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	return tagsRequests
}

func resourceVirtualServerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	virtualServerInfo, _, err := inst.Client.VirtualServer.GetVirtualServer(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	nicInfo, err := inst.Client.VirtualServer.GetNicList(ctx, virtualServerInfo.VirtualServerId)
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, virtualServerInfo.ProductGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	productIdToNameMapper := make(map[string]string)
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			productIdToNameMapper[productInfo.ProductId] = productInfo.ProductName
		}
	}

	rd.Set("image_id", virtualServerInfo.ImageId)
	rd.Set("virtual_server_name", virtualServerInfo.VirtualServerName)
	rd.Set("delete_protection", virtualServerInfo.DeletionProtectionEnabled)
	//rd.Set("service_level", virtualServerInfo.ServiceLevel)
	rd.Set("contract_discount", virtualServerInfo.Contract)
	if len(virtualServerInfo.NextContractId) > 0 {
		info, err := inst.Client.Product.GetProductsUsingGET(ctx, virtualServerInfo.NextContractId)
		if err != nil {
			return diag.FromErr(err)
		}
		rd.Set("next_contract_discount", info.Name)
	} else {
		rd.Set("next_contract_discount", "None")
	}
	rd.Set("vpc_id", virtualServerInfo.VpcId)
	rd.Set("use_dns", virtualServerInfo.DnsEnabled)
	rd.Set("initial_script_content", virtualServerInfo.InitialScriptContent)

	sgIds := common.HclListObject{}
	for _, sg := range virtualServerInfo.SecurityGroupIds {
		sgIds = append(sgIds, sg.SecurityGroupId)
	}
	rd.Set("security_group_ids", sgIds)

	//extStorages := common.HclListObject{}
	//for _, blockId := range virtualServerInfo.BlockStorageIds {
	//
	//	var blockInfo blockstorage2.BlockStorageResponse
	//	blockInfo, _, err = inst.Client.BlockStorage.ReadBlockStorage(ctx, blockId)
	//	if err != nil {
	//		continue
	//	}
	//	if *blockInfo.IsBootDisk {
	//		rd.Set("os_storage_name", blockInfo.BlockStorageName)
	//		rd.Set("os_storage_size_gb", int(blockInfo.BlockStorageSize))
	//		rd.Set("os_storage_encrypted", blockInfo.EncryptEnabled)
	//	} else {
	//		extStorageInfo := common.HclKeyValueObject{}
	//		extStorageInfo["block_storage_id"] = blockInfo.BlockStorageId
	//		extStorageInfo["name"] = blockInfo.BlockStorageName
	//		extStorageInfo["storage_size_gb"] = int(blockInfo.BlockStorageSize)
	//		extStorageInfo["encrypted"] = blockInfo.EncryptEnabled
	//		if blockProductName, ok := productIdToNameMapper[blockInfo.ProductId]; ok {
	//			extStorageInfo["product_name"] = blockProductName
	//		} else {
	//			extStorageInfo["product_name"] = "UNKNOWN"
	//		}
	//		extStorageInfo["product_id"] = blockInfo.ProductId
	//		extStorageInfo["shared_type"] = blockInfo.SharedType
	//		extStorages = append(extStorages, extStorageInfo)
	//	}
	//}
	// detail Virtual Server 내의 blockId 리스트 정보를 가지고, api 호출을 하여
	// detail BlockStorageReponse 정보를 리스트로 구성
	blockStorageResponseList := getBlockStorageResponseList(ctx, virtualServerInfo.BlockStorageIds, inst)

	for _, blockStorageResponse := range blockStorageResponseList {
		if *blockStorageResponse.IsBootDisk {
			rd.Set("os_storage_name", blockStorageResponse.BlockStorageName)
			rd.Set("os_storage_size_gb", int(blockStorageResponse.BlockStorageSize))
			rd.Set("os_storage_encrypted", blockStorageResponse.EncryptEnabled)
		}
	}

	prevExternalStorageList := rd.Get("external_storage").([]interface{})
	extStorages := make([]map[string]interface{}, 0)
	// 이전 externa_storage 의 순서를 유지하기 위한 처리
	for _, prevExternalStorage := range prevExternalStorageList {
		mapPrevExternalStorage := prevExternalStorage.(map[string]interface{})
		for _, blockStorageResponse := range blockStorageResponseList {
			if strings.Compare(mapPrevExternalStorage["name"].(string), blockStorageResponse.BlockStorageName) == 0 {
				extStorages = append(extStorages, getExternalStorageMapSchema(blockStorageResponse, productIdToNameMapper, mapPrevExternalStorage))
				break
			}
		}
	}

	for _, blockStorageResponse := range blockStorageResponseList {
		if *blockStorageResponse.IsBootDisk {
			continue
		}
		var idx int
		for idx = 0; idx < len(extStorages); idx++ {
			extStorage := extStorages[idx]
			if blockStorageResponse.BlockStorageId == extStorage["block_storage_id"] {
				break
			}
		}
		if idx == len(extStorages) {
			extStorage := getExtStorage(blockStorageResponse, productIdToNameMapper)
			extStorages = append(extStorages, extStorage)
		}
	}

	rd.Set("external_storage", extStorages)

	// Set cpu / memory
	scale, err := client.FindProductById(ctx, inst.Client, virtualServerInfo.ProductGroupId, virtualServerInfo.ServerTypeId)
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

	ipv4 := virtualServerInfo.Ip
	var localSubnetInfos []common.HclKeyValueObject
	//var localSubnetInfos []LocalSubnetInfo
	//var localSubnetId string
	//var localSubnetIpv4 string
	var subnetId string
	//var subnetIpv4 string
	var natIpv4 string
	//var ipv4 string
	//var natIpv4 string
	//var subnetId string
	//var publicIpId string
	for _, vsNicId := range virtualServerInfo.NicIds {
		for _, nic := range nicInfo.Contents {
			if nic.NicId == vsNicId {
				if nic.SubnetType == "VM" || nic.SubnetType == "BM" {
					localSubnetInfos = append(localSubnetInfos, common.HclKeyValueObject{
						"id":        nic.NicId,
						"subnet_id": nic.SubnetId,
						"ipv4":      nic.Ip,
					})
				} else if nic.SubnetType == "PUBLIC" {
					subnetId = nic.SubnetId
					natIpv4 = nic.NatIp
				} else {
					subnetId = nic.SubnetId
				}
				break
			}
		}
	}

	rd.Set("ipv4", ipv4)
	rd.Set("subnet_id", subnetId)
	rd.Set("local_subnet", localSubnetInfos)
	rd.Set("nat_ipv4", natIpv4)

	if natIpv4 != "" {
		publicIpInfo, err := inst.Client.PublicIp.GetPublicIps(ctx,
			&publicip2.PublicIpOpenApiV3ControllerApiListPublicIpsV3Opts{
				IpAddress:     optional.NewString(natIpv4),
				PublicIpState: optional.String{},
				UplinkType:    optional.String{},
				VpcId:         optional.NewString(virtualServerInfo.VpcId),
				CreatedBy:     optional.String{},
				Page:          optional.Int32{},
				Size:          optional.Int32{},
				Sort:          optional.Interface{},
			})
		if err != nil {
			diagnostics = diag.FromErr(err)
			return
		}

		if len(publicIpInfo.Contents) != 0 {
			rd.Set("public_ip_id", publicIpInfo.Contents[0].PublicIpAddressId)
		}

	} else {
		rd.Set("public_ip_id", "")
	}
	rd.Set("state", virtualServerInfo.VirtualServerState)
	rd.Set("key_pair_id", virtualServerInfo.KeyPairId)
	rd.Set("placement_group_id", virtualServerInfo.PlacementGroupId)

	return nil
}

func getExternalStorageMapSchema(blockStorageResponse blockstorage2.BlockStorageResponse, productIdToNameMapper map[string]string, mapPrevExternalStorage map[string]interface{}) map[string]interface{} {
	extStorage := getExtStorage(blockStorageResponse, productIdToNameMapper)
	tags := mapPrevExternalStorage["tags"].(map[string]interface{})
	tagElem := make(map[string]string)
	for key, value := range tags {
		tagElem[key] = value.(string)
	}
	extStorage["tags"] = tagElem
	return extStorage
}

func getExtStorage(blockStorageResponse blockstorage2.BlockStorageResponse, productIdToNameMapper map[string]string) map[string]interface{} {
	extStorage := make(map[string]interface{})
	extStorage["block_storage_id"] = blockStorageResponse.BlockStorageId
	extStorage["name"] = blockStorageResponse.BlockStorageName
	extStorage["storage_size_gb"] = int(blockStorageResponse.BlockStorageSize)
	extStorage["encrypted"] = blockStorageResponse.EncryptEnabled
	if blockProductName, ok := productIdToNameMapper[blockStorageResponse.ProductId]; ok {
		extStorage["product_name"] = blockProductName
	} else {
		extStorage["product_name"] = "UNKNOWN"
	}
	extStorage["product_id"] = blockStorageResponse.ProductId
	extStorage["shared_type"] = blockStorageResponse.SharedType
	return extStorage
}

func getBlockStorageResponseList(ctx context.Context, blockStorageIds []string, inst *client.Instance) []blockstorage2.BlockStorageResponse {
	blockStorageResponseList := make([]blockstorage2.BlockStorageResponse, 0)
	for _, id := range blockStorageIds {
		blockInfo, _, err := inst.Client.BlockStorage.ReadBlockStorage(ctx, id)
		if err != nil {
			continue
		}
		blockStorageResponseList = append(blockStorageResponseList, blockInfo)
	}
	return blockStorageResponseList
}

func resourceVirtualServerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	virtualServerInfo, _, err := inst.Client.VirtualServer.GetVirtualServer(ctx, rd.Id())
	if err != nil {
		return
	}

	targetProductGroupId := virtualServerInfo.ProductGroupId

	if rd.HasChanges("cpu_count", "memory_size_gb") {
		cpuCount := rd.Get("cpu_count").(int)
		memorySizeGB := rd.Get("memory_size_gb").(int)

		// Find VM scaling
		productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
		if err != nil {
			return diag.FromErr(err)
		}
		scaleProductInfo := getScaleProductInfoFromProductGroup(productGroup, cpuCount, memorySizeGB)
		if scaleProductInfo == nil {
			return diag.Errorf("Server Type Not Found !!")
		}

		// Update scale
		_, err = inst.Client.VirtualServer.UpdateScale(ctx, virtualServerInfo.VirtualServerId, scaleProductInfo.ProductName)
		if err != nil {
			return
		}

		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("delete_protection") {
		isDeleteProtectionEnabled := rd.Get("delete_protection").(bool)
		_, _, err = inst.Client.VirtualServer.UpdateDeleteProtectionEnabled(ctx, rd.Id(), isDeleteProtectionEnabled)
		if err != nil {
			return
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("security_group_ids") {
		securityGroupIds := getSecurityGroupIds(rd)
		currentSGData := make(map[string]bool)
		for _, sgId := range securityGroupIds {
			currentSGData[sgId] = true
		}

		var deletedIds []string
		for _, sgResult := range virtualServerInfo.SecurityGroupIds {
			if _, ok := currentSGData[sgResult.SecurityGroupId]; ok {
				currentSGData[sgResult.SecurityGroupId] = false
			} else {
				deletedIds = append(deletedIds, sgResult.SecurityGroupId)
			}
		}

		for _, sgId := range deletedIds {
			_, err = inst.Client.VirtualServer.DeleteSecurityGroup(ctx, virtualServerInfo.VirtualServerId, sgId)
			if err != nil {
				continue
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}

		for sgId, v := range currentSGData {
			if !v {
				// Already attached security-group
				continue
			}
			_, err = inst.Client.VirtualServer.AddSecurityGroup(ctx, virtualServerInfo.VirtualServerId, sgId)
			if err != nil {
				return
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
	}

	if rd.HasChanges("local_subnet") {

		var nicInfoList virtualserver2.ListResponseOfNicResponse
		nicInfoList, err = inst.Client.VirtualServer.GetNicList(ctx, rd.Id())
		if err != nil {
			return
		}

		mapSubnetId2NicId := make(map[string]string)
		for _, nic := range nicInfoList.Contents {
			mapSubnetId2NicId[nic.SubnetId] = nic.NicId
		}

		o, n := rd.GetChange("local_subnet")
		oldVal := o.(common.HclListObject)
		newVal := n.(common.HclListObject)
		var oldList []LocalSubnetInfo
		var newList []LocalSubnetInfo
		oldList, err = convertLocalSubnet(oldVal)
		newList, err = convertLocalSubnet(newVal)

		type matchInfo struct {
			Trigger  int
			NicId    string
			SubnetId string
			Ipv4     string
		}
		match := make(map[string]matchInfo)
		for _, oSubInfo := range oldList {
			okey := oSubInfo.NicId + "|" + oSubInfo.SubnetId + "|" + oSubInfo.Ipv4
			match[okey] = matchInfo{
				Trigger:  -1, // Default detach
				NicId:    oSubInfo.NicId,
				SubnetId: oSubInfo.SubnetId,
				Ipv4:     oSubInfo.Ipv4,
			}
		}
		for _, nSubInfo := range newList {
			nkey := nSubInfo.NicId + "|" + nSubInfo.SubnetId + "|" + nSubInfo.Ipv4
			if _, ok := match[nkey]; ok {
				match[nkey] = matchInfo{
					Trigger:  0, // Preserve
					NicId:    nSubInfo.NicId,
					SubnetId: nSubInfo.SubnetId,
					Ipv4:     nSubInfo.Ipv4,
				}
			} else {
				match[nkey] = matchInfo{
					Trigger:  1, // Attach
					NicId:    nSubInfo.NicId,
					SubnetId: nSubInfo.SubnetId,
					Ipv4:     nSubInfo.Ipv4,
				}
			}
		}

		for _, mi := range match {
			if mi.Trigger > 0 {
				// Attach
				_, _, err = inst.Client.VirtualServer.AttachLocalSubnet(ctx, virtualServerInfo.VirtualServerId, mi.SubnetId, mi.Ipv4)
				if err != nil {
					return
				}
			} else if mi.Trigger < 0 {
				// Detach
				if nicId, ok := mapSubnetId2NicId[mi.SubnetId]; ok {
					_, _, err = inst.Client.VirtualServer.DetachLocalSubnet(ctx, virtualServerInfo.VirtualServerId, nicId)
					if err != nil {
						return
					}
				} else {
					diagnostics = diag.Errorf("detach target local subnet network interface not found.")
					return
				}
			} else {
				// No-op
				continue
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
	}

	if rd.HasChanges("public_ip_id") || rd.HasChanges("nat_enabled") {

		var nicInfoList virtualserver2.ListResponseOfNicResponse
		nicInfoList, err = inst.Client.VirtualServer.GetNicList(ctx, rd.Id())
		if err != nil {
			return
		}

		// Find public subnet
		subnetId := rd.Get("subnet_id")
		var nicId string
		for _, nicInfo := range nicInfoList.Contents {
			if nicInfo.SubnetId == subnetId {
				nicId = nicInfo.NicId
			}
		}
		if len(nicId) == 0 {
			diagnostics = diag.Errorf("public subnet network interface not found")
			return
		}

		natEnabled := rd.Get("nat_enabled").(bool)
		publicIpId := rd.Get("public_ip_id").(string)
		if rd.HasChanges("nat_enabled") {
			if rd.HasChanges("public_ip_id") {
				err = detachAndAttachPublicIpId(ctx, rd, inst, virtualServerInfo.VirtualServerId, nicId, natEnabled)
				if err != nil {
					return diag.FromErr(err)
				}
			} else {
				if len(publicIpId) == 0 {
					if natEnabled == false {
						_, _, err = inst.Client.VirtualServer.DetachPublicIp(ctx, virtualServerInfo.VirtualServerId, nicId)
					} else {
						_, _, err = inst.Client.VirtualServer.AttachPublicIp(ctx, virtualServerInfo.VirtualServerId, nicId, "")
					}
					if err != nil {
						return
					}
				}
			}
		} else {
			if rd.HasChanges("public_ip_id") {
				err = detachAndAttachPublicIpId(ctx, rd, inst, virtualServerInfo.VirtualServerId, nicId, natEnabled)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("internal_ip_address") || rd.HasChanges("subnet_id") {
		newSubnetId := rd.Get("subnet_id").(string)
		newInternalIpAddress := rd.Get("internal_ip_address").(string)
		if len(newInternalIpAddress) != 0 {
			res, err := inst.Client.Subnet.CheckAvailableSubnetIp(ctx, newSubnetId, newInternalIpAddress)
			if err != nil {
				return diag.FromErr(err)
			}
			if *res.Result == false {
				return diag.Errorf("Not Available Internal Ip Address")
			}
		}
		_, err := inst.Client.VirtualServer.UpdateVirtualServerSubnetIp(ctx, rd.Id(), virtualserver.VirtualServerSubnetIpUpdateRequest{
			SubnetId:          newSubnetId,
			InternalIpAddress: newInternalIpAddress,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("contract_discount") {
		oldContractDiscount, newContractDiscount := rd.GetChange("contract_discount")
		strOldContractDiscount := oldContractDiscount.(string)
		strNewContractDiscount := newContractDiscount.(string)
		if strings.Compare(strOldContractDiscount, "None") == 0 {
			_, err := inst.Client.VirtualServer.UpdateVirtualServerContract(ctx, rd.Id(), virtualserver.VirtualServerContractUpdateRequest{
				ContractDiscount: strNewContractDiscount,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.Errorf("Once the Contract Discount created, it can be changed after the contract discount period expires.")
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("next_contract_discount") && strings.Compare(rd.Get("contract_discount").(string), "None") != 0 {
		_, err := inst.Client.VirtualServer.UpdateVirtualServerNextContract(ctx, rd.Id(), virtualserver.VirtualServerContractUpdateRequest{
			ContractDiscount: rd.Get("next_contract_discount").(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("external_storage") {
		oldExternalStorages, newExternalStorages := rd.GetChange("external_storage")
		tempStrOldExternalStorages := make([]map[string]interface{}, 0)
		tempStrNewExternalStorages := make([]map[string]interface{}, 0)
		for _, oldExternalStorage := range oldExternalStorages.([]interface{}) {
			tempStrOldExternalStorages = append(tempStrOldExternalStorages, oldExternalStorage.(map[string]interface{}))
		}
		for _, newExternalStorage := range newExternalStorages.([]interface{}) {
			tempStrNewExternalStorages = append(tempStrNewExternalStorages, newExternalStorage.(map[string]interface{}))
		}
		//oldExternalStorageList := getExternalStorageStructArray(oldExternalStorages.([]map[string]interface{}))
		//newExternalStorageList := getExternalStorageStructArray(newExternalStorages.([]map[string]interface{}))
		oldExternalStorageList := getExternalStorageStructArray(tempStrOldExternalStorages)
		newExternalStorageList := getExternalStorageStructArray(tempStrNewExternalStorages)
		externalStorageListWithChangedSize := getExternalStorageListWithChangedSize(oldExternalStorageList, newExternalStorageList)

		addedExternalStorageList := getDiff(oldExternalStorageList, newExternalStorageList)
		deletedExternalStorageList := getDiff(newExternalStorageList, oldExternalStorageList)

		for _, externalStorageWithChangeSize := range externalStorageListWithChangedSize {
			inst.Client.BlockStorage.ResizeBlockStorage(ctx, blockstorage.UpdateBlockStorageRequest{
				BlockStorageId:   externalStorageWithChangeSize.BlockStorageId,
				BlockStorageSize: externalStorageWithChangeSize.StorageSizeGb,
				ProductId:        externalStorageWithChangeSize.ProductId,
			})
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}

		productGroup, err := inst.Client.Product.GetProductGroup(ctx, virtualServerInfo.ProductGroupId)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, addedExternalStorage := range addedExternalStorageList {
			_, err = inst.Client.BlockStorage.CreateBlockStorage(ctx, getCreatedBlockStorageRequest(addedExternalStorage, productGroup, rd.Id()))
			if err != nil {
				return diag.FromErr(err)
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
		for _, deletedExternalStorage := range deletedExternalStorageList {
			if strings.Compare(deletedExternalStorage.SharedType, "SHARED") == 0 {
				_, err := inst.Client.BlockStorage.DetachBlockStorage(ctx, deletedExternalStorage.BlockStorageId, blockstorage.BlockStorageDetachRequest{
					VirtualServerId: rd.Id(),
				})
				if err != nil {
					return diag.FromErr(err)
				}
			} else {
				_, err := inst.Client.BlockStorage.DeleteBlockStorage(ctx, deletedExternalStorage.BlockStorageId)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
	}

	if rd.HasChanges("state") {
		var VmState string
		if strings.Compare(strings.ToUpper(rd.Get("state").(string)), "STOPPED") == 0 {
			_, err = inst.Client.VirtualServer.StopVirtualServer(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			VmState = common.StoppedState
		}

		if strings.Compare(strings.ToUpper(rd.Get("state").(string)), "RUNNING") == 0 {
			_, err = inst.Client.VirtualServer.StartVirtualServer(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			VmState = common.RunningState
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{VmState}, true)
		if err != nil {
			return
		}
	}

	return resourceVirtualServerRead(ctx, rd, meta)
}

func detachAndAttachPublicIpId(ctx context.Context, rd *schema.ResourceData, inst *client.Instance, virtualServerId string, nicId string, natEnabled bool) error {
	_, n := rd.GetChange("public_ip_id")
	newPublicId := n.(string)

	_, _, err := inst.Client.VirtualServer.DetachPublicIp(ctx, virtualServerId, nicId)
	// skip 'Nic's nat is empty.' error
	const ErrorCodeNicNatIsEmpty = "PRODUCT-VIRTUALSERVER-INTERNAL-00022"
	if err != nil && !strings.Contains(err.Error(), ErrorCodeNicNatIsEmpty) {
		return err
	}
	err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return err
	}

	if natEnabled == true || len(newPublicId) > 0 {
		_, _, err = inst.Client.VirtualServer.AttachPublicIp(ctx, virtualServerId, nicId, newPublicId)
		if err != nil {
			return err
		}
		err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCreatedBlockStorageRequest(addedExternalStorage virtualserver.ExternalStorage, productGroup product.ProductGroupDetailResponse, virtualServerId string) blockstorage.CreateBlockStorageRequest {
	createBlockStorageRequest := blockstorage.CreateBlockStorageRequest{}
	createBlockStorageRequest.Tags = make([]blockstorage.TagRequest, 0)
	for _, tagRequest := range addedExternalStorage.Tags {
		createBlockStorageRequest.Tags = append(createBlockStorageRequest.Tags, blockstorage.TagRequest{
			TagKey:   tagRequest.TagKey,
			TagValue: tagRequest.TagValue,
		})
	}
	createBlockStorageRequest.BlockStorageName = addedExternalStorage.BlockStorageName
	createBlockStorageRequest.BlockStorageSize = addedExternalStorage.StorageSizeGb
	createBlockStorageRequest.EncryptEnabled = addedExternalStorage.EncryptEnabled
	createBlockStorageRequest.ProductId = getProductId(productGroup, addedExternalStorage.ProductName)
	createBlockStorageRequest.SharedType = addedExternalStorage.SharedType
	createBlockStorageRequest.VirtualServerId = virtualServerId
	return createBlockStorageRequest
}

func getProductId(productGroup product.ProductGroupDetailResponse, productName string) string {
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			if strings.Compare(productName, productInfo.ProductName) == 0 {
				return productInfo.ProductId
			}
		}
	}
	return ""
}

func getDiff(orgList []virtualserver.ExternalStorage, destList []virtualserver.ExternalStorage) []virtualserver.ExternalStorage {
	diffEntities := make([]virtualserver.ExternalStorage, 0)
	for _, dest := range destList {
		var i int
		for i = 0; i < len(orgList); i++ {
			if strings.Compare(dest.BlockStorageId, orgList[i].BlockStorageId) == 0 {
				break
			}
		}
		if i == len(orgList) {
			diffEntities = append(diffEntities, dest)
		}
	}
	return diffEntities
}

func getExternalStorageListWithChangedSize(oldExternalStorageList []virtualserver.ExternalStorage, newExternalStorageList []virtualserver.ExternalStorage) []virtualserver.ExternalStorage {
	externalStorageListWithChangedSize := make([]virtualserver.ExternalStorage, 0)
	for _, oldExternalStorage := range oldExternalStorageList {
		for _, newExternalStorage := range newExternalStorageList {
			if strings.Compare(oldExternalStorage.BlockStorageId, newExternalStorage.BlockStorageId) == 0 && oldExternalStorage.StorageSizeGb != newExternalStorage.StorageSizeGb {
				externalStorageListWithChangedSize = append(externalStorageListWithChangedSize, virtualserver.ExternalStorage{
					BlockStorageId: newExternalStorage.BlockStorageId,
					StorageSizeGb:  newExternalStorage.StorageSizeGb,
					ProductId:      newExternalStorage.ProductId,
				})
			}
		}
	}
	return externalStorageListWithChangedSize
}

func getExternalStorageStructArray(mapExternalStorages []map[string]interface{}) []virtualserver.ExternalStorage {
	externalStorageList := make([]virtualserver.ExternalStorage, 0)
	externalStorage := virtualserver.ExternalStorage{}
	for _, mapOldExternalStorage := range mapExternalStorages {
		//var strBlockStorageId string
		//var intStorageSizeGb int32
		//var boolEncryptEnabled bool
		//var strProductName string
		//var strSharedType string
		//var strProductId string
		//var strBlockStorageName string

		if blockStorageId, ok := mapOldExternalStorage["block_storage_id"]; ok {
			//strBlockStorageId = blockStorageId.(string)
			externalStorage.BlockStorageId = blockStorageId.(string)
		}
		if storageSizeGb, ok := mapOldExternalStorage["storage_size_gb"]; ok {
			//intStorageSizeGb = int32(storageSizeGb.(int))
			externalStorage.StorageSizeGb = int32(storageSizeGb.(int))
		}
		if encryptEnabled, ok := mapOldExternalStorage["encrypted"]; ok {
			//boolEncryptEnabled = encryptEnabled.(bool)
			externalStorage.EncryptEnabled = encryptEnabled.(bool)
		}
		if productName, ok := mapOldExternalStorage["product_name"]; ok {
			//strProductName = productName.(string)
			externalStorage.ProductName = productName.(string)
		}
		if sharedType, ok := mapOldExternalStorage["shared_type"]; ok {
			//strSharedType = sharedType.(string)
			externalStorage.SharedType = sharedType.(string)
		}
		if productId, ok := mapOldExternalStorage["product_id"]; ok {
			//strProductId = productId.(string)
			externalStorage.ProductId = productId.(string)
		}

		if name, ok := mapOldExternalStorage["name"]; ok {
			//strBlockStorageName = name.(string)
			externalStorage.BlockStorageName = name.(string)
		}

		tags := make([]virtualserver.TagRequest, 0)
		//tagsMap := mapOldExternalStorage["tags"].(map[string]interface{})
		if tagsMap, ok := mapOldExternalStorage["tags"]; ok {
			for key, value := range tagsMap.(map[string]interface{}) {
				tags = append(tags, virtualserver.TagRequest{
					TagKey:   key,
					TagValue: value.(string),
				})
			}
			log.Println("externalStorage.Tags : ", tags)
			externalStorage.Tags = tags
		}

		//externalStorageList = append(externalStorageList, virtualserver.ExternalStorage{
		//	BlockStorageId:   strBlockStorageId,
		//	StorageSizeGb:    intStorageSizeGb,
		//	EncryptEnabled:   boolEncryptEnabled,
		//	ProductName:      strProductName,
		//	SharedType:       strSharedType,
		//	ProductId:        strProductId,
		//	BlockStorageName: strBlockStorageName,
		//})
		externalStorageList = append(externalStorageList, externalStorage)
	}
	return externalStorageList
}

//func getOldAndNewInternalIpAddress(data *schema.ResourceData) (string, string) {
//	oldValue, newValue := data.GetChange("internal_ip_address")
//	oldInternalIpAddress := oldValue.(string)
//	newInternalIpAddress := newValue.(string)
//
//	return oldInternalIpAddress, newInternalIpAddress
//}

func resourceVirtualServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	error := WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState, common.ErrorState}, false)
	if error != nil {
		return diag.FromErr(error)
	}

	_, err := inst.Client.VirtualServer.DeleteVirtualServer(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func WaitForVirtualServerStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.VirtualServer.GetVirtualServer(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, info.VirtualServerState, nil
	})
}
