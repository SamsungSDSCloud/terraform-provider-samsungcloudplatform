package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/kubernetesengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	samsungcloudplatform.RegisterResource("Kubernetes", "samsungcloudplatform_kubernetes_node_pool", ResourceKubernetesNodePool())
}

func ResourceKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: createNodePool,
		ReadContext:   readNodePool,
		UpdateContext: updateNodePool,
		DeleteContext: deleteNodePool,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ComputedIf("desired_node_count", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.Get("auto_scale").(bool)
			}),
			/*
				customdiff.ComputedIf("min_node_count", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
					return !diff.Get("auto_scale").(bool)
				}),
				customdiff.ComputedIf("max_node_count", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
					return !diff.Get("auto_scale").(bool)
				}),
			*/
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if diff.Get("auto_scale").(bool) {
					if _, ok := diff.GetOk("desired_node_count"); ok {
						return errors.New("desired_node_count is not supported when auto_scale is enabled")
					}

					/*if _, ok := diff.GetOk("min_node_count"); !ok {
						return errors.New("min_node_count should be given when auto_scale is enabled")
					}

					if _, ok := diff.GetOk("max_node_count"); !ok {
						return errors.New("max_node_count should be given when auto_scale is enabled")
					}*/
				} else {
					if _, ok := diff.GetOk("desired_node_count"); !ok {
						return errors.New("desired_node_count should be given when auto_scale is disabled")
					}

					/*if _, ok := diff.GetOk("min_node_count"); ok {
						return errors.New("min_node_count is not supported when auto_scale is disabled")
					}

					if _, ok := diff.GetOk("max_node_count"); ok {
						return errors.New("max_node_count is not supported when auto_scale is disabled")
					}*/
				}

				return nil
			},
		),

		Importer: &schema.ResourceImporter{
			// Composite import id: "<engineId>/<nodePoolId>". The Read function
			// depends on engine_id being populated; passthrough alone doesn't set
			// it, which is why upstream `terraform import` returned 400 from the
			// API (URL ended up with an empty {kubernetesEngineId} segment).
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), "/", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("import id must be in format <engineId>/<nodePoolId>, got %q", d.Id())
				}
				if err := d.Set("engine_id", parts[0]); err != nil {
					return nil, err
				}
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Node pool name",
				ValidateFunc: validation.StringLenBetween(3, 20),
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of scp_kubernetes_engine resource",
			},
			"availability_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Availability zone name.",
			},
			"auto_recovery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable auto recovery",
			},
			"auto_scale": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable auto scale",
			},
			"encrypt_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Encrypt enabled",
			},
			"desired_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Desired node count in the pool (Desired node count is 0 when auto_scale is enabled)",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image ID (use scp_standard_image data source)",
			},
			"max_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum node count",
			},
			"min_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Minimum node count",
			},
			"scale_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "s1v2m4",
				Description: "Scale name",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "SSD",
				Description: "Storage name (Currently only SSD is supported)",
			},
			"storage_size_gb": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "100",
				Description: "Storage size in GB (default 100)",
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label Key",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label Value",
						},
					},
				},
				Description: "labels",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Taint Effect",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Taint Key",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Taint Value",
						},
					},
				},
				Description: "Taints",
			},
			"advanced_settings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Performing subjects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_unsafe_sysctls": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allowed Unsafe Sysctls",
						},
						"container_log_max_files": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Container Log Max Files",
						},
						"container_log_max_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Container Log Max Size",
						},
						"image_gc_high_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ImageGc High Threshold",
						},
						"image_gc_low_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Image Gc Low Threshold",
						},
						"max_pods": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max Pods",
						},
						"pod_max_pids": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Pod Max Pids",
						},
					},
				},
			},
		},
		Description: "Provides a K8s Node Pool resource.",
	}
}

func createNodePool(ctx context.Context, data *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil

	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	engineId := data.Get("engine_id").(string)
	scaleName := data.Get("scale_name").(string)
	var serverType string

	if strings.HasPrefix(scaleName, "h") {
		serverType = "High Capacity"
	} else if strings.HasPrefix(scaleName, "g") {
		serverType = "GPU"
	} else {
		serverType = "Standard"
	}

	if err != nil {
		return
	}

	response, _, err := inst.Client.KubernetesEngine.CreateNodePool(
		ctx,
		engineId,
		kubernetesengine.CreateNodePoolRequest{
			AvailabilityZoneName: data.Get("availability_zone_name").(string),
			AutoRecovery:         data.Get("auto_recovery").(bool),
			AutoScale:            data.Get("auto_scale").(bool),
			ContractName:         "None",
			DesiredNodeCount:     int32(data.Get("desired_node_count").(int)),
			EncryptEnabled:       data.Get("encrypt_enabled").(bool),
			ImageId:              data.Get("image_id").(string),
			MaxNodeCount:         int32(data.Get("max_node_count").(int)),
			MinNodeCount:         int32(data.Get("min_node_count").(int)),
			NodePoolName:         data.Get("name").(string),
			ServerType:           serverType,
			ScaleName:            data.Get("scale_name").(string),
			ServiceLevelName:     "None",
			StorageName:          data.Get("storage_name").(string),
			StorageSize:          data.Get("storage_size_gb").(string),
			Labels:               toLabelRequestList(data.Get("labels").([]interface{})),
			Taints:               toTaintRequestList(data.Get("taints").([]interface{})),
			AdvancedSettings:     toAdvancedSettingsRequest(data.Get("advanced_settings").(interface{})),
		})

	if err != nil {
		return
	}

	data.SetId(response.ResourceId)

	time.Sleep(5 * time.Second)

	//FAIL, ERROR, NOT READY, RUNNING
	err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))
	if err != nil {
		return
	}

	diagnostics = readNodePool(ctx, data, meta)
	return
}

func readNodePool(ctx context.Context, data *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil

	defer func() {
		if err != nil {
			data.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	nodePool, _, err := inst.Client.KubernetesEngine.ReadNodePool(ctx, data.Get("engine_id").(string), data.Id())

	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	scale, err := client.FindProductById(ctx, inst.Client, nodePool.ProductGroupId, nodePool.ScaleId)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuFound := false
	memoryFound := false

	for _, item := range scale.Item {
		if item.ItemType == "cpu" {
			cpuCount, err := strconv.Atoi(item.ItemValue)

			if err != nil {
				continue
			}

			data.Set("cpu_count", cpuCount)
			cpuFound = true
		} else if item.ItemType == "memory" {
			memorySize, err := strconv.Atoi(item.ItemValue)

			if err != nil {
				continue
			}

			data.Set("memory_size_gb", memorySize)
			memoryFound = true
		}
	}

	if !cpuFound || !memoryFound {
		return diag.FromErr(fmt.Errorf("failed to find scale product"))
	}

	serviceLevel, err := client.FindProductById(ctx, inst.Client, nodePool.ProductGroupId, nodePool.ServiceLevelId)
	if err != nil {
		return diag.FromErr(err)
	}

	storage, err := client.FindProductById(ctx, inst.Client, nodePool.ProductGroupId, nodePool.StorageId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Populate min/max/desired unconditionally so the imported state matches
	// reality regardless of which knobs the user has set in HCL. (Original
	// upstream only set min/max when auto_scale=true, which surfaced as
	// "(known after apply)" drift on imported pools with auto_scale=false.)
	data.Set("max_node_count", nodePool.MaxNodeCount)
	data.Set("min_node_count", nodePool.MinNodeCount)
	data.Set("desired_node_count", nodePool.DesiredNodeCount)

	// The schema declares "scale_name" + "storage_name" (Optional+ForceNew),
	// not "service_level" + "storage_type" — those keys were silent no-ops.
	// Populating the right ones avoids destroy/recreate during a refresh of
	// an imported pool. Also fixes the "taints" key case and adds the
	// previously-missing "encrypt_enabled" attribute.
	data.Set("scale_name", scale.ProductName)
	data.Set("auto_recovery", nodePool.AutoRecovery)
	data.Set("auto_scale", nodePool.AutoScale)
	data.Set("encrypt_enabled", nodePool.EncryptEnabled)
	data.Set("image_id", nodePool.ImageId)
	data.Set("storage_name", storage.ProductName)
	data.Set("storage_size_gb", nodePool.StorageSize)
	data.Set("name", nodePool.NodePoolName)
	data.Set("labels", nodePool.Labels)
	data.Set("taints", nodePool.Taints)
	_ = serviceLevel // intentionally unused; "service_level" not in schema

	return
}

func updateNodePool(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("labels") {
		engineId := data.Get("engine_id").(string)

		if len(data.Get("labels").([]interface{})) == 0 {
			var request []kubernetesengine.LabelRequestToUpdate
			_, _, err := inst.Client.KubernetesEngine.UpdateNodePoolLabels(ctx, engineId, data.Id(), kubernetesengine.UpdateNodePoolLabelsRequest{
				LabelRequestToUpdate: request,
			})

			time.Sleep(5 * time.Second)

			//FAIL, ERROR, NOT READY, RUNNING
			err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdateNodePoolLabels(ctx, engineId, data.Id(), kubernetesengine.UpdateNodePoolLabelsRequest{
				LabelRequestToUpdate: toLabelRequestListToUpdate(data.Get("labels").([]interface{})),
			})

			time.Sleep(5 * time.Second)

			//FAIL, ERROR, NOT READY, RUNNING
			err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("taints") {
		engineId := data.Get("engine_id").(string)

		if len(data.Get("taints").([]interface{})) == 0 {
			var request []kubernetesengine.TaintRequestToUpdate
			_, _, err := inst.Client.KubernetesEngine.UpdateNodePoolTaints(ctx, engineId, data.Id(), kubernetesengine.UpdateNodePoolTatintsRequest{
				TaintRequestToUpdate: request,
			})

			time.Sleep(5 * time.Second)

			//FAIL, ERROR, NOT READY, RUNNING
			err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdateNodePoolTaints(ctx, engineId, data.Id(), kubernetesengine.UpdateNodePoolTatintsRequest{
				TaintRequestToUpdate: toTaintRequestListToUpdate(data.Get("taints").([]interface{})),
			})

			time.Sleep(5 * time.Second)

			//FAIL, ERROR, NOT READY, RUNNING
			err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("auto_recovery", "auto_scale", "desired_node_count", "max_node_count", "min_node_count") {

		engineId := data.Get("engine_id").(string)

		_, _, err := inst.Client.KubernetesEngine.UpdateNodePool(ctx, engineId, data.Id(), kubernetesengine.NodePoolUpdateRequest{
			AutoRecovery:     data.Get("auto_recovery").(bool),
			AutoScale:        data.Get("auto_scale").(bool),
			DesiredNodeCount: int32(data.Get("desired_node_count").(int)),
			MaxNodeCount:     int32(data.Get("max_node_count").(int)),
			MinNodeCount:     int32(data.Get("min_node_count").(int)),
		})

		time.Sleep(5 * time.Second)

		//FAIL, ERROR, NOT READY, RUNNING
		err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChanges("image_id") {
		engineId := data.Get("engine_id").(string)
		nodePool, _, err := inst.Client.KubernetesEngine.ReadNodePool(ctx, data.Get("engine_id").(string), data.Id())
		beforeImageid := nodePool.ImageId
		afterImageid := data.Get("image_id").(string)

		beforeImageType, beforeOsType, beforeK8sVersion, _, beforeProductId, err := getImageInfo(ctx, beforeImageid, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		afterImageType, afterOsType, afterK8sVersion, _, afterProductId, err := getImageInfo(ctx, afterImageid, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		if afterImageType == "" || afterImageType == "" || afterImageType == "" || afterK8sVersion == "" {
			return diag.Errorf("It is not kubernetes image (" + afterImageid + ")")
		}

		beforeK8sMinorVersion, _ := strconv.ParseInt(strings.Split(beforeK8sVersion, ".")[1], 10, 64)
		afterK8sMinorVersion, _ := strconv.ParseInt(strings.Split(afterK8sVersion, ".")[1], 10, 64)

		if beforeImageType != afterImageType {
			return diag.Errorf("Image type not match (berfor : " + beforeImageType + " / after : " + afterImageType + ")")
		}

		if beforeOsType != afterOsType {
			return diag.Errorf("OS type not match (berfor : " + beforeOsType + " / after : " + afterOsType + ")")
		}

		if beforeProductId != afterProductId {
			return diag.Errorf("Product Id not match (berfor : " + beforeProductId + " / after : " + afterProductId + ")")
		}

		// K8s minor version 차이가 1만 나야함
		if afterK8sMinorVersion-beforeK8sMinorVersion != 1 {
			return diag.Errorf("Cannot upgrade from " + beforeK8sVersion + " to " + afterK8sVersion)
		}

		_, _, err = inst.Client.KubernetesEngine.UpgradeNodePool(ctx, engineId, data.Id(), kubernetesengine.NodePoolUgradeRequest{
			UpgradeImageId: data.Get("image_id").(string),
		})

		if err != nil {
			return diag.FromErr(err)
		}
		time.Sleep(5 * time.Second)

		//FAIL, ERROR, NOT READY, RUNNING
		err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"Running"}, refreshNodePool(ctx, meta, engineId, data.Id(), true))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readNodePool(ctx, data, meta)
}

func deleteNodePool(ctx context.Context, data *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil

	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)

	_, err = inst.Client.KubernetesEngine.DeleteNodePool(ctx, engineId, data.Id())
	if err != nil && !common.IsDeleted(err) {
		return
	}

	time.Sleep(5 * time.Second)

	err = client.WaitForStatus(ctx, inst.Client, []string{"Deleting"}, []string{"DELETED"}, refreshNodePool(ctx, meta, engineId, data.Id(), false))
	if err != nil {
		return
	}

	return
}

func refreshNodePool(ctx context.Context, meta interface{}, engineId string, nodePoolId string, errorOnNotFound bool) func() (interface{}, string, error) {
	inst := meta.(*client.Instance)

	return func() (interface{}, string, error) {
		//nodePool, httpStatus, err := inst.Client.KubernetesEngine.ReadNodePool(ctx, engineId, nodePoolId)

		nodePool, httpStatus, err := inst.Client.KubernetesEngine.GetNodePoolList(ctx, engineId, &kubernetesengine2.NodePoolV2ApiListNodePoolsV2Opts{
			NodePoolName: optional.String{},
			CreatedBy:    optional.String{},
			Page:         optional.NewInt32(0),
			Size:         optional.NewInt32(100),
			Sort:         optional.String{},
		})

		for _, node := range nodePool.Contents {
			if node.NodePoolId == nodePoolId {
				if httpStatus == 200 {
					return node, node.NodePoolState, nil
				} else if httpStatus == 404 {
					if errorOnNotFound {
						return nil, "", fmt.Errorf("kubernetes nodepool with id=%s not found", nodePoolId)
					}

					return nodePool, "DELETED", nil
				} else if err != nil {
					return nil, "", err
				}
			}
		}

		if httpStatus == 200 {
			return nodePool, "DELETED", nil
		}
		return nil, "", fmt.Errorf("failed to read kubernetes nodepool(%s) status:%d", nodePoolId, httpStatus)
	}
}

func getImageInfo(ctx context.Context, imageId string, meta interface{}) (string, string, string, string, string, error) {
	var err error = nil
	inst := meta.(*client.Instance)

	standardImage, err := inst.Client.Image.GetStandardImageInfo(ctx, imageId)

	if err != nil {
		return "", "", "", "", "", err
	}

	k8sVersion := standardImage.Properties["k8s.version"]
	osVersion := standardImage.Properties["os.version"]
	productId := standardImage.Products[0].ProductId
	return standardImage.ImageType, standardImage.OsType, k8sVersion, osVersion, productId, err
}

func toLabelRequestList(list []interface{}) []kubernetesengine.LabelRequest {
	if len(list) == 0 {
		return nil
	}
	for _, val := range list {
		if val == nil {
			return nil
		}
	}

	var result []kubernetesengine.LabelRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)

		result = append(result, kubernetesengine.LabelRequest{
			Key:   kv["key"].(string),
			Value: kv["value"].(string),
		})
	}
	return result
}

func toLabelRequestListToUpdate(list []interface{}) []kubernetesengine.LabelRequestToUpdate {
	if len(list) == 0 {
		return nil
	}
	var result []kubernetesengine.LabelRequestToUpdate

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)

		if kv["key"] == "" {
			break
		}

		result = append(result, kubernetesengine.LabelRequestToUpdate{
			Key:   kv["key"].(string),
			Value: kv["value"].(string),
		})
	}
	return result
}

func toTaintRequestList(list []interface{}) []kubernetesengine.TaintRequest {
	if len(list) == 0 {
		return nil
	}
	for _, val := range list {
		if val == nil {
			return nil
		}
	}

	var result []kubernetesengine.TaintRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)

		result = append(result, kubernetesengine.TaintRequest{
			Effect: kv["effect"].(string),
			Key:    kv["key"].(string),
			Value:  kv["value"].(string),
		})
	}
	return result
}

func toTaintRequestListToUpdate(list []interface{}) []kubernetesengine.TaintRequestToUpdate {
	if len(list) == 0 {
		return nil
	}
	var result []kubernetesengine.TaintRequestToUpdate

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)

		if kv["key"] == "" {
			break
		}

		result = append(result, kubernetesengine.TaintRequestToUpdate{
			Effect: kv["effect"].(string),
			Key:    kv["key"].(string),
			Value:  kv["value"].(string),
		})
	}
	return result
}

func toAdvancedSettingsRequest(advancedSettings interface{}) *kubernetesengine.AdvancedSettingsRequest {

	if advancedSettings == nil {
		return nil
	}

	// Terraform schema 에서는 보통 *schema.Set 로 전달됩니다.
	set, ok := advancedSettings.(*schema.Set)
	if !ok {
		// 타입이 맞지 않을 경우 panic(또는 에러 반환) – 상황에 맞게 처리하세요.
		panic("advanced_settings must be a *schema.Set")
	}

	// Set 안에 하나의 맵이 들어 있다고 가정합니다.
	if set.Len() == 0 {
		return nil
	}

	inst := set.List() // []interface{}
	// 첫 번째(그리고 유일한) 요소를 map[string]interface{} 로 변환
	m := inst[0].(map[string]interface{})

	result := &kubernetesengine.AdvancedSettingsRequest{
		AllowedUnsafeSysctls: m["allowed_unsafe_sysctls"].(string),
		ContainerLogMaxFiles: int32(m["container_log_max_files"].(int)),
		ContainerLogMaxSize:  int32(m["container_log_max_size"].(int)),
		ImageGcHighThreshold: int32(m["image_gc_high_threshold"].(int)),
		ImageGcLowThreshold:  int32(m["image_gc_low_threshold"].(int)),
		MaxPods:              int32(m["max_pods"].(int)),
		PodMaxPids:           int32(m["pod_max_pids"].(int)),
	}

	return result
}
