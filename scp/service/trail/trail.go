package trail

import (
	"context"
	"fmt"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/client/loggingaudit"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTrail() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTrail,
		ReadContext:   readTrail,
		UpdateContext: updateTrail,
		DeleteContext: deleteTrail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"is_logging_target_all_resource": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"is_logging_target_all_user": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"logging_target_resource_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "",
			},
			"logging_target_users": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
					},
				},
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"save_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
		},
	}
}

func convertUsers(list common.HclListObject) ([]loggingaudit.UserResponse, error) {
	var result []loggingaudit.UserResponse
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := loggingaudit.UserResponse{}
		if userId, ok := itemObject["user_id"]; ok {
			info.UserId = userId.(string)
		}

		result = append(result, info)
	}
	return result, nil
}

func getResourceIds(rd *schema.ResourceData) []string {
	resourceIds := rd.Get("logging_target_resource_ids").([]interface{})
	ids := make([]string, len(resourceIds))
	for i, valueIpv4 := range resourceIds {
		ids[i] = valueIpv4.(string)
	}
	return ids
}

func createTrail(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	users, err := convertUsers(data.Get("logging_target_users").(common.HclListObject))
	if err != nil {
		return nil
	}

	response, err := inst.Client.Loggingaudit.CreateTrail(ctx, loggingaudit.CreateTrailRequest{
		TrailName:                  data.Get("name").(string),
		ObsBucketId:                data.Get("obs_bucket_id").(string),
		IsLoggingTargetAllUser:     data.Get("is_logging_target_all_user").(string),
		LoggingTargetUsers:         users,
		IsLoggingTargetAllResource: data.Get("is_logging_target_all_resource").(string),
		LoggingTargetResourceIds:   getResourceIds(data),
		TrailSaveType:              data.Get("save_type").(string),
		TrailDescription:           data.Get("description").(string),
	})

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	err = waitForTrailStatus(ctx, inst.Client, response.ObsBucketId, []string{}, []string{"Active"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.TrailId)

	return readTrail(ctx, data, meta)
}

func readTrail(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.Loggingaudit.ReadTrail(ctx, data.Id())

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("name", info.TrailName)
	data.Set("is_logging_target_all_resource", info.IsLoggingTargetAllResource)
	data.Set("is_logging_target_all_user", info.IsLoggingTargetAllUser)
	data.Set("logging_target_resource_ids", info)
	data.Set("logging_target_users", info)
	data.Set("obs_bucket_id", info.ObsBucketId)
	data.Set("save_type", info.TrailSaveType)
	data.Set("description", info.TrailDescription)

	return nil
}

func updateTrail(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteTrail(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
