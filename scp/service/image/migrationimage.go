package image

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_migration_image", ResourceMigrationImage())
}

func ResourceMigrationImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationImageCreate,
		ReadContext:   resourceMigrationImageRead,
		UpdateContext: resourceMigrationImageUpdate,
		DeleteContext: resourceMigrationImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "access key for ova",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "secret key for ova",
			},
			"az_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Availability Zone Name",
			},
			"image_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Migration Image Name",
				ValidateDiagFunc: common.ValidateName3to60AlphaNumericWithSpaceDashUnderscoreStartsWithLowerAlpha,
			},
			"ova_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ova url",
			},
			"os_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OS User Id",
			},
			"os_user_password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Os User Password",
			},
			"original_image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Original Image Id",
			},
			"image_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image Description",
			},
			"icon": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image state (ACTIVE)",
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image type (STANDARD, CUSTOM, MIGRATION)",
			},
			"origin_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS type (Windows, Ubuntu, ..)",
			},
			"product_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{Schema: elemProduct()},
			},
			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_dt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_dt": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Description: "Provides Migration Image resource.",
	}
}
func resourceMigrationImageCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	AccessKey := rd.Get("access_key").(string)
	SecretKey := rd.Get("secret_key").(string)
	AvailabilityZoneName := rd.Get("az_name").(string)
	ImageName := rd.Get("image_name").(string)
	OriginalImageId := rd.Get("original_image_id").(string)
	OsAdminCredential := image2.VirtualServerCreateOsCredentialRequest{
		OsUserId:       rd.Get("os_user_id").(string),
		OsUserPassword: rd.Get("os_user_password").(string),
	}
	OvaUrl := rd.Get("ova_url").(string)
	ServiceZoneId := rd.Get("service_zone_id").(string)
	ImageDescription := rd.Get("image_description").(string)

	createRequest := image2.MigrationImageCreateRequest{
		AccessKey:            AccessKey,
		SecretKey:            SecretKey,
		AvailabilityZoneName: AvailabilityZoneName,
		ImageName:            ImageName,
		OriginalImageId:      OriginalImageId,
		OsAdminCredential:    &OsAdminCredential,
		OvaUrl:               OvaUrl,
		ServiceZoneId:        ServiceZoneId,
		ImageDescription:     ImageDescription,
	}
	response, err := inst.Client.MigrationImage.CreateMigrationImage(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	err = waitMigrationImageCreating(ctx, inst.Client, response.ResourceId)
	rd.SetId(response.ResourceId)
	return resourceMigrationImageRead(ctx, rd, meta)
}

func resourceMigrationImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	MigrationImageInfo, _, err := inst.Client.MigrationImage.GetMigrationImageInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}
	rd.Set("name", MigrationImageInfo.ImageName)

	return nil
}
func resourceMigrationImageUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	if rd.HasChanges("image_description") {
		_, err := inst.Client.MigrationImage.UpdateMigrationImage(ctx, rd.Id(), rd.Get("image_description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMigrationImageRead(ctx, rd, meta)
}
func resourceMigrationImageDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	error := WaitMigrationImageStatus(ctx, inst.Client, rd.Id(), common.VirtualServerProcessingStates(), []string{common.ActiveState, common.StoppedState}, false)
	if error != nil {
		return diag.FromErr(error)
	}

	_, err := inst.Client.MigrationImage.DeleteMigrationImage(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = WaitMigrationImageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitMigrationImageCreating(ctx context.Context, scpClient *client.SCPClient, vpcId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		info, _, err := scpClient.MigrationImage.GetMigrationImageInfo(ctx, vpcId)
		if err != nil {
			return nil, "", err
		}
		return info, info.ImageState, nil
	})
}

func WaitMigrationImageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.MigrationImage.GetMigrationImageInfo(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, info.ImageState, nil
	})
}
