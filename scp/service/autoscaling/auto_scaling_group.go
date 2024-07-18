package autoscaling

import (
	"context"
	"errors"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
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
	scp.RegisterResource("scp_auto_scaling_group", ResourceAutoScalingGroup())
}

func ResourceAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutoScalingGroupCreate,
		ReadContext:   resourceAutoScalingGroupRead,
		UpdateContext: resourceAutoScalingGroupUpdate,
		DeleteContext: resourceAutoScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asg_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Auto-Scaling Group name. (3 to 20 using English letters, numbers and -)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]*$`), "Must be 3 to 20 using English letters, numbers and -."),
				),
			},
			"server_name_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server name prefix. (3 to 26 characters, starts with a lowercase letter, and uses lowercase letters, numbers, and -)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 26),
					validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9-]*$`), "Must be 3 to 26 characters, starts with a lowercase letter, and uses lowercase letters, numbers, and -."),
				),
			},
			"lc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Launch Configuration ID",
			},
			"min_server_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Min server count",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_server_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Max server count",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"desired_server_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Desired server count",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"desired_server_count_editable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Desired server count editable",
			},
			"vpc_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "VPC information",
				MaxItems:    1,
				Elem:        resourceVpcInfo(),
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Security Group ID list",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Availability zone name",
			},
			"multi_availability_zone_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable multi availability zone feature for this Auto-Scaling Group.",
			},
			"file_storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "File Storage ID",
			},
			"asg_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto-Scaling Group ID",
			},
			"asg_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto-Scaling Group state",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID",
			},
			"block_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block ID",
			},
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service ID",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service zone ID",
			},
			"lc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Launch Configuration name",
			},
			"is_terminating": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is terminating",
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
			"dns_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "DNS enabled",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Auto-Scaling Group resource.",
	}
}

func resourceVpcInfo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID",
			},
			"local_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local subnet ID",
			},
			"nat_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "NAT enabled",
			},
		},
	}
}

func resourceAutoScalingGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	info, _, err := inst.Client.AutoScaling.GetAutoScalingGroupDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(info, rd, "DeploymentEnvType", "ServiceLevelProductId", "SubnetId", "VpcId", "LocalSubnetId", "tags")
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceAutoScalingGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	// Get values from schema
	desiredServerCount := int32(rd.Get("desired_server_count").(int))
	desiredServerCountEditable := rd.Get("desired_server_count_editable").(bool)
	lcId := rd.Get("lc_id").(string)
	maxServerCount := int32(rd.Get("max_server_count").(int))
	minServerCount := int32(rd.Get("min_server_count").(int))
	multiAvailabilityZoneEnabled := rd.Get("multi_availability_zone_enabled").(bool)
	securityGroupIds := getActualStringList(rd, "security_group_ids")
	var vpcInfoReq autoscaling2.AsgVpcInfoRequest
	vpcInfoList := rd.Get("vpc_info").(common.HclListObject)
	vpcInfo := vpcInfoList[0].(map[string]interface{})
	vpcInfoReq.VpcId = vpcInfo["vpc_id"].(string)
	vpcInfoReq.SubnetId = vpcInfo["subnet_id"].(string)
	vpcInfoReq.LocalSubnetId = vpcInfo["local_subnet_id"].(string)
	natEnabled := vpcInfo["nat_enabled"].(bool)
	vpcInfoReq.NatEnabled = &natEnabled
	fileStorageId := rd.Get("file_storage_id").(string)

	// Validation
	validateServerCount(minServerCount, desiredServerCount, maxServerCount)
	if _, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcInfoReq.VpcId); err != nil {
		return diag.FromErr(err)
	}
	if _, _, err := inst.Client.Subnet.GetSubnet(ctx, vpcInfoReq.SubnetId); err != nil {
		return diag.FromErr(err)
	}
	if _, _, err := inst.Client.Subnet.GetSubnet(ctx, vpcInfoReq.LocalSubnetId); err != nil {
		return diag.FromErr(err)
	}
	if _, _, err := inst.Client.AutoScaling.GetLaunchConfigurationDetail(ctx, lcId); err != nil {
		return diag.FromErr(err)
	}

	createRequest := autoscaling2.AutoScalingGroupCreateV4Request{
		AsgName:                      rd.Get("asg_name").(string),
		AvailabilityZoneName:         rd.Get("availability_zone_name").(string),
		DesiredServerCount:           &desiredServerCount,
		DesiredServerCountEditable:   &desiredServerCountEditable,
		LcId:                         lcId,
		MaxServerCount:               &maxServerCount,
		MinServerCount:               &minServerCount,
		MultiAvailabilityZoneEnabled: &multiAvailabilityZoneEnabled,
		SecurityGroupIds:             securityGroupIds,
		ServerNamePrefix:             rd.Get("server_name_prefix").(string),
		VpcInfo:                      &vpcInfoReq,
		FileStorageId:                fileStorageId,
	}
	result, _, err := inst.Client.AutoScaling.CreateAutoScalingGroup(ctx, createRequest, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = waitForAutoScalingGroupStatus(ctx, inst.Client, result.AsgId, []string{}, []string{"In Service"}, true)
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.SetId(result.AsgId)

	return resourceAutoScalingGroupRead(ctx, rd, meta)
}

func validateServerCount(min int32, desired int32, max int32) error {
	if min > desired {
		return errors.New("min_server_count must not be greater than desired_server_count")
	}
	if desired > max {
		return errors.New("desired_server_count must not be greater than max_server_count")
	}
	return nil
}

func getActualStringList(rd *schema.ResourceData, key string) []string {
	rawStrings := rd.Get(key).([]interface{})
	actualStrings := make([]string, len(rawStrings))
	for i, v := range rawStrings {
		actualStrings[i] = v.(string)
	}
	return actualStrings
}

func resourceAutoScalingGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("desired_server_count_editable", "lc_id", "security_group_ids") {
		// update
		desiredServerCountEditable := rd.Get("desired_server_count_editable").(bool)
		lcId := rd.Get("lc_id").(string)
		securityGroupIds := getActualStringList(rd, "security_group_ids")
		updateRequest := autoscaling2.AutoScalingGroupUpdateRequest{
			DesiredServerCountEditable: &desiredServerCountEditable,
			LcId:                       lcId,
			SecurityGroupIds:           securityGroupIds,
		}
		_, _, err := inst.Client.AutoScaling.UpdateAutoScalingGroup(ctx, rd.Id(), updateRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("min_server_count", "desired_server_count", "max_server_count") {
		min := int32(rd.Get("min_server_count").(int))
		desired := int32(rd.Get("desired_server_count").(int))
		max := int32(rd.Get("max_server_count").(int))
		updateRequest := autoscaling2.AsgServerCountUpdateRequest{
			MinServerCount:     &min,
			DesiredServerCount: &desired,
			MaxServerCount:     &max,
		}

		_, _, err := inst.Client.AutoScaling.UpdateAutoScalingGroupServerCount(ctx, rd.Id(), updateRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		// wait for async job complete
		err = waitForAutoScalingGroupStatus(ctx, inst.Client, rd.Id(), []string{"Scale Out", "Scale In", "Attach to LB", "Detach from LB"}, []string{"In Service"}, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAutoScalingGroupRead(ctx, rd, meta)
}

func resourceAutoScalingGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.AutoScaling.DeleteAutoScalingGroup(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	err = waitForAutoScalingGroupStatus(ctx, inst.Client, rd.Id(), []string{"Terminating", "Detach from LB"}, []string{"Terminated"}, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForAutoScalingGroupStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.AutoScaling.GetAutoScalingGroupDetail(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "Terminated", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "Terminated", nil
			}
			return nil, "", err
		}
		return info, info.AsgState, nil
	})
}
