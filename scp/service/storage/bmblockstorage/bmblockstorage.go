package bmblockstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/bmblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	baremetalblockstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/bare-metal-block-storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	scp.RegisterResource("scp_bm_block_storage", ResourceBmBlockStorage())
}

func ResourceBmBlockStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createBmBlockStorage,
		ReadContext:   readBmBlockStorage,
		UpdateContext: updateBmBlockStorage,
		DeleteContext: deleteBmBlockStorage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The block storage name to create. (3 to 28 characters with -)",
				ValidateDiagFunc: common.ValidateName3to28Dash,
			},
			"storage_size_gb": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				Description:      "The storage size(GB) of the block storage to create. (10 to  16384 GB)",
				ValidateDiagFunc: validateStorageSize10to16384,
			},
			"bm_server_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:    5,
				Description: "Baremetal server IDs to which you want to assign the block storage.",
			},
			"product_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "You can use by selecting SSD or HDD based storage.",
			},
			"encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Encrypt the volume to be created and create it. When encryption is applied, performance degradation of around 10% occurs.",
			},
			"snapshot_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Use an additional 100-300% of the Block Storage capacity you created. If auto-creation is set, snapshots are created and saved automatically according to the specified cycle. You can restore using the saved snapshot.",
			},
			"snapshot_capacity_rate": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validateSnapshotCapacityRate100to500,
				Description:      "snapshot capacity rate(100 ~ 500)",
			},
			"snap_shot_schedule": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateDiagFunc: validateSnapShotSchedule,
				Description:      "schedule for snapshot",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a BM Block Storage resource.",
	}
}

func validateStorageSize10to16384(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := int32(v.(int))

	err := common.CheckInt32Range(value, 10, 16384)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func validateSnapshotCapacityRate100to500(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := int32(v.(int))

	err := common.CheckInt32Range(value, 100, 500)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("capacity rate must be between 100 and 500"),
			AttributePath: path,
		})
	}

	if value%50 != 0 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("capacity rate must be multiple of 50"),
			AttributePath: path,
		})
	}

	return diags
}

func validateSnapShotSchedule(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(map[string]interface{})

	dayOfWeek := value["day_of_week"].(string)
	frequency := value["frequency"].(string)
	hour, _ := strconv.Atoi(value["hour"].(string))

	if !regexp.MustCompile("^(|SUN|MON|TUE|WED|THU|FRI|SAT)$").MatchString(dayOfWeek) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("value must be SUN|MON|TUE|WED|THU|FRI|SAT"),
			AttributePath: path,
		})
	}

	if !regexp.MustCompile("^(NONE|DAILY|WEEKLY)$").MatchString(frequency) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("value must be NONE|DAILY|WEEKLY"),
			AttributePath: path,
		})
	}

	err := common.CheckInt32Range(int32(hour), 0, 23)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("hour must be between 0 and 23"),
			AttributePath: path,
		})
	}

	return diags
}

func createBmBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	baremetalServerIds := make([]string, 0)
	for _, baremetalServerId := range data.Get("bm_server_ids").([]interface{}) {
		baremetalServerIds = append(baremetalServerIds, baremetalServerId.(string))
	}

	snapshotScheduleInfo := data.Get("snap_shot_schedule").(map[string]interface{})

	snapshotSchedule := bmblockstorage.SnapshotSchedule{}

	if len(snapshotScheduleInfo) != 0 {

		snapshotSchedule.DayOfWeek = snapshotScheduleInfo["day_of_week"].(string)
		snapshotSchedule.Frequency = snapshotScheduleInfo["frequency"].(string)

		if snapshotScheduleInfo["hour"] != nil {

			hour, err := strconv.Atoi(snapshotScheduleInfo["hour"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

			snapshotSchedule.Hour = int32(hour)
		}

	}

	serverInfo, _, err := inst.Client.BareMetal.GetBareMetalServerDetail(ctx, baremetalServerIds[0])

	if err != nil {
		return diag.FromErr(err)
	}

	productId, _ := client.FindProductId(ctx, inst.Client, serverInfo.ProductGroupId, common.ProductTypeDisk, data.Get("product_name").(string))

	if len(productId) == 0 {
		return diag.Errorf("ERROR productId")
	}

	response, _, err := inst.Client.BareMetalBlockStorage.CreateBareMetalBlockStorage(ctx, bmblockstorage.BmBlockStorageCreateRequest{
		BareMetalBlockStorageName: data.Get("name").(string),
		BareMetalBlockStorageSize: (int32)(data.Get("storage_size_gb").(int)),
		EncryptionEnabled:         data.Get("encrypted").(bool),
		IsSnapshotPolicy:          data.Get("snapshot_policy").(bool),
		SnapshotSchedule:          snapshotSchedule,
		SnapshotCapacityRate:      (int32)(data.Get("snapshot_capacity_rate").(int)),
		BareMetalServerIds:        baremetalServerIds,
		ServiceZoneId:             serverInfo.ServiceZoneId,
		ProductId:                 productId,
		Tags:                      data.Get("tags").(map[string]interface{}),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForBmBlockStorageStatus(ctx, inst.Client, response.ResourceId, []string{common.CreatingState}, []string{common.ActiveState})
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ResourceId)

	return readBmBlockStorage(ctx, data, meta)
}

func readBmBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	serverInfo, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, data.Id())
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	bmServerIds := getBmServerIds(serverInfo)
	data.Set("bm_server_ids", bmServerIds)

	data.Set("name", serverInfo.BareMetalBlockStorageName)
	data.Set("storage_size_gb", serverInfo.BareMetalBlockStorageSize)
	data.Set("encrypted", serverInfo.EncryptionEnabled)
	data.Set("storage_size_gb", serverInfo.BareMetalBlockStorageSize)

	snapshotInfo, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageSnapshotList(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if len(snapshotInfo.Contents) != 0 {
		data.Set("snapshot_policy", snapshotInfo.Contents[0].IsSnapshotPolicy)
		data.Set("snapshot_capacity_rate", snapshotInfo.Contents[0].SnapshotCapacityRate)

		if len(snapshotInfo.Contents[0].Snapshots) != 0 {
			scheduleInfo, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageScheduleList(ctx, data.Id())

			if err != nil {
				return diag.FromErr(err)
			}

			if len(scheduleInfo.Contents) != 0 {
				data.Set("snap_shot_schedule", scheduleInfo.Contents[0].SnapshotSchedule)
			}
		}
	}

	tfTags.SetTags(ctx, data, meta, data.Id())
	return nil
}

func getBmServerIds(info baremetalblockstorage.BmBlockStorageDetailResponse) []string {
	bmServerIds := make([]string, 0)
	for _, bmServerId := range info.Servers {
		bmServerIds = append(bmServerIds, bmServerId.ServerId)
	}
	return bmServerIds
}

// 삭제하려면 detach를 해야해서 attach, detach만 구현
func updateBmBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("bm_server_ids") {
		deleteBmIds, addedBmIds := getDetachedAttachedVmIds(data)

		if len(deleteBmIds) != 0 {
			_, _, err := inst.Client.BareMetalBlockStorage.DetachBareMetalBlockStorage(ctx, data.Id(), deleteBmIds)

			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForBmBlockStorageStatus(ctx, inst.Client, data.Id(), []string{common.EditingState}, []string{common.ActiveState})
		}

		if len(addedBmIds) != 0 {
			_, _, err := inst.Client.BareMetalBlockStorage.AttachBareMetalBlockStorage(ctx, data.Id(), addedBmIds)

			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForBmBlockStorageStatus(ctx, inst.Client, data.Id(), []string{common.EditingState}, []string{common.ActiveState})
		}

	}

	tfTags.UpdateTags(ctx, data, meta, data.Id())

	return readBmBlockStorage(ctx, data, meta)
}

func getDetachedAttachedVmIds(data *schema.ResourceData) ([]string, []string) {
	oldValue, newValue := data.GetChange("bm_server_ids")
	oldValues := oldValue.([]interface{})
	newValues := newValue.([]interface{})
	oldBmIds := make([]string, len(oldValues))
	newBmIds := make([]string, len(newValues))
	for i, oldVmId := range oldValues {
		oldBmIds[i] = oldVmId.(string)
	}
	for i, newVmId := range newValues {
		newBmIds[i] = newVmId.(string)
	}

	// oldBmIds bmId가 newBmIds에 그대로 있는지 확인
	deletedBmIds := make([]string, 0)
	for _, oldBmId := range oldBmIds {
		var i int
		for i = 0; i < len(newBmIds); i++ {
			if strings.Compare(newBmIds[i], oldBmId) == 0 {
				break
			}
		}
		if i == len(newBmIds) {
			deletedBmIds = append(deletedBmIds, oldBmId)
		}
	}

	// oldBmIds bmId가 newBmIds에 그대로 있는지 확인
	addedBmIds := make([]string, 0)
	for _, newBmId := range newBmIds {
		var i int
		for i = 0; i < len(oldBmIds); i++ {
			if strings.Compare(newBmId, oldBmIds[i]) == 0 {
				break
			}
		}
		if i == len(oldBmIds) && len(newBmId) != 0 {
			addedBmIds = append(addedBmIds, newBmId)
		}
	}

	return deletedBmIds, addedBmIds

}

// Mount 상태를 Unmount 상태로 진행한 후에 삭제 진행해야한다.
func deleteBmBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, _, err := inst.Client.BareMetalBlockStorage.DeleteBareMetalBlockStorage(ctx, data.Id())

	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForBmBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{common.DeletedState})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForBmBlockStorageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, id)

		if err != nil {
			if c == 404 {
				return "", "DELETED", nil
			}

			if c == 403 {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, strings.ToUpper(info.BareMetalBlockStorageState), nil
	})
}
