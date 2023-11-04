package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/iam"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_iam_access_keys", DatasourceAccessKeys())
}

func DatasourceAccessKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAccessKeysRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":                  common.DatasourceFilter(),
			"project_id":              {Type: schema.TypeString, Optional: true, Description: "Project ID"},
			"access_key_project_type": {Type: schema.TypeString, Optional: true, Description: "Access key's project type"},
			"access_key_state":        {Type: schema.TypeString, Optional: true, Description: "Access key state (ACTIVATED or DEACTIVATED)"},
			"active_yn":               {Type: schema.TypeBool, Optional: true, Description: "Whether the key is activated or not"},
			"project_name":            {Type: schema.TypeString, Optional: true, Description: "Access key's project name"},
			"contents":                {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceAccessKeysElem()},
			"total_count":             {Type: schema.TypeInt, Computed: true, Description: "Total count"},
		},
	}
}

func datasourceAccessKeysRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	careActive := false
	var activeYn bool
	projectId := rd.Get("project_id").(string)
	projectType := rd.Get("access_key_project_type").(string)
	state := rd.Get("access_key_state").(string)
	projectName := rd.Get("project_name").(string)

	ayn, careActive := rd.GetOk("active_yn")
	if careActive {
		activeYn = ayn.(bool)
	}

	var err error
	var response iam.PageResponseV2OfAccessKeysResponse
	if careActive == true {
		response, err = inst.Client.Iam.ListAccessKeys(ctx, projectId, projectType, state, optional.NewBool(activeYn), projectName)
	} else {
		response, err = inst.Client.Iam.ListAccessKeys(ctx, projectId, projectType, state, optional.Bool{}, projectName)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(response.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceAccessKeys().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", response.TotalCount)

	return nil
}

func datasourceAccessKeysElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":           {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"access_key":           {Type: schema.TypeString, Computed: true, Description: "Access key"},
			"access_key_activated": {Type: schema.TypeBool, Computed: true, Description: "Access key activated"},
			"access_key_id":        {Type: schema.TypeString, Computed: true, Description: "Access key ID"},
			"access_key_state":     {Type: schema.TypeString, Computed: true, Description: "Access key state"},
			"expired_dt":           {Type: schema.TypeString, Computed: true, Description: "Expiration date"},
			"project_name":         {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"created_by":           {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":      {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":     {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":           {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":          {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":     {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email":    {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":          {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}
