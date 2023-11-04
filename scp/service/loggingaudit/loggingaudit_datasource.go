package loggingaudit

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	loggingaudit2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/logging-audit"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	scp.RegisterDataSource("scp_loggingaudits", DatasourceLoggingAudits())
	scp.RegisterDataSource("scp_loggingaudit", DatasourceLoggingAudit())
	scp.RegisterDataSource("scp_loggingaudit_users", DatasourceLoggingAuditUsers())
}

func DatasourceLoggingAudits() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLoggingAuditsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Logging object ID",
			},
			"object_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Logging object name",
			},
			"target_product_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Name of logging target products",
			},
			"target_regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target regions",
			},
			"target_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Logging target resources",
			},
			"product_offering": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALL", "PUBLIC", "PRIVATE", "GOV"}, false)),
				Description:      "Offering scope. One of ALL, PUBLIC, PRIVATE, GOV.",
			},
			"request_client_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Console", "Api", "System"}, false)),
				Description:      "requesting client type. One of Console, Api, System.",
			},
			"request_start_dt": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Request start date. Default : current date",
			},
			"request_end_dt": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Request start date. Default : date 3 months ago",
			},
			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Success", "Fail"}, false)),
				Description:      "Job result's state. Success or Fail",
			},
			"page": {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Request page number"},
			"size": {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Page count"},
			"sort": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User name",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count",
			},
			"contents": {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceLogElem()},
		},
		Description: "Provides the activity log list",
	}
}

func datasourceLogElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":           {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"audit_content":        {Type: schema.TypeString, Computed: true, Description: "Audit content"},
			"audit_detail_content": {Type: schema.TypeString, Computed: true, Description: "Audit details"},
			"cluster_id":           {Type: schema.TypeString, Computed: true, Description: "Cluster ID"},
			"cluster_namespace_id": {Type: schema.TypeString, Computed: true, Description: "Cluster namespace ID"},
			"event_topic_name":     {Type: schema.TypeString, Computed: true, Description: "Event topic name"},
			"id":                   {Type: schema.TypeString, Computed: true, Description: "Logging ID"},
			"log_error_message":    {Type: schema.TypeString, Computed: true, Description: "Log error message"},
			"object_id":            {Type: schema.TypeString, Computed: true, Description: "Object ID"},
			"object_name":          {Type: schema.TypeString, Computed: true, Description: "Object name"},
			"product_name":         {Type: schema.TypeString, Computed: true, Description: "Product name"},
			"project_name":         {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"region":               {Type: schema.TypeString, Computed: true, Description: "Region"},
			"requested_by":         {Type: schema.TypeString, Computed: true, Description: "Requester"},
			"request_client_type":  {Type: schema.TypeString, Computed: true, Description: "requesting client type"},
			"request_dt":           {Type: schema.TypeString, Computed: true, Description: "Request date"},
			"resource_type":        {Type: schema.TypeString, Computed: true, Description: "Resource type"},
			"state":                {Type: schema.TypeString, Computed: true, Description: "Job's result state"},
			"user_email":           {Type: schema.TypeString, Computed: true, Description: "User email"},
			"user_name":            {Type: schema.TypeString, Computed: true, Description: "User name"},
		},
	}
}

func datasourceLoggingAuditsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	objectId := rd.Get("object_id").(string)
	objectName := rd.Get("object_name").(string)
	userName := rd.Get("user_name").(string)
	targetProductNames := common.ToStringList(rd.Get("target_product_names").(*schema.Set).List())
	targetRegions := common.ToStringList(rd.Get("target_regions").(*schema.Set).List())
	targetResources := common.ToStringList(rd.Get("target_resources").(*schema.Set).List())
	productOffering := rd.Get("product_offering").(string)
	requestClientType := rd.Get("request_client_type").(string)
	requestStartDt := rd.Get("request_start_dt").(string)
	requestEndDt := rd.Get("request_end_dt").(string)
	state := rd.Get("state").(string)
	page := rd.Get("page").(int)
	size := rd.Get("size").(int)
	sort := common.ToStringList(rd.Get("sort").(*schema.Set).List())

	startDt, err := time.Parse("2021-09-30T00:00:00.000Z", requestStartDt)
	endDt, err := time.Parse("2021-09-30T00:00:00.000Z", requestEndDt)

	// time은 빈 값으로 omit이 되질 않아서 값을 넣어줌, api가 수정되어야 할 것 같다만...
	if requestEndDt == "" || requestStartDt == "" {
		endDt = time.Now()
		startDt = endDt.AddDate(0, -3, 0)
	}

	// utc가 아니면 bad request가 뜸 ㅠㅠ
	startDt = startDt.UTC()
	endDt = endDt.UTC()

	response, _, err := inst.Client.Loggingaudit.ListLoggings(ctx, loggingaudit2.LoggingSearchCriteria{
		LoggingObjectId:          objectId,
		LoggingTargetProductName: targetProductNames,
		LoggingTargetRegion:      targetRegions,
		LoggingTargetResource:    targetResources,
		ObjectName:               objectName,
		Page:                     int32(page),
		ProductOffering:          productOffering,
		RequestClientType:        requestClientType,
		RequestEndDt:             endDt,
		RequestStartDt:           startDt,
		Size:                     int32(size),
		Sort:                     sort,
		State:                    state,
		UserName:                 userName,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertTrailListToHclSet(response)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceLoggingAudits().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", response.TotalCount)

	return nil
}

func convertTrailListToHclSet(logs loggingaudit2.PageResponseV2OfLoggingsResponse) common.HclSetObject {
	var logList common.HclSetObject
	for _, log := range logs.Contents {
		if len(log.Id) == 0 {
			continue
		}

		kv := common.HclKeyValueObject{
			"project_id":           log.ProjectId,
			"audit_content":        log.AuditContent,
			"audit_detail_content": log.AuditDetailContent,
			"cluster_id":           log.ClusterId,
			"cluster_namespace_id": log.ClusterNamespaceId,
			"event_topic_name":     log.EventTopicName,
			"id":                   log.Id,
			"log_error_message":    log.LogErrorMessage,
			"object_id":            log.ObjectId,
			"object_name":          log.ObjectName,
			"product_name":         log.ProductName,
			"project_name":         log.ProjectName,
			"region":               log.Region,
			"requested_by":         log.RequestBy,
			"request_client_type":  log.RequestClientType,
			"request_dt":           log.RequestDt.String(),
			"resource_type":        log.ResourceType,
			"state":                log.State,
			"user_email":           log.UserEmail,
			"user_name":            log.UserName,
		}
		logList = append(logList, kv)
	}
	return logList
}

func DatasourceLoggingAudit() *schema.Resource {
	var logResource schema.Resource
	logResource.ReadContext = datasourceLoggingAuditRead
	logResource.Schema = datasourceLogElem().Schema

	logResource.Schema["logging_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Logging ID",
	}
	logResource.Description = "Provides detailed logging information for a given logging id"

	return &logResource
}

func datasourceLoggingAuditRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	id := rd.Get("logging_id").(string)
	logInfo, _, err := inst.Client.Loggingaudit.DetailLogging(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(id)
	rd.Set("project_id", logInfo.ProjectId)
	rd.Set("audit_content", logInfo.AuditContent)
	rd.Set("audit_detail_content", logInfo.AuditDetailContent)
	rd.Set("cluster_id", logInfo.ClusterId)
	rd.Set("cluster_namespace_id", logInfo.ClusterNamespaceId)
	rd.Set("event_topic_name", logInfo.EventTopicName)
	rd.Set("id", logInfo.Id)
	rd.Set("log_error_message", logInfo.LogErrorMessage)
	rd.Set("object_id", logInfo.ObjectId)
	rd.Set("object_name", logInfo.ObjectName)
	rd.Set("product_name", logInfo.ProductName)
	rd.Set("project_name", logInfo.ProjectName)
	rd.Set("region", logInfo.Region)
	rd.Set("requested_by", logInfo.RequestBy)
	rd.Set("request_client_type", logInfo.RequestClientType)
	rd.Set("request_dt", logInfo.RequestDt)
	rd.Set("resource_type", logInfo.ResourceType)
	rd.Set("state", logInfo.State)
	rd.Set("user_email", logInfo.UserEmail)
	rd.Set("user_name", logInfo.UserName)

	return nil
}

func DatasourceLoggingAuditUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLoggingAuditUsersRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User name",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count",
			},
			"contents": {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceLogUserElem()},
		},
		Description: "Provides a datasource that retrieves a list of project target users by project ID",
	}
}

func datasourceLogUserElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email":     {Type: schema.TypeString, Computed: true, Description: "Email"},
			"user_id":   {Type: schema.TypeString, Computed: true, Description: "User ID"},
			"user_name": {Type: schema.TypeString, Computed: true, Description: "User Name"},
		},
	}
}

func datasourceLoggingAuditUsersRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	response, _, err := inst.Client.Loggingaudit.ListUsers(ctx, rd.Get("user_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var userList common.HclSetObject
	for _, user := range response.Contents {
		if len(user.UserId) == 0 {
			continue
		}

		kv := common.HclKeyValueObject{
			"user_id":   user.UserId,
			"user_name": user.UserName,
			"email":     user.Email,
		}
		userList = append(userList, kv)
	}

	if f, ok := rd.GetOk("filter"); ok {
		userList = common.ApplyFilter(DatasourceLoggingAuditUsers().Schema, f.(*schema.Set), userList)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", userList)
	rd.Set("total_count", len(userList))

	return nil
}
