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
	scp.RegisterDataSource("scp_kubernetes_engine_versions", DatasourceEngineVersions())
}

func DatasourceEngineVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVersionList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"contents": {Type: schema.TypeList, Optional: true, Description: "K8s engine list", Elem: datasourceVersionElem()},
			"page":     {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":     {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
		},
		Description: "Provides list of K8s versions",
	}
}

func dataSourceVersionList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.KubernetesEngine.GetEngineVersionList(ctx, &kubernetesengine2.K8sTemplateV2ApiListKubernetesVersionV21Opts{
		Page: optional.NewInt32((int32)(rd.Get("page").(int))),
		Size: optional.NewInt32((int32)(rd.Get("size").(int))),
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

func datasourceVersionElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":  {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"k8s_version": {Type: schema.TypeString, Computed: true, Description: "K8s version"},
		},
	}
}
