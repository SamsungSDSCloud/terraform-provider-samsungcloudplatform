package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/kubernetesengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/kubernetes-engine2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

func init() {
	scp.RegisterResource("scp_kubernetes_node_pool", ResourceKubernetesNodePool())
}

func ResourceKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: createNodePool,
		ReadContext:   readNodePool,
		UpdateContext: updateNodePool,
		DeleteContext: deleteNodePool,

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
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew:    false,
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
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SSD",
				Description: "Storage type (Currently only SSD is supported)",
			},
			"storage_size_gb": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "100",
				Description: "Storage size in GB (default 100)",
			},
			/*"service_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service level (Use None for now)",
			},*/
			"cpu_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "CPU count for node VMs (default 2)",
			},
			"memory_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "Memory size in GB for node VMs (default 4)",
			},
		},
		Description: "Provides a K8s Node Pool resource.",
	}
}

func getRequiredProducts(ctx context.Context, scpClient *client.SCPClient, engineId string, cpuCount int, memorySizeGB int, serviceLevel string, storageType string) (productGroupId string, contractId string, scaleId string, serviceLevelId string, storageProductId string, err error) {
	engineInfo, _, err := scpClient.KubernetesEngine.ReadEngine(ctx, engineId)
	if err != nil {
		return
	}

	productGroupId, err = client.FindProductGroupId(ctx, scpClient, engineInfo.ZoneId, common.ContainerProductGroup, common.KubernetesEngineVmProductName)
	if err != nil {
		return
	}

	contractId, err = client.FindProductId(ctx, scpClient, productGroupId, common.ContractProductType, "None")
	if err != nil {
		return
	}

	scaleId, err = client.FindScaleProduct(ctx, scpClient, productGroupId, cpuCount, memorySizeGB)
	if err != nil {
		return
	}

	serviceLevelId, err = client.FindProductId(ctx, scpClient, productGroupId, "SERVICE_LEVEL", serviceLevel)
	if err != nil {
		return
	}

	storageProductId, err = client.FindProductId(ctx, scpClient, productGroupId, "DEFAULT_DISK", storageType)
	if err != nil {
		return
	}

	return
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
	cpuCount := data.Get("cpu_count").(int)
	memorySizeGB := data.Get("memory_size_gb").(int)
	storageType := data.Get("storage_type").(string)
	//serviceLevel := data.Get("service_level").(string)
	serviceLevel := "None"

	productGroupId, contractId, scaleId, serviceLevelId, storageProductId, err := getRequiredProducts(ctx, inst.Client, engineId, cpuCount, memorySizeGB, serviceLevel, storageType)
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
			ContractId:           contractId,
			DesiredNodeCount:     int32(data.Get("desired_node_count").(int)),
			ImageId:              data.Get("image_id").(string),
			MaxNodeCount:         int32(data.Get("max_node_count").(int)),
			MinNodeCount:         int32(data.Get("min_node_count").(int)),
			NodePoolName:         data.Get("name").(string),
			ProductGroupId:       productGroupId,
			ScaleId:              scaleId,
			ServiceLevelId:       serviceLevelId,
			StorageId:            storageProductId,
			//StorageSize:          storageSize,
			StorageSize: data.Get("storage_size_gb").(string),
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

	//nodePool, _, err := inst.Client.KubernetesEngine.ReadNodePool(ctx, data.Get("engine_id").(string), data.Id())
	nodePoolList, _, err := inst.Client.KubernetesEngine.GetNodePoolList(ctx, data.Get("engine_id").(string), &kubernetesengine2.NodePoolV2ControllerApiListNodePoolsV2Opts{
		NodePoolName: optional.String{},
		CreatedBy:    optional.String{},
		Page:         optional.NewInt32(0),
		Size:         optional.NewInt32(100),
		Sort:         optional.String{},
	})

	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	var nodePool kubernetesengine2.NodePoolsResponse
	for _, item := range nodePoolList.Contents {
		if item.NodePoolId == data.Id() {
			nodePool = item
		}
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

	if *nodePool.AutoScale {
		data.Set("max_node_count", nodePool.MaxNodeCount)
		data.Set("min_node_count", nodePool.MinNodeCount)
	} else {
		data.Set("desired_node_count", nodePool.DesiredNodeCount)
	}

	data.Set("service_level", serviceLevel.ProductName)
	data.Set("auto_recovery", nodePool.AutoRecovery)
	data.Set("auto_scale", nodePool.AutoScale)
	data.Set("image_id", nodePool.ImageId)
	data.Set("storage_type", storage.ProductName)
	data.Set("storage_size_gb", nodePool.StorageSize)
	data.Set("name", nodePool.NodePoolName)

	return
}

func updateNodePool(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("auto_recovery", "auto_scale", "contract_id", "desired_node_count",
		"image_id", "max_node_count", "min_node_count", "name", "product_id_group", "scale_id",
		"service_level_id", "storage_id", "storage_size") {

		engineId := data.Get("engine_id").(string)
		cpuCount := data.Get("cpu_count").(int)
		memorySizeGB := data.Get("memory_size_gb").(int)
		//serviceLevel := data.Get("service_level").(string)
		serviceLevel := "None"
		storageType := data.Get("storage_type").(string)
		storageSize := "100"

		productGroupId, contractId, scaleId, serviceLevelId, storageProductId, err := getRequiredProducts(ctx, inst.Client, engineId, cpuCount, memorySizeGB, serviceLevel, storageType)

		_, _, err = inst.Client.KubernetesEngine.UpdateNodePool(ctx, engineId, data.Id(), kubernetesengine.NodePoolUpdateRequest{
			AutoRecovery:     data.Get("auto_recovery").(bool),
			AutoScale:        data.Get("auto_scale").(bool),
			ContractId:       contractId,
			DesiredNodeCount: int32(data.Get("desired_node_count").(int)),
			ImageId:          data.Get("image_id").(string),
			MaxNodeCount:     int32(data.Get("max_node_count").(int)),
			MinNodeCount:     int32(data.Get("min_node_count").(int)),
			NodePoolName:     data.Get("name").(string),
			ProductGroupId:   productGroupId,
			ScaleId:          scaleId,
			ServiceLevelId:   serviceLevelId,
			StorageId:        storageProductId,
			StorageSize:      storageSize,
			//StorageSize:      data.Get("storage_size_gb").(string),
		})

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

		nodePool, httpStatus, err := inst.Client.KubernetesEngine.GetNodePoolList(ctx, engineId, &kubernetesengine2.NodePoolV2ControllerApiListNodePoolsV2Opts{
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
