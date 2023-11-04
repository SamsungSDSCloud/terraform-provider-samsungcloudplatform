package firewall

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/firewall2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterDataSource("scp_firewalls", DatasourceFirewalls())
}

func DatasourceFirewalls() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceFirewallsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC id",
			},
			"target_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target firewall resource id. (e.g. Internet Gateway, NAT Gateway, Load Balancer, ...)",
			},
			"firewalls": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Firewall list",
				Elem:        common.GetDatasourceItemsSchema(DatasourceFirewall()),
			},
		},
		Description: "Provides list of firewalls",
	}
}

func convertFirewallListToHclSet(firewalls []firewall2.FirewallListItemResponse) (common.HclSetObject, []string) {
	var firewallList common.HclSetObject
	var ids []string
	for _, fw := range firewalls {
		if len(fw.FirewallId) == 0 {
			continue
		}
		ids = append(ids, fw.FirewallId)
		kv := common.HclKeyValueObject{
			"id":          fw.FirewallId,
			"name":        fw.FirewallName,
			"state":       fw.FirewallState,
			"vpc_id":      fw.VpcId,
			"target_id":   fw.ObjectId,
			"target_type": fw.ObjectType,
		}
		firewallList = append(firewallList, kv)
	}
	return firewallList, ids
}

func datasourceFirewallsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	vpcId := rd.Get("vpc_id").(string)
	targetId := rd.Get("target_id").(string)

	firewalls, _, err := inst.Client.Firewall.GetFirewallList(ctx, vpcId, targetId, "")
	if err != nil {
		return diag.FromErr(err)
	}

	setFirewalls, ids := convertFirewallListToHclSet(firewalls.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		setFirewalls = common.ApplyFilter(DatasourceFirewall().Schema, f.(*schema.Set), setFirewalls)
	}
	rd.SetId(common.GenerateHash(ids))
	rd.Set("ids", ids)
	rd.Set("firewalls", setFirewalls)

	return nil
}
