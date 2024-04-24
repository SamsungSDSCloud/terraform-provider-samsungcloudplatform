package kubernetes

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_kubernetes_node_pools", DatasourceNodePools())
	scp.RegisterDataSource("scp_kubernetes_node_pool", DatasourceNodePool())
}

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

	responses, _, err := inst.Client.KubernetesEngine.GetNodePoolList(ctx, engineId, &kubernetesengine2.NodePoolV2ApiListNodePoolsV2Opts{
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
			"region":             {Type: schema.TypeString, Computed: true, Description: "Modification Date"},
		},
	}
}

func DatasourceNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: nodePoolDetail, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("KubernetesEngineId"): {Type: schema.TypeString, Required: true, Description: "Kubernetes Engine Id"},
			common.ToSnakeCase("NodePoolId"):         {Type: schema.TypeString, Required: true, Description: "Node Pool Id"},
			common.ToSnakeCase("ProjectId"):          {Type: schema.TypeString, Computed: true, Description: "Project Id"},

			common.ToSnakeCase("NodePoolName"):     {Type: schema.TypeString, Computed: true, Description: "NodePoolName"},
			common.ToSnakeCase("K8sVersion"):       {Type: schema.TypeString, Computed: true, Description: "K8sVersion"},
			common.ToSnakeCase("ImageId"):          {Type: schema.TypeString, Computed: true, Description: "ImageId"},
			common.ToSnakeCase("NodePoolStatus"):   {Type: schema.TypeString, Computed: true, Description: "NodePoolStatus"},
			common.ToSnakeCase("Upgradable"):       {Type: schema.TypeBool, Computed: true, Description: "Upgradable"},
			common.ToSnakeCase("AutoScale"):        {Type: schema.TypeBool, Computed: true, Description: "AutoScale"},
			common.ToSnakeCase("MinNodeCount"):     {Type: schema.TypeInt, Computed: true, Description: "MinNodeCount"},
			common.ToSnakeCase("MaxNodeCount"):     {Type: schema.TypeInt, Computed: true, Description: "MaxNodeCount"},
			common.ToSnakeCase("AutoRecovery"):     {Type: schema.TypeBool, Computed: true, Description: "AutoRecovery"},
			common.ToSnakeCase("EncryptEnabled"):   {Type: schema.TypeBool, Computed: true, Description: "EncryptEnabled"},
			common.ToSnakeCase("ProductGroupId"):   {Type: schema.TypeString, Computed: true, Description: "ProductGroupId"},
			common.ToSnakeCase("StorageId"):        {Type: schema.TypeString, Computed: true, Description: "StorageId"},
			common.ToSnakeCase("StorageSize"):      {Type: schema.TypeString, Computed: true, Description: "StorageSize"},
			common.ToSnakeCase("ScaleId"):          {Type: schema.TypeString, Computed: true, Description: "ScaleId"},
			common.ToSnakeCase("ContractId"):       {Type: schema.TypeString, Computed: true, Description: "ContractId"},
			common.ToSnakeCase("ServiceLevelId"):   {Type: schema.TypeString, Computed: true, Description: "ServiceLevelId"},
			common.ToSnakeCase("CurrentNodeCount"): {Type: schema.TypeInt, Computed: true, Description: "CurrentNodeCount"},
			common.ToSnakeCase("DesiredNodeCount"): {Type: schema.TypeInt, Computed: true, Description: "DesiredNodeCount"},
			common.ToSnakeCase("CreatedBy"):        {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):        {Type: schema.TypeString, Computed: true, Description: "Created Dt"},
			common.ToSnakeCase("ModifiedBy"):       {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):       {Type: schema.TypeString, Computed: true, Description: "Modified Dt"},
		},
		Description: "Provides Kubernetes Engine Detail",
	}
}

func nodePoolDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	engineId := rd.Get("kubernetes_engine_id").(string)
	nodePoolId := rd.Get("node_pool_id").(string)

	response, _, err := inst.Client.KubernetesEngine.ReadNodePool(ctx, engineId, nodePoolId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())

	rd.Set(common.ToSnakeCase("ProjectId"), response.ProjectId)
	rd.Set(common.ToSnakeCase("NodePoolName"), response.NodePoolName)
	rd.Set(common.ToSnakeCase("K8sVersion"), response.K8sVersion)
	rd.Set(common.ToSnakeCase("ImageId"), response.ImageId)
	rd.Set(common.ToSnakeCase("NodePoolStatus"), response.NodePoolStatus)
	rd.Set(common.ToSnakeCase("Upgradable"), response.Upgradable)
	rd.Set(common.ToSnakeCase("AutoScale"), response.AutoScale)
	rd.Set(common.ToSnakeCase("MinNodeCount"), response.MinNodeCount)
	rd.Set(common.ToSnakeCase("MaxNodeCount"), response.MaxNodeCount)
	rd.Set(common.ToSnakeCase("AutoRecovery"), response.AutoRecovery)
	rd.Set(common.ToSnakeCase("EncryptEnabled"), response.EncryptEnabled)
	rd.Set(common.ToSnakeCase("ProductGroupId"), response.ProductGroupId)
	rd.Set(common.ToSnakeCase("StorageId"), response.StorageId)
	rd.Set(common.ToSnakeCase("StorageSize"), response.StorageSize)
	rd.Set(common.ToSnakeCase("ScaleId"), response.ScaleId)
	rd.Set(common.ToSnakeCase("ContractId"), response.ContractId)
	rd.Set(common.ToSnakeCase("ServiceLevelId"), response.ServiceLevelId)
	rd.Set(common.ToSnakeCase("CurrentNodeCount"), response.CurrentNodeCount)
	rd.Set(common.ToSnakeCase("DesiredNodeCount"), response.DesiredNodeCount)
	rd.Set(common.ToSnakeCase("CreatedBy"), response.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), response.CreatedDt.String())
	rd.Set(common.ToSnakeCase("ModifiedBy"), response.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), response.ModifiedDt.String())

	return nil
}
