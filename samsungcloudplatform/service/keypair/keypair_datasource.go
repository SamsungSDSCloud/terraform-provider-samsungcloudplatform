package keypair

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/keypair"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_key_pairs", DatasourceKeyPairs())
}

func DatasourceKeyPairs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {Type: schema.TypeString, Optional: true, Description: "Key Pair Name"},
			"key_pair_id":   {Type: schema.TypeString, Optional: true, Description: "Key Pair Id"},
			"created_by":    {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":          {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":          {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":          {Type: schema.TypeString, Optional: true, Description: "Sorting"},
			"contents":      {Type: schema.TypeList, Optional: true, Description: "Key Pair list", Elem: datasourceElem()},
			"total_count":   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of key pairs",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if len(rd.Get("key_pair_id").(string)) > 0 {

		response, _, err := inst.Client.KeyPair.GetKeyPair(ctx, rd.Get("key_pair_id").(string))
		if err != nil {
			diag.FromErr(err)
		}

		contents := make([]map[string]interface{}, 0)

		content := common.ToMap(response)
		content["virtual_server_id_list"] = response.VirtualServerIdList
		content["launch_configuration_id_list"] = response.LaunchConfigurationIdList
		contents = append(contents, content)

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", 1)
	} else {

		requestParam := keypair.ListKeyPairsRequestParam{
			KeyPairName: rd.Get("key_pair_name").(string),
			CreatedBy:   rd.Get("created_by").(string),
			Page:        (int32)(rd.Get("page").(int)),
			Size:        (int32)(rd.Get("size").(int)),
			Sort:        rd.Get("sort").(string),
		}

		responses, err := inst.Client.KeyPair.ListKeyPairs(ctx, requestParam)
		if err != nil {
			return diag.FromErr(err)
		}

		contents := common.ConvertStructToMaps(responses.Contents)

		for i, resContent := range responses.Contents {
			contents[i]["virtual_server_id_list"] = resContent.VirtualServerIdList
			contents[i]["launch_configuration_id_list"] = resContent.LaunchConfigurationIdList
		}

		var totalCount int32 = 0
		if responses.TotalCount > (int32)(rd.Get("size").(int)) {
			totalCount = (int32)(rd.Get("size").(int))
		} else {
			totalCount = responses.TotalCount
		}

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", totalCount)
	}
	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":     {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"key_pair_id":    {Type: schema.TypeString, Computed: true, Description: "Key Pair Id"},
			"key_pair_name":  {Type: schema.TypeString, Computed: true, Description: "Key Pair Name"},
			"key_pair_state": {Type: schema.TypeString, Computed: true, Description: "Key Pair State"},
			"virtual_server_id_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Virtual Server Id List",
			},
			"launch_configuration_id_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Launch Configuration Id List",
			},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "Person who created the resource"},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "Person who modified the resource"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
