package baremetal

import (
	"context"
	"errors"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
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
	"net"
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
			"bm_server_name": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateName3to24LowerAlphaDashStartsWithLowerCase,
				},
				Description: "Bare-metal server name",
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
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Public IP id of this bare-metal server. Public-IP must be a valid public-ip resource which is attached to the VPC.",
			},
			"use_dns": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:    schema.TypeBool,
					Default: false,
				},
				Description: "Enable DNS feature for this bare-metal server.",
			},
			"use_hyper_threading": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:    schema.TypeString,
					Default: "N",
				},
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
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    5,
				Description: "IP address of this bare-metal server",
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateIpv4InList,
				},
			},
			"nat_enabled": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:    schema.TypeBool,
					Default: false,
				},
				Description: "Enable NAT feature for this bare-metal server.",
			},
			"local_subnet_enabled": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Description: "Enable local subnet for this bare-metal server",
			},
			"local_subnet_id": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Local Subnet id of this bare-metal server. Subnet must be a valid subnet resource which is attached to the VPC.",
			},
			"local_subnet_ipv4": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Local IP address of this bare-metal server",
			},
			"state": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateBmState,
				},
				Description: "Baremetal Server State(ex. RUNNING, STOPPED)",
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
	attr := path[0].(cty.GetAttrStep)
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

func validateIpv4InList(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[0].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)
	if value == "" {
		return diags
	}

	trial := net.ParseIP(value)
	if trial.To4() == nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q is not IP address", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validateBmState(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[0].(cty.GetAttrStep)
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

func convertInterfaceToStringList(input []interface{}) []string {

	output := make([]string, len(input))

	for i, v := range input {
		if v == nil {
			output[i] = ""
		} else if str, ok := v.(string); ok {
			output[i] = str
		}
	}

	return output
}

func convertInterfaceToBoolList(input []interface{}) []bool {

	output := make([]bool, len(input))

	for i, v := range input {
		if v == nil {
			output[i] = false
		} else if b, ok := v.(bool); ok {
			output[i] = b
		}
	}

	return output
}

// update 시에 달라진 index와 새로운 값을 반환시켜준다.
func getChangeDiffIndex(rd *schema.ResourceData, key string) (idxList []int, newValue []string) {
	o, n := rd.GetChange(key)
	oldValue := convertInterfaceToStringList(o.([]interface{}))
	newValue = convertInterfaceToStringList(n.([]interface{}))

	for index, _ := range oldValue {
		if oldValue[index] != newValue[index] {
			idxList = append(idxList, index)
		}
	}
	return idxList, newValue
}

func resourceBareMetalServerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	serverName := convertInterfaceToStringList(rd.Get("bm_server_name").([]interface{}))

	isDeleteProtected := rd.Get("delete_protection").(bool)
	cpuCount := rd.Get("cpu_count").(int)
	memorySizeGB := rd.Get("memory_size_gb").(int)

	contractDiscount := rd.Get("contract_discount").(string)
	blockStorageList := rd.Get("block_storages").(common.HclListObject)

	subnetId := rd.Get("subnet_id").(string)
	publicIpId := convertInterfaceToStringList(rd.Get("public_ip_id").([]interface{}))

	adminAccount := rd.Get("admin_account").(string)
	adminPassword := rd.Get("admin_password").(string)
	initialScript := rd.Get("initial_script").(string)

	states := convertInterfaceToStringList(rd.Get("state").([]interface{}))

	for _, state := range states {
		if strings.Compare(state, common.StoppedState) == 0 {
			return diag.Errorf("state value must be RUNNING")
		}
	}

	vpcId := rd.Get("vpc_id").(string)
	imageId := rd.Get("image_id").(string)
	useDNS := convertInterfaceToBoolList(rd.Get("use_dns").([]interface{}))
	ipAddr := convertInterfaceToStringList(rd.Get("ipv4").([]interface{}))
	natEnabled := convertInterfaceToBoolList(rd.Get("nat_enabled").([]interface{}))
	useHyperThreading := convertInterfaceToStringList(rd.Get("use_hyper_threading").([]interface{}))
	localSubnetEnabled := convertInterfaceToBoolList(rd.Get("local_subnet_enabled").([]interface{}))
	localSubnetIds := convertInterfaceToStringList(rd.Get("local_subnet_id").([]interface{}))
	localSubnetIp := convertInterfaceToStringList(rd.Get("local_subnet_ipv4").([]interface{}))

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

	serverDetailList := make([]baremetal.BMServerDetailsRequest, 0)

	for index, _ := range serverName {

		localSubnetId := localSubnetIds[index]
		if !localSubnetEnabled[index] {
			localSubnetId = "" // to bypass local-subnet data manipulation bug
		}

		if !localSubnetEnabled[index] && localSubnetId != "" {
			return diag.Errorf("local subnet is disabled, but has subnet id")
		}

		serverDetails := baremetal.BMServerDetailsRequest{
			BareMetalLocalSubnetEnabled:   localSubnetEnabled[index],
			BareMetalLocalSubnetId:        localSubnetId,
			BareMetalLocalSubnetIpAddress: localSubnetIp[index],
			BareMetalServerName:           serverName[index],
			DnsEnabled:                    useDNS[index],
			IpAddress:                     ipAddr[index],
			NatEnabled:                    natEnabled[index],
			PublicIpAddressId:             publicIpId[index],
			ServerTypeId:                  scaleId,
			StorageDetails:                blockStorageInfoList,
			UseHyperThreading:             useHyperThreading[index],
		}
		serverDetailList = append(serverDetailList, serverDetails)
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

	serverNames := make([]string, 0)
	ipv4List := make([]string, 0)
	useDnsList := make([]bool, 0)
	useHyperThreadingList := make([]string, 0)
	natEnabledList := make([]bool, 0)
	publicIpIdList := make([]string, 0)
	localSubnetEnabledList := make([]bool, 0)
	localSubnetIdList := make([]string, 0)
	localSubnetIpv4List := make([]string, 0)
	states := make([]string, 0)

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

	serverNames = append(serverNames, bmServerInfo.BareMetalServerName)
	ipv4List = append(ipv4List, bmServerInfo.IpAddress)
	useDnsList = append(useDnsList, bmServerInfo.DnsEnabled == "Y")
	useHyperThreadingList = append(useHyperThreadingList, bmServerInfo.UseHyperThreading)
	natEnabledList = append(natEnabledList, bmServerInfo.NatIpAddress != "")
	localSubnetEnabledList = append(localSubnetEnabledList, bmServerInfo.BareMetalLocalSubnetId != "")
	localSubnetIdList = append(localSubnetIdList, bmServerInfo.BareMetalLocalSubnetId)
	localSubnetIpv4List = append(localSubnetIpv4List, bmServerInfo.BareMetalLocalSubnetIpAddress)
	states = append(states, strings.ToUpper(bmServerInfo.BareMetalServerState))

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
			publicIpIdList = append(publicIpIdList, "")
		} else {
			publicIpIdList = append(publicIpIdList, publicIpInfo.Contents[0].PublicIpAddressId)
		}
	} else {
		publicIpIdList = append(publicIpIdList, "")
	}

	for i := 1; i < len(baremetalIds); i++ {
		bmServerInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, baremetalIds[i])
		if err != nil {
			rd.SetId("")
			if common.IsDeleted(err) {
				return nil
			}
			return diag.FromErr(err)
		}

		serverNames = append(serverNames, bmServerInfo.BareMetalServerName)
		ipv4List = append(ipv4List, bmServerInfo.IpAddress)
		useDnsList = append(useDnsList, bmServerInfo.DnsEnabled == "Y")
		useHyperThreadingList = append(useHyperThreadingList, bmServerInfo.UseHyperThreading)
		natEnabledList = append(natEnabledList, bmServerInfo.NatIpAddress != "")
		localSubnetEnabledList = append(localSubnetEnabledList, bmServerInfo.BareMetalLocalSubnetId != "")
		localSubnetIdList = append(localSubnetIdList, bmServerInfo.BareMetalLocalSubnetId)
		localSubnetIpv4List = append(localSubnetIpv4List, bmServerInfo.BareMetalLocalSubnetIpAddress)
		states = append(states, strings.ToUpper(bmServerInfo.BareMetalServerState))

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
				publicIpIdList = append(publicIpIdList, "")
			} else {
				publicIpIdList = append(publicIpIdList, publicIpInfo.Contents[0].PublicIpAddressId)
			}
		} else {
			publicIpIdList = append(publicIpIdList, "")
		}
	}
	rd.Set("bm_server_name", serverNames)
	rd.Set("ipv4", ipv4List)
	rd.Set("local_subnet_id", localSubnetIdList)
	rd.Set("local_subnet_ipv4", localSubnetIpv4List)
	rd.Set("local_subnet_enabled", localSubnetEnabledList)
	rd.Set("public_ip_id", publicIpIdList)
	rd.Set("nat_enabled", natEnabledList)
	rd.Set("use_dns", useDnsList)
	rd.Set("use_hyper_threading", useHyperThreadingList)
	rd.Set("state", states)
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
		!rd.HasChanges("local_subnet_enabled") && !rd.HasChanges("nat_enabled") &&
		!rd.HasChanges("block_storages") && !rd.HasChanges("state") && !rd.HasChanges("tags") {
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

	if rd.HasChanges("state") {
		idxList, newValue := getChangeDiffIndex(rd, "state")

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
	}

	if err != nil {
		return
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
