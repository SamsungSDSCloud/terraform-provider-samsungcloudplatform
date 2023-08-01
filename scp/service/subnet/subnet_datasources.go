package subnet

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/subnet2"
	"github.com/antihax/optional"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_subnets", DatasourceSubnets())
	scp.RegisterDataSource("scp_subnet_resources", DatasourceSubnetResources())
	scp.RegisterDataSource("scp_subnet_ip_check", DatasourceSubnetCheckIp())
}

func DatasourceSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_cidr_block": {Type: schema.TypeString, Optional: true, Description: "Subnet CIDR block"},
			"subnet_id":         {Type: schema.TypeString, Optional: true, Description: "Subnet id"},
			"subnet_name":       {Type: schema.TypeString, Optional: true, Description: "Subnet name"},
			"subnet_states":     {Type: schema.TypeString, Optional: true, Description: "Subnet states (ACTIVE, ERROR)"},
			"subnet_types":      {Type: schema.TypeString, Optional: true, Description: "Subnet types (PUBLIC, PRIVATE)"},
			"vpc_id":            {Type: schema.TypeString, Optional: true, Description: "VPC id"},
			"created_by":        {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":          {Type: schema.TypeList, Optional: true, Description: "Subnet list", Elem: datasourceElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true, Description: "Subnet list size"},
		},
		Description: "Provides list of subnets.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := &subnet2.SubnetOpenApiControllerApiListSubnetV2Opts{
		SubnetCidrBlock: optional.NewString(rd.Get("subnet_cidr_block").(string)),
		SubnetId:        optional.NewString(rd.Get("subnet_id").(string)),
		SubnetName:      optional.NewString(rd.Get("subnet_name").(string)),
		SubnetTypes:     optional.NewInterface(rd.Get("subnet_types").(string)),
		SubnetStates:    optional.NewInterface(rd.Get("subnet_states").(string)),
		VpcId:           optional.NewString(rd.Get("vpc_id").(string)),
		CreatedBy:       optional.NewString(rd.Get("created_by").(string)),
		Page:            optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:            optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:            optional.Interface{},
	}

	responses, _, err := inst.Client.Subnet.GetSubnetList(ctx, requestParam)
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
			"gateway_ip_address": {Type: schema.TypeString, Computed: true, Description: "Ip address of gateway"},
			"subnet_cidr_block":  {Type: schema.TypeString, Computed: true, Description: "Subnet CIDR block"},
			"subnet_id":          {Type: schema.TypeString, Computed: true, Description: "Subnet id"},
			"subnet_name":        {Type: schema.TypeString, Computed: true, Description: "Subnet name"},
			"subnet_purpose":     {Type: schema.TypeString, Computed: true, Description: "Purpose of subnet (GENERAL)"},
			"subnet_state":       {Type: schema.TypeString, Computed: true, Description: "Subnet status"},
			"subnet_type":        {Type: schema.TypeString, Computed: true, Description: "Subnet type (PUBLIC, PRIVATE)"},
			"vpc_id":             {Type: schema.TypeString, Computed: true, Description: "VPC id"},
			"created_dt":         {Type: schema.TypeString, Computed: true, Description: "Creation date"},
		},
	}
}

func DatasourceSubnetResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetResourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id":          {Type: schema.TypeString, Required: true, Description: "Subnet id"},
			"ip_address":         {Type: schema.TypeString, Optional: true, Description: "Ip address"},
			"linked_object_type": {Type: schema.TypeString, Optional: true, Description: "Type of object linked by subnet"},
			"page":               {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":               {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":           {Type: schema.TypeList, Optional: true, Description: "Subnet resource list size", Elem: datasourceSubnetResourceElem()},
			"total_count":        {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of subnet resources",
	}
}

func dataSourceSubnetResourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)

	requestParam := subnet2.SubnetVipOpenApiControllerApiListSubnetResourcesV2Opts{
		IpAddress:        optional.NewString(rd.Get("ip_address").(string)),
		LinkedObjectType: optional.NewString(rd.Get("linked_object_type").(string)),
		Page:             optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:             optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:             optional.Interface{},
	}

	responses, _, err := inst.Client.Subnet.GetSubnetResourcesV2List(ctx, subnetId, &requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceSubnetResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_address":         {Type: schema.TypeString, Computed: true, Description: "Ip address"},
			"ip_id":              {Type: schema.TypeString, Computed: true, Description: "ip id"},
			"ip_state":           {Type: schema.TypeString, Computed: true, Description: "Ip status"},
			"linked_object_id":   {Type: schema.TypeString, Computed: true, Description: "Id of object linked by subnet"},
			"linked_object_name": {Type: schema.TypeString, Computed: true, Description: "Name of object linked by subnet"},
			"linked_object_type": {Type: schema.TypeString, Computed: true, Description: "Type of object linked by subnet"},
			"ip_description":     {Type: schema.TypeString, Computed: true, Description: "Description of ip"},
			"created_by":         {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":         {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":        {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":        {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func DatasourceSubnetCheckIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSubnetIpCheck,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip_address": {Type: schema.TypeString, Required: true, Description: "Subnet ip Address"},
			"subnet_id":  {Type: schema.TypeString, Required: true, Description: "Subnet id"},
			"result":     {Type: schema.TypeBool, Optional: true, Description: "Subnet ip Check Usable"},
		},
		Description: "Provides list of subnets.",
	}
}

func datasourceSubnetIpCheck(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, err := inst.Client.Subnet.CheckAvailableSubnetIp(ctx, rd.Get("subnet_id").(string), rd.Get("ip_address").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("result", responses.Result)

	return nil
}
