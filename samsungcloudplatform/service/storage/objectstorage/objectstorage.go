package objectstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"regexp"
	"strings"
)

var accessControlRuleTypes = []string{
	"IP_ADDRESS_RANGE",
	"VIRTUAL_SERVER",
	"BARE_METAL_SERVER",
	"MULTI_GPU_CLUSTER",
	"HADOOP",
	"GPU_SERVER",
	"VPC_ENDPOINT",
	"MY_SQL",
	"MARIA_DB",
	"TIBERO",
	"SQL_SERVER",
	"POSTGRE_SQL",
	"EPAS",
}

var objectStorageBucketUserPurposes = []string{
	"PUBLIC",
	"PRIVATE",
}

var objectStorageProductNames = []string{
	"Object Storage",
}

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_obs_bucket", ResourceObjectStorageBucket())
}

func ResourceObjectStorageBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: createBucket,
		ReadContext:   readBucket,
		UpdateContext: updateBucket,
		DeleteContext: deleteBucket,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"object_storage_bucket_access_control_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Object Storage Bucket Access Control Enabled",
			},
			"access_control_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Object Storage Bucket Access Control Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access Control Rule Type",
						},
						"rule_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access Control Rule Value",
						},
					},
				},
			},
			"object_storage_bucket_file_encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Object Storage Bucket File Encryption Enabled",
			},
			"object_storage_bucket_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Object Storage Bucket Name",
				ValidateDiagFunc: ValidateName3To63LowerAlphaNumberDashNoLastDash,
			},
			"object_storage_bucket_version_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Object Storage Bucket Version Enabled",
			},
			"object_storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Object Storage ID",
			},
			"product_names": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Product Names",
			},
			"object_storage_bucket_dr_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Object Storage Bucket DR Enabled",
			},
			"object_storage_bucket_dr_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object Storage Bucket DR Type",
			},
			"object_storage_bucket_purpose": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object Storage Bucket Purpose",
			},
			"object_storage_bucket_user_purpose": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "PUBLIC",
				ForceNew:         true,
				Description:      "Object Storage Bucket User Purpose",
				ValidateDiagFunc: ValidateObjectStorageBucketUserPurpose,
			},
			"sync_object_storage_bucket_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sync Object Storage Bucket ID",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone ID",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides an Object Storage Bucket Resource.",
	}
}

func convertAccessIpAddressRanges(list common.HclListObject) ([]objectstorage.AccessControlRule, error) {
	var result []objectstorage.AccessControlRule

	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := objectstorage.AccessControlRule{}
		if v, ok := itemObject["rule_value"]; ok {
			info.RuleValue = v.(string)
		}
		if t, ok := itemObject["rule_type"]; ok {
			err := CheckStringInStringList(accessControlRuleTypes, t.(string))
			if err != nil {
				if err.Error() == "input string is not in the string list" {
					return nil, fmt.Errorf("invalid access control rule type")
				}
				return nil, err
			}
			info.RuleType = t.(string)
		}
		result = append(result, info)
	}
	return result, nil
}

func countAccessControlRulesByType(rules []objectstorage.AccessControlRule, ruleType string) int {
	count := 0
	for _, rule := range rules {
		if rule.RuleType == ruleType {
			count++
		}
	}
	return count
}

func convertToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, 0)
	for _, interfaceElem := range interfaceArray {
		stringArray = append(stringArray, interfaceElem.(string))
	}
	return stringArray
}

func convertStringListToFormattedString(stringList []string) string {
	switch len(stringList) {
	case 1:
		return fmt.Sprintf("%q", stringList[0])
	case 2:
		return fmt.Sprintf(`%q or %q`, stringList[0], stringList[1])
	default:
		var sb strings.Builder
		for _, s := range stringList[:len(stringList)-1] {
			sb.WriteString(fmt.Sprintf("%q or ", s))
		}
		sb.WriteString(fmt.Sprintf("%q", stringList[len(stringList)-1]))
		return sb.String()
	}
}

func createBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	// input access control rules conversion and validation
	// TODO: add invalid acces control rule value case
	accessControlRules, err := convertAccessIpAddressRanges(rd.Get("access_control_rules").(common.HclListObject))
	if err != nil {
		if err.Error() == "invalid access control rule type" {
			return diag.Errorf("Input \"Access Control Rule Type\" is invalid. It must be one of " + convertStringListToFormattedString(accessControlRuleTypes) + ".")
		}
		return diag.Errorf("\"Access Control Rule Type\" validation check failed.")
	}

	CountIpAddressRange := countAccessControlRulesByType(accessControlRules, "IP_ADDRESS_RANGE")
	if CountIpAddressRange > 201 {
		return diag.Errorf("\"IP_ADDRESS_RANGE\" type of rule cannot exceed 200")
	}

	// input service zone id validation
	ServiceZoneId := rd.Get("service_zone_id").(string)
	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var ServiceZoneIds []string
	for _, response := range projectInfo.ServiceZones {
		ServiceZoneIds = append(ServiceZoneIds, response.ServiceZoneId)
	}
	err = CheckStringInStringList(ServiceZoneIds, ServiceZoneId)
	if err != nil {
		if err.Error() == "input string is not in the string list" {
			return diag.Errorf("Input \"Service Zone ID\" is invalid. Check if the Service Zone ID is valid in the project.")
		}
		return diag.Errorf("\"Service Zone ID\" validation check failed.")
	}

	// input object storage id validation
	ObjectStorageId := rd.Get("object_storage_id").(string)
	readObjectStorageListResponses, err := inst.Client.ObjectStorage.ReadObjectStorageList(ctx, ServiceZoneId, objectstorage.ReadObjectStorageListRequest{})
	if err != nil {
		return diag.Errorf("\"Object Storage ID\" validation check failed.")
	} else {
		var ObjectStorageIds []string
		for _, response := range readObjectStorageListResponses.Contents {
			ObjectStorageIds = append(ObjectStorageIds, response.ObjectStorageId)
		}
		err = CheckStringInStringList(ObjectStorageIds, ObjectStorageId)

		if err != nil {
			if err.Error() == "input string is not in the string list" {
				return diag.Errorf("Input \"Object Storage ID\" is invalid. Check if the Object Storage ID is valid in the project.")
			}
			return diag.Errorf("\"Object Storage ID\" validation check failed.")
		}
	}

	// input product names validation
	ProductNames := convertToStringArray(rd.Get("product_names").([]interface{}))
	if !reflect.DeepEqual(objectStorageProductNames, ProductNames) {
		return diag.Errorf("Input \"Product Names\" is invalid. It must be [\"Object Storage\"].")
	}

	// pre-check conditions for using DR
	objectStorageBucketDrEnabled := rd.Get("object_storage_bucket_dr_enabled").(bool)
	objectStorageBucketVersionEnabled := rd.Get("object_storage_bucket_version_enabled").(bool)
	syncObjectStorageBucketId := rd.Get("sync_object_storage_bucket_id").(string)

	if objectStorageBucketDrEnabled && !objectStorageBucketVersionEnabled {
		return diag.Errorf("To use DR, the object storage bucket versioning must be enabled.")
	} else if objectStorageBucketDrEnabled && objectStorageBucketVersionEnabled {
		/*
			TODO: add DR available service zone check
			- cannot check the sync object storage bucket is in DR capable service zone.
		*/
		if syncObjectStorageBucketId != "" {
			syncBucketInfo, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, syncObjectStorageBucketId)
			if err != nil {
				return diag.Errorf("Can not get the information of input sync object storage bucket.")
			} else {
				if !*syncBucketInfo.ObjectStorageBucketVersionEnabled {
					return diag.Errorf("To use DR, an object storage bucket versioning of input sync object storage bucket must be enabled.")
				} else if *syncBucketInfo.ObjectStorageBucketDrEnabled {
					return diag.Errorf("To use DR, a DR control of input sync object storage bucket must be disabled.")
				}
			}
		} else {
			return diag.Errorf("\"Sync Object Storage Bucket ID\" is required to use DR.")
		}
	}

	/*
		TODO: add object storage bucket name duplication check
		- cannot get other user's private bucket info. by ReadBucketList
		- cannot get exact object storage bucket name search result by ReadBucketList(only support 'LIKE')
	*/

	response, err := inst.Client.ObjectStorage.CreateBucket(ctx, objectstorage.CreateBucketRequest{
		ObjectStorageBucketAccessControlEnabled:  rd.Get("object_storage_bucket_access_control_enabled").(bool),
		AccessControlRules:                       accessControlRules,
		ObjectStorageBucketFileEncryptionEnabled: rd.Get("object_storage_bucket_file_encryption_enabled").(bool),
		ObjectStorageBucketName:                  rd.Get("object_storage_bucket_name").(string),
		ObjectStorageBucketVersionEnabled:        objectStorageBucketVersionEnabled,
		ObjectStorageBucketUserPurpose:           rd.Get("object_storage_bucket_user_purpose").(string),
		ObjectStorageId:                          ObjectStorageId,
		ServiceZoneId:                            ServiceZoneId,
		ProductNames:                             ProductNames,
		Tags:                                     rd.Get("tags").(map[string]interface{}),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ObjectStorageBucketId)

	// enable DR
	if objectStorageBucketDrEnabled && objectStorageBucketVersionEnabled {
		err = inst.Client.ObjectStorage.UpdateBucketDr(ctx, rd.Id(), objectStorageBucketDrEnabled, syncObjectStorageBucketId)
		if err != nil {
			// if dr failed, rollback create action
			deleteBucket(ctx, rd, meta)
			return diag.FromErr(err)
		}
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"UPDATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return readBucket(ctx, rd, meta)
}

func readBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, rd.Id())

	objectStorageBucketAccessControlRules := common.HclSetObject{}
	for _, accessControlRules := range info.ObjectStorageBucketAccessControlRules {
		objectStorageBucketAccessControlRules = append(objectStorageBucketAccessControlRules, common.HclKeyValueObject{
			"rule_type":  accessControlRules.RuleType,
			"rule_value": accessControlRules.RuleValue,
		})
	}

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("object_storage_bucket_access_control_enabled", info.ObjectStorageBucketAccessControlEnabled)
	rd.Set("access_control_rules", objectStorageBucketAccessControlRules)
	rd.Set("object_storage_bucket_file_encryption_enabled", info.ObjectStorageBucketFileEncryptionEnabled)
	rd.Set("object_storage_bucket_name", info.ObjectStorageBucketName)
	rd.Set("object_storage_bucket_version_enabled", info.ObjectStorageBucketVersionEnabled)
	rd.Set("object_storage_id", info.ObjectStorageId)
	rd.Set("service_zone_id", info.ServiceZoneId)
	rd.Set("object_storage_bucket_purpose", info.ObjectStorageBucketPurpose)
	rd.Set("object_storage_bucket_user_purpose", info.ObjectStorageBucketUserPurpose)
	rd.Set("object_storage_bucket_dr_enabled", info.ObjectStorageBucketDrEnabled)
	rd.Set("sync_object_storage_bucket_id", info.SyncObjectStorageBucketId)
	rd.Set("object_storage_bucket_dr_type", info.ObjectStorageBucketDrType)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func updateBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("object_storage_bucket_file_encryption_enabled") {
		_, err := inst.Client.ObjectStorage.UpdateBucketEncryption(ctx, rd.Id(), rd.Get("object_storage_bucket_file_encryption_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("object_storage_bucket_access_control_enabled") ||
		(rd.Get("object_storage_bucket_access_control_enabled").(bool) && rd.HasChanges("access_control_rules")) {
		// input access control rules conversion and validation
		// TODO: add invalid access control rule value case
		accessControlRules, err := convertAccessIpAddressRanges(rd.Get("access_control_rules").(common.HclListObject))
		if err != nil {
			if err.Error() == "invalid access control rule type" {
				return diag.Errorf("Input \"Access Control Rule Type\" is invalid. It must be one of " + convertStringListToFormattedString(accessControlRuleTypes) + ".")
			}
			return diag.Errorf("\"Access Control Rule Type\" validation check failed.")
		}
		CountIpAddressRange := countAccessControlRulesByType(accessControlRules, "IP_ADDRESS_RANGE")
		if CountIpAddressRange > 200 {
			return diag.Errorf("\"IP_ADDRESS_RANGE\" type of rule cannot exceed 200")
		}
		_, err = inst.Client.ObjectStorage.CreateBucketIps(ctx, rd.Id(), rd.Get("object_storage_bucket_access_control_enabled").(bool), accessControlRules)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// TODO: add unadapted change error(access control rules)
	/*
		if !rd.Get("object_storage_bucket_access_control_enabled").(bool) && rd.HasChanges("access_control_rules") {
			diags := append(diag.Diagnostics{}, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Change of \"Access Control Rules\" is not applied",
				Detail:   "If \"Object Storage Bucket Access Control Enabled\" is FALSE or Null, you cannot apply access control rules.",
			})
			return diags
		}
	*/

	objectStorageBucketVersionChanged := rd.HasChanges("object_storage_bucket_version_enabled")
	objectStorageBucketDrChanged := rd.HasChanges("object_storage_bucket_dr_enabled")
	objectStorageBucketVersionEnabled := rd.Get("object_storage_bucket_version_enabled").(bool)
	objectStorageBucketDrEnabled := rd.Get("object_storage_bucket_dr_enabled").(bool)
	syncObjectStorageBucketId := rd.Get("sync_object_storage_bucket_id").(string)

	switch rd.Get("object_storage_bucket_dr_type") {
	case "":
		if objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> true
				changeVersionControl(inst, ctx, rd, objectStorageBucketVersionEnabled)
				changeDRControl(inst, ctx, rd, syncObjectStorageBucketId, objectStorageBucketDrEnabled)
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR false -> true
				return diag.Errorf("To use DR, the object storage bucket versioning must be enabled.")
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> false
			} else {
				// versioning true -> false && DR true -> false
			}
		} else if !objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> true && DR false -> true
				changeDRControl(inst, ctx, rd, syncObjectStorageBucketId, objectStorageBucketDrEnabled)
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> false && DR false -> true
				return diag.Errorf("To use DR, the object storage bucket versioning must be enabled.")
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning true -> true && DR true -> false
			} else {
				// versioning false -> false && DR true -> false
			}
		} else if objectStorageBucketVersionChanged && !objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR true -> true
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> false
				// versioning enable
				changeVersionControl(inst, ctx, rd, objectStorageBucketVersionEnabled)
			} else {
				// versioning true -> false && DR false -> false
				// versioning disable
				changeVersionControl(inst, ctx, rd, objectStorageBucketVersionEnabled)
			}
		} else {
			// !object_storage_bucket_version_changed && !object_storage_bucket_dr_changed
			// nothing to update
		}
	case "ORIGIN":
		// TODO: add unadapted change error(same sync bucket id validation)
		if objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR false -> true
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> false
			} else {
				// versioning true -> false && DR true -> false
				changeDRControl(inst, ctx, rd, syncObjectStorageBucketId, objectStorageBucketDrEnabled)
				changeVersionControl(inst, ctx, rd, objectStorageBucketVersionEnabled)
			}
		} else if !objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> true && DR false -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> false && DR false -> true
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning true -> true && DR true -> false
				changeDRControl(inst, ctx, rd, syncObjectStorageBucketId, objectStorageBucketDrEnabled)
			} else {
				// versioning false -> false && DR true -> false
			}
		} else if objectStorageBucketVersionChanged && !objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR true -> true
				return diag.Errorf("While DR is enabled, you cannot disable versioning for an object storage bucket.")
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> false
			} else {
				// versioning true -> false && DR false -> false
			}
		} else {
			// !object_storage_bucket_version_changed && !object_storage_bucket_dr_changed
			// nothing to update
		}
	case "CLONE":
		// TODO: add unadapted change error(same sync bucket id validation)
		if objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR false -> true
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> false
			} else {
				// versioning true -> false && DR true -> false
				// TODO: add unadapted change error(Clone bucket)
				/*
					diags := append(diag.Diagnostics{}, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Change of \"Object Storage Bucket DR Enabled\" is not applied",
						Detail:   "If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable version control.",
					})
					diags = append(diag.Diagnostics{}, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Change of \"Object Storage Bucket Version Enabled\" is not applied",
						Detail:   "If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable version control.",
					})
					return diags
				*/
			}
		} else if !objectStorageBucketVersionChanged && objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> true && DR false -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> false && DR false -> true
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning true -> true && DR true -> false
				// return diag.Errorf("If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable DR sync.")
				// TODO: add unadapted change error(Clone bucket)
				/*
					diags := append(diag.Diagnostics{}, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Change of \"Object Storage Bucket DR Enabled\" is not applied",
						Detail:   "If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable version control.",
					})
					return diags
				*/
			} else {
				// versioning false -> false && DR true -> false
			}
		} else if objectStorageBucketVersionChanged && !objectStorageBucketDrChanged {
			if objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning false -> true && DR true -> true
			} else if !objectStorageBucketVersionEnabled && objectStorageBucketDrEnabled {
				// versioning true -> false && DR true -> true
				// return diag.Errorf("If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable version control.")
				diags := append(diag.Diagnostics{}, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Change of \"Object Storage Bucket Version Enabled\" is not applied",
					Detail:   "If \"Object Storage Bucket DR Type\" is \"CLONE\", you cannot disable version control.",
				})
				return diags
			} else if objectStorageBucketVersionEnabled && !objectStorageBucketDrEnabled {
				// versioning false -> true && DR false -> false
			} else {
				// versioning true -> false && DR false -> false
			}
		} else {
			// !object_storage_bucket_version_changed && !object_storage_bucket_dr_changed
			// nothing to update
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return readBucket(ctx, rd, meta)
}

func deleteBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.ObjectStorage.DeleteBucket(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"DELETING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForObjectStorageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.ObjectStorage.ReadBucket(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		if info.ObjectStorageBucketId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.ObjectStorageBucketState, nil
	})
}

func CheckStringLength(str string, min int, max int) error {
	if len(str) < min {
		return fmt.Errorf("input must be longer than %v characters", min)
	} else if len(str) > max {
		return fmt.Errorf("input must be shorter than %v characters", max)
	} else {
		return nil
	}
}

func CheckStringInStringList(strList []string, str string) error {
	if len(strList) == 0 {
		return fmt.Errorf("input string list is empty")
	}
	for _, s := range strList {
		if s == str {
			return nil
		}
	}
	return fmt.Errorf("input string is not in the string list")
}

func ValidateName3To63LowerAlphaNumberDashNoLastDash(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := CheckStringLength(value, 3, 63)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
		return diags
	}

	// Check characters
	if !regexp.MustCompile("^[a-z0-9][a-z0-9-]+[a-z0-9]$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must start with a lowercase and number, and enter using lowercase, number, and -. However, it does not end with -.", attrKey),
			AttributePath: path,
		})
		return diags
	}

	return diags
}

func ValidateObjectStorageBucketUserPurpose(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	err := CheckStringInStringList(objectStorageBucketUserPurposes, value)
	if err != nil {
		if err.Error() == "input string is not in the string list" {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Attribute %q must be one of "+convertStringListToFormattedString(objectStorageBucketUserPurposes)+".", attrKey),
				AttributePath: path,
			})
			return diags
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
				AttributePath: path,
			})
			return diags
		}
	}

	return diags
}

func changeDRControl(inst *client.Instance, ctx context.Context, rd *schema.ResourceData, syncBucketId string, drEnabled bool) diag.Diagnostics {
	if syncBucketId != "" {
		syncBucketInfo, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, syncBucketId)
		if err != nil {
			return diag.Errorf("Can not get the information of input sync object storage bucket.")
		} else {
			if drEnabled { // drEnabled == true
				/*
					TODO: add DR available service zone check
					- cannot check the sync object storage bucket is in DR capable service zone.
				*/
				if !*syncBucketInfo.ObjectStorageBucketVersionEnabled {
					return diag.Errorf("To use DR, an object storage bucket versioning of input sync object storage bucket must be enabled.")
				} else if *syncBucketInfo.ObjectStorageBucketDrEnabled {
					return diag.Errorf("To use DR, a DR control of input sync object storage bucket must be disabled.")
				} else {
					err = inst.Client.ObjectStorage.UpdateBucketDr(ctx, rd.Id(), drEnabled, syncBucketId)
					if err != nil {
						return diag.FromErr(err)
					}
					err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"UPDATING"}, []string{"ACTIVE"}, true)
					if err != nil {
						return diag.FromErr(err)
					}
				}

			} else { // drEnabled == false
				if syncBucketInfo.SyncObjectStorageBucketId == rd.Id() && syncBucketInfo.ObjectStorageBucketDrType == "CLONE" {
					err = inst.Client.ObjectStorage.UpdateBucketDr(ctx, rd.Id(), drEnabled, syncBucketId)
					if err != nil {
						return diag.FromErr(err)
					}
					// clear sync_object_storage_bucket_id
					rd.Set("sync_object_storage_bucket_id", "")
					err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"UPDATING"}, []string{"ACTIVE"}, true)
					if err != nil {
						return diag.FromErr(err)
					}
				} else {
					return diag.Errorf("\"Sync Object Storage Bucket ID\" is not valid.")
				}
			}
		}
	} else {
		return diag.Errorf("\"Sync Object Storage Bucket ID\" is required to use DR.")
	}
	return nil
}

func changeVersionControl(inst *client.Instance, ctx context.Context, rd *schema.ResourceData, versioningEnabled bool) diag.Diagnostics {
	_, err := inst.Client.ObjectStorage.UpdateVersioning(ctx, rd.Id(), versioningEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
