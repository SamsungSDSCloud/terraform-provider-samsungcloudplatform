package objectstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
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
			"is_obs_bucket_ip_address_filter_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ip Address Filter Is Enabled",
			},
			"obs_bucket_access_ip_address_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Object Storage Bucket Access IP Address Ranges",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs_bucket_access_ip_address_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object Storage Bucket Access IP Address Range",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Range Type",
						},
					},
				},
			},
			"obs_bucket_file_encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object Storage Bucket File Encryption Algorithm (AES256)",
			},
			"obs_bucket_file_encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable File Encryption for Object Storage Bucket",
			},
			"obs_bucket_file_encryption_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object Storage Bucket File Encryption Type (SSE-S3)",
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Object Storage Bucket Name",
				// name validation 필요
			},
			"obs_bucket_version_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Object Storage Bucket Version usage",
			},
			"obs_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object Storage Id",
			},
			"is_obs_bucket_dr_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable Object Storage Bucket DR",
			},
			"replica_obs_bucket_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Replica Object Storage Bucket ID",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone ID",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags",
			},
		},
		Description: "Provides an Object Storage resource.",
	}
}

func convertAccessIpAddressRanges(list common.HclListObject) ([]objectstorage.ObsBucketAccessIpAddressInfo, error) {
	var result []objectstorage.ObsBucketAccessIpAddressInfo
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := objectstorage.ObsBucketAccessIpAddressInfo{}
		if ipAddressRange, ok := itemObject["obs_bucket_access_ip_address_range"]; ok {
			info.ObsBucketAccessIpAddressRange = ipAddressRange.(string)
		}
		if t, ok := itemObject["type"]; ok {
			info.Type = t.(string)
		}
		result = append(result, info)
	}
	return result, nil
}

func createBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	isNameDuplicated, err := inst.Client.ObjectStorage.CheckBucketName(ctx, rd.Get("obs_id").(string), rd.Get("obs_bucket_name").(string))
	if err != nil {
		return diag.FromErr(err)
	} else if isNameDuplicated {
		return diag.Errorf("Bucket Name is duplicated")
	}

	obsBucketAccessIpAddressInfos, err := convertAccessIpAddressRanges(rd.Get("obs_bucket_access_ip_address_ranges").(common.HclListObject))
	if err != nil {
		return diag.Errorf("Bucket Access IP Address Range is not valid")
	}

	response, err := inst.Client.ObjectStorage.CreateBucket(ctx, objectstorage.CreateBucketRequest{
		IsObsBucketIpAddressFilterEnabled: rd.Get("is_obs_bucket_ip_address_filter_enabled").(bool),
		ObsBucketAccessIpAddressRanges:    obsBucketAccessIpAddressInfos,
		ObsBucketFileEncryptionAlgorithm:  rd.Get("obs_bucket_file_encryption_algorithm").(string),
		ObsBucketFileEncryptionEnabled:    rd.Get("obs_bucket_file_encryption_enabled").(bool),
		ObsBucketFileEncryptionType:       rd.Get("obs_bucket_file_encryption_type").(string),
		ObsBucketName:                     rd.Get("obs_bucket_name").(string),
		ObsBucketVersionEnabled:           rd.Get("obs_bucket_version_enabled").(bool),
		ObsId:                             rd.Get("obs_id").(string),
		ZoneId:                            rd.Get("zone_id").(string),
		Tags:                              getTagRequestArray(rd),
	})

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, response.ObsBucketId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ObsBucketId)

	// if dr == enabled //
	if rd.Get("is_obs_bucket_dr_enabled").(bool) {
		err := inst.Client.ObjectStorage.UpdateBucketDr(ctx, response.ObsBucketId, rd.Get("is_obs_bucket_dr_enabled").(bool), rd.Get("replica_obs_bucket_id").(string))
		if err != nil {
			// if dr failed, rollback create action
			return deleteBucket(ctx, rd, meta)
		}
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, response.ObsBucketId, []string{"UPDATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return readBucket(ctx, rd, meta)
}

func readBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.ObjectStorage.ReadBucket(ctx, rd.Id())

	obsBucketAccessIpAddressRanges := common.HclSetObject{}
	for _, ipAddressRange := range info.ObsBucketAccessIpAddressRanges {
		obsBucketAccessIpAddressRanges = append(obsBucketAccessIpAddressRanges, common.HclKeyValueObject{
			"obs_bucket_access_ip_address_range": ipAddressRange.ObsBucketAccessIpAddressRange,
			"type":                               ipAddressRange.Type_,
		})
	}

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("is_obs_bucket_ip_address_filter_enabled", info.IsObsBucketIpAddressFilterEnabled)
	rd.Set("obs_bucket_access_ip_address_ranges", obsBucketAccessIpAddressRanges)
	rd.Set("obs_bucket_file_encryption_algorithm", info.ObsBucketFileEncryptionAlgorithm)
	rd.Set("obs_bucket_file_encryption_enabled", info.ObsBucketFileEncryptionEnabled)
	rd.Set("obs_bucket_file_encryption_type", info.ObsBucketFileEncryptionType)
	rd.Set("obs_bucket_name", info.ObsBucketName)
	rd.Set("obs_bucket_version_enabled", info.ObsBucketVersionEnabled)
	rd.Set("obs_id", info.ObsId)
	rd.Set("zone_id", info.ZoneId)
	rd.Set("dr_enable", info.IsObsBucketDrEnabled)
	rd.Set("replica_obs_bucket_id", info.ObsSyncBucketId)

	return nil
}

func updateBucket(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("obs_bucket_version_enabled") {
		_, err := inst.Client.ObjectStorage.UpdateVersioning(ctx, rd.Id(), rd.Get("obs_bucket_version_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("obs_bucket_file_encryption_enabled") {
		_, err := inst.Client.ObjectStorage.UpdateBucketEncryption(ctx, rd.Id(), objectstorage.UpdateBucketRequest{
			ObsBucketFileEncryptionAlgorithm: rd.Get("obs_bucket_file_encryption_algorithm").(string),
			ObsBucketFileEncryptionEnabled:   rd.Get("obs_bucket_file_encryption_enabled").(bool),
			ObsBucketFileEncryptionType:      rd.Get("obs_bucket_file_encryption_type").(string),
			ObsBucketVersionEnabled:          rd.Get("obs_bucket_version_enabled").(bool),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("is_obs_bucket_ip_address_filter_enabled") ||
		(rd.Get("is_obs_bucket_ip_address_filter_enabled").(bool) && rd.HasChanges("obs_bucket_access_ip_address_ranges")) {
		obsBucketAccessIpAddressInfos, err := convertAccessIpAddressRanges(rd.Get("obs_bucket_access_ip_address_ranges").(common.HclListObject))
		inst.Client.ObjectStorage.CreateBucketIps(ctx, rd.Id(), rd.Get("is_obs_bucket_ip_address_filter_enabled").(bool), obsBucketAccessIpAddressInfos)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("is_obs_bucket_dr_enabled") {
		err := inst.Client.ObjectStorage.UpdateBucketDr(ctx, rd.Id(), rd.Get("is_obs_bucket_dr_enabled").(bool), rd.Get("replica_obs_bucket_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForObjectStorageStatus(ctx, inst.Client, rd.Id(), []string{"UPDATING"}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
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
		if info.ObsBucketId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.ObsBucketState, nil
	})
}

func getTagRequestArray(rd *schema.ResourceData) []objectstorage.TagRequest {
	tags := rd.Get("tags").(map[string]interface{})
	tagsRequests := make([]objectstorage.TagRequest, 0)
	for key, value := range tags {
		tagsRequests = append(tagsRequests, objectstorage.TagRequest{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	return tagsRequests
}
