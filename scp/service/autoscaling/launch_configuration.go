package autoscaling

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/autoscaling/autoscaling_common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_launch_configuration", ResourceLaunchConfiguration())
}

func ResourceLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLaunchConfigurationCreate,
		ReadContext:   resourceLaunchConfigurationRead,
		UpdateContext: resourceLaunchConfigurationUpdate,
		DeleteContext: resourceLaunchConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"block_storages": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        resourceLaunchConfigurationBlockStorageElem(),
				Description: "Block Storage list",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image ID",
			},
			"initial_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Virtual Server's initial script",
			},
			"key_pair_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key pair ID",
			},
			"lc_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Launch Configuration name",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-z0-9-]*$`), "Must be 3 to 20 using English letters, numbers and -."),
				),
			},
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server type",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service zone ID",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID",
			},
			"asg_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Auto-Scaling Group ID list",
			},
			"block_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block ID",
			},
			"contract_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contract product ID",
			},
			"lc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Launch Configuration ID",
			},
			"os_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS product ID",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS type",
			},
			"product_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product group ID",
			},
			"scale_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scale product ID",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The person who created the resource",
			},
			"created_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The person who modified the resource",
			},
			"modified_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification date",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Launch Configuration resource.",
	}
}

func resourceLaunchConfigurationBlockStorageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "Block Storage size (GB)",
				ValidateFunc: validation.IntAtLeast(4),
			},
			"disk_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Block Storage product (default value : SSD)",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`SSD|HDD`), "Must be one of \"SSD\" or \"HDD\"."),
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Encryption enabled",
			},
			"is_boot_disk": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Is boot disk or not",
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product ID",
			},
		},
	}
}

func resourceLaunchConfigurationCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	response, _, err := inst.Client.AutoScaling.CreateLaunchConfigurationGroup(ctx, autoscaling2.LaunchConfigCreateV6Request{
		BlockStorages: convertBlockStorages(rd.Get("block_storages").(common.HclListObject)),
		ImageId:       rd.Get("image_id").(string),
		InitialScript: rd.Get("initial_script").(string),
		KeyPairId:     rd.Get("key_pair_id").(string),
		LcName:        rd.Get("lc_name").(string),
		ServerType:    rd.Get("server_type").(string),
		ServiceZoneId: rd.Get("service_zone_id").(string),
	}, rd.Get("tags").(map[string]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.LcId)

	return resourceLaunchConfigurationRead(ctx, rd, meta)
}

func convertBlockStorages(list common.HclListObject) []autoscaling2.LaunchConfigBlockStorageV2Request {
	var result []autoscaling2.LaunchConfigBlockStorageV2Request
	for _, itemObject := range list {
		item := itemObject.(common.HclKeyValueObject)
		var request autoscaling2.LaunchConfigBlockStorageV2Request
		if v, ok := item["block_storage_size"]; ok {
			size := int32(v.(int))
			request.BlockStorageSize = &size
		}
		if v, ok := item["disk_type"]; ok {
			request.DiskType = v.(string)
		}
		if v, ok := item["encryption_enabled"]; ok {
			request.EncryptionEnabled = new(bool)
			*request.EncryptionEnabled = v.(bool)
		}
		if v, ok := item["is_boot_disk"]; ok {
			request.IsBootDisk = new(bool)
			*request.IsBootDisk = v.(bool)
		}
		result = append(result, request)
	}
	return result
}

func resourceLaunchConfigurationRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, _, err := inst.Client.AutoScaling.GetLaunchConfigurationDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(response, rd, "tags")
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceLaunchConfigurationUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return
	}

	return resourceLaunchConfigurationRead(ctx, rd, meta)
}

func resourceLaunchConfigurationDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.AutoScaling.DeleteLaunchConfigurationGroup(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	return nil
}
