package securitygroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_security_group", ResourceSecurityGroup())
	scp.RegisterResource("scp_security_group_user_ip", ResourceSecurityGroupUserIp())
}

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC id",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Subnet name. (3 to 20 lowercase characters with - and _)",
				ValidateDiagFunc: common.ValidateName3to20DashUnderscore,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Subnet description",
				ValidateDiagFunc: common.ValidateDescriptionMaxlength50,
			},
			"is_loggable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
		},
		Description: "Provides a Security Group resource.",
	}
}

func resourceSecurityGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	vpcId := rd.Get("vpc_id").(string)
	name := rd.Get("name").(string)
	description := rd.Get("description").(string)
	isLoggable := rd.Get("is_loggable").(bool)

	inst := meta.(*client.Instance)

	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	isNameInvalid, err := inst.Client.SecurityGroup.CheckSecurityGroupName(ctx, name, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input security group name is invalid (maybe duplicated) : " + name)
	}

	productGroupId, err := client.FindProductGroupId(ctx, inst.Client, vpcInfo.ServiceZoneId, common.NetworkProductGroup, common.SecurityGroupProductName)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := inst.Client.SecurityGroup.CreateSecurityGroup(ctx, productGroupId, vpcInfo.ServiceZoneId, vpcId, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSecurityGroupStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	if isLoggable {
		_, err = inst.Client.SecurityGroup.UpdateSecurityGroupIsLoggable(ctx, response.ResourceId, isLoggable)
		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForSecurityGroupStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	rd.SetId(response.ResourceId)

	return resourceSecurityGroupRead(ctx, rd, meta)
}

func resourceSecurityGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	info, _, err := inst.Client.SecurityGroup.GetSecurityGroup(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("vpc_id", info.VpcId)
	rd.Set("name", info.SecurityGroupName)
	rd.Set("description", info.SecurityGroupDescription)
	rd.Set("is_loggable", info.IsLoggable)

	return nil
}

func resourceSecurityGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("description") {
		_, err := inst.Client.SecurityGroup.UpdateSecurityGroupDescription(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rd.HasChanges("is_loggable") {
		_, err := inst.Client.SecurityGroup.UpdateSecurityGroupIsLoggable(ctx, rd.Id(), rd.Get("is_loggable").(bool))
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForSecurityGroupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceSecurityGroupRead(ctx, rd, meta)
}

func resourceSecurityGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	isLoggable := rd.Get("is_loggable").(bool)

	if isLoggable {
		_, err := inst.Client.SecurityGroup.UpdateSecurityGroupIsLoggable(ctx, rd.Id(), false)
		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForSecurityGroupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := inst.Client.SecurityGroup.DeleteSecurityGroup(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForSecurityGroupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSecurityGroupStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.SecurityGroup.GetSecurityGroup(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.SecurityGroupState, nil
	})
}

func ResourceSecurityGroupUserIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupUserIpCreate,
		ReadContext:   resourceSecurityGroupUserIpRead,
		DeleteContext: resourceSecurityGroupUserIpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target SecurityGroup id",
			},
			"user_ip_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of Directly Attached IP",
			},
			"user_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Address of Directly Attached IP",
			},
			"user_ip_description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Description of Directly Attached IP",
			},
		},
		Description: "Provides a Security Group User IP Attachment resource.",
	}
}

func resourceSecurityGroupUserIpCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	sgId := rd.Get("security_group_id").(string)
	userIpType := rd.Get("user_ip_type").(string)
	userIpAddress := rd.Get("user_ip_address").(string)
	userIpDescription := rd.Get("user_ip_description").(string)

	inst := meta.(*client.Instance)

	response, err := inst.Client.SecurityGroup.AttachUserIpToSecurityGroup(ctx, sgId, userIpType, userIpAddress, userIpDescription)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.Set("security_group_id", response.ResourceId)
	rd.SetId(userIpAddress)

	return nil
}

func resourceSecurityGroupUserIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityGroupUserIpDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.SecurityGroup.DetachUserIpFromSecurityGroup(ctx, rd.Get("security_group_id").(string), rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
