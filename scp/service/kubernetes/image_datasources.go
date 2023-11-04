package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/kubernetesapps"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_kubernetes_apps_images", DatasourceKubernetesAppsImages())
}

func DatasourceKubernetesAppsImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceKubernetesAppsImageList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("Category"):         {Type: schema.TypeString, Optional: true, Description: "App image category"},
			common.ToSnakeCase("ImageId"):          {Type: schema.TypeString, Optional: true, Description: "App image id"},
			common.ToSnakeCase("ImageName"):        {Type: schema.TypeString, Optional: true, Description: "App image name"},
			common.ToSnakeCase("IsCarepack"):       {Type: schema.TypeString, Optional: true, Description: "Check whether it is carepack or not "},
			common.ToSnakeCase("IsNew"):            {Type: schema.TypeString, Optional: true, Description: "Check whether it is new image or not"},
			common.ToSnakeCase("IsRecommended"):    {Type: schema.TypeString, Optional: true, Description: "Check whether it is recommendation image or not"},
			common.ToSnakeCase("PricePolicy"):      {Type: schema.TypeString, Optional: true, Description: "Price policy"},
			common.ToSnakeCase("ProductGroupName"): {Type: schema.TypeString, Optional: true, Description: "Product group name"},
			common.ToSnakeCase("Page"):             {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):             {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                             {Type: schema.TypeList, Optional: true, Description: "K8s app image list", Elem: datasourceStandardImageElem()},
			"total_count":                          {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of K8s app images",
	}
}

func datasourceKubernetesAppsImageList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := kubernetesapps.ListStandardImageRequest{
		Category:         rd.Get(common.ToSnakeCase("Category")).(string),
		ImageId:          rd.Get(common.ToSnakeCase("ImageId")).(string),
		ImageName:        rd.Get(common.ToSnakeCase("ImageName")).(string),
		IsCarepack:       rd.Get(common.ToSnakeCase("IsCarepack")).(string),
		IsNew:            rd.Get(common.ToSnakeCase("IsNew")).(string),
		IsRecommended:    rd.Get(common.ToSnakeCase("IsRecommended")).(string),
		PricePolicy:      rd.Get(common.ToSnakeCase("PricePolicy")).(string),
		ProductGroupName: rd.Get(common.ToSnakeCase("ProductGroupName")).(string),
		Page:             (int32)(rd.Get(common.ToSnakeCase("Page")).(int)),
		Size:             (int32)(rd.Get(common.ToSnakeCase("Size")).(int)),
	}

	responses, err := inst.Client.KubernetesApps.GetImageList(ctx, requestParam)

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceStandardImageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("BaseImage"):        {Type: schema.TypeString, Computed: true, Description: "Base image of Apps"},
			common.ToSnakeCase("Category"):         {Type: schema.TypeString, Computed: true, Description: "App image category"},
			common.ToSnakeCase("IconFile"):         {Type: schema.TypeString, Computed: true, Description: "Icon file encoded in base64"},
			common.ToSnakeCase("IconFileName"):     {Type: schema.TypeString, Computed: true, Description: "Icon file name"},
			common.ToSnakeCase("ImageAttr"):        {Type: schema.TypeMap, Computed: true, Description: "App image attributes"},
			common.ToSnakeCase("ImageId"):          {Type: schema.TypeString, Computed: true, Description: "App image id"},
			common.ToSnakeCase("ImageName"):        {Type: schema.TypeString, Computed: true, Description: "App image name"},
			common.ToSnakeCase("PoolId"):           {Type: schema.TypeString, Computed: true, Description: "Block id of this region"},
			common.ToSnakeCase("ProductGroupAttr"): {Type: schema.TypeMap, Computed: true, Description: "Product group attributes"},
			common.ToSnakeCase("ProductGroupId"):   {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			common.ToSnakeCase("ProductGroupName"): {Type: schema.TypeString, Computed: true, Description: "Product group name"},
			common.ToSnakeCase("ZoneId"):           {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("Description"):      {Type: schema.TypeString, Computed: true, Description: "Description"},
			common.ToSnakeCase("ProjectId"):        {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("CreatedBy"):        {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):        {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			common.ToSnakeCase("ModifiedBy"):       {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):       {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
