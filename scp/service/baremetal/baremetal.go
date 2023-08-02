package baremetal

import (
	"context"
	"errors"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/storage/bmblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/service/image"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/public-ip2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_bm_server", ResourceBareMetalServer())
}

func ResourceBareMetalServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBareMetalServerCreate,
		ReadContext:   resourceBareMetalServerRead,
		UpdateContext: resourceBareMetalServerUpdate,
		DeleteContext: resourceBareMetalServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bm_server_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateName3to28AlphaDashStartsWithLowerCase,
				Description:      "Bare-metal server name",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable delete protection for this bare-metal server",
			},
			"cpu_count": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "CPU core count(8, 16, ..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"memory_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Memory size in gigabytes(16, 32,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"contract_discount": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contract : None, 1-year, 3-year",
			},
			"block_storages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "block storages",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Block storage name",
						},
						"product_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "SSD",
							Description: "Storage product name : SSD",
						},
						"storage_size_gb": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Storage size in gigabytes",
						},
						"encrypted": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Use encryption for this storage",
						},
					},
				},
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id of this bare-metal server",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image id of this bare-metal server",
			},
			"initial_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Initialization script",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id of this bare-metal server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Public IP id of this bare-metal server. Public-IP must be a valid public-ip resource which is attached to the VPC.",
			},
			"use_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DNS feature for this bare-metal server.",
			},
			"use_hyper_threading": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "Enable hyper-threading feature for this bare-metal server.",
			},
			"admin_account": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateName3to20DashUnderscore,
				Description:      "Admin account for this bare-metal server OS. For linux, this must be 'root'. For Windows, this must not be 'administrator'.",
			},
			"admin_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: common.ValidatePassword8to20,
				Description:      "Admin account password for this bare-metal server OS. (CAUTION) The actual plain-text password will be sent to your email.",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address of this bare-metal server",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable NAT feature for this bare-metal server.",
			},
			"local_subnet_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable local subnet for this bare-metal server",
			},
			"local_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local Subnet id of this bare-metal server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"local_subnet_ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local IP address of this bare-metal server",
			},
		},
		Description: "Provides a Bare-metal Server resource.",
	}
}

// ConvertBlockStorageList : this routine is based on common.ConvertExternalStorageList
func ConvertBlockStorageList(list common.HclListObject, diskToId map[string]string) ([]baremetal.BMAdditionalBlockStorageCreateRequest, error) {
	var result []baremetal.BMAdditionalBlockStorageCreateRequest
	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		var info baremetal.BMAdditionalBlockStorageCreateRequest
		if v, ok := item["name"]; ok {
			info.BareMetalBlockStorageName = v.(string)
		}
		if v, ok := item["product_name"]; ok {
			info.BareMetalBlockStorageType = v.(string)
		}
		if v, ok := item["storage_size_gb"]; ok {
			info.BareMetalBlockStorageSize = (int32)(v.(int))
		}
		if v, ok := item["encrypted"]; ok {
			info.EncryptionEnabled = v.(bool)
		}
		if extProductId, ok := diskToId[info.BareMetalBlockStorageType]; ok {
			info.BareMetalBlockStorageTypeId = extProductId
		} else {
			return nil, errors.New("disk type is not matched")
		}
		result = append(result, info)
	}
	return result, nil
}

func resourceBareMetalServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	serverName := rd.Get("bm_server_name").(string)
	isDeleteProtected := rd.Get("delete_protection").(bool)
	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	contractDiscount := rd.Get("contract_discount").(string)
	blockStorageList := rd.Get("block_storages").(common.HclListObject)

	subnetId := rd.Get("subnet_id").(string)
	publicIpId := rd.Get("public_ip_id").(string)

	adminAccount := rd.Get("admin_account").(string)
	adminPassword := rd.Get("admin_password").(string)
	initialScript := rd.Get("initial_script").(string)

	vpcId := rd.Get("vpc_id").(string)
	imageId := rd.Get("image_id").(string)
	useDNS := rd.Get("use_dns").(bool)
	//ipAddr := rd.Get("ipv4").(string)
	natEnabled := rd.Get("nat_enabled").(bool)
	useHyperThreading := rd.Get("use_hyper_threading").(string)
	localSubnetEnabled := rd.Get("local_subnet_enabled").(bool)
	localSubnetId := rd.Get("local_subnet_id").(string)
	localSubnetIp := rd.Get("local_subnet_ipv4").(string)
	if !localSubnetEnabled {
		localSubnetId = "" // to bypass local-subnet data manipulation bug
	}

	if !localSubnetEnabled && localSubnetId != "" {
		return diag.Errorf("local subnet is disabled, but has subnet id")
	}
	// Get vpc info
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return
	}

	blockId := ""
	projectDetails, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, zone := range projectDetails.ServiceZones {
		if zone.ServiceZoneId == vpcInfo.ServiceZoneId {
			blockId = zone.BlockId
		}
	}
	if blockId == "" {
		return diag.FromErr(errors.New("vpc info is not valid"))
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

	diskProductNameToId := common.ProductToIdMap(common.ProductDisk, &productGroup)
	if len(diskProductNameToId) == 0 {
		diagnostics = diag.Errorf("Failed to find external disk product")
		return
	}

	contractToId := common.ProductToIdMap(common.ProductContractDiscount, &productGroup)
	if len(contractToId) == 0 {
		diagnostics = diag.Errorf("Failed to find contract info")
		return
	}

	// Find bare-metal scaling
	scaleId, err := client.FindScaleProduct(ctx, inst.Client, targetProductGroupId, cpuCount, memorySizeGB)
	if err != nil {
		return
	}

	// Find contract
	contractId, ok := contractToId[contractDiscount]
	if !ok {
		diagnostics = diag.Errorf("Invalid contract")
		return
	}

	blockStorageInfoList, err := ConvertBlockStorageList(blockStorageList, diskProductNameToId)
	if err != nil {
		return diag.FromErr(err)
	}

	serverDetails := baremetal.BMServerDetailsRequest{
		BareMetalLocalSubnetEnabled:   localSubnetEnabled,
		BareMetalLocalSubnetId:        localSubnetId,
		BareMetalLocalSubnetIpAddress: localSubnetIp,
		BareMetalServerName:           serverName,
		DnsEnabled:                    useDNS,
		//IpAddress:                     ipAddr,
		NatEnabled:        natEnabled,
		PublicIpAddressId: publicIpId,
		ServerTypeId:      scaleId,
		StorageDetails:    blockStorageInfoList,
		UseHyperThreading: useHyperThreading,
	}
	createRequest := baremetal.BMServerCreateRequest{
		BlockId:                   blockId,
		ContractId:                contractId,
		DeletionProtectionEnabled: isDeleteProtected,
		ImageId:                   imageId,
		InitScript:                initialScript,
		SubnetId:                  subnetId,
		OsUserId:                  adminAccount,
		OsUserPassword:            adminPassword,
		ProductGroupId:            targetProductGroupId,
		ServiceZoneId:             vpcInfo.ServiceZoneId,
		ServerDetails:             []baremetal.BMServerDetailsRequest{serverDetails},
		Tags:                      []baremetal.TagRequest{},
		VpcId:                     vpcId,
	}
	createResponse, err := inst.Client.BareMetal.CreateBareMetalServer(ctx, createRequest)
	if err != nil {
		return
	}

	err = WaitForBMServerStatus(ctx, inst.Client, createResponse.ResourceId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
	if err != nil {
		return
	}

	rd.SetId(createResponse.ResourceId)

	return resourceBareMetalServerRead(ctx, rd, meta)
}

func getImageInfo(ctx context.Context, serviceZoneId string, imageId string, meta interface{}) (bool, string, error) {
	var err error = nil
	inst := meta.(*client.Instance)

	isOsWindows := false
	var targetProductGroupId string

	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, serviceZoneId, image.ActiveState, common.ServicedGroupCompute, common.ServicedForBaremetalServer)
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
	return isOsWindows, targetProductGroupId, err
}

func resourceBareMetalServerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	bmServerInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, bmServerInfo.ProductGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	productIdToNameMapper := make(map[string]string)
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			productIdToNameMapper[productInfo.ProductId] = productInfo.ProductName
		}
	}

	rd.Set("image_id", bmServerInfo.ImageId)
	rd.Set("bm_server_name", bmServerInfo.BareMetalServerName)
	rd.Set("delete_protection", bmServerInfo.DeletionProtectionEnabled)
	rd.Set("contract_discount", bmServerInfo.Contract)
	rd.Set("vpc_id", bmServerInfo.VpcId)
	rd.Set("use_dns", bmServerInfo.DnsEnabled)
	rd.Set("initial_script", bmServerInfo.InitialScriptContent)

	blockStorages := common.HclListObject{}
	for _, blockId := range bmServerInfo.BareMetalBlockStorageIds {
		blockInfo, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, blockId)
		if err != nil {
			continue
		}
		blockStorageInfo := common.HclKeyValueObject{}
		blockStorageInfo["name"] = blockInfo.BareMetalBlockStorageName
		blockStorageInfo["storage_size_gb"] = int(blockInfo.BareMetalBlockStorageSize)
		blockStorageInfo["encrypted"] = blockInfo.EncryptionEnabled
		if blockProductName, ok := productIdToNameMapper[blockInfo.ProductId]; ok {
			blockStorageInfo["product_name"] = blockProductName
		} else {
			blockStorageInfo["product_name"] = "UNKNOWN"
		}
		blockStorages = append(blockStorages, blockStorageInfo)
	}
	rd.Set("block_storages", blockStorages)

	// Set cpu / memory
	scale, err := client.FindProductById(ctx, inst.Client, bmServerInfo.ProductGroupId, bmServerInfo.ServerTypeId)
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

	ipv4 := bmServerInfo.IpAddress
	localSubnetId := bmServerInfo.BareMetalLocalSubnetId
	localSubnetIp := bmServerInfo.BareMetalLocalSubnetIpAddress
	subnetId := bmServerInfo.SubnetId
	natIpv4 := bmServerInfo.NatIpAddress
	natEnabled := bmServerInfo.PublicNatStatus

	rd.Set("ipv4", ipv4)
	rd.Set("subnet_id", subnetId)
	rd.Set("local_subnet_id", localSubnetId)
	rd.Set("local_subnet_ipv4", localSubnetIp)
	rd.Set("local_subnet_enabled", localSubnetId != "")
	rd.Set("nat_ipv4", natIpv4)
	rd.Set("nat_enabled", natEnabled == "SUCCESS")

	if natIpv4 != "" {
		publicIpInfo, err := inst.Client.PublicIp.GetPublicIpList(ctx,
			bmServerInfo.ServiceZoneId, &publicip2.PublicIpOpenApiControllerApiListPublicIpsV2Opts{
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
			// this case is found on auto assign mode
			rd.Set("public_ip_id", "")
		} else {
			rd.Set("public_ip_id", publicIpInfo.Contents[0].PublicIpAddressId)
		}
	}
	return nil
}

func resourceBareMetalServerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	if !rd.HasChanges("delete_protection") && !rd.HasChanges("contract_discount") && !rd.HasChanges("local_subnet_enabled") && !rd.HasChanges("nat_enabled") && !rd.HasChanges("block_storages") {
		return diag.Errorf("nothing to update")
	}

	serverInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, rd.Id())
	if err != nil {
		return
	}
	targetProductGroupId := serverInfo.ProductGroupId

	if rd.HasChanges("delete_protection") {
		deleteProtection := rd.Get("delete_protection").(bool)
		isDeleteProtectionEnabled := "Y"
		if !deleteProtection {
			isDeleteProtectionEnabled = "N"
		}
		_, err = inst.Client.BareMetal.ChangeBMDeletePolicy(ctx, rd.Id(), isDeleteProtectionEnabled)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("contract_discount") {
		contractDiscount := rd.Get("contract_discount").(string)

		// Get product group information
		productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
		//productGroup, err := inst.Client.Product.GetProducesList(ctx, vpcInfo.ServiceZoneId, targetProductGroupId, "")
		if err != nil {
			return diag.FromErr(err)
		}

		contractToId := common.ProductToIdMap(common.ProductContractDiscount, &productGroup)
		if len(contractToId) == 0 {
			diagnostics = diag.Errorf("Failed to find contract info")
			return
		}

		contractId, ok := contractToId[contractDiscount]
		if !ok {
			diagnostics = diag.Errorf("Invalid contract")
			return
		}

		_, err = inst.Client.BareMetal.ChangeBMContract(ctx, rd.Id(), contractId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("local_subnet_enabled") {
		localSubnetEnabled := rd.Get("local_subnet_enabled").(bool)
		localSubnetIp := rd.Get("local_subnet_ipv4").(string)
		localSubnetId := rd.Get("local_subnet_id").(string)

		if localSubnetEnabled {
			_, err = inst.Client.BareMetal.AttachBMLocalSubnet(ctx, rd.Id(), localSubnetId, localSubnetIp)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err = inst.Client.BareMetal.DetachBMLocalSubnet(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// wait for server state change to editing //
		time.Sleep(3 * time.Second)

		err = WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("nat_enabled") {
		natEnabled := rd.Get("nat_enabled").(bool)
		publicIpId := rd.Get("public_ip_id").(string)

		addressType := "AUTO"
		if publicIpId != "" {
			addressType = "MANUAL"
		}
		if natEnabled {
			_, err = inst.Client.BareMetal.EnableBMNat(ctx, rd.Id(), addressType, publicIpId)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err = inst.Client.BareMetal.DisableBMNat(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		err = WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("block_storages") {
		storageList, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorages(ctx)

		// Get product group information
		productGroup, err := inst.Client.Product.GetProductGroup(ctx, targetProductGroupId)
		if err != nil {
			return diag.FromErr(err)
		}

		diskProductNameToId := common.ProductToIdMap(common.ProductDisk, &productGroup)
		if len(diskProductNameToId) == 0 {
			diagnostics = diag.Errorf("Failed to find external disk product")
			return
		}

		o, n := rd.GetChange("block_storages")
		oldVal := o.(common.HclListObject)
		newVal := n.(common.HclListObject)
		var oldList []baremetal.BMAdditionalBlockStorageCreateRequest
		var newList []baremetal.BMAdditionalBlockStorageCreateRequest
		oldList, err = ConvertBlockStorageList(oldVal, diskProductNameToId)
		newList, err = ConvertBlockStorageList(newVal, diskProductNameToId)

		if len(oldList) == len(newList) {
			return diag.Errorf("there is no change in number of storages")
		}

		for _, oldRequest := range oldList {
			isExist := false
			for _, newRequest := range newList {
				if oldRequest.BareMetalBlockStorageName == newRequest.BareMetalBlockStorageName {
					isExist = true
				}
			}
			if isExist {
				continue
			}
			storageId := ""
			for _, storageInfo := range storageList.Contents {
				if storageInfo.BareMetalBlockStorageName == oldRequest.BareMetalBlockStorageName {
					storageId = storageInfo.BareMetalBlockStorageId
				}
			}
			if storageId == "" {
				return diag.Errorf("storage name is not exist")
			}
			// if not exist, this old request has to be deleted

			_, _, err := inst.Client.BareMetalBlockStorage.DetachBareMetalBlockStorage(ctx, storageId, []string{rd.Id()})
			if err != nil {
				return diag.FromErr(err)
			}
			err = WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return diag.FromErr(err)
			}

			// if storage is orphan, delete storage
			result, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, storageId)
			if err != nil {
				return diag.FromErr(err)
			}
			if len(result.Servers) == 0 {
				_, _, _ = inst.Client.BareMetalBlockStorage.DeleteBareMetalBlockStorage(ctx, storageId)
			}
		}
		for _, newRequest := range newList {
			isExist := false
			for _, oldRequest := range oldList {
				if oldRequest.BareMetalBlockStorageName == newRequest.BareMetalBlockStorageName {
					isExist = true
				}
			}
			if isExist {
				continue
			}
			// Get vpc info
			vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, rd.Get("vpc_id").(string))
			if err != nil {
				return
			}

			request := bmblockstorage.BmBlockStorageCreateRequest{
				BareMetalBlockStorageName: newRequest.BareMetalBlockStorageName,
				BareMetalServerIds:        []string{rd.Id()},
				BareMetalBlockStorageSize: newRequest.BareMetalBlockStorageSize,
				EncryptionEnabled:         newRequest.EncryptionEnabled,
				IsSnapshotPolicy:          false,
				SnapshotCapacityRate:      0,
				ServiceZoneId:             vpcInfo.ServiceZoneId,
				SnapshotSchedule:          bmblockstorage.SnapshotSchedule{},
				Tags:                      []bmblockstorage.TagRequest{},
				ProductId:                 newRequest.BareMetalBlockStorageTypeId,
			}
			// if not exist, this new request has to be created
			_, _, err = inst.Client.BareMetalBlockStorage.CreateBareMetalBlockStorage(ctx, request)
			if err != nil {
				return diag.FromErr(err)
			}
			err = WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceBareMetalServerRead(ctx, rd, meta)
}

func resourceBareMetalServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	error := WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
	if error != nil {
		return diag.FromErr(error)
	}

	_, err := inst.Client.BareMetal.DeleteBareMetalServer(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = WaitForBMServerStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func WaitForBMServerStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.BareMetal.GetBareMetalServerDetail(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, strings.ToUpper(info.BareMetalServerState), nil
	})
}
