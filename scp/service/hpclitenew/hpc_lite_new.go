package hpclitenew

import (
	context2 "context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/hpclitenew"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_hpc_lite_new", ResourceHpcLiteNew())
}

func ResourceHpcLiteNew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHpcLiteNewCreate,
		ReadContext:   resourceHpcLiteNewRead,
		UpdateContext: resourceHpcLiteNewUpdate,
		DeleteContext: resourceHpcLiteNewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		// TODO Validation 추가
		Schema: map[string]*schema.Schema{
			"co_service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) CO Pool ID",
			},
			"contract": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Contract",
			},
			"hyper_threading_enabled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) HT Enabled",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Image ID",
			},
			"init_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HPC Lite(New) Init Script",
			},
			"os_user_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "HPC Lite(New) OS User ID",
				ValidateDiagFunc: validateOsUserId,
			},
			"os_user_password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: common.ValidatePassword8to20,
				Description:      "HPC Lite(New) OS User PWD",
			},
			"product_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Product Group ID",
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) block Id",
			},
			"server_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Server Type",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Service Zone ID",
			},
			"vlan_pool_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HPC Lite(New) Vlan Pool CIDR",
			},
			"server_details": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HPC Lite(New) ID",
						},
						"server_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "HPC Lite(New) Server Detail Name",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "HPC Lite(New) Server Detail ip address",
						},
					},
				},
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Hpc Lite(New) resource.",
		CustomizeDiff: func(ctx context2.Context, diff *schema.ResourceDiff, i interface{}) error {
			if diff.Id() == "" {
				//create
			} else {
				//update
				if diff.HasChange("co_service_zone_id") {
					return fmt.Errorf("co_service_zone_id can't be modified.")
				}
				if diff.HasChange("contract") {
					return fmt.Errorf("contract can't be modified.")
				}
				if diff.HasChange("hyper_threading_enabled") {
					return fmt.Errorf("hyper_threading_enabled can't be modified.")
				}
				if diff.HasChange("image_id") {
					return fmt.Errorf("image_id can't be modified.")
				}
				if diff.HasChange("init_script") {
					return fmt.Errorf("init_script can't be modified.")
				}
				if diff.HasChange("os_user_id") {
					return fmt.Errorf("os_user_id can't be modified.")
				}
				if diff.HasChange("os_user_password") {
					return fmt.Errorf("os_user_password can't be modified.")
				}
				if diff.HasChange("product_group_id") {
					return fmt.Errorf("product_group_id can't be modified.")
				}
				if diff.HasChange("resource_pool_id") {
					return fmt.Errorf("resource_pool_id can't be modified.")
				}
				if diff.HasChange("server_type") {
					return fmt.Errorf("server_type can't be modified.")
				}
				if diff.HasChange("service_zone_id") {
					return fmt.Errorf("service_zone_id can't be modified.")
				}
				if diff.HasChange("vlan_pool_cidr") {
					return fmt.Errorf("vlan_pool_cidr can't be modified.")
				}
			}
			return nil
		},
	}
}

func validateOsUserId(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	osUserId := v.(string)
	if osUserId != "root" {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("OS User ID must be root"),
			AttributePath: path,
		})
	}
	return diags
}

func resourceHpcLiteNewCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()
	inst := meta.(*client.Instance)

	var serverDetailsRequestList []hpclitenew.ServerDetailRequest
	for _, server := range rd.Get("server_details").([]interface{}) {
		serverDetail := server.(map[string]interface{})
		serverDetailsRequestList = append(serverDetailsRequestList, hpclitenew.ServerDetailRequest{
			ServerName: serverDetail["server_name"].(string),
			IpAddress:  serverDetail["ip_address"].(string),
		})
	}

	request := hpclitenew.HpcLiteNewCreateRequest{
		CoServiceZoneId:       rd.Get("co_service_zone_id").(string),
		Contract:              rd.Get("contract").(string),
		HyperThreadingEnabled: rd.Get("hyper_threading_enabled").(string),
		ImageId:               rd.Get("image_id").(string),
		InitScript:            rd.Get("init_script").(string),
		OsUserId:              rd.Get("os_user_id").(string),
		OsUserPassword:        rd.Get("os_user_password").(string),
		ProductGroupId:        rd.Get("product_group_id").(string),
		ResourcePoolId:        rd.Get("resource_pool_id").(string),
		ServerDetails:         serverDetailsRequestList,
		ServerType:            rd.Get("server_type").(string),
		ServiceZoneId:         rd.Get("service_zone_id").(string),
		Tags:                  rd.Get("tags").(map[string]interface{}),
		VlanPoolCidr:          rd.Get("vlan_pool_cidr").(string),
	}

	response, _, err := inst.Client.HpcLiteNew.CreateHpcLiteNew(ctx, request)
	if err != nil {
		diag.FromErr(err)
	}

	for _, serverId := range response.ResourceIdList {
		err = waitForAllHpcLiteNewStatus(ctx, inst.Client, serverId, []string{common.CreatingState}, []string{common.RunningState}, true)
		if err != nil {
			diag.FromErr(err)
		}
	}

	setResourceId(rd, response.ResourceIdList)

	return resourceHpcLiteNewRead(ctx, rd, meta)
}

func resourceHpcLiteNewRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	serverNameToServerDetail := make(map[string]interface{})
	for _, server := range rd.Get("server_details").([]interface{}) {
		serverDetail := server.(map[string]interface{})
		serverNameToServerDetail[serverDetail["server_name"].(string)] = serverDetail
	}

	var serverDetails []interface{}
	serverIds := getServerIds(rd)
	for _, serverId := range serverIds {
		res, _, err := inst.Client.HpcLiteNew.GetHpcLiteNewDetail(ctx, serverId)
		if err != nil {
			rd.SetId("")
			return diag.FromErr(err)
		}
		// not import
		if _, hasKey := serverNameToServerDetail[res.ServerName]; hasKey {
			serverDetail := serverNameToServerDetail[res.ServerName].(map[string]interface{})
			serverDetail["id"] = serverId
			if serverDetail["ip_address"] == nil || serverDetail["ip_address"] == "" {
				serverDetail["ip_address"] = res.IpAddress
			}
		} else {
			serverDetail := make(map[string]interface{})
			serverDetail["id"] = serverId
			serverDetail["server_name"] = res.ServerName
			serverDetail["ip_address"] = res.IpAddress
			serverDetails = append(serverDetails, serverDetail)
		}
	}
	if len(serverDetails) > 0 {
		rd.Set("server_details", serverDetails)
	}
	res, _, err := inst.Client.HpcLiteNew.GetHpcLiteNewDetail(ctx, serverIds[0])
	if _, exists := rd.GetOk("co_service_zone_id"); !exists {
		rd.Set("co_service_zone_id", res.CoServiceZone)
	}
	if _, exists := rd.GetOk("contract"); !exists {
		rd.Set("contract", res.Contract)
	}
	if _, exists := rd.GetOk("hyper_threading_enabled"); !exists {
		rd.Set("hyper_threading_enabled", res.HyperThreading)
	}
	if _, exists := rd.GetOk("image_id"); !exists {
		rd.Set("image_id", res.ImageId)
	}
	if _, exists := rd.GetOk("init_script"); !exists {
		rd.Set("init_script", res.InitScript)
	}
	if _, exists := rd.GetOk("os_user_id"); !exists {
		rd.Set("os_user_id", "root")
	}
	//rd.Set("os_user_password", res.us)
	//rd.Set("product_group_id", res.pro)
	//rd.Set("resource_pool_id",
	if _, exists := rd.GetOk("server_type"); !exists {
		rd.Set("server_type", strings.Split(res.ServerType, " ")[0])
	}
	if _, exists := rd.GetOk("service_zone_id"); !exists {
		rd.Set("service_zone_id", res.ZoneId)
	}
	//rd.Set("vlan_pool_cidr", res.)

	tfTags.SetTags(ctx, rd, meta, rd.Id())
	return nil
}

func resourceHpcLiteNewUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()
	if rd.HasChanges("server_details") {
		ov, nv := rd.GetChange("server_details")
		oldServerDetails := ov.([]interface{})
		newServerDetails := nv.([]interface{})

		oldDetails, _ := mapServerNameToDetail(oldServerDetails)
		newDetails, hasDuplicatedName := mapServerNameToDetail(newServerDetails)
		if hasDuplicatedName {
			errMsg := "Server Name is duplicated."
			err := rd.Set("server_details", ov)
			if err != nil {
				errMsg += err.Error()
			}
			return diag.Errorf(errMsg)
		}

		if len(oldDetails) > len(newDetails) {
			var deleteServerIds []string
			var notDeleteServerIds []string
			for serverName := range oldDetails {
				serverId := oldDetails[serverName].(map[string]interface{})["id"].(string)
				if _, hasKey := newDetails[serverName]; hasKey {
					notDeleteServerIds = append(notDeleteServerIds, serverId)
				} else {
					deleteServerIds = append(deleteServerIds, serverId)
				}
			}

			deleteHpcLiteNewServers(ctx, rd, deleteServerIds, inst)
			setResourceId(rd, notDeleteServerIds)
		} else if len(oldDetails) < len(newDetails) {
			var serverDetailsRequestList []hpclitenew.ServerDetailRequest
			for serverName := range newDetails {
				if _, hasKey := oldDetails[serverName]; !hasKey {
					serverDetailsRequestList = append(serverDetailsRequestList, hpclitenew.ServerDetailRequest{
						ServerName: serverName,
						IpAddress:  newDetails[serverName].(map[string]interface{})["ip_address"].(string),
					})
				}
			}
			request := hpclitenew.HpcLiteNewCreateRequest{
				CoServiceZoneId:       rd.Get("co_service_zone_id").(string),
				Contract:              rd.Get("contract").(string),
				HyperThreadingEnabled: rd.Get("hyper_threading_enabled").(string),
				ImageId:               rd.Get("image_id").(string),
				InitScript:            rd.Get("init_script").(string),
				OsUserId:              rd.Get("os_user_id").(string),
				OsUserPassword:        rd.Get("os_user_password").(string),
				ProductGroupId:        rd.Get("product_group_id").(string),
				ResourcePoolId:        rd.Get("resource_pool_id").(string),
				ServerDetails:         serverDetailsRequestList,
				ServerType:            rd.Get("server_type").(string),
				ServiceZoneId:         rd.Get("service_zone_id").(string),
				Tags:                  rd.Get("tags").(map[string]interface{}),
				VlanPoolCidr:          rd.Get("vlan_pool_cidr").(string),
			}
			response, _, err := inst.Client.HpcLiteNew.CreateHpcLiteNew(ctx, request)
			if err != nil {
				diag.FromErr(err)
			}
			currentServerIds := getServerIds(rd)
			for _, serverId := range response.ResourceIdList {
				err = waitForAllHpcLiteNewStatus(ctx, inst.Client, serverId, []string{common.CreatingState}, []string{common.RunningState}, true)
				if err != nil {
					diag.FromErr(err)
				}
				currentServerIds = append(currentServerIds, serverId)
			}
			setResourceId(rd, currentServerIds)
		}
	}
	if rd.HasChanges("tags") {
		serverIds := getServerIds(rd)
		for _, serverId := range serverIds {
			tfTags.UpdateTags(ctx, rd, meta, serverId)
		}
	}
	return resourceHpcLiteNewRead(ctx, rd, meta)
}

func mapServerNameToDetail(serverDetails []interface{}) (map[string]interface{}, bool) {
	nameToDetails := make(map[string]interface{})
	for _, v := range serverDetails {
		serverDetail := v.(map[string]interface{})
		nameToDetails[serverDetail["server_name"].(string)] = serverDetail
	}
	hasDuplicatedName := len(nameToDetails) < len(serverDetails)
	return nameToDetails, hasDuplicatedName
}

func resourceHpcLiteNewDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	deleteServerIds := strings.Split(rd.Id(), ",")
	deleteHpcLiteNewServers(ctx, rd, deleteServerIds, inst)

	return nil
}

func deleteHpcLiteNewServers(ctx context.Context, rd *schema.ResourceData, deleteServerIds []string, inst *client.Instance) {
	request := hpclitenew.HpcLiteNewDeleteRequest{
		ServerIds:     deleteServerIds,
		ServiceZoneId: rd.Get("service_zone_id").(string),
	}

	for _, serverId := range deleteServerIds {
		err := waitForAllHpcLiteNewStatus(ctx, inst.Client, serverId, common.VirtualServerProcessingStates(), []string{common.RunningState}, true)
		if err != nil {
			diag.FromErr(err)
		}
	}

	_, _, err := inst.Client.HpcLiteNew.DeleteHpcLiteNew(ctx, request)
	if err != nil {
		diag.FromErr(err)
	}

	for _, serverId := range deleteServerIds {
		err = waitForAllHpcLiteNewStatus(ctx, inst.Client, serverId, common.VirtualServerProcessingStates(), []string{common.DeletedState}, false)
		if err != nil {
			diag.FromErr(err)
		}
	}
}

func waitForAllHpcLiteNewStatus(ctx context.Context, scpClient *client.SCPClient, serverId string, pendingStates []string, targetStates []string, checkNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		responseVo, c, err := scpClient.HpcLiteNew.GetHpcLiteNewDetail(ctx, serverId)

		if err != nil {
			if c == 404 && !checkNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !checkNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return responseVo, strings.ToUpper(responseVo.ServerState), nil
	})
}

func setResourceId(rd *schema.ResourceData, resourceIdList []string) {
	rd.SetId(strings.Join(resourceIdList, ","))
}

func getServerIds(rd *schema.ResourceData) []string {
	return strings.Split(rd.Id(), ",")
}
