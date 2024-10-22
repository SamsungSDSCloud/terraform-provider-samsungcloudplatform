package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_kubernetes_subnet", DatasourceSubnet())
}

func DatasourceSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetCheck,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {Type: schema.TypeString, Required: true, Description: "Subnet Id"},
			"vpc_id":    {Type: schema.TypeString, Required: true, Description: "Vpc Id"},
			"result":    {Type: schema.TypeBool, Computed: true, Description: "Result"},
		},
		Description: "Check whether Subnet is usable for Kubernetes Engine or not (usable : true, not usable : false)\n\n",
	}
}

func dataSourceSubnetCheck(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.KubernetesEngine.CheckUsableSubnet(ctx, rd.Get("subnet_id").(string), rd.Get("vpc_id").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("result", responses.Result)

	return nil
}
