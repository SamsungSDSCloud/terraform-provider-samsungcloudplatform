package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_kubernetes_engines", DatasourceEngines())
}

func DatasourceEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"k8s_version":              {Type: schema.TypeString, Optional: true, Description: "K8s cluster version"},
			"kubernetes_engine_name":   {Type: schema.TypeString, Optional: true, Description: "K8s engine name"},
			"kubernetes_engine_status": {Type: schema.TypeString, Optional: true, Description: "K8s engine status"},
			"region":                   {Type: schema.TypeString, Optional: true, Description: "Region"},
			"created_by":               {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":                     {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                     {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                 {Type: schema.TypeList, Optional: true, Description: "K8s engine list", Elem: datasourceElem()},
			"total_count":              {Type: schema.TypeInt, Computed: true, Description: "Content list size"},
		},
		Description: "Provides list of K8s engines",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.KubernetesEngine.GetEngineList(ctx, &kubernetesengine2.K8sEngineV2ApiListKubernetesEnginesV2Opts{
		K8sVersion:             optional.NewInterface(rd.Get("k8s_version").(string)),
		KubernetesEngineName:   optional.NewString(rd.Get("kubernetes_engine_name").(string)),
		KubernetesEngineStatus: optional.NewInterface(rd.Get("kubernetes_engine_status").(string)),
		Region:                 optional.NewInterface(rd.Get("region").(string)),
		CreatedBy:              optional.NewString(rd.Get("created_by").(string)),
		Page:                   optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                   optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                   optional.String{},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":               {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"k8s_version":              {Type: schema.TypeString, Computed: true, Description: "K8s version"},
			"kubernetes_engine_id":     {Type: schema.TypeString, Computed: true, Description: "K8s engine id"},
			"kubernetes_engine_name":   {Type: schema.TypeString, Computed: true, Description: "K8s engine name"},
			"kubernetes_engine_status": {Type: schema.TypeString, Computed: true, Description: "K8s engine status"},
			"node_count":               {Type: schema.TypeInt, Computed: true, Description: "K8s node count"},
			"region":                   {Type: schema.TypeString, Computed: true, Description: "Region name"},
			"security_group_id":        {Type: schema.TypeString, Computed: true, Description: "Security group id"},
			"subnet_id":                {Type: schema.TypeString, Computed: true, Description: "Subnet id"},
			"volume_id":                {Type: schema.TypeString, Computed: true, Description: "File storage volume id"},
			"vpc_id":                   {Type: schema.TypeString, Computed: true, Description: "Vpc id"},
			"created_by":               {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":               {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			"modified_by":              {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":              {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
