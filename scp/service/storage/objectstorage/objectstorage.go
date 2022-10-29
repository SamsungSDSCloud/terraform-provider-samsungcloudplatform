package objectstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceObjectStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createObjectStorage,
		ReadContext:   readObjectStorage,
		UpdateContext: updateObjectStorage,
		DeleteContext: deleteObjectStorage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip_address_filter_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"access_ip_address_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address_range": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
					},
				},
			},
			"file_encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"file_encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"file_encryption_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"version_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"obs_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
		},
	}
}

func convertAccessIpAddressRanges(list common.HclListObject) ([]objectstorage.ObsBucketAccessIpAddressInfo, error) {
	var result []objectstorage.ObsBucketAccessIpAddressInfo
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := objectstorage.ObsBucketAccessIpAddressInfo{}
		if ip_address_range, ok := itemObject["ip_address_range"]; ok {
			info.ObsBucketAccessIpAddressRange = ip_address_range.(string)
		}
		if t, ok := itemObject["type"]; ok {
			info.Type = t.(string)
		}

		result = append(result, info)
	}
	return result, nil
}

func createObjectStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	obsBucketAccessIpAddressInfos, err := convertAccessIpAddressRanges(data.Get("access_ip_address_ranges").(common.HclListObject))
	if err != nil {
		return nil
	}

	response, err := inst.Client.ObjectStorage.CreateObjectStorage(ctx, objectstorage.CreateObjectStorageRequest{
		IsObsBucketIpAddressFilterEnabled: data.Get("ip_address_filter_enabled").(bool),
		ObsBucketAccessIpAddressRanges:    obsBucketAccessIpAddressInfos,
		ObsBucketFileEncryptionAlgorithm:  data.Get("file_encryption_algorithm").(string),
		ObsBucketFileEncryptionEnabled:    data.Get("file_encryption_enabled").(bool),
		ObsBucketFileEncryptionType:       data.Get("file_encryption_type").(string),
		ObsBucketName:                     data.Get("name").(string),
		ObsBucketVersionEnabled:           data.Get("version_enabled").(bool),
		ObsId:                             data.Get("obs_id").(string),
		ZoneId:                            data.Get("zone_id").(string),
	})

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	err = waitForObjectStorageStatus(ctx, inst.Client, response.ObsBucketId, []string{}, []string{"Active"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ObsBucketId)

	return readObjectStorage(ctx, data, meta)
}

func readObjectStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.ObjectStorage.ReadObjectStorage(ctx, data.Id())

	s := common.HclSetObject{}
	for _, svc := range info.ObsBucketAccessIpAddressRanges {
		s = append(s, common.HclKeyValueObject{
			"obs_bucket_access_ip_address_range": svc.ObsBucketAccessIpAddressRange,
			"type":                               svc.Type_,
		})
	}
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("ip_address_filter_enabled", info.IsObsBucketIpAddressFilterEnabled)
	data.Set("access_ip_address_ranges", s)
	data.Set("file_encryption_algorithm", info.ObsBucketFileEncryptionAlgorithm)
	data.Set("file_encryption_enabled", info.ObsBucketFileEncryptionEnabled)
	data.Set("file_encryption_type", info.ObsBucketFileEncryptionType)
	data.Set("name", info.ObsBucketName)
	data.Set("version_enabled", info.ObsBucketVersionEnabled)
	data.Set("obs_id", info.ObsId)
	data.Set("zone_id", info.ZoneId)

	return nil
}

func updateObjectStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("version_enabled") {
		inst.Client.ObjectStorage.UpdateVersioning(ctx, data.Id(), data.Get("version_enabled").(bool))
	}

	if data.HasChanges("file_encryption_enabled") {
		inst.Client.ObjectStorage.UpdateEncryption(ctx, data.Id(), objectstorage.S3BucketUpdateRequest{
			ObsBucketFileEncryptionAlgorithm: data.Get("file_encryption_algorithm").(string),
			ObsBucketFileEncryptionEnabled:   data.Get("file_encryption_enabled").(bool),
			ObsBucketFileEncryptionType:      data.Get("file_encryption_type").(string),
			ObsBucketVersionEnabled:          data.Get("version_enabled").(bool),
		})
	}

	if data.HasChanges("ip_address_filter_enabled") ||
		(data.Get("ip_address_filter_enabled").(bool) && data.HasChanges("access_ip_address_ranges")) {
		obsBucketAccessIpAddressInfos, err := convertAccessIpAddressRanges(data.Get("access_ip_address_ranges").(common.HclListObject))
		inst.Client.ObjectStorage.CreateBucketIps(ctx, data.Id(), data.Get("ip_address_filter_enabled").(bool), obsBucketAccessIpAddressInfos)
		if err != nil {
			return nil
		}
	}

	return readObjectStorage(ctx, data, meta)
}

func deleteObjectStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.ObjectStorage.DeleteObjectStorage(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForObjectStorageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.ObjectStorage.ReadObjectStorage(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}

			return nil, "", err
		}
		if info.ObsId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.ObsBucketState, nil
	})
}
