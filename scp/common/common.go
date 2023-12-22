package common

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"hash/fnv"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/product"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	NetworkProductGroup      string = "NETWORKING"
	VpcProductName           string = "VPC Traffic"
	DirectConnectProductName string = "Direct Connect"
	PublicIpProductName      string = "Reserved IP"
	SecurityGroupProductName string = "Security Group"
	FileStorageProductName   string = "File Storage(New)"

	StorageProductGroup           string = "STORAGE"
	ContainerProductGroup         string = "CONTAINER"
	KubernetesEngineVmProductName string = "Kubernetes Engine VM"

	ContractProductType string = "CONTRACT_DISCOUNT"

	CreatingState    string = "CREATING"
	ReservedState    string = "RESERVED"
	ActiveState      string = "ACTIVE"
	InActiveState    string = "INACTIVE"
	DeployingState   string = "DEPLOYING"
	DeletedState     string = "DELETED"
	TerminatingState string = "TERMINATING"
	RunningState     string = "RUNNING"
	AvailableState   string = "AVAILABLE"
	UnavailableState string = "UNAVAILABLE"
	UnknownState     string = "UNKNOWN"
	ErrorState       string = "ERROR"
	EditingState     string = "EDITING"
	StartingState    string = "STARTING"
	StoppingState    string = "STOPPING"
	StoppedState     string = "STOPPED"
	RestartingState  string = "RESTARTING"
	SoftDeletedState string = "SOFT_DELETED"
	UpgradingState   string = "UPGRADING"

	VpcPublicIpPurpose            string = "NAT"
	VpcPublicIpNetworkServiceType string = "VPC"

	ServicedGroupCompute       string = "COMPUTE"
	ServicedForVirtualServer   string = "Virtual Server"
	ServicedForGpuServer       string = "GPU Server"
	ServicedForBaremetalServer string = "Baremetal Server"
	ServicedGroupDatabase      string = "DATABASE"
	ServicedForPostgresql      string = "PostgreSQL"
	ServicedForMariadb         string = "Mariadb"
	ServicedForMySql           string = "MySql"
	ServicedForEpas            string = "EPAS"
	ServicedForSqlServer       string = "Microsoft SQL Server"
	ServicedForTibero          string = "Tibero"

	ProductTypeDisk string = "DISK"

	// Product & Product Group state

	ProductActiveState    string = "ACTIVE"
	ProductAvailableState string = "AVAILABLE"

	// Product & Product Group key

	ProductDefaultDisk      string = "DEFAULT_DISK"
	ProductDisk             string = "DISK"
	ProductScale            string = "SCALE"
	ProductIP               string = "IP"
	ProductContractDiscount string = "CONTRACT_DISCOUNT"
	ProductCloudType        string = "CLOUD_TYPE"
	ProductMonitoringTool   string = "MONITORING_TOOL"
	ProductOS               string = "OS"
	ProductPGLevel1         string = "PG_LEVEL1"
	ProductPGLevel2         string = "PG_LEVEL2"
	ProductServiceLevel     string = "SERVICE_LEVEL"

	OsTypeWindows string = "WINDOWS"

	LinuxAdminAccount   string = "root"
	WindowsAdminAccount string = "Administrator"

	DeploymentEnvironmentDev string = "DEV"
	DeploymentEnvironmentPrd string = "PRD"

	BlockStorageTypeOS      string = "OS"
	BlockStorageTypeData    string = "DATA"
	BlockStorageTypeArchive string = "ARCHIVE"
)

func DatabaseProcessingStates() []string {
	return []string{CreatingState, EditingState, StartingState, RestartingState, StoppingState, TerminatingState, UpgradingState}
}

func VirtualServerProcessingStates() []string {
	return []string{CreatingState, EditingState, StartingState, RestartingState, StoppingState, TerminatingState, StoppedState}
}

func NetworkProcessingStates() []string {
	return []string{CreatingState, EditingState, StartingState, TerminatingState}
}

// ResourceMetaData Resource information meta data structure
type ResourceMetaData = []map[string]interface{}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

type HclKeyValueObject = map[string]interface{}
type HclSetObject = []HclKeyValueObject
type HclListObject = []interface{}

//func flattenStringList(list []*string) []interface{} {
//	vs := make([]interface{}, 0, len(list))
//	for _, v := range list {
//		vs = append(vs, *v)
//	}
//	return vs
//}
//
//func flattenStringSet(list []*string) *schema.Set {
//	return schema.NewSet(schema.HashString, flattenStringList(list))
//}

func GenerateHash(values []string) string {
	hasher := fnv.New64a()

	for _, v := range values {
		hasher.Write([]byte(v))
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, hasher.Sum64())
	return base64.StdEncoding.EncodeToString(b)
}
func GetDatasourceItemsSchema(rs *schema.Resource) *schema.Resource {
	if _, ok := rs.Schema["id"]; !ok {
		rs.Schema["id"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	// Ensure Create/Read are not set for nested sub-resource schemas. Otherwise, terraform will validate them
	// as though they were resources.
	rs.CreateContext = nil
	rs.ReadContext = nil

	return convertResourceFieldsToDatasourceFields(rs)
}

// This is mainly used to ensure that fields of a datasource item are compliant with Terraform schema validation
// All datasource return items should have computed-only fields; and not require Diff, Validation, or Default settings.
func convertResourceFieldsToDatasourceFields(resourceSchema *schema.Resource) *schema.Resource {
	resultSchema := map[string]*schema.Schema{}
	for k, fieldSchema := range resourceSchema.Schema {
		isComputed := fieldSchema.Required || fieldSchema.Computed
		fieldSchema.Computed = true
		fieldSchema.Required = false
		fieldSchema.Optional = false
		fieldSchema.DiffSuppressFunc = nil
		fieldSchema.ValidateFunc = nil
		fieldSchema.ValidateDiagFunc = nil
		fieldSchema.ConflictsWith = nil
		fieldSchema.Default = nil
		fieldSchema.MaxItems = 0
		fieldSchema.MaxItems = 0
		if fieldSchema.Type == schema.TypeSet {
			fieldSchema.Type = schema.TypeList
			fieldSchema.Set = nil
		}

		if fieldSchema.Elem != nil {
			if resource, ok := fieldSchema.Elem.(*schema.Resource); ok {
				fieldSchema.Elem = convertResourceFieldsToDatasourceFields(resource)
			}
		}

		if isComputed {
			resultSchema[k] = fieldSchema
		}
	}

	resourceSchema.Schema = resultSchema
	return resourceSchema
}

type AdditionalStorageInfo struct {
	//Id           string
	ProductName  string
	StorageUsage string
	StorageSize  int
}

func ConvertAdditionalStorageList(list HclListObject) []AdditionalStorageInfo {
	var result []AdditionalStorageInfo
	for _, itemObject := range list {
		item := itemObject.(HclKeyValueObject)
		var info AdditionalStorageInfo
		//if v, ok := item["id"]; ok {
		//	info.Id = v.(string)
		//}
		if v, ok := item["product_name"]; ok {
			info.ProductName = v.(string)
		}
		if v, ok := item["storage_usage"]; ok {
			info.StorageUsage = v.(string)
		}
		if v, ok := item["storage_size_gb"]; ok {
			info.StorageSize = v.(int)
		}
		result = append(result, info)
	}
	return result
}

type ExternalStorageInfo struct {
	Name        string
	ProductName string
	StorageSize int
	Encrypted   bool
	SharedType  string
}

func ConvertExternalStorageList(list HclListObject) []ExternalStorageInfo {
	var result []ExternalStorageInfo
	for _, itemObject := range list {
		item := itemObject.(HclKeyValueObject)
		var info ExternalStorageInfo
		if v, ok := item["name"]; ok {
			info.Name = v.(string)
		}
		if v, ok := item["product_name"]; ok {
			info.ProductName = v.(string)
		}
		if v, ok := item["storage_size_gb"]; ok {
			info.StorageSize = v.(int)
		}
		if v, ok := item["encrypted"]; ok {
			info.Encrypted = v.(bool)
		}
		result = append(result, info)
	}
	return result
}

func ExternalStorageResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_storage_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Id",
			},
			"shared_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "DEDICATED",
				Description: "SHARED/DEDICATED",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External storage name",
			},
			"product_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SSD",
				Description: "Storage product name : SSD",
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage product Id",
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
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags",
			},
		},
	}
}

func ProductToIdMap(productType string, productGroup *product.ProductGroupDetailResponse) map[string]string {
	result := make(map[string]string)
	if productInfos, ok := productGroup.Products[productType]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState != ProductAvailableState {
				continue
			}
			result[productInfo.ProductName] = productInfo.ProductId
		}
	}
	return result
}

func ProductIdToNameMap(productType string, productGroup *product.ProductGroupDetailResponse) map[string]string {
	result := make(map[string]string)
	if productInfos, ok := productGroup.Products[productType]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState != ProductAvailableState {
				continue
			}
			result[productInfo.ProductId] = productInfo.ProductName
		}
	}
	return result
}

func FindProductId(productType string, productName string, productGroup *product.ProductGroupDetailResponse) (string, error) {
	if productInfos, ok := productGroup.Products[productType]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState != ProductAvailableState {
				continue
			}
			if productInfo.ProductName == productName {
				return productInfo.ProductId, nil
			}
		}
	}
	return "", fmt.Errorf("product for type '%s' not found", productType)
}

func FirstProductId(productType string, productGroup *product.ProductGroupDetailResponse) (product.ProductForCalculatorResponse, error) {
	if productInfos, ok := productGroup.Products[productType]; ok {
		for _, productInfo := range productInfos {
			if productInfo.ProductState != ProductAvailableState {
				continue
			}
			return productInfo, nil
		}
	}
	return product.ProductForCalculatorResponse{}, fmt.Errorf("product for type '%s' not found", productType)
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToMap(in any) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)

	m := map[string]interface{}{}

	for field, val := range inInterface {
		field = ToSnakeCase(field)

		log.Println("KV Pair: ", field, val)

		var typeOfValue []interface{} //TODO:

		log.Println("=================================================")
		log.Println("typeOfValue: ", reflect.TypeOf(val))
		log.Println("=================================================")
		if reflect.TypeOf(val) == reflect.TypeOf(typeOfValue) {
			m[field] = ConvertStructToMaps(val.([]interface{}))
		} else {
			m[field] = val
		}

	}
	log.Println("=================================================")
	log.Println("ToMap: m: ", m)
	log.Println("=================================================")
	return m
}

func ConvertStructToMaps[T any](contents []T) []map[string]interface{} {
	var contentMaps []map[string]interface{}

	log.Println("=================================================")
	log.Println("contents: ", contents)
	log.Println("=================================================")

	for _, content := range contents {
		contentMaps = append(contentMaps, ToMap(content))
	}

	log.Println("=================================================")
	log.Println("contentMaps: ", contentMaps)
	log.Println("=================================================")

	return contentMaps
}

func IsSubnetContainsIp(subnetString, ipString string) bool {
	_, subnet, _ := net.ParseCIDR(subnetString)
	ip := net.ParseIP(ipString)

	return subnet.Contains(ip)
}

func IsDeleted(err error) bool {
	if strings.HasPrefix(err.Error(), "404") {
		//except resource deleted error
		log.Println("deleted resource")
		return true
	}
	return false
}

func ExpandInterfaceToStringList(origin []interface{}) []string {
	strList := make([]string, len(origin))
	for i, v := range origin {
		strList[i] = v.(string)
	}
	return strList
}

func ToStringList(interfaceList []interface{}) []string {
	if len(interfaceList) == 0 {
		return nil
	}
	stringList := make([]string, len(interfaceList))
	for i, iVal := range interfaceList {
		stringList[i] = iVal.(string)
	}
	return stringList
}

func getAddRemoveItemStringListFromStringList(oldList []string, newList []string) ([]string, []string) {
	oldSet := make(map[string]struct{})
	newSet := make(map[string]struct{})

	for _, item := range oldList {
		oldSet[item] = struct{}{}
	}

	for _, item := range newList {
		newSet[item] = struct{}{}
	}

	var ok bool
	var addList, removeList []string

	for _, item := range oldList {
		if _, ok = newSet[item]; !ok {
			removeList = append(removeList, item)
		}
	}

	for _, item := range newList {
		if _, ok = oldSet[item]; !ok {
			addList = append(addList, item)
		}
	}

	return addList, removeList
}

func GetAddRemoveItemsStringList(rd *schema.ResourceData, key string) ([]string, []string) {
	o, n := rd.GetChange(key)

	oldList := ToStringList(o.([]interface{}))
	newList := ToStringList(n.([]interface{}))

	return getAddRemoveItemStringListFromStringList(oldList, newList)
}

func GetAddRemoveItemsStringListFromSet(rd *schema.ResourceData, key string) ([]string, []string) {
	o, n := rd.GetChange(key)

	oRaw := o.(*schema.Set).List()
	nRaw := n.(*schema.Set).List()

	oldList := ToStringList(oRaw)
	newList := ToStringList(nRaw)

	return getAddRemoveItemStringListFromStringList(oldList, newList)
}

func GetKeyString(rd *schema.ResourceData, key string) optional.String {
	if len(rd.Get(key).(string)) > 0 {
		return optional.NewString(rd.Get(key).(string))
	} else {
		return optional.String{}
	}
}
