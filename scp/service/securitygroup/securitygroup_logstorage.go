package securitygroup

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	securitygroup2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/security-group2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_security_group_logstorage", resourceSecurityGroupLogStorage())
}

func resourceSecurityGroupLogStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupLogStorageCreate,
		ReadContext:   resourceSecurityGroupLogStorageRead,
		UpdateContext: resourceSecurityGroupLogStorageUpdate,
		DeleteContext: resourceSecurityGroupLogStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id",
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object storage bucket id to save Security Group log",
			},
		},
		Description: "Set up Security Group log storage.",
	}
}

func resourceSecurityGroupLogStorageCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	vpcId := rd.Get("vpc_id").(string)
	obsBucketId := rd.Get("obs_bucket_id").(string)

	SecurityGroupLogStorageCreatRequest := securitygroup2.SecurityGroupLogStorageCreatRequest{
		ObsBucketId: obsBucketId,
		VpcId:       vpcId,
	}
	logStorage, _, err := inst.Client.SecurityGroup.CreateSecurityGroupLogStorage(ctx, SecurityGroupLogStorageCreatRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(logStorage.LogStorageId)

	return resourceSecurityGroupLogStorageRead(ctx, rd, meta)
}

func resourceSecurityGroupLogStorageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	logStageInfo, _, err := inst.Client.SecurityGroup.GetSecurityGroupLogStorage(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set("vpc_id", logStageInfo.VpcId)
	rd.Set("obs_bucket_id", logStageInfo.ObsBucketId)

	return nil
}

func resourceSecurityGroupLogStorageUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("obs_bucket_id") {
		result, err := inst.Client.SecurityGroup.ListSecurityGroupsByLoggable(ctx, rd.Get("vpc_id").(string), true)
		if err != nil {
			return diag.FromErr(err)
		}
		if result.TotalCount > 0 {
			return diag.FromErr(fmt.Errorf("the bucket cannot be changed while logs are in use"))
		}

		obsBucketId := rd.Get("obs_bucket_id").(string)

		if _, _, err := inst.Client.SecurityGroup.UpdateSecurityGroupLogStorage(ctx, rd.Id(), obsBucketId); err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceSecurityGroupLogStorageRead(ctx, rd, meta)
}

func resourceSecurityGroupLogStorageDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	result, err := inst.Client.SecurityGroup.ListSecurityGroupsByLoggable(ctx, rd.Get("vpc_id").(string), true)
	if err != nil {
		return diag.FromErr(err)
	}
	if result.TotalCount > 0 {
		return diag.FromErr(fmt.Errorf("logs are in use and cannot be deleted"))
	}
	if err := inst.Client.SecurityGroup.DeleteSecurityGroupLogStorage(ctx, rd.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
