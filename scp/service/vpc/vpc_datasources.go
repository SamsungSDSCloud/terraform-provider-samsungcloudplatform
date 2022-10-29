package vpc

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServiceZoneId"): {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("VpcId"):         {Type: schema.TypeString, Optional: true, Description: "VPC id"},
			common.ToSnakeCase("VpcName"):       {Type: schema.TypeString, Optional: true, Description: "VPC name"},
			common.ToSnakeCase("VpcStates"):     {Type: schema.TypeString, Optional: true, Description: "VPC status"},
			common.ToSnakeCase("CreatedBy"):     {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			common.ToSnakeCase("Page"):          {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):          {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                          {Type: schema.TypeList, Optional: true, Description: "VPC list", Elem: datasourceElem()},
			"total_count":                       {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of vpcs.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := vpc.ListVpcRequest{
		ServiceZoneId: rd.Get(common.ToSnakeCase("ServiceZoneId")).(string),
		VpcId:         rd.Get(common.ToSnakeCase("VpcId")).(string),
		VpcName:       rd.Get(common.ToSnakeCase("VpcName")).(string),
		VpcStates:     rd.Get(common.ToSnakeCase("VpcStates")).(string),
		CreatedBy:     rd.Get(common.ToSnakeCase("CreatedBy")).(string),
		Page:          (int32)(rd.Get(common.ToSnakeCase("Page")).(int)),
		Size:          (int32)(rd.Get(common.ToSnakeCase("Size")).(int)),
	}

	responses, err := inst.Client.Vpc.GetVpcListV2(ctx, requestParam)
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
			common.ToSnakeCase("BlockId"):       {Type: schema.TypeString, Computed: true, Description: "Block id"},
			common.ToSnakeCase("CreatedBy"):     {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):     {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):    {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):    {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("ProjectId"):     {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ServiceZoneId"): {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("VpcId"):         {Type: schema.TypeString, Computed: true, Description: "VPC id"},
			common.ToSnakeCase("VpcName"):       {Type: schema.TypeString, Computed: true, Description: "VPC name"},
			common.ToSnakeCase("VpcState"):      {Type: schema.TypeString, Computed: true, Description: "VPC status"},
		},
	}
}
