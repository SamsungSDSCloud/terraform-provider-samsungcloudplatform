package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceNodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodePoolList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"kubernetes_engine_id": {Type: schema.TypeString, Required: true, Description: "K8s engine id"},
			"node_pool_name":       {Type: schema.TypeString, Optional: true, Description: "K8s NodePool name"},
			"created_by":           {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":                 {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                 {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":             {Type: schema.TypeList, Optional: true, Description: "K8s Node Pool list", Elem: datasourceNodePoolElem()},
			"total_count":          {Type: schema.TypeInt, Computed: true, Description: "Total list count"},
		},
		Description: "Provides list of K8s node pools",
	}
}

func dataSourceNodePoolList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	engineId := rd.Get("kubernetes_engine_id").(string)
	if len(engineId) == 0 {
		return diag.Errorf("kubernetes engine id not found")
	}

	responses, _, err := inst.Client.KubernetesEngine.GetNodePoolList(ctx, engineId, &kubernetesengine2.NodePoolV2ControllerApiListNodePoolsV2Opts{
		NodePoolName: optional.NewString(rd.Get("node_pool_name").(string)),
		CreatedBy:    optional.NewString(rd.Get("created_by").(string)),
		Page:         optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:         optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:         optional.String{},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	ids := uuid.NewV4().String()
	//rd.SetId(ids)
	rd.SetId(ids)
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceNodePoolElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":         {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"auto_recovery":      {Type: schema.TypeBool, Computed: true, Description: "Enable auto recovery"},
			"auto_scale":         {Type: schema.TypeBool, Computed: true, Description: "Enable auto scale"},
			"contract_id":        {Type: schema.TypeString, Computed: true, Description: "Contract id"},
			"current_node_count": {Type: schema.TypeInt, Computed: true, Description: "Current node count in the pool"},
			"desired_node_count": {Type: schema.TypeInt, Computed: true, Description: "Desired node count in the pool"},
			"image_id":           {Type: schema.TypeString, Computed: true, Description: "Image id"},
			"in_progress":        {Type: schema.TypeBool, Computed: true, Description: "Check inProgress status"},
			"k8s_version":        {Type: schema.TypeString, Computed: true, Description: "K8s version"},
			"max_node_count":     {Type: schema.TypeInt, Computed: true, Description: "Maximum node count"},
			"min_node_count":     {Type: schema.TypeInt, Computed: true, Description: "Minimum node count"},
			"node_pool_id":       {Type: schema.TypeString, Computed: true, Description: "NodePool id"},
			"node_pool_name":     {Type: schema.TypeString, Computed: true, Description: "NodePool name"},
			"node_pool_state":    {Type: schema.TypeString, Computed: true, Description: "NodePool status"},
			"product_group_id":   {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			"scale_id":           {Type: schema.TypeString, Computed: true, Description: "Scale id"},
			"service_level_id":   {Type: schema.TypeString, Computed: true, Description: "Service level id"},
			"storage_id":         {Type: schema.TypeString, Computed: true, Description: "Storage id"},
			"storage_size":       {Type: schema.TypeString, Computed: true, Description: "Storage size in GB"},
			"os_type":            {Type: schema.TypeString, Computed: true, Description: "Host OS type (Ubuntu, Window,..)"},
			"upgradable":         {Type: schema.TypeBool, Computed: true, Description: "Where to enable upgrade"},
			"created_by":         {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":         {Type: schema.TypeString, Computed: true, Description: "Creation Date"},
			"modified_by":        {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":        {Type: schema.TypeString, Computed: true, Description: "Modification Date"},
		},
	}
}
