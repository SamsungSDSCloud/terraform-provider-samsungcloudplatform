package kubernetes

import (
	"context"
	"fmt"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/client/kubernetesengine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

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
				Computed:    true,
				ForceNew:    true,
				Description: "Load balancer ID",
			},
			"public_acl_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
			/*
				"cifs_volume_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "CIFS volume id",
				},
			*/
		},
		Description: "Provides a K8s Engine resource.",
	}
}

func createEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	vpcId := data.Get("vpc_id").(string)
	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)

	//cifsVolumeId := data.Get("cifs_volume_id").(string)

	response, _, err := inst.Client.KubernetesEngine.CreateEngine(ctx, kubernetesengine.CreateEngineRequest{
		CloudLoggingEnabled:  data.Get("cloud_logging_enabled").(bool),
		K8sVersion:           data.Get("kubernetes_version").(string),
		KubernetesEngineName: data.Get("name").(string),
		LbId:                 data.Get("load_balancer_id").(string),
		PublicAclIpAddress:   data.Get("public_acl_ip_address").(string),
		SecurityGroupId:      data.Get("security_group_id").(string),
		SubnetId:             data.Get("subnet_id").(string),
		VolumeId:             data.Get("volume_id").(string),
		//CifsVolumeId:         cifsVolumeId,
		VpcId:  vpcId,
		ZoneId: vpcInfo.ServiceZoneId,
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
		return diag.FromErr(err)
	}

	kubeConfig, _, err := inst.Client.KubernetesEngine.GetKubeConfig(ctx, data.Id())

	data.Set("name", engine.KubernetesEngineName)
	data.Set("cloud_logging_enabled", engine.CloudLoggingEnabled)
	data.Set("kubernetes_version", engine.K8sVersion)
	data.Set("load_balancer_id", engine.LbId)
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
		_, _, err := inst.Client.KubernetesEngine.UpdateEngine(ctx, data.Id(), kubernetesengine.UpdateEngineRequest{
			K8sVersion:         data.Get("kubernetes_version").(string),
			PublicAclIpAddress: data.Get("public_acl_ip_address").(string),
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readEngine(ctx, data, meta)
}

func deleteEngine(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// TODO: Delete node pool first

	_, err := inst.Client.KubernetesEngine.DeleteEngine(ctx, data.Id())
	if err != nil {
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
