package gslb

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/gslb"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/gslb2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_gslb", ResourceGslb())
}

func ResourceGslb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGslbCreate,
		ReadContext:   resourceGslbRead,
		UpdateContext: resourceGslbUpdate,
		DeleteContext: resourceGslbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"gslb_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "GSLB Name",
				ValidateDiagFunc: validateName4to40LowerAlphaAndNumeric,
			},
			"gslb_env_usage": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "GSLB Environment Usage",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PUBLIC"}, false)),
			},
			"gslb_algorithm": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "GSLB Algorithm. One of RATIO, RTT",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RATIO", "RTT"}, false)),
			},
			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "GSLB Health Check Protocol. One of ICMP, TCP, HTTP, HTTPS, NONE",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ICMP", "TCP", "HTTP", "HTTPS", "NONE"}, false)),
			},
			"gslb_health_check_interval": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "GSLB Health Check Interval. (5 to 300)",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 300)),
			},
			"gslb_health_check_timeout": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 300)),
				Description:      "GSLB Health Check Timeout. (5 to 300, greater than gslb_health_check_interval)",
			},
			"probe_timeout": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "GSLB Health Check Probe Timeout. (5 to 300),  It must be greater than the Heath Check Interval.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 300)),
			},
			"service_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "GSLB Health Check Service Port. (5 to 300),  It must be greater than the Heath Check Interval.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
			},
			"gslb_health_check_user_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "GSLB Health Check User Id",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 60)),
			},
			"gslb_health_check_user_password": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "GSLB Health Check User Password",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 250)),
			},
			"gslb_send_string": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "GSLB Health Check Send String",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 300)),
			},
			"gslb_response_string": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "GSLB Health Check Response String",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 300)),
			},
			"gslb_resources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "GSLB Resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gslb_destination": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Gslb Resource Destination",
							ValidateDiagFunc: validateIpv4,
						},
						"gslb_region": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Gslb Resource Region",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"KR-EAST-1", "KR-WEST-1", "KR-WEST-2"}, false)),
						},
						"gslb_resource_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Gslb Resource Weight",
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"gslb_resource_description": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Gslb Resource Description",
							ValidateFunc: validation.StringLenBetween(0, 100),
						},
					},
				},
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Gslb resource.",
	}
}

func validateName4to40LowerAlphaAndNumeric(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 4, 40)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z\\d]*$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q: Enter 4-40 char (lowercase, number).", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validateIpv4(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check characters
	if !regexp.MustCompile("^(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]\\d|[1-9])(\\.(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]\\d|\\d)){3}$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q: Supports only IPv4 format.", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validateGslbTimeResourceCount(rd *schema.ResourceData) error {
	gslbHealthCheckInterval := rd.Get("gslb_health_check_interval").(int)
	gslbHealthCheckTimeout := rd.Get("gslb_health_check_timeout").(int)

	if gslbHealthCheckInterval >= gslbHealthCheckTimeout {
		return fmt.Errorf("gslb_health_check_interval must be greater than gslb_health_check_timeout")
	} else {
		return nil
	}

	gslbResourceList := rd.Get("gslb_resources").(*schema.Set).List()
	if len(gslbResourceList) < 1 {
		return fmt.Errorf("There must be at least one gslb_resource in GSLB Serivce.")
	} else {
		return nil
	}
}

func checkStringLength(str string, min int, max int) error {
	if len(str) < min {
		return fmt.Errorf("input must be longer than %v characters", min)
	} else if len(str) > max {
		return fmt.Errorf("input must be shorter than %v characters", max)
	} else {
		return nil
	}
}

func convertGslbResources(list common.HclListObject) ([]gslb.GslbResourceRequest, error) {
	var result []gslb.GslbResourceRequest
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := gslb.GslbResourceRequest{}
		if gslbDestination, ok := itemObject["gslb_destination"]; ok {
			info.GslbDestination = gslbDestination.(string)
		}
		if gslbRegion, ok := itemObject["gslb_region"]; ok {
			info.GslbRegion = gslbRegion.(string)
		}
		if gslbResourceWeight, ok := itemObject["gslb_resource_weight"]; ok {
			info.GslbResourceWeight = int32(gslbResourceWeight.(int))
		}
		if gslbResourceDescription, ok := itemObject["gslb_resource_description"]; ok {
			info.GslbResourceDescription = gslbResourceDescription.(string)
		}

		result = append(result, info)
	}

	return result, nil
}

func resourceGslbCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	gslbName := rd.Get("gslb_name").(string)
	gslbEnvUsage := rd.Get("gslb_env_usage").(string)
	gslbAlgorithm := rd.Get("gslb_algorithm").(string)
	gslbHealthCheck := gslb.GslbHealthCheckRequest{
		Protocol:                    rd.Get("protocol").(string),
		GslbHealthCheckInterval:     int32(rd.Get("gslb_health_check_interval").(int)),
		GslbHealthCheckTimeout:      int32(rd.Get("gslb_health_check_timeout").(int)),
		ProbeTimeout:                int32(rd.Get("probe_timeout").(int)),
		ServicePort:                 int32(rd.Get("service_port").(int)),
		GslbHealthCheckUserId:       rd.Get("gslb_health_check_user_id").(string),
		GslbHealthCheckUserPassword: rd.Get("gslb_health_check_user_password").(string),
		GslbSendString:              rd.Get("gslb_send_string").(string),
		GslbResponseString:          rd.Get("gslb_response_string").(string),
	}

	gslbResources, err := convertGslbResources(rd.Get("gslb_resources").(*schema.Set).List())

	createRequest := gslb.CreateGslbRequest{
		GslbName:        gslbName,
		GslbEnvUsage:    gslbEnvUsage,
		GslbAlgorithm:   gslbAlgorithm,
		GslbHealthCheck: gslbHealthCheck,
		GslbResources:   gslbResources,
		Tags:            rd.Get("tags").(map[string]interface{}),
	}

	validateErr := validateGslbTimeResourceCount(rd)
	if validateErr != nil {
		return diag.FromErr(validateErr)
	}

	result, _, err := inst.Client.Gslb.CreateGslb(ctx, createRequest)

	if err != nil {
		return
	}

	err = waitForGslbStatus(ctx, inst.Client, result.ResourceId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return
	}

	rd.SetId(result.ResourceId)

	return resourceGslbRead(ctx, rd, meta)
}

func resourceGslbRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	gslbInfo, _, err := inst.Client.Gslb.GetGslb(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("gslb_id", gslbInfo.GslbId)
	rd.Set("gslb_name", strings.Split(gslbInfo.GslbName, ".")[0])
	rd.Set("gslb_env_usage", gslbInfo.GslbEnvUsage)
	rd.Set("gslb_algorithm", gslbInfo.GslbAlgorithm)
	rd.Set("protocol", gslbInfo.GslbHealthCheck.Protocol)
	rd.Set("gslb_health_check_interval", gslbInfo.GslbHealthCheck.GslbHealthCheckInterval)
	rd.Set("gslb_health_check_timeout", gslbInfo.GslbHealthCheck.GslbHealthCheckTimeout)
	rd.Set("probe_timeout", gslbInfo.GslbHealthCheck.ProbeTimeout)
	rd.Set("service_port", gslbInfo.GslbHealthCheck.ServicePort)
	rd.Set("gslb_health_check_user_id", gslbInfo.GslbHealthCheck.GslbHealthCheckUserId)
	rd.Set("gslb_send_string", gslbInfo.GslbHealthCheck.GslbSendString)
	rd.Set("gslb_response_string", gslbInfo.GslbHealthCheck.GslbResponseString)

	var gslbResourceInfo gslb2.ListResponseGslbResourceMappingResponse
	gslbResourceInfo, err = inst.Client.Gslb.GetGslbResource(ctx, rd.Id())

	var gslbResources []common.HclKeyValueObject

	for _, gslbResource := range gslbResourceInfo.Contents {
		gslbResources = append(gslbResources, common.HclKeyValueObject{
			"gslb_destination":          gslbResource.GslbDestination,
			"gslb_region":               gslbResource.GslbRegion,
			"gslb_resource_weight":      gslbResource.GslbResourceWeight,
			"gslb_resource_description": gslbResource.GslbResourceDescription,
		})
	}

	rd.Set("gslb_resources", gslbResources)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceGslbUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges("gslb_algorithm") {
		_, _, err := inst.Client.Gslb.UpdateGslbAlgorithm(ctx, rd.Id(), rd.Get("gslb_algorithm").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForGslbStatus(ctx, inst.Client, rd.Id(), []string{"EDITING"}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if rd.HasChanges("protocol", "gslb_health_check_interval", "gslb_health_check_timeout", "probe_timeout", "service_port", "gslb_health_check_user_id", "gslb_send_string", "gslb_response_string") {
		protocol := rd.Get("protocol").(string)
		gslbHealthCheckInterval := int32(rd.Get("gslb_health_check_interval").(int))
		gslbHealthCheckTimeout := int32(rd.Get("gslb_health_check_timeout").(int))
		servicePort := int32(rd.Get("service_port").(int))
		probeTimeout := int32(rd.Get("probe_timeout").(int))
		gslbHealthCheckUserId := rd.Get("gslb_health_check_user_id").(string)
		gslbHealthCheckUserPassword := rd.Get("gslb_health_check_user_password").(string)
		gslbSendString := rd.Get("gslb_send_string").(string)
		gslbResponseString := rd.Get("gslb_response_string").(string)
		updateHealthCheckRequest := gslb2.ChangeGslbHealthCheckRequest{
			Protocol:                    protocol,
			GslbHealthCheckInterval:     gslbHealthCheckInterval,
			GslbHealthCheckTimeout:      gslbHealthCheckTimeout,
			ServicePort:                 servicePort,
			ProbeTimeout:                probeTimeout,
			GslbHealthCheckUserId:       gslbHealthCheckUserId,
			GslbHealthCheckUserPassword: gslbHealthCheckUserPassword,
			GslbSendString:              gslbSendString,
			GslbResponseString:          gslbResponseString,
		}

		validateErr := validateGslbTimeResourceCount(rd)
		if validateErr != nil {
			return diag.FromErr(validateErr)
		}

		_, _, err := inst.Client.Gslb.UpdateGslbHealthCheck(ctx, rd.Id(), updateHealthCheckRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForGslbStatus(ctx, inst.Client, rd.Id(), []string{"EDITING"}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if rd.HasChanges("gslb_resources") {
		gslbResourceList := rd.Get("gslb_resources").(*schema.Set).List()

		gslbResources := make([]gslb2.GslbResourceMappingRequest, len(gslbResourceList))

		for i, gslbResource := range gslbResourceList {
			r := gslbResource.(map[string]interface{})
			if t, ok := r["gslb_destination"]; ok {
				gslbResources[i].GslbDestination = t.(string)
			}
			if t, ok := r["gslb_region"]; ok {
				gslbResources[i].GslbRegion = t.(string)
			}
			if v, ok := r["gslb_resource_weight"]; ok {
				gslbResources[i].GslbResourceWeight = int32(v.(int))
			}
			if t, ok := r["gslb_resource_description"]; ok {
				gslbResources[i].GslbResourceDescription = t.(string)
			}
		}

		validateErr := validateGslbTimeResourceCount(rd)
		if validateErr != nil {
			return diag.FromErr(validateErr)
		}

		_, _, err := inst.Client.Gslb.UpdateGslbResources(ctx, rd.Id(), gslbResources)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForGslbStatus(ctx, inst.Client, rd.Id(), []string{"EDITING"}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGslbRead(ctx, rd, meta)
}

func resourceGslbDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Gslb.DeleteGslb(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	time.Sleep(10 * time.Second)

	err = waitForGslbStatus(ctx, inst.Client, rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForGslbStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Gslb.GetGslb(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.GslbState, nil
	})
}
