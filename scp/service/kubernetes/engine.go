package kubernetes

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/kubernetesengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func init() {
	scp.RegisterResource("scp_kubernetes_engine", ResourceKubernetesEngine())
}

func ResourceKubernetesEngine() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEngine,
		ReadContext:   readEngine,
		UpdateContext: updateEngine,
		DeleteContext: deleteEngine,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Kubernetes engine name",
				ValidateFunc: validation.StringLenBetween(3, 30),
			},
			"cloud_logging_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     nil,
				Computed:    true,
				Description: "Enable cloud logging",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kubernetes version (Contact administrator to check supported version)",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Load balancer ID",
			},
			"private_acl_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag key",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value",
						},
						"value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag value",
						},
					},
				},
				Description: "Tag list",
			},
			"public_acl_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of comma separated IP addresses (CIDR or Single IP) for access control",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security group ID",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID",
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "File storage volume ID",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID",
			},
			"public_endpoint": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "Public endpoint URL for the kubernetes cluster",
			},
			"kube_config": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "Kube config of the kubernetes cluster",
			},
			"cifs_volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CIFS volume id",
			},
		},
		Description: "Provides a K8s Engine resource.",
	}
}

func createEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	vpcId := data.Get("vpc_id").(string)
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)

	response, _, err := inst.Client.KubernetesEngine.CreateEngine(ctx, kubernetesengine.CreateEngineRequest{
		CloudLoggingEnabled:  data.Get("cloud_logging_enabled").(bool),
		K8sVersion:           data.Get("kubernetes_version").(string),
		KubernetesEngineName: data.Get("name").(string),
		LbId:                 data.Get("load_balancer_id").(string),
		PrivateAclResources:  toPrivateAclResourcesRequestList(data.Get("private_acl_resources").([]interface{})),
		PublicAclIpAddress:   data.Get("public_acl_ip_address").(string),
		SecurityGroupId:      data.Get("security_group_id").(string),
		SubnetId:             data.Get("subnet_id").(string),
		VolumeId:             data.Get("volume_id").(string),
		CifsVolumeId:         data.Get("cifs_volume_id").(string),
		VpcId:                vpcId,
		ZoneId:               vpcInfo.ServiceZoneId,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ResourceId)

	time.Sleep(10 * time.Second)

	err = client.WaitForStatus(ctx, inst.Client, []string{"CREATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
	if err != nil {
		return diag.FromErr(err)
	}

	return readEngine(ctx, data, meta)
}

func readEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	engine, _, err := inst.Client.KubernetesEngine.ReadEngine(ctx, data.Id())
	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	kubeConfig, _, err := inst.Client.KubernetesEngine.GetKubeConfig(ctx, data.Id())

	data.Set("name", engine.KubernetesEngineName)
	data.Set("cloud_logging_enabled", engine.CloudLoggingEnabled)
	data.Set("kubernetes_version", engine.K8sVersion)
	data.Set("load_balancer_id", engine.LbId)
	data.Set("private_acl_resources", engine.PrivateAclResources)
	data.Set("public_acl_ip_address", engine.PublicAclIpAddress)
	data.Set("security_group_id", engine.SecurityGroupId)
	data.Set("subnet_id", engine.SubnetId)
	data.Set("volume_id", engine.VolumeId)
	data.Set("vpc_id", engine.VpcId)
	data.Set("zone_id", engine.ZoneId)
	data.Set("public_endpoint", engine.PublicEndpointUrl)
	data.Set("kube_config", kubeConfig)

	return nil
}

func updateEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("kubernetes_version", "public_acl_ip_address") {
		if data.HasChanges("public_acl_ip_address") && !data.HasChanges("kubernetes_version") {
			_, _, err := inst.Client.KubernetesEngine.UpdateEngine(ctx, data.Id(), kubernetesengine.UpdateEngineRequest{
				K8sVersion:         "",
				PublicAclIpAddress: data.Get("public_acl_ip_address").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else if data.HasChanges("kubernetes_version") && !data.HasChanges("public_acl_ip_address") {
			_, _, err := inst.Client.KubernetesEngine.UpdateEngine(ctx, data.Id(), kubernetesengine.UpdateEngineRequest{
				K8sVersion:         data.Get("kubernetes_version").(string),
				PublicAclIpAddress: "",
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdateEngine(ctx, data.Id(), kubernetesengine.UpdateEngineRequest{
				K8sVersion:         data.Get("kubernetes_version").(string),
				PublicAclIpAddress: data.Get("public_acl_ip_address").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("cloud_logging_enabled") {
		_, _, err := inst.Client.KubernetesEngine.UpdateLoggingEngine(ctx, data.Id(), kubernetesengine.UpdateEngineLoggingRequest{
			CloudLoggingEnabled: data.Get("cloud_logging_enabled").(bool),
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChanges("load_balancer_id") {
		if len(data.Get("load_balancer_id").(string)) == 0 {
			_, _, err := inst.Client.KubernetesEngine.UpdateLoadBalancerEngine(ctx, data.Id(), kubernetesengine.UpdateEngineLoadBalancerRequest{
				LoadBalancerEnabled: false,
				LbId:                data.Get("load_balancer_id").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdateLoadBalancerEngine(ctx, data.Id(), kubernetesengine.UpdateEngineLoadBalancerRequest{
				LoadBalancerEnabled: true,
				LbId:                data.Get("load_balancer_id").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("cifs_volume_id") {
		if len(data.Get("cifs_volume_id").(string)) == 0 {
			_, _, err := inst.Client.KubernetesEngine.UpdateCifsVolumeEngine(ctx, data.Id(), kubernetesengine.UpdateEngineCifsVolumeRequest{
				CifsVolumeIdEnabled: false,
				CifsVolumeId:        data.Get("cifs_volume_id").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdateCifsVolumeEngine(ctx, data.Id(), kubernetesengine.UpdateEngineCifsVolumeRequest{
				CifsVolumeIdEnabled: true,
				CifsVolumeId:        data.Get("cifs_volume_id").(string),
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("private_acl_resources") {
		if len(data.Get("private_acl_resources").([]interface{})) == 0 {
			var request []kubernetesengine.PrivateAclResourcesRequest
			_, _, err := inst.Client.KubernetesEngine.UpdatePrivateAclEngine(ctx, data.Id(), kubernetesengine.UpdateEnginePrivateAclRequest{
				PrivateAclResources: request,
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := inst.Client.KubernetesEngine.UpdatePrivateAclEngine(ctx, data.Id(), kubernetesengine.UpdateEnginePrivateAclRequest{
				PrivateAclResources: toPrivateAclResourcesRequestList(data.Get("private_acl_resources").([]interface{})),
			})

			if err != nil {
				return diag.FromErr(err)
			}

			time.Sleep(10 * time.Second)
			err = client.WaitForStatus(ctx, inst.Client, []string{"UPDATING"}, []string{"RUNNING"}, refreshEngine(ctx, meta, data.Id(), true))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return readEngine(ctx, data, meta)
}

func deleteEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// TODO: Delete node pool first

	_, err := inst.Client.KubernetesEngine.DeleteEngine(ctx, data.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = client.WaitForStatus(ctx, inst.Client, []string{"DELETING"}, []string{"DELETED"}, refreshEngine(ctx, meta, data.Id(), false))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func refreshEngine(ctx context.Context, meta interface{}, id string, errorOnNotFound bool) func() (interface{}, string, error) {
	inst := meta.(*client.Instance)

	return func() (interface{}, string, error) {
		engine, httpStatus, err := inst.Client.KubernetesEngine.ReadEngine(ctx, id)

		if httpStatus == 200 {
			return engine, engine.KubernetesEngineStatus, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("kubernetes engine with id=%s not found", id)
			}

			return engine, "DELETED", nil
		} else if err != nil {
			return nil, "", err
		}

		return nil, "", fmt.Errorf("failed to read kubernetes engine(%s) status:%d", id, httpStatus)
	}
}

func toPrivateAclResourcesRequestList(list []interface{}) []kubernetesengine.PrivateAclResourcesRequest {
	if len(list) == 0 {
		return nil
	}
	var result []kubernetesengine.PrivateAclResourcesRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)

		result = append(result, kubernetesengine.PrivateAclResourcesRequest{
			Id:    kv["id"].(string),
			Type:  kv["type"].(string),
			Value: kv["value"].(string),
		})
	}
	return result
}
