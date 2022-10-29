package virtualserver

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/service/image"
	blockstorage2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/block-storage2"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/public-ip2"
	virtualserver2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/virtual-server2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Actual virtual-server name",
			},
			"name_prefix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "VirtualServer name prefix",
				ValidateDiagFunc: common.ValidateName3to20NoSpecials,
			},
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contract : None, 1-year, 3-year",
			},
			"external_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "External block storage.",
				Elem:        common.ExternalStorageResourceSchema(),
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server timezone",
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
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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
				Default:          "Write administrator's account",
				ValidateDiagFunc: common.ValidateName3to20DashUnderscore,
				Description:      "Admin account for this virtual server OS. For linux, this must be 'root'. For Windows, this must not be 'administrator'.",
			},
			"admin_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: common.ValidatePassword8to20,
				Description:      "Admin account password for this virtual server OS. (CAUTION) The actual plain-text password will be sent to your email.",
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

	namePrefix := rd.Get("name_prefix").(string)
	isDeleteProtected := rd.Get("delete_protection").(bool)
	timezone := rd.Get("timezone").(string)
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

	antiAffinity := rd.Get("anti_affinity").(bool)

	vpcId := rd.Get("vpc_id").(string)

	imageId := rd.Get("image_id").(string)

	useDNS := rd.Get("use_dns").(bool)

	// Get vpc info
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}

	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, vpcInfo.ServiceZoneId, image.ActiveState, common.ServicedGroupCompute, common.ServicedForVirtualServer)
	if err != nil {
		return
	}

	isOsWindows := false
	var targetProductGroupId string
	for _, c := range standardImages.Contents {
		if c.ImageId == (imageId + "dr") { // code for "stage" environment
			imageId = imageId + "dr"
		}
		if c.ImageId == imageId {
			targetProductGroupId = c.ProductGroupId
			if c.OsType == common.OsTypeWindows {
				isOsWindows = true
			}
		}
	}

	if isOsWindows && adminAccount != common.LinuxAdminAccount {
		diagnostics = diag.Errorf("Linux admin account must be root")
		return
	}

	if len(targetProductGroupId) == 0 {
		diagnostics = diag.Errorf("Product group id not found from image")
		return
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find OS disk product
	osDiskProductId, err := common.FirstProductId(common.ProductDefaultDisk, &productGroup)
	if err != nil {
		return
	}

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

	contractToId := common.ProductToIdMap(common.ProductContractDiscount, &productGroup)
	if len(contractToId) == 0 {
		diagnostics = diag.Errorf("Failed to find contract info")
		return
	}

	// Find VM scaling
	scaleId, err := client.FindScaleProduct(ctx, inst.Client, targetProductGroupId, cpuCount, memorySizeGB)
	if err != nil {
		return
	}

	// Find service level
	serviceLevelId, ok := serviceLevelToId["None"]
	if !ok {
		diagnostics = diag.Errorf("Invalid service level")
		return
	}

	// Find contract
	contractId, ok := contractToId[contractDiscount]
	if !ok {
		diagnostics = diag.Errorf("Invalid contract")
		return
	}

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

	var extStorages []virtualserver.BlockStorageInfo
	externalStorageInfoList := common.ConvertExternalStorageList(externalStorageList)
	for _, extStorageInfo := range externalStorageInfoList {
		if extProductId, ok := externalDiskProductNameToId[extStorageInfo.ProductName]; ok {
			extStorages = append(extStorages, virtualserver.BlockStorageInfo{
				BlockStorageName: extStorageInfo.Name,
				DiskSize:         int32(extStorageInfo.StorageSize),
				EncryptEnabled:   extStorageInfo.Encrypted,
				ProductId:        extProductId,
			})
		}
	}

	initialScriptShell := "bash"
	if isOsWindows {
		initialScriptShell = "powershell"
	}
	initialScript := virtualserver.InitialScriptInfo{
		EncodingType:         "plain",
		InitialScriptContent: rd.Get("initial_script_content").(string),
		InitialScriptShell:   initialScriptShell,
		InitialScriptType:    "text",
	}

	var name string
	for i := 1; ; i++ {
		name = fmt.Sprintf("%s-%03d", namePrefix, i)
		namedServerList, _, err := inst.Client.VirtualServer.GetVirtualServerList(ctx, name)
		if err != nil {
			continue
		}
		if len(namedServerList.Contents) == 0 {
			break
		}
	}

	createRequest := virtualserver.CreateRequest{
		BlockStorage: virtualserver.BlockStorageInfo{
			BlockStorageName: osStorageName,
			DiskSize:         int32(osStorageSize),
			EncryptEnabled:   osStorageEncrypted,
			ProductId:        osDiskProductId,
		},
		ContractId:                contractId,
		DeletionProtectionEnabled: isDeleteProtected,
		DnsEnabled:                useDNS,
		ExtraBlockStorages:        extStorages,
		ImageId:                   imageId,
		InitialScript:             initialScript,
		LocalSubnet:               virtualserver.LocalSubnetInfo{},
		Nic: virtualserver.NicInfo{
			InternalIpAddress: "",
			NatEnabled:        useNAT,
			PublicIpAddressId: publicIpId,
			SubnetId:          subnetId,
		},
		OsAdmin: virtualserver.OsAdminInfo{
			OsUserId:       adminAccount,
			OsUserPassword: adminPassword,
		},
		ProductGroupId:    targetProductGroupId,
		SecurityGroupIds:  getSecurityGroupIds(rd),
		ServerGroupId:     serverGroupId,
		ServerTypeId:      scaleId,
		ServiceLevelId:    serviceLevelId,
		ServiceZoneId:     vpcInfo.ServiceZoneId,
		Timezone:          timezone,
		VirtualServerName: name,
	}
	createResponse, err := inst.Client.VirtualServer.CreateVirtualServer(ctx, createRequest)
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

	rd.SetId(createResponse.ResourceId)

	return resourceVirtualServerRead(ctx, rd, meta)
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
		return
	}

	nicInfo, err := inst.Client.VirtualServer.GetNicList(ctx, virtualServerInfo.VirtualServerId)
	if err != nil {
		return
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
	rd.Set("name", virtualServerInfo.VirtualServerName)
	rd.Set("delete_protection", virtualServerInfo.DeletionProtectionEnabled)
	rd.Set("timezone", virtualServerInfo.Timezone)
	rd.Set("service_level", virtualServerInfo.ServiceLevel)
	rd.Set("contract_discount", virtualServerInfo.Contract)
	rd.Set("vpc_id", virtualServerInfo.VpcId)
	rd.Set("use_dns", virtualServerInfo.DnsEnabled)
	rd.Set("initial_script_content", virtualServerInfo.InitialScriptContent)

	sgIds := common.HclListObject{}
	for _, sg := range virtualServerInfo.SecurityGroupIds {
		sgIds = append(sgIds, sg.SecurityGroupId)
	}
	rd.Set("security_group_ids", sgIds)

	extStorages := common.HclListObject{}
	for _, blockId := range virtualServerInfo.BlockStorageIds {

		var blockInfo blockstorage2.BlockStorageResponse
		blockInfo, _, err = inst.Client.BlockStorage.ReadBlockStorage(ctx, blockId)
		if err != nil {
			continue
		}
		if blockInfo.IsBootDisk {
			rd.Set("os_storage_name", blockInfo.BlockStorageName)
			rd.Set("os_storage_size_gb", int(blockInfo.BlockStorageSize))
			rd.Set("os_storage_encrypted", blockInfo.EncryptEnabled)
		} else {
			extStorageInfo := common.HclKeyValueObject{}
			extStorageInfo["name"] = blockInfo.BlockStorageName
			extStorageInfo["storage_size_gb"] = int(blockInfo.BlockStorageSize)
			extStorageInfo["encrypted"] = blockInfo.EncryptEnabled
			if blockProductName, ok := productIdToNameMapper[blockInfo.ProductId]; ok {
				extStorageInfo["product_name"] = blockProductName
			} else {
				extStorageInfo["product_name"] = "UNKNOWN"
			}
			extStorages = append(extStorages, extStorageInfo)
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
			if nic.SubnetType == "VM" {
				if nic.NicId == vsNicId {
					localSubnetInfos = append(localSubnetInfos, common.HclKeyValueObject{
						"id":        nic.NicId,
						"subnet_id": nic.SubnetId,
						"ipv4":      nic.Ip,
					})
					break
				}
			}
			if nic.SubnetType == "PUBLIC" {
				if nic.NicId == vsNicId {
					subnetId = nic.SubnetId
					//subnetIpv4 = nic.Ip
					natIpv4 = nic.NatIp
					break
				}
			}
		}
	}

	rd.Set("ipv4", ipv4)
	rd.Set("subnet_id", subnetId)
	rd.Set("local_subnet", localSubnetInfos)
	rd.Set("nat_ipv4", natIpv4)

	if natIpv4 != "" {
		publicIpInfo, err := inst.Client.PublicIp.GetPublicIpList(ctx,
			virtualServerInfo.ServiceZoneId, &publicip2.PublicIpOpenApiControllerApiListPublicIpsV21Opts{
				IpAddress:       optional.NewString(natIpv4),
				IsBillable:      optional.Bool{},
				IsViewable:      optional.Bool{},
				PublicIpPurpose: optional.String{},
				PublicIpState:   optional.String{},
				UplinkType:      optional.String{},
				CreatedBy:       optional.String{},
				Page:            optional.Int32{},
				Size:            optional.Int32{},
				Sort:            optional.Interface{},
			})
		if err != nil {
			diagnostics = diag.FromErr(err)
			return
		}

		if len(publicIpInfo.Contents) == 0 {
			diagnostics = diag.Errorf("public ip information not found")
			return
		}

		rd.Set("public_ip_id", publicIpInfo.Contents[0].PublicIpAddressId)
	}

	return nil
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
		var scaleId string
		scaleId, err = client.FindScaleProduct(ctx, inst.Client, targetProductGroupId, cpuCount, memorySizeGB)
		if err != nil {
			return
		}

		// Update scale
		_, err = inst.Client.VirtualServer.UpdateScale(ctx, virtualServerInfo.VirtualServerId, scaleId)
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

	if rd.HasChanges("local_subnet_ids") {

		var nicInfoList virtualserver2.ListResponseOfNicResponse
		nicInfoList, err = inst.Client.VirtualServer.GetNicList(ctx, rd.Id())
		if err != nil {
			return
		}

		mapSubnetId2NicId := make(map[string]string)
		for _, nic := range nicInfoList.Contents {
			mapSubnetId2NicId[nic.SubnetId] = nic.NicId
		}

		o, n := rd.GetChange("local_subnet_ids")
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

	if rd.HasChanges("public_ip_id") {

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

		o, n := rd.GetChange("public_ip_id")
		oldPublicId := o.(string)
		newPublicId := n.(string)

		if len(oldPublicId) > 0 {
			_, _, err = inst.Client.VirtualServer.DetachPublicIp(ctx, virtualServerInfo.VirtualServerId, nicId)
			if err != nil {
				return
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
		if len(newPublicId) > 0 {
			_, _, err = inst.Client.VirtualServer.AttachPublicIp(ctx, virtualServerInfo.VirtualServerId, nicId, newPublicId)
			if err != nil {
				return
			}
			err = WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return
			}
		}
	}

	return resourceVirtualServerRead(ctx, rd, meta)
}

func resourceVirtualServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	error := WaitForVirtualServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
	if error != nil {
		return diag.FromErr(error)
	}

	_, err := inst.Client.VirtualServer.DeleteVirtualServer(ctx, rd.Id())
	if err != nil {
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
