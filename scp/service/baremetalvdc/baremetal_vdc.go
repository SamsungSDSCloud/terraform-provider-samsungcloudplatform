package baremetalvdc

import (
	"context"
	"errors"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/baremetalvdc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	baremetal2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/image"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/product"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_bm_server_vdc", ResourceBareMetalServerVDC())
}

func ResourceBareMetalServerVDC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVxLanBareMetalServerCreate,
		ReadContext:   resourceVxLanBareMetalServerRead,
		UpdateContext: resourceVxLanBareMetalServerUpdate,
		DeleteContext: resourceVxLanBareMetalServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"block_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "BLOCK ID of this bare-metal server",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SERVICE ZONE ID of this bare-metal server",
			},
			"vdc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VDC ID of this bare-metal server",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id of this bare-metal server. Subnet must be a valid subnet resource which is attached to the VDC.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image name of this bare-metal server",
			},
			"contract_discount": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contract : None, 1 Year, 3 Year",
				ValidateDiagFunc: common.ValidateContract,
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable delete protection for this bare-metal server",
			},
			"admin_account": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateName3to20LowerAlphaAndNumberOnly,
				Description:      "Admin account for this bare-metal server OS. For linux, this must be 'root'. For Windows, this must not be 'administrator'.",
			},
			"admin_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: common.ValidatePassword8to20,
				Description:      "Admin account password for this bare-metal server OS. (CAUTION) The actual plain-text password will be sent to your email.",
			},
			"initial_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Initialization script",
			},
			"cpu_count": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "CPU core count(8, 16, ..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"memory_size": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Memory size in gigabytes(16, 32,..)",
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"servers": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Server name",
							ValidateDiagFunc: common.ValidateName3to28AlphaNumberDash,
						},
						"dns_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable DNS feature for this bare-metal server.",
						},
						"ip_address": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							ValidateDiagFunc: common.ValidateIpv4WithEmptyValue,
							Description:      "IP address of this bare-metal server",
						},
						"use_hyper_threading": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "N",
							Description: "Enable hyper-threading feature for this bare-metal server.(ex Y, N)",
						},
					},
				},
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
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Bare-metal Server(VDC) resource.",
	}
}

func resourceVxLanBareMetalServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vdcId := rd.Get("vdc_id").(string)
	imageName := rd.Get("image_name").(string)

	servers := rd.Get("servers").(common.HclListObject)
	blockStorageList := rd.Get("block_storages").(common.HclListObject)

	cpuCount := rd.Get("cpu_count").(int)
	memorySize := rd.Get("memory_size").(int)
	contractDiscount := rd.Get("contract_discount").(string)

	deleteProtection := rd.Get("delete_protection").(bool)
	adminAccount := rd.Get("admin_account").(string)
	adminPassword := rd.Get("admin_password").(string)
	initialScript := rd.Get("initial_script").(string)

	serviceZoneId := rd.Get("service_zone_id").(string)
	blockId := rd.Get("block_id").(string)

	// Get imageId
	isOsWindows, imageId, targetProductGroupId, err := getImageInfo(ctx, serviceZoneId, imageName, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(imageId) == 0 {
		diagnostics = diag.Errorf("Image id not found from imageName")
		return
	}

	if len(targetProductGroupId) == 0 {
		diagnostics = diag.Errorf("Product group id not found from image")
		return
	}

	// validate adminAccount
	if !isOsWindows && adminAccount != common.LinuxAdminAccount {
		adminAccount = common.LinuxAdminAccount
		log.Println("Linux admin account must be root")
	}

	if isOsWindows && (adminAccount == common.WindowsAdminAccount || len(adminAccount) < 5) {
		diagnostics = diag.Errorf("Windows admin account must be 5 to 20 alpha-numeric characters with special character and not be 'administrator'.")
		return
	}

	// window 서버일 때, 이름 길이 확인
	if isOsWindows {
		err := validateServerNameInWindowImage(servers)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get product information
	productGroup, err := inst.Client.Product.GetProducesList(ctx, serviceZoneId, targetProductGroupId, "")
	if err != nil {
		return diag.FromErr(err)
	}

	diskProductNameToId := ProductGroupDetailToIdMap(productToTypeList(common.ProductDisk, &productGroup))
	if len(diskProductNameToId) == 0 {
		diagnostics = diag.Errorf("Failed to find external disk product")
		return
	}

	// Find contract
	contractToId := productToTypeList(common.ProductContractDiscount, &productGroup)
	if len(contractToId) == 0 {
		diagnostics = diag.Errorf("Failed to find contract info")
		return
	}

	var contractId string

	for _, c := range contractToId {
		if c.ProductName == contractDiscount {
			contractId = c.ProductId
			break
		}
	}

	if len(contractId) == 0 {
		diagnostics = diag.Errorf("Invalid contract")
		return
	}

	// Find bare-metal scaling
	scaleToId := productToTypeList(common.ProductScale, &productGroup)
	if len(scaleToId) == 0 {
		diagnostics = diag.Errorf("Failed to find scale info")
		return
	}

	var scaleId string

	for _, s := range scaleToId {
		cpu, memory, err := client.FindScaleInfo(ctx, inst.Client, targetProductGroupId, s.ProductId)
		if err != nil {
			return
		}

		if cpuCount == cpu && memorySize == memory {
			scaleId = s.ProductId
			break
		}
	}

	if len(scaleId) == 0 {
		diagnostics = diag.Errorf("Invalid scale")
		return
	}

	// Get BS(임시 disable 처리)
	_, err = baremetal2.ConvertBlockStorageList(blockStorageList, diskProductNameToId)
	if err != nil {
		return diag.FromErr(err)
	}

	serverDetails := ConvertServerList(servers)

	for idx, _ := range serverDetails {
		serverDetails[idx].ServerTypeId = scaleId
		//serverDetails[idx].StorageDetails = blockStorageInfoList
	}

	createRequest := baremetalvdc.BMVDCServerCreateRequest{
		BlockId:                   blockId,
		ContractId:                contractId,
		DeletionProtectionEnabled: deleteProtection,
		ImageId:                   imageId,
		InitScript:                initialScript,
		OsUserId:                  adminAccount,
		OsUserPassword:            adminPassword,
		ProductGroupId:            targetProductGroupId,
		ServerDetails:             serverDetails,
		ServiceZoneId:             serviceZoneId,
		SubnetId:                  subnetId,
		VdcId:                     vdcId,
	}

	createResponse, err := inst.Client.BareMetalVdc.CreateBareMetalServerVDC(ctx, createRequest, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		return
	}

	resourceIds := strings.Split(createResponse.ResourceId, ",")

	for _, resourceId := range resourceIds {
		err = WaitForBMServerStatus(ctx, inst.Client, resourceId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			return
		}
	}

	rd.SetId(createResponse.ResourceId)

	return resourceVxLanBareMetalServerRead(ctx, rd, meta)
}

func getImageInfo(ctx context.Context, serviceZoneId string, imageName string, meta interface{}) (bool, string, string, error) {
	var err error = nil
	inst := meta.(*client.Instance)

	isOsWindows := false
	var targetProductGroupId string
	var imageId string

	standardImages, err := inst.Client.Image.GetStandardImageList(ctx, serviceZoneId, image.ActiveState, common.ServicedGroupCompute, common.ServicedForBaremetalServer)
	if err != nil {
		return false, "", "", err
	}

	for _, c := range standardImages.Contents {
		if result := strings.Contains(c.ImageName, imageName); result {
			targetProductGroupId = c.ProductGroupId
			imageId = c.ImageId
			if c.OsType == common.OsTypeWindows {
				isOsWindows = true
			}
		}
	}
	return isOsWindows, imageId, targetProductGroupId, err
}

func productToTypeList(productType string, products *product.ListResponseV2ProductsResponse1) []product.ProductsResponse1 {
	result := make([]product.ProductsResponse1, 0)

	for _, p := range products.Contents {
		if p.ProductType == productType {
			result = append(result, p)
		}
	}
	return result
}

func ProductGroupDetailToIdMap(productList []product.ProductsResponse1) map[string]string {
	result := make(map[string]string)

	for _, p := range productList {
		result[p.ProductName] = p.ProductId
	}
	return result
}

func ConvertServerList(list common.HclListObject) []baremetalvdc.BMServerDetailsRequest {
	result := make([]baremetalvdc.BMServerDetailsRequest, 0)

	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		var info baremetalvdc.BMServerDetailsRequest
		if v, ok := item["name"]; ok {
			info.BareMetalServerName = v.(string)
		}
		if v, ok := item["dns_enabled"]; ok {
			info.DnsEnabled = v.(bool)
		}
		if v, ok := item["ip_address"]; ok {
			info.IpAddress = v.(string)
		}
		if v, ok := item["use_hyper_threading"]; ok {
			info.UseHyperThreading = v.(string)
		}
		result = append(result, info)
	}

	return result
}

func validateServerNameInWindowImage(list common.HclListObject) error {

	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		v, ok := item["name"]

		if !ok {
			return errors.New("There is no name attribute")
		}

		if len(v.(string)) > 15 {
			return errors.New("Servers using Windows must be 3 to 15 characters long. ")
		}

	}

	return nil
}

func resourceVxLanBareMetalServerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	baremetalIds := strings.Split(rd.Id(), ",")

	bmVdcServerInfo, _, err := inst.Client.BareMetalVdc.GetBareMetalServerDetailVDC(ctx, baremetalIds[0])
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("block_id", bmVdcServerInfo.BlockId)
	rd.Set("service_zone_id", bmVdcServerInfo.ServiceZoneId)

	rd.Set("vdc_id", bmVdcServerInfo.VdcId)
	rd.Set("subnet_id", bmVdcServerInfo.SubnetId)
	rd.Set("image_name", strings.ToUpper(bmVdcServerInfo.OsType)+" "+bmVdcServerInfo.OsVersion)

	rd.Set("contract_discount", bmVdcServerInfo.Contract)
	rd.Set("delete_protection", bmVdcServerInfo.DeletionProtectionEnabled == "Y")
	rd.Set("initial_script", bmVdcServerInfo.InitialScriptContent)

	// Set cpu / memory
	scale, err := client.FindProductById(ctx, inst.Client, bmVdcServerInfo.ProductGroupId, bmVdcServerInfo.ServerTypeId)
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
			rd.Set("memory_size", memorySize)
			memoryFound = true
		}
	}
	if !cpuFound || !memoryFound {
		return
	}

	// Get product group information
	productGroup, err := inst.Client.Product.GetProductGroup(ctx, bmVdcServerInfo.ProductGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get disk information for BS
	productIdToNameMapper := make(map[string]string)
	if productInfos, ok := productGroup.Products[common.ProductDisk]; ok {
		for _, productInfo := range productInfos {
			productIdToNameMapper[productInfo.ProductId] = productInfo.ProductName
		}
	}

	blockStorages := common.HclListObject{}
	for _, blockId := range bmVdcServerInfo.BareMetalBlockStorageIds {
		blockInfo, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, blockId)
		if err != nil {
			continue
		}
		blockStorageInfo := common.HclKeyValueObject{}
		blockStorageInfo["name"] = blockInfo.BareMetalBlockStorageName
		blockStorageInfo["storage_size_gb"] = blockInfo.BareMetalBlockStorageSize
		blockStorageInfo["encrypted"] = blockInfo.EncryptionEnabled
		if blockProductName, ok := productIdToNameMapper[blockInfo.ProductId]; ok {
			blockStorageInfo["product_name"] = blockProductName
		} else {
			blockStorageInfo["product_name"] = "UNKNOWN"
		}
		blockStorages = append(blockStorages, blockStorageInfo)
	}
	rd.Set("block_storages", blockStorages)

	servers := common.HclListObject{}

	for _, baremetalId := range baremetalIds {
		bmVdcServerInfo, _, err := inst.Client.BareMetalVdc.GetBareMetalServerDetailVDC(ctx, baremetalId)
		if err != nil {
			rd.SetId("")
			if common.IsDeleted(err) {
				return nil
			}
			return diag.FromErr(err)
		}
		serverInfo := common.HclKeyValueObject{}

		serverInfo["name"] = bmVdcServerInfo.BareMetalServerName
		serverInfo["ip_address"] = bmVdcServerInfo.IpAddress
		serverInfo["use_hyper_threading"] = bmVdcServerInfo.UseHyperThreading
		serverInfo["dns_enabled"] = bmVdcServerInfo.DnsEnabled == "Y"

		servers = append(servers, serverInfo)
	}

	rd.Set("servers", servers)

	tfTags.SetTags(ctx, rd, meta, baremetalIds[0])

	return nil
}

// TODO: 추후 구현
func resourceVxLanBareMetalServerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	return resourceVxLanBareMetalServerRead(ctx, rd, meta)
}

func resourceVxLanBareMetalServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	baremetalIds := strings.Split(rd.Id(), ",")

	for _, baremetalId := range baremetalIds {
		err := WaitForBMServerStatus(ctx, inst.Client, baremetalId, common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err := inst.Client.BareMetalVdc.DeleteBareMetalServersVDC(ctx, baremetalIds)
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	for _, baremetalId := range baremetalIds {
		err = WaitForBMServerStatus(ctx, inst.Client, baremetalId, common.VirtualServerProcessingStates(), []string{common.DeletedState}, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func WaitForBMServerStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.BareMetalVdc.GetBareMetalServerDetailVDC(ctx, id)
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
