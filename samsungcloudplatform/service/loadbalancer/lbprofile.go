package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_lb_profile", ResourceLbProfile())
}

func ResourceLbProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLbProfileCreate,
		ReadContext:   resourceLbProfileRead,
		UpdateContext: resourceLbProfileUpdate,
		DeleteContext: resourceLbProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target Load Balancer id",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateName3to20DashInMiddle,
				Description:      "Name of profile. (3 to 20 with dash in middle)",
			},
			"category": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateLbProfileCategory,
				Description:      "Category of profile. (PERSISTENCE or APPLICATION)",
			},
			"persistence_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateLbProfilePersistenceType,
				Description:      "Persistence type. (SOURCE_IP, COOKIE) (Only persistence category)",
			},
			"layer_type": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:         true,
				ValidateDiagFunc: ValidateLbProfileLayerType,
				Description:      "Protocol layer type (Only application category). (L4, L7)",
			},
			"redirect_type": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:         true,
				Description: "HTTP redirection option.",
			},
			"request_header_size": {
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
				Description: "Request header size (Only application category with L7 layer. Recommend: 1024). (1 to 65536)",
			},
			"response_header_size": {
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
				Description: "Response header size (Only application category with L7 layer. Recommend: 4096). (1 to 65536)",
			},
			"response_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
				Description: "Request header size (Only application category with L7 layer. Recommend: 60). (1 to 2147483647)",
			},
			"session_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
				Description: "Session timeout value (Only application category. Recommend: 300). (30 to 5400)",
			},
			"x_forwarded_for": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:         true,
				ValidateDiagFunc: ValidateLbProfileForwardedFor,
				Description:      "Forwarded for value (Only application category with L7 layer). (None, INSERT, REPLACE)",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Load Balancer Profile resource.",
	}
}

func resourceLbProfileCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	lbId := rd.Get("lb_id").(string)
	name := rd.Get("name").(string)
	persistenceType := rd.Get("persistence_type").(string)
	category := rd.Get("category").(string)
	layerType := rd.Get("layer_type").(string)
	redirectType := rd.Get("redirect_type").(string)
	requestHeaderSize := rd.Get("request_header_size").(int)
	responseHeaderSize := rd.Get("response_header_size").(int)
	responseTimeout := rd.Get("response_timeout").(int)
	sessionTimeout := rd.Get("session_timeout").(int)
	xForwardedFor := rd.Get("x_forwarded_for").(string)

	// Extra validation
	if category == "PERSISTENCE" {
		if len(layerType) != 0 ||
			requestHeaderSize != 0 ||
			responseHeaderSize != 0 ||
			responseTimeout != 0 ||
			sessionTimeout != 0 ||
			len(xForwardedFor) != 0 {
			return diag.Errorf("In 'PERSISTENCE' category, following fields must not be set : layer_type, request_header_size, response_header_size, response_timeout, session_timeout, x_forwarded_for")
		}
		if len(persistenceType) == 0 {
			return diag.Errorf("In 'PERSISTENCE' category, 'persistence_type' must be set.")
		}
	} else if category == "APPLICATION" {
		if layerType == "L4" {
			if requestHeaderSize != 0 ||
				responseHeaderSize != 0 ||
				responseTimeout != 0 ||
				len(xForwardedFor) != 0 {
				return diag.Errorf("In 'APPLICATION' category 'L4' layer, following fields must not be set : request_header_size, response_header_size, response_timeout, x_forwarded_for")
			}
		}
	}

	inst := meta.(*client.Instance)

	isNameInvalid, err := inst.Client.LoadBalancer.CheckLbProfileName(ctx, lbId, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input profile name is invalid (maybe duplicated) : " + name)
	}

	protocol := "TCP"
	if category == "APPLICATION" {
		if sessionTimeout < 30 || sessionTimeout > 5400 {
			return diag.Errorf("Session duration time is out of bounds. (30 ~ 5,400 sec). ")
		}

		if layerType == "L7" {
			protocol = "HTTP"
			persistenceType = "HTTP"
			if requestHeaderSize < 1 || requestHeaderSize >= 65536 {
				return diag.Errorf("Request header's size should be 1 ~ 65,536 bytes or less. ")
			}
			if responseHeaderSize < 1 || responseHeaderSize >= 65536 {
				return diag.Errorf("Response header's size should be 1 ~ 65,536 bytes or less. ")
			}
			if responseTimeout < 1 || responseTimeout > 2147483647 {
				return diag.Errorf("Response latency out of range. (1 ~ 2,147,483,647 sec)")
			}
		} else if layerType == "L4" {
			persistenceType = "FAST_TCP"
		}
	}

	tags := rd.Get("tags").(map[string]interface{})
	result, err := inst.Client.LoadBalancer.CreateLbProfile(ctx, lbId, layerType, category, name, persistenceType, protocol, redirectType, requestHeaderSize, responseHeaderSize, responseTimeout, sessionTimeout, xForwardedFor, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForLbProfileStatus(ctx, inst.Client, result.ResourceId, lbId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceId)

	return resourceLbProfileRead(ctx, rd, meta)
}

func resourceLbProfileRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.LoadBalancer.GetLbProfile(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", info.LbProfileName)
	rd.Set("category", info.LbProfileCategory)
	if info.LbProfileCategory == "PERSISTENCE" {
		rd.Set("persistence_type", info.LbProfileType)
	} else if info.LbProfileCategory == "APPLICATION" {
		rd.Set("layer_type", info.LayerType)
		rd.Set("redirect_type", info.LbProfileAttrs.RedirectType)
		rd.Set("request_header_size", info.LbProfileAttrs.RequestHeaderSize)
		rd.Set("response_header_size", info.LbProfileAttrs.ResponseHeaderSize)
		rd.Set("response_timeout", info.LbProfileAttrs.ResponseTimeout)
		rd.Set("session_timeout", info.LbProfileAttrs.SessionTimeout)
		rd.Set("x_forwarded_for", info.LbProfileAttrs.XforwardedFor)
	} else {
		rd.SetId("")
		return diag.Errorf("Invalid profile category found")
	}

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceLbProfileUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	hasChange := rd.HasChanges("request_header_size")
	hasChange = hasChange || rd.HasChanges("response_header_size")
	hasChange = hasChange || rd.HasChanges("response_timeout")
	hasChange = hasChange || rd.HasChanges("session_timeout")
	hasChange = hasChange || rd.HasChanges("x_forwarded_for")

	if hasChange {
		requestHeaderSize := rd.Get("request_header_size").(int)
		responseHeaderSize := rd.Get("response_header_size").(int)
		responseTimeout := rd.Get("response_timeout").(int)
		sessionTimeout := rd.Get("session_timeout").(int)

		_, err := inst.Client.LoadBalancer.UpdateLbProfile(ctx, rd.Id(), rd.Get("lb_id").(string), requestHeaderSize, responseHeaderSize, responseTimeout, sessionTimeout, rd.Get("x_forwarded_for").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := waitForLbProfileStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), []string{}, []string{"ACTIVE"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLbProfileRead(ctx, rd, meta)
}

func resourceLbProfileDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.LoadBalancer.DeleteLbProfile(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForLbProfileStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ValidateLbProfilePersistenceType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbProfileIP     string = "SOURCE_IP"
		LbProfileCookie string = "COOKIE"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := v.(string)

	// Check size string
	if (value != LbProfileIP) && (value != LbProfileCookie) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateLbProfileCategory(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbApp         string = "APPLICATION"
		LbPersistence string = "PERSISTENCE"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := v.(string)

	// Check size string
	if (value != LbApp) && (value != LbPersistence) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateLbProfileLayerType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbL4 string = "L4"
		LbL7 string = "L7"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := v.(string)

	// Check size string
	if (value != LbL4) && (value != LbL7) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateLbProfileForwardedFor(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbNone    string = "None"
		LbInsert  string = "INSERT"
		LbReplace string = "REPLACE"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := v.(string)

	// Check size string
	if (value != LbNone) && (value != LbInsert) && (value != LbReplace) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func waitForLbProfileStatus(ctx context.Context, scpClient *client.SCPClient, id string, loadBalancer string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.LoadBalancer.GetLbProfile(ctx, id, loadBalancer)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.LbProfileState, nil
	})
}
