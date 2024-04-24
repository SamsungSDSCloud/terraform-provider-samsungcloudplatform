package baremetal

import (
	"context"
	"errors"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/image"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/public-ip2"
	"github.com/antihax/optional"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"regexp"
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
			"servers": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bm_server_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateName3to28AlphaNumberDash,
							Description:      "Bare-metal server name",
						},
						"ipv4": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							ValidateDiagFunc: common.ValidateIpv4WithEmptyValue,
							Description:      "IP address of this bare-metal server",
						},
						"nat_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable NAT feature for this bare-metal server.",
						},
						"public_ip_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Public IP id of this bare-metal server. Public-IP must be a valid public-ip resource which is attached to the VPC.",
						},
						"local_subnet_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable local subnet for this bare-metal server",
						},
						"local_subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Local Subnet id of this bare-metal server. Subnet must be a valid subnet resource which is attached to the VPC.",
						},
						"local_subnet_ipv4": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							ValidateDiagFunc: common.ValidateIpv4WithEmptyValue,
							Description:      "Local IP address of this bare-metal server",
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
							Default:     "N",
							Description: "Enable hyper-threading feature for this bare-metal server.(ex. Y, N)",
						},
						"state": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateBmState,
							Description:      "Baremetal Server State(ex. RUNNING, STOPPED)",
						},
					},
				},
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
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Bare-metal Server resource.",
	}
}

func checkStringLength(str string, min int, max int) error {
	if len(str) < min {
		return fmt.Errorf("input must be longer than %v characters", min)
	} else if len(str) > max {
		return fmt.Errorf("input must be shorter than %v characters", max)
	} else {
		return nil
	}
}

func validateName3to24LowerAlphaDashStartsWithLowerCase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 24)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z][a-z0-9\\-]+[a-z0-9]$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must start with lower case character and contain only lower case, -, numerical characters and end with lower case, numerical character", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validateBmState(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	if !regexp.MustCompile("^(RUNNING|STOPPED)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Only RUNNING or STOPPED value of Attribute %q is allowed ", attrKey),
			AttributePath: path,
		})
	}

	return diags
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

func ConvertBaremetalServerList(list common.HclListObject, storageList []baremetal.BMAdditionalBlockStorageCreateRequest, scaleId string) ([]baremetal.BMServerDetailsRequest, error) {
	var result []baremetal.BMServerDetailsRequest
	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		var info baremetal.BMServerDetailsRequest
		if v, ok := item["bm_server_name"]; ok {
			info.BareMetalServerName = v.(string)
		}
		if v, ok := item["ipv4"]; ok {
			info.IpAddress = v.(string)
		}
		if v, ok := item["nat_enabled"]; ok {
			info.NatEnabled = v.(bool)
		}
		if v, ok := item["public_ip_id"]; ok {
			info.PublicIpAddressId = v.(string)
		}
		if v, ok := item["local_subnet_enabled"]; ok {
			info.BareMetalLocalSubnetEnabled = v.(bool)
		}
		if v, ok := item["local_subnet_id"]; ok {
			info.BareMetalLocalSubnetId = v.(string)
		}
		if v, ok := item["local_subnet_ipv4"]; ok {
			info.BareMetalLocalSubnetIpAddress = v.(string)
		}
		if v, ok := item["use_dns"]; ok {
			info.DnsEnabled = v.(bool)
		}
		if v, ok := item["use_hyper_threading"]; ok {
			info.UseHyperThreading = v.(string)
		}

		if v, ok := item["state"]; ok {
			if strings.Compare(strings.ToUpper(v.(string)), common.StoppedState) == 0 {
				return nil, errors.New("state value must be RUNNING")
			}
		}

		if !info.BareMetalLocalSubnetEnabled {
			info.BareMetalLocalSubnetId = ""
		}

		info.ServerTypeId = scaleId
		info.StorageDetails = storageList

		result = append(result, info)
	}
	return result, nil
}

func convertInterfaceToStringList(list common.HclListObject, key string) ([]string, error) {

	output := make([]string, 0)

	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		if v, ok := item[key]; ok {
			output = append(output, v.(string))
		} else {
			return nil, fmt.Errorf("%q value must be exist", key)
		}
	}

	return output, nil
}

// servers block에서 update 시에 달라진 index와 새로운 값을 반환시켜준다.
func getChangeDiffIndex(rd *schema.ResourceData, key string) ([]int, []string, error) {
	o, n := rd.GetChange("servers")

	oldValue, err := convertInterfaceToStringList(o.(common.HclListObject), key)
	if err != nil {
		return nil, nil, err
	}

	newValue, err := convertInterfaceToStringList(n.(common.HclListObject), key)
	if err != nil {
		return nil, nil, err
	}

	idxList := make([]int, 0)

	for index, _ := range oldValue {
		if oldValue[index] != newValue[index] {
			idxList = append(idxList, index)
		}
	}
	return idxList, newValue, nil
}

func resourceBareMetalServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	isDeleteProtected := rd.Get("delete_protection").(bool)
	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	contractDiscount := rd.Get("contract_discount").(string)
	blockStorageList := rd.Get("block_storages").(common.HclListObject)

	subnetId := rd.Get("subnet_id").(string)

	adminAccount := rd.Get("admin_account").(string)
	adminPassword := rd.Get("admin_password").(string)
	initialScript := rd.Get("initial_script").(string)

	vpcId := rd.Get("vpc_id").(string)
	imageId := rd.Get("image_id").(string)

	servers := rd.Get("servers").(common.HclListObject)

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

	if isOsWindows {
		err := common.ValidateServerNameInWindowImage(servers)
		if err != nil {
			return diag.FromErr(err)
		}
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

	serverDetailList, err := ConvertBaremetalServerList(servers, blockStorageInfoList, scaleId)
	if err != nil {
		return diag.FromErr(err)
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
		ServerDetails:             serverDetailList,
		VpcId:                     vpcId,
	}

	createResponse, err := inst.Client.BareMetal.CreateBareMetalServer(ctx, createRequest, rd.Get("tags").(map[string]interface{}))
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
	baremetalIds := strings.Split(rd.Id(), ",")

	bmServerInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, baremetalIds[0])
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

	rd.Set("subnet_id", bmServerInfo.SubnetId)
	rd.Set("image_id", bmServerInfo.ImageId)
	rd.Set("delete_protection", bmServerInfo.DeletionProtectionEnabled == "Y")
	rd.Set("contract_discount", bmServerInfo.Contract)
	rd.Set("vpc_id", bmServerInfo.VpcId)
	rd.Set("initial_script", bmServerInfo.InitialScriptContent)

	blockStorages := common.HclListObject{}
	for _, blockId := range bmServerInfo.BareMetalBlockStorageIds {
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

	servers := common.HclListObject{}

	for _, baremetalId := range baremetalIds {
		bmServerInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, baremetalId)
		if err != nil {
			rd.SetId("")
			if common.IsDeleted(err) {
				return nil
			}
			return diag.FromErr(err)
		}

		serverInfo := common.HclKeyValueObject{}

		serverInfo["bm_server_name"] = bmServerInfo.BareMetalServerName
		serverInfo["ipv4"] = bmServerInfo.IpAddress
		serverInfo["nat_enabled"] = bmServerInfo.NatIpAddress != ""
		serverInfo["local_subnet_enabled"] = bmServerInfo.BareMetalLocalSubnetId != ""
		serverInfo["local_subnet_id"] = bmServerInfo.BareMetalLocalSubnetId
		serverInfo["local_subnet_ipv4"] = bmServerInfo.BareMetalLocalSubnetIpAddress
		serverInfo["use_hyper_threading"] = bmServerInfo.UseHyperThreading
		serverInfo["use_dns"] = bmServerInfo.DnsEnabled == "Y"
		serverInfo["state"] = strings.ToUpper(bmServerInfo.BareMetalServerState)

		natIpv4 := bmServerInfo.NatIpAddress

		if natIpv4 != "" {
			publicIpInfo, err := inst.Client.PublicIp.GetPublicIps(ctx, &publicip2.PublicIpOpenApiV3ControllerApiListPublicIpsV3Opts{
				IpAddress:     optional.NewString(natIpv4),
				VpcId:         optional.NewString(bmServerInfo.VpcId),
				PublicIpState: optional.String{},
				UplinkType:    optional.String{},
				CreatedBy:     optional.String{},
				Page:          optional.Int32{},
				Size:          optional.Int32{},
				Sort:          optional.Interface{},
			})
			if err != nil {
				diagnostics = diag.FromErr(err)
				return
			}

			if len(publicIpInfo.Contents) == 0 {
				// this case is found on auto assign mode
				serverInfo["public_ip_id"] = ""
			} else {
				serverInfo["public_ip_id"] = publicIpInfo.Contents[0].PublicIpAddressId
			}
		} else {
			serverInfo["public_ip_id"] = ""
		}

		servers = append(servers, serverInfo)
	}

	rd.Set("servers", servers)

	tfTags.SetTags(ctx, rd, meta, baremetalIds[0])

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

	if !rd.HasChanges("delete_protection") && !rd.HasChanges("contract_discount") &&
		!rd.HasChanges("block_storages") && !rd.HasChanges("servers") && !rd.HasChanges("tags") {
		return diag.Errorf("nothing to update")
	}

	serverIds := strings.Split(rd.Id(), ",")

	/*
		serverInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, serverIds[0])
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
	*/

	if rd.HasChanges("servers") {
		idxList, newValue, err := getChangeDiffIndex(rd, "state")
		if err != nil {
			return diag.FromErr(err)
		}

		stopBaremetalIds := make([]string, 0)
		startBaremetalIds := make([]string, 0)

		for _, value := range idxList {
			if strings.Compare(newValue[value], common.StoppedState) == 0 {
				stopBaremetalIds = append(stopBaremetalIds, serverIds[value])
			} else {
				startBaremetalIds = append(startBaremetalIds, serverIds[value])
			}
		}

		// 실행(RUNNING) -> 중지(STOPPED)
		if len(stopBaremetalIds) != 0 {
			_, err := inst.Client.BareMetal.StopBareMetalServer(ctx, baremetal.BMStartStopRequest{
				BareMetalServerIds: stopBaremetalIds,
			})

			if err != nil {
				return diag.FromErr(err)
			}

			// wait for server state change to editing
			for _, id := range stopBaremetalIds {
				err = WaitForBMServerStatus(ctx, inst.Client, id, common.VirtualServerProcessingStates(), []string{common.StoppedState}, true)

				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// 중지(STOPPED) -> 시작(RUNNING)
		if len(startBaremetalIds) != 0 {
			_, err := inst.Client.BareMetal.StartBareMetalServer(ctx, baremetal.BMStartStopRequest{
				BareMetalServerIds: startBaremetalIds,
			})

			if err != nil {
				return diag.FromErr(err)
			}

			for _, id := range startBaremetalIds {
				// wait for server state change to editing
				err = WaitForBMServerStatus(ctx, inst.Client, id, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)

				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	/*
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
	*/

	for _, bmId := range serverIds {
		err = tfTags.UpdateTags(ctx, rd, meta, bmId)
		if err != nil {
			return
		}
	}

	return resourceBareMetalServerRead(ctx, rd, meta)
}

func resourceBareMetalServerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	baremetalIds := strings.Split(rd.Id(), ",")

	for _, baremetalId := range baremetalIds {
		err := WaitForBMServerStatus(ctx, inst.Client, baremetalId, common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err := inst.Client.BareMetal.DeleteBareMetalServers(ctx, baremetalIds)

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
