package placementgroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/placementgroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/virtualserver"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
)

func init() {
	scp.RegisterResource("scp_placement_group", ResourcePlacementGroup())
}

func ResourcePlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlacementGroupCreate,
		ReadContext:   resourcePlacementGroupRead,
		UpdateContext: resourcePlacementGroupUpdate,
		DeleteContext: resourcePlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"placement_group_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validatePlacementGroupName,
				Description:      "Placement Group Name",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone Id",
			},
			"virtual_server_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validatePlacementGroupServerType,
				Description:      "Virtual Server Type",
			},
			"virtual_server_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Virtual Server Id List",
			},
			"availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Availability Zone Name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description",
			},
			"tags": tfTags.TagsSchema(),
		},
	}
}

func resourcePlacementGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	placementGroupName := rd.Get("placement_group_name").(string)
	serviceZoneId := rd.Get("service_zone_id").(string)
	virtualServerType := rd.Get("virtual_server_type").(string)
	virtualServerIds := rd.Get("virtual_server_ids").([]interface{})
	arrVirtualServerIds := make([]string, len(virtualServerIds))
	for i, virtualServerId := range virtualServerIds {
		arrVirtualServerIds[i] = virtualServerId.(string)
	}
	availabilityZoneName := rd.Get("availability_zone_name").(string)
	description := rd.Get("description").(string)
	tags := rd.Get("tags").(map[string]interface{})
	tagsRequests := make([]placementgroup.TagRequest, 0)
	for key, value := range tags {
		tagsRequests = append(tagsRequests, placementgroup.TagRequest{
			TagKey:   key,
			TagValue: value.(string),
		})
	}

	createResponse, err := inst.Client.PlacementGroup.CreatePlacementGroup(ctx, placementgroup.CreateRequest{
		AvailabilityZoneName:      availabilityZoneName,
		PlacementGroupName:        placementGroupName,
		ServiceZoneId:             serviceZoneId,
		Tags:                      rd.Get("tags").(map[string]interface{}),
		VirtualServerType:         virtualServerType,
		PlacementGroupDescription: description,
	})
	if err != nil {
		return
	}

	err = WaitForPlacementGroupState(ctx, inst.Client, createResponse.PlacementGroupId, []string{}, []string{common.ActiveState}, true)
	if err != nil {
		return
	}

	rd.SetId(createResponse.PlacementGroupId)
	return resourcePlacementGroupRead(ctx, rd, meta)
}

func resourcePlacementGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()
	inst := meta.(*client.Instance)
	placementGroupInfo, _, err := inst.Client.PlacementGroup.GetPlacementGroup(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("placement_group_name", placementGroupInfo.PlacementGroupName)
	rd.Set("description", placementGroupInfo.PlacementGroupDescription)
	rd.Set("availability_zone_name", placementGroupInfo.AvailabilityZoneName)
	rd.Set("virtual_server_type", placementGroupInfo.VirtualServerType)
	rd.Set("service_zone_id", placementGroupInfo.ServiceZoneId)
	rd.Set("virtual_server_ids", placementGroupInfo.VirtualServerIdList)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourcePlacementGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("virtual_server_ids") {
		oldVmIds, newVmIds := getOldAndNewVmIds(rd)
		deletedVmIds := getDeletedVmIds(oldVmIds, newVmIds)
		addedVmIds := getAddedVmIds(oldVmIds, newVmIds)

		for _, deletedVmId := range deletedVmIds {
			err := inst.Client.PlacementGroup.RemovePlacementGroupMember(ctx, rd.Id(), deletedVmId)
			if err != nil {
				return diag.FromErr(err)
			}

			err = virtualserver.WaitForVirtualServerStatus(ctx, inst.Client, deletedVmId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, addedVmId := range addedVmIds {
			err := inst.Client.PlacementGroup.AddPlacementGroupMember(ctx, rd.Id(), addedVmId)
			if err != nil {
				return diag.FromErr(err)
			}

			err = virtualserver.WaitForVirtualServerStatus(ctx, inst.Client, addedVmId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if rd.HasChanges("description") {
		description := rd.Get("description").(string)
		err := inst.Client.PlacementGroup.UpdatePlacementGroupDescription(ctx, rd.Id(), placementgroup.UpdatePlacementGroupDescriptionRequest{
			PlacementGroupDescription: description,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil

}

func resourcePlacementGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.PlacementGroup.DeletePlacementGroup(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = WaitForPlacementGroupState(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func WaitForPlacementGroupState(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.PlacementGroup.GetPlacementGroup(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, info.PlacementGroupState, nil
	})
}

func getOldAndNewVmIds(data *schema.ResourceData) ([]string, []string) {
	oldValue, newValue := data.GetChange("virtual_server_ids")
	oldValues := oldValue.([]interface{})
	newValues := newValue.([]interface{})
	oldVmIds := make([]string, len(oldValues))
	newVmIds := make([]string, len(newValues))
	for i, oldVmId := range oldValues {
		oldVmIds[i] = oldVmId.(string)
	}
	for i, newVmId := range newValues {
		newVmIds[i] = newVmId.(string)
	}
	return oldVmIds, newVmIds
}

func getDeletedVmIds(oldVmIds []string, newVmIds []string) []string {
	deletedVmIds := make([]string, 0)

	for _, oldVmId := range oldVmIds {
		var i int
		for i = 0; i < len(newVmIds); i++ {
			if strings.Compare(newVmIds[i], oldVmId) == 0 {
				break
			}
		}
		if i == len(newVmIds) {
			deletedVmIds = append(deletedVmIds, oldVmId)
		}
	}
	return deletedVmIds
}

func getAddedVmIds(oldVmIds []string, newVmIds []string) []string {
	addedVmIds := make([]string, 0)
	for _, newVmId := range newVmIds {
		var i int
		for i = 0; i < len(oldVmIds); i++ {
			if strings.Compare(newVmId, oldVmIds[i]) == 0 {
				break
			}
		}
		if i == len(oldVmIds) {
			addedVmIds = append(addedVmIds, newVmId)
		}
	}
	return addedVmIds
}

func validatePlacementGroupName(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	var err error = nil
	if len(value) < 3 {
		err = fmt.Errorf("input must be longer than %v characters", 3)
	} else if len(value) > 30 {
		err = fmt.Errorf("input must be shorter than %v characters", 30)
	} else {
		err = nil
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\-0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters with dash", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validatePlacementGroupServerType(v interface{}, path cty.Path) diag.Diagnostics {

	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	if value != "s1" && value != "h1" {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must be s1 or h1 value", attrKey),
			AttributePath: path,
		})
	}

	return diags
}
