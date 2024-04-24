package image

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	image "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_custom_image", ResourceCustomImage())
}

func ResourceCustomImage() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceCustomImageCreate,
		ReadContext:   resourceCustomImageRead,
		UpdateContext: resourceCustomImageUpdate,
		DeleteContext: resourceCustomImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Custom image name.",
				ValidateDiagFunc: common.ValidateName3to60AlphaNumericWithSpaceDashUnderscoreStartsWithLowerAlpha,
			},
			"origin_virtual_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Origin virtual server id.",
			},
			"image_description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Custom image description.",
				ValidateDiagFunc: common.ValidateDescriptionMaxlength50,
			},
			"tags":                   tfTags.TagsSchema(),
			"project_id":             {Type: schema.TypeString, Computed: true},
			"availability_zone_name": {Type: schema.TypeString, Computed: true},
			"base_image":             {Type: schema.TypeString, Computed: true},
			"block_id":               {Type: schema.TypeString, Computed: true},
			"category":               {Type: schema.TypeString, Computed: true},
			"default_disk_size":      {Type: schema.TypeInt, Computed: true},
			"disk_size":              {Type: schema.TypeInt, Computed: true, Description: "Extra disk size."},
			"disks":                  {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: elemDisk()}},
			"icon":                   {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"image_id":               {Type: schema.TypeString, Computed: true},
			"image_state":            {Type: schema.TypeString, Computed: true, Description: "Image state (ACTIVE)"},
			"image_type":             {Type: schema.TypeString, Computed: true, Description: "Image type (STANDARD, CUSTOM, MIGRATION)"},
			"origin_image_id":        {Type: schema.TypeString, Computed: true},
			"origin_image_name":      {Type: schema.TypeString, Computed: true},
			"os_type":                {Type: schema.TypeString, Computed: true, Description: "OS type (Windows, Ubuntu, ..)"},
			"product_group_id":       {Type: schema.TypeString, Computed: true},
			"products":               {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: elemProduct()}},
			"properties":             {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"service_zone_id":        {Type: schema.TypeString, Computed: true},
			"created_by":             {Type: schema.TypeString, Computed: true},
			"created_dt":             {Type: schema.TypeString, Computed: true},
			"modified_by":            {Type: schema.TypeString, Computed: true},
			"modified_dt":            {Type: schema.TypeString, Computed: true},
		},
		Description: "Provides a Custom Image resource.",
	}
}

func resourceCustomImageCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, _, err := inst.Client.CustomImage.CreateCustomImage(ctx, image.CustomImageCreateRequest{
		ImageName:        rd.Get("image_name").(string),
		VirtualServerId:  rd.Get("origin_virtual_server_id").(string),
		ImageDescription: rd.Get("image_description").(string),
	}, rd.Get("tags").(map[string]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForCustomImageStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(response.ResourceId)
	return resourceCustomImageRead(ctx, rd, meta)
}

func resourceCustomImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	responseCustomImage, _, err := inst.Client.CustomImage.GetCustomImage(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	if rd.Set("image_id", responseCustomImage.ImageId) != nil {
		return nil
	}

	// response를 map형태로 변환 후 스키마에 반영
	mapCustomImageDetail := common.ToMap(responseCustomImage)
	for k, v := range mapCustomImageDetail {
		if rd.Set(k, v) != nil {
			return nil
		}
	}
	if rd.Set("disks", common.ConvertStructToMaps(responseCustomImage.Disks)) != nil {
		return nil
	}
	if rd.Set("products", common.ConvertStructToMaps(responseCustomImage.Products)) != nil {
		return nil
	}
	tfTags.SetTags(ctx, rd, meta, rd.Id())
	return nil
}

func resourceCustomImageUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if rd.HasChanges("image_description") {
		desc := rd.Get("image_description").(string)

		inst := meta.(*client.Instance)
		_, err := inst.Client.CustomImage.UpdateCustomImageDescription(ctx, rd.Id(), desc)

		if err != nil {
			diag.Errorf(err.Error())
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCustomImageRead(ctx, rd, meta)
}

func resourceCustomImageDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	err := inst.Client.CustomImage.DeleteCustomImage(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForCustomImageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCustomImageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.CustomImage.GetCustomImage(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.ImageState, nil
	})
}
