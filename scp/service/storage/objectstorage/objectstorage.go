package objectstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_obs_bucket", ResourceObjectStorageBucket())
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Object Storage Bucket Name",
				// name validation 필요
			},
			"object_storage_bucket_version_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Object Storage Bucket Version Enabled",
			},
			"object_storage_id": {
				Type:        schema.TypeString,
				Required:    true,
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
				Optional:    true,
				Description: "Object Storage Bucket DR Type",
			},
			"object_storage_bucket_purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object Storage Bucket Purpose",
			},
			"object_storage_bucket_user_purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object Storage Bucket User Purpose",
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
		Description: "Provides an Object Storage resource.",
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
			info.RuleType = t.(string)
		}
		result = append(result, info)
	}
	return result, nil
}

func convertToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, 0)
	for _, interfaceElem := range interfaceArray {
		stringArray = append(stringArray, interfaceElem.(string))
	}
	return stringArray
}

func createBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	isNameDuplicated, err := inst.Client.ObjectStorage.CheckBucketName(ctx, rd.Get("object_storage_id").(string), rd.Get("object_storage_bucket_name").(string))
	if err != nil {
		return diag.FromErr(err)
	} else if isNameDuplicated {
		return diag.Errorf("Bucket Name is duplicated")
	}

	accessControlRules, err := convertAccessIpAddressRanges(rd.Get("access_control_rules").(common.HclListObject))
	if err != nil {
		return diag.Errorf("Access Control Rule is not valid")
	}
	productNames := convertToStringArray(rd.Get("product_names").([]interface{}))

	response, err := inst.Client.ObjectStorage.CreateBucket(ctx, objectstorage.CreateBucketRequest{
		ObjectStorageBucketAccessControlEnabled:  rd.Get("object_storage_bucket_access_control_enabled").(bool),
		AccessControlRules:                       accessControlRules,
		ObjectStorageBucketFileEncryptionEnabled: rd.Get("object_storage_bucket_file_encryption_enabled").(bool),
		ObjectStorageBucketName:                  rd.Get("object_storage_bucket_name").(string),
		ObjectStorageBucketVersionEnabled:        rd.Get("object_storage_bucket_version_enabled").(bool),
		ObjectStorageBucketUserPurpose:           rd.Get("object_storage_bucket_user_purpose").(string),
		ObjectStorageId:                          rd.Get("object_storage_id").(string),
		ServiceZoneId:                            rd.Get("service_zone_id").(string),
		ProductNames:                             productNames,
		Tags:                                     rd.Get("tags").(map[string]interface{}),
	})

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, response.ObjectStorageBucketId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ObjectStorageBucketId)

	// if dr == enabled //
	if rd.Get("object_storage_bucket_dr_enabled").(bool) {
		err := inst.Client.ObjectStorage.UpdateBucketDr(ctx, response.ObjectStorageBucketId, rd.Get("object_storage_bucket_dr_enabled").(bool), rd.Get("sync_object_storage_bucket_id").(string))
		if err != nil {
			// if dr failed, rollback create action
			return deleteBucket(ctx, rd, meta)
		}
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, response.ObjectStorageBucketId, []string{"UPDATING"}, []string{"ACTIVE"}, true)
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
	rd.Set("object_storage_bucket_dr_enabled", info.ObjectStorageBucketDrEnabled)
	rd.Set("object_storage_bucket_dr_type", info.ObjectStorageBucketDrType)
	rd.Set("object_storage_bucket_purpose", info.ObjectStorageBucketPurpose)
	rd.Set("object_storage_bucket_user_purpose", info.ObjectStorageBucketUserPurpose)
	rd.Set("sync_object_storage_bucket_id", info.SyncObjectStorageBucketId)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func updateBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("object_storage_bucket_version_enabled") {
		_, err := inst.Client.ObjectStorage.UpdateVersioning(ctx, rd.Id(), rd.Get("object_storage_bucket_version_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("object_storage_bucket_file_encryption_enabled") {
		_, err := inst.Client.ObjectStorage.UpdateBucketEncryption(ctx, rd.Id(), rd.Get("object_storage_bucket_file_encryption_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("object_storage_bucket_access_control_enabled") ||
		(rd.Get("object_storage_bucket_access_control_enabled").(bool) && rd.HasChanges("access_control_rules")) {
		accessControlRules, err := convertAccessIpAddressRanges(rd.Get("access_control_rules").(common.HclListObject))
		inst.Client.ObjectStorage.CreateBucketIps(ctx, rd.Id(), rd.Get("object_storage_bucket_access_control_enabled").(bool), accessControlRules)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("object_storage_bucket_dr_enabled") {

		syncObjectStorageBucketId, ok := rd.Get("sync_object_storage_bucket_id").(string)
		if ok && rd.Get("sync_object_storage_bucket_id").(string) != "" {
			err := inst.Client.ObjectStorage.UpdateBucketDr(ctx, rd.Id(), rd.Get("object_storage_bucket_dr_enabled").(bool), syncObjectStorageBucketId)
			if err != nil {
				return diag.FromErr(err)
			}
			err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"UPDATING"}, []string{"ACTIVE"}, true)
			if err != nil {
				return diag.FromErr(err)
			}
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
