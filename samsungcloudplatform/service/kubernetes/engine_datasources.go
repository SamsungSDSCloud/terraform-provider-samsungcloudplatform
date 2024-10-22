package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_kubernetes_engines", DatasourceEngines())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_kubernetes_engine", DatasourceEngine())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_kubernetes_kubeconfig", DatasourceKubeConfig())
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

func DatasourceEngine() *schema.Resource {
	return &schema.Resource{
		ReadContext: engineDetail, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("KubernetesEngineId"):     {Type: schema.TypeString, Required: true, Description: "Engine Id"},
			common.ToSnakeCase("ProjectId"):              {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			common.ToSnakeCase("ZoneId"):                 {Type: schema.TypeString, Computed: true, Description: "Zone Id"},
			common.ToSnakeCase("Region"):                 {Type: schema.TypeString, Computed: true, Description: "Region"},
			common.ToSnakeCase("KubernetesEngineName"):   {Type: schema.TypeString, Computed: true, Description: "Kubernetes Engine Name"},
			common.ToSnakeCase("ClusterPrefix"):          {Type: schema.TypeString, Computed: true, Description: "Cluster Prefix"},
			common.ToSnakeCase("KubernetesEngineStatus"): {Type: schema.TypeString, Computed: true, Description: "Kubernetes Engine Status"},
			common.ToSnakeCase("K8sVersion"):             {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			common.ToSnakeCase("PrivateEndpointUrl"):     {Type: schema.TypeString, Computed: true, Description: "PrivateEndpoint Url"},
			common.ToSnakeCase("PrivateEndpointAccessControlResourceList"): {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag key",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag value",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "Tag list"},
			common.ToSnakeCase("PublicEndpointUrl"):             {Type: schema.TypeString, Computed: true, Description: "Public Endpoint Url"},
			common.ToSnakeCase("PublicEndpointAccessControlIp"): {Type: schema.TypeString, Computed: true, Description: "Public Endpoint Access Control Ip"},
			common.ToSnakeCase("VpcId"):                         {Type: schema.TypeString, Computed: true, Description: "Vpc Id"},
			common.ToSnakeCase("SubnetId"):                      {Type: schema.TypeString, Computed: true, Description: "Subnet Id"},
			common.ToSnakeCase("SecurityGroupId"):               {Type: schema.TypeString, Computed: true, Description: "Security Group Id"},
			common.ToSnakeCase("LoadBalancerId"):                {Type: schema.TypeString, Computed: true, Description: "Load Balancer Id"},
			common.ToSnakeCase("VolumeId"):                      {Type: schema.TypeString, Computed: true, Description: "Volume Id"},
			common.ToSnakeCase("CifsVolumeId"):                  {Type: schema.TypeString, Computed: true, Description: "Cifs Volume Id"},
			common.ToSnakeCase("CreatedBy"):                     {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                     {Type: schema.TypeString, Computed: true, Description: "Created Dt"},
			common.ToSnakeCase("ModifiedBy"):                    {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                    {Type: schema.TypeString, Computed: true, Description: "Modified Dt"},
		},
		Description: "Provides Kubernetes Engine Detail",
	}
}

func engineDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	engineId := rd.Get("kubernetes_engine_id").(string)

	response, _, err := inst.Client.KubernetesEngine.ReadEngine(ctx, engineId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())

	rd.Set(common.ToSnakeCase("ProjectId"), response.ProjectId)
	rd.Set(common.ToSnakeCase("ZoneId"), response.ZoneId)
	rd.Set(common.ToSnakeCase("Region"), response.Region)
	rd.Set(common.ToSnakeCase("KubernetesEngineName"), response.KubernetesEngineName)
	rd.Set(common.ToSnakeCase("ClusterPrefix"), response.ClusterPrefix)
	rd.Set(common.ToSnakeCase("KubernetesEngineStatus"), response.KubernetesEngineStatus)
	rd.Set(common.ToSnakeCase("K8sVersion"), response.K8sVersion)
	rd.Set(common.ToSnakeCase("PrivateEndpointUrl"), response.PrivateEndpointUrl)
	var privateEndpointAccessControlResourceList common.HclSetObject
	for _, privateEndpointAccessControlResource := range response.PrivateEndpointAccessControlResourceList {
		kv := common.HclKeyValueObject{
			"resource_id":   privateEndpointAccessControlResource.ResourceID,
			"resource_name": privateEndpointAccessControlResource.ResourceName,
			"resource_type": privateEndpointAccessControlResource.ResourceType,
		}
		privateEndpointAccessControlResourceList = append(privateEndpointAccessControlResourceList, kv)
	}
	rd.Set(common.ToSnakeCase("PrivateEndpointAccessControlResourceList"), privateEndpointAccessControlResourceList)
	rd.Set(common.ToSnakeCase("PublicEndpointUrl"), response.PublicEndpointUrl)
	rd.Set(common.ToSnakeCase("PublicEndpointAccessControlIp"), response.PublicEndpointAccessControlIp)
	rd.Set(common.ToSnakeCase("VpcId"), response.VpcId)
	rd.Set(common.ToSnakeCase("SubnetId"), response.SubnetId)
	rd.Set(common.ToSnakeCase("SecurityGroupId"), response.SecurityGroupId)
	rd.Set(common.ToSnakeCase("LoadBalancerId"), response.LoadBalancerId)
	rd.Set(common.ToSnakeCase("VolumeId"), response.VolumeId)
	rd.Set(common.ToSnakeCase("CifsVolumeId"), response.CifsVolumeId)
	rd.Set(common.ToSnakeCase("CreatedBy"), response.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), response.CreatedDt.String())
	rd.Set(common.ToSnakeCase("ModifiedBy"), response.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), response.ModifiedDt.String())

	return nil
}

func DatasourceKubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: engineKubeConfig, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("KubernetesEngineId"): {Type: schema.TypeString, Required: true, Description: "Engine Id"},
			common.ToSnakeCase("kubeconfigType"):     {Type: schema.TypeString, Required: true, Description: "kubeconfig Type"},
			common.ToSnakeCase("KubeConfig"):         {Type: schema.TypeString, Computed: true, Description: "KubeConfig"},
		},
		Description: "Provides Kubernetes Engine Detail",
	}
}

func engineKubeConfig(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	engineId := rd.Get("kubernetes_engine_id").(string)
	kubeconfigType := rd.Get("kubeconfig_type").(string)

	response, _, err := inst.Client.KubernetesEngine.GetKubeConfig(ctx, engineId, kubeconfigType)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())

	rd.Set(common.ToSnakeCase("KubeConfig"), response)

	return nil
}
