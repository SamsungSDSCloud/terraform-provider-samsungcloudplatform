package trail

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_trail", ResourceTrail())
}

func ResourceTrail() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrailCreate,
		ReadContext:   resourceTrailRead,
		UpdateContext: resourceTrailUpdate,
		DeleteContext: resourceTrailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Trail name",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(5, 20),
					validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z_-]+$`), "must contain only alphanumeric,-,_ characters"),
				)),
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Object storage bucket ID",
			},
			"is_logging_target_all_user": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether for all users",
			},
			"logging_target_users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target user ID list",
			},
			"is_logging_target_all_resource": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to log all resources",
			},
			"logging_target_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target resource ID list",
			},
			"is_logging_target_all_region": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to target all regions",
			},
			"logging_target_regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target regions list",
			},
			"save_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trail save type. JSON or CSV",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"JSON",
					"CSV",
				}, false)),
			},
			"use_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use trail verification",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         false,
				Description:      "",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 400)),
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"ACTIVE",
					"STOPPED",
				}, false)),
			},
			"obs_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage name",
			},
			"obs_folder_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage folder name",
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage bucket name",
			},
			"batch_end_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Batch processing end date and time",
			},
			"batch_first_start_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Batch processing first start date and time",
			},
			"batch_last_success_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time of last successful batch processing",
			},
			"batch_start_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Batch processing start date and time",
			},
			"batch_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "batch processing status",
			},
			"tags": tfTags.TagsSchema(),
		},
	}
}

func resourceTrailCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	result, err := inst.Client.Loggingaudit.CheckTrailName(ctx, rd.Get("name").(string))
	if !result {
		return diag.Errorf("trailName is already exists.")
	}
	if err != nil {
		return
	}

	tags := rd.Get("tags").(map[string]interface{})
	request := loggingaudit.CreateTrailRequest{
		TrailName:                  rd.Get("name").(string),
		ObsBucketId:                rd.Get("obs_bucket_id").(string),
		IsLoggingTargetAllUser:     rd.Get("is_logging_target_all_user").(bool),
		IsLoggingTargetAllResource: rd.Get("is_logging_target_all_resource").(bool),
		TrailSaveType:              rd.Get("save_type").(string),
		TrailDescription:           rd.Get("description").(string),
		IsLoggingTargetAllRegion:   rd.Get("is_logging_target_all_region").(bool),
		UseVerification:            rd.Get("use_verification").(bool),
	}
	if !request.IsLoggingTargetAllUser {
		request.LoggingTargetUsers = common.ToStringList(rd.Get("logging_target_users").(*schema.Set).List())
	}
	if !request.IsLoggingTargetAllResource {
		request.LoggingTargetResourceIds = common.ToStringList(rd.Get("logging_target_resource_ids").(*schema.Set).List())
	}
	if !request.IsLoggingTargetAllRegion {
		request.LoggingTargetRegions = common.ToStringList(rd.Get("logging_target_regions").(*schema.Set).List())
	}

	response, err := inst.Client.Loggingaudit.CreateTrail(ctx, tags, request)

	if err != nil {
		return
	}

	err = waitForTrailStatus(ctx, inst.Client, response.TrailId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return
	}

	rd.SetId(response.TrailId)

	return resourceTrailRead(ctx, rd, meta)
}

func resourceTrailRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.Loggingaudit.ReadTrail(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", info.TrailName)
	rd.Set("is_logging_target_all_user", getBoolFromYn(info.IsLoggingTargetAllUser))
	rd.Set("is_logging_target_all_region", getBoolFromYn(info.IsLoggingTargetAllRegion))
	rd.Set("is_logging_target_all_resource", getBoolFromYn(info.IsLoggingTargetAllResource))
	rd.Set("state", info.TrailState)

	if v := info.LoggingTargetResourceIds; v != nil {
		rd.Set("logging_target_resource_ids", info.LoggingTargetResourceIds)
	}

	if v := info.LoggingTargetUsers; v != nil {
		var users []string
		for _, itemObject := range info.LoggingTargetUsers {
			users = append(users, itemObject.UserId)
		}
		rd.Set("logging_target_users", users)
	}

	if v := info.LoggingTargetRegions; v != nil {
		rd.Set("logging_target_regions", info.LoggingTargetRegions)
	}

	rd.Set("obs_bucket_id", info.ObsBucketId)
	rd.Set("save_type", info.TrailSaveType)
	rd.Set("description", info.TrailDescription)
	rd.Set("use_verification", getBoolFromYn(info.ValidationYn))

	rd.Set("obs_name", info.ObjectStorageName)
	rd.Set("obs_folder_name", info.ObjectStorageFolderName)
	rd.Set("obs_bucket_name", info.ObsBucketName)
	rd.Set("batch_end_dt", info.TrailBatchEndDt)
	rd.Set("batch_first_start_dt", info.TrailBatchFirstStartDt)
	rd.Set("batch_last_success_dt", info.TrailBatchLastSuccessDt)
	rd.Set("batch_start_dt", info.TrailBatchStartDt)
	rd.Set("batch_state", info.TrailBatchState)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceTrailUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var updateTrail loggingaudit.UpdateTrailRequest

	if rd.HasChanges("description") {
		updateTrail.TrailUpdateType = "description"
		updateTrail.TrailDescription = rd.Get("description").(string)

		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("logging_target_resource_ids") || rd.HasChanges("is_logging_target_all_resource") {
		updateTrail.TrailUpdateType = "resource"
		allTarget := rd.Get("is_logging_target_all_resource").(bool)
		updateTrail.IsLoggingTargetAllResource = getYnFromBool(allTarget)
		if !allTarget {
			updateTrail.LoggingTargetResourceIds = common.ToStringList(rd.Get("logging_target_resource_ids").(*schema.Set).List())
		}

		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("logging_target_users") || rd.HasChanges("is_logging_target_all_user") {
		updateTrail.TrailUpdateType = "user"
		allTarget := rd.Get("is_logging_target_all_user").(bool)
		updateTrail.IsLoggingTargetAllUser = getYnFromBool(allTarget)
		if !allTarget {
			updateTrail.LoggingTargetUsers = common.ToStringList(rd.Get("logging_target_users").(*schema.Set).List())
		}

		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("save_type") {
		updateTrail.TrailUpdateType = "saveType"
		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("logging_target_regions") || rd.HasChanges("is_logging_target_all_region") {
		updateTrail.TrailUpdateType = "region"
		allTarget := rd.Get("is_logging_target_all_region").(bool)
		updateTrail.IsLoggingTargetAllRegion = getYnFromBool(allTarget)

		if !allTarget {
			updateTrail.LoggingTargetRegions = common.ToStringList(rd.Get("logging_target_regions").(*schema.Set).List())
		}

		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("use_verification") {
		updateTrail.TrailUpdateType = "validation"
		updateTrail.UseVerification = rd.Get("use_verification").(bool)

		_, err := inst.Client.Loggingaudit.UpdateTrail(ctx, rd.Id(), updateTrail)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("state") {
		state := rd.Get("state").(string)
		if state == "STOPPED" {
			_, _, err := inst.Client.Loggingaudit.StopTrail(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		} else if state == "ACTIVE" {
			_, _, err := inst.Client.Loggingaudit.StartTrail(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.Errorf("Wrong state : %s", state)
		}

		err := waitForTrailStatus(ctx, inst.Client, rd.Id(), []string{}, []string{state}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTrailRead(ctx, rd, meta)
}

func resourceTrailDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, _, err := inst.Client.Loggingaudit.DeleteTrail(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForTrailStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Loggingaudit.ReadTrail(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}

			return nil, "", err
		}
		if info.TrailId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.TrailState, nil
	})
}

func getBoolFromYn(flag string) bool {
	if flag == "Y" {
		return true
	} else if flag == "N" {
		return false
	} else {
		fmt.Errorf("wrong string flag. (Not Y or N)")
		return false
	}
}

func getYnFromBool(flag bool) string {
	if flag {
		return "Y"
	} else {
		return "N"
	}
}
