package securitygroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	securitygroup2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/security-group2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	scp.RegisterDataSource("scp_security_group", DatasourceSecurityGroupInfo())
	scp.RegisterDataSource("scp_security_groups", DatasourceSecurityGroups())
	scp.RegisterDataSource("scp_security_group_user_ips", DatasourceSecurityGroupUserIps())
	scp.RegisterDataSource("scp_security_group_rule", DatasourceSecurityGroupRuleInfo())
	scp.RegisterDataSource("scp_security_group_rules", DatasourceSecurityGroupRules())
	scp.RegisterDataSource("scp_security_group_log_storage", DatasourceSecurityGroupLogStorageInfo())
	scp.RegisterDataSource("scp_security_group_log_storages", DatasourceSecurityGroupLogStorages())
}

func DatasourceSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"contents": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security Group list", Elem: datasourceSecurityGroupsElem(),
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total list size",
			},
		},
		Description: "Provides list of Security Groups",
	}
}

func datasourceSecurityGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	option := securitygroup2.SecurityGroupOpenApiControllerV2ApiListSecurityGroupV2Opts{
		IsLoggable:          optional.Bool{},
		Scopes:              optional.Interface{},
		SecurityGroupName:   optional.String{},
		SecurityGroupStates: optional.Interface{},
		VpcId:               optional.String{},
		CreatedBy:           optional.String{},
		Page:                optional.Int32{},
		Size:                optional.NewInt32(10000),
		Sort:                optional.Interface{},
	}

	responses, err := inst.Client.SecurityGroup.ListSecurityGroups(ctx, &option)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceSecurityGroups().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func datasourceSecurityGroupsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("AttachedObjectCount"): {Type: schema.TypeInt, Computed: true, Description: "The number of Objects which is attached with Security Group"},
			common.ToSnakeCase("IsLoggable"):          {Type: schema.TypeBool, Computed: true, Description: "Is loggable"},
			common.ToSnakeCase("RuleCount"):           {Type: schema.TypeInt, Computed: true, Description: "The number of Rules"},
			common.ToSnakeCase("Scope"):               {Type: schema.TypeString, Computed: true, Description: "Security Group Scope of Use"},
			common.ToSnakeCase("SecurityGroupId"):     {Type: schema.TypeString, Computed: true, Description: "Security Group ID"},
			common.ToSnakeCase("SecurityGroupName"):   {Type: schema.TypeString, Computed: true, Description: "Security Group name"},
			common.ToSnakeCase("SecurityGroupState"):  {Type: schema.TypeString, Computed: true, Description: "Security Group state"},
			common.ToSnakeCase("ZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Zone ID of Resource"},
			common.ToSnakeCase("CreatedBy"):           {Type: schema.TypeString, Computed: true, Description: "Resource creator"},
			common.ToSnakeCase("CreatedDt"):           {Type: schema.TypeString, Computed: true, Description: "Resource created date"},
			common.ToSnakeCase("ModifiedBy"):          {Type: schema.TypeString, Computed: true, Description: "Resource last modifier"},
			common.ToSnakeCase("ModifiedDt"):          {Type: schema.TypeString, Computed: true, Description: "Resource last modified date"},
		},
	}
}

func DatasourceSecurityGroupInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupInfoRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("SecurityGroupId"):          {Type: schema.TypeString, Required: true, Description: "Security Group ID"},
			common.ToSnakeCase("IsLoggable"):               {Type: schema.TypeBool, Computed: true, Description: "Is loggable"},
			common.ToSnakeCase("RuleCount"):                {Type: schema.TypeInt, Computed: true, Description: "The number of Rules"},
			common.ToSnakeCase("Scope"):                    {Type: schema.TypeString, Computed: true, Description: "Security Group Scope of Use"},
			common.ToSnakeCase("SecurityGroupName"):        {Type: schema.TypeString, Computed: true, Description: "Security Group name"},
			common.ToSnakeCase("SecurityGroupState"):       {Type: schema.TypeString, Computed: true, Description: "Security Group state"},
			common.ToSnakeCase("VendorObjectId"):           {Type: schema.TypeString, Computed: true, Description: "Vendor Object ID"},
			common.ToSnakeCase("VpcId"):                    {Type: schema.TypeString, Computed: true, Description: "VPC ID"},
			common.ToSnakeCase("ZoneId"):                   {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			common.ToSnakeCase("SecurityGroupDescription"): {Type: schema.TypeString, Computed: true, Description: "Security Group description"},
			common.ToSnakeCase("CreatedBy"):                {Type: schema.TypeString, Computed: true, Description: "creator"},
			common.ToSnakeCase("CreatedDt"):                {Type: schema.TypeString, Computed: true, Description: "created datetime"},
			common.ToSnakeCase("ModifiedBy"):               {Type: schema.TypeString, Computed: true, Description: "last modified user"},
			common.ToSnakeCase("ModifiedDt"):               {Type: schema.TypeString, Computed: true, Description: "Resource modified datetime"},
			common.ToSnakeCase("ProjectId"):                {Type: schema.TypeString, Computed: true, Description: "Project ID"},
		},
		Description: "Provides Security Group Info",
	}
}

func datasourceSecurityGroupInfoRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	SecurityGroupId := rd.Get("security_group_id").(string)

	info, _, err := inst.Client.SecurityGroup.GetSecurityGroup(ctx, SecurityGroupId)

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set(common.ToSnakeCase("IsLoggable"), info.IsLoggable)
	rd.Set(common.ToSnakeCase("RuleCount"), info.RuleCount)
	rd.Set(common.ToSnakeCase("Scope"), info.Scope)
	rd.Set(common.ToSnakeCase("SecurityGroupId"), info.SecurityGroupId)
	rd.Set(common.ToSnakeCase("SecurityGroupName"), info.SecurityGroupName)
	rd.Set(common.ToSnakeCase("SecurityGroupState"), info.SecurityGroupState)
	rd.Set(common.ToSnakeCase("VendorObjectId"), info.VendorObjectId)
	rd.Set(common.ToSnakeCase("VpcId"), info.VpcId)
	rd.Set(common.ToSnakeCase("ZoneId"), info.ZoneId)
	rd.Set(common.ToSnakeCase("SecurityGroupDescription"), info.SecurityGroupDescription)
	rd.Set(common.ToSnakeCase("CreatedBy"), info.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), time.Time.String(info.CreatedDt))
	rd.Set(common.ToSnakeCase("ModifiedBy"), info.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), time.Time.String(info.ModifiedDt))
	rd.Set(common.ToSnakeCase("ProjectId"), info.ProjectId)

	rd.SetId(uuid.NewV4().String())

	return nil
}

func DatasourceSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupRulesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security Group ID",
			},
			"contents": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security Group Rule list", Elem: datasourceSecurityGroupRulesElem(),
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total list size",
			},
		},
		Description: "Provides list of Security Group Rules",
	}
}

func datasourceSecurityGroupRulesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	SecurityGroupId := rd.Get("security_group_id").(string)

	option := securitygroup2.SecurityGroupOpenApiControllerV2ApiListSecurityGroupRuleV2Opts{
		IpAddress:       optional.String{},
		RuleDescription: optional.String{},
		RuleDirection:   optional.String{},
		RuleStates:      optional.Interface{},
		Page:            optional.Int32{},
		Size:            optional.NewInt32(10000),
		Sort:            optional.Interface{},
	}

	responses, err := inst.Client.SecurityGroup.ListSecurityGroupRules(ctx, SecurityGroupId, &option)
	if err != nil {
		return diag.FromErr(err)
	}

	setSecurityGroupRules := convertSecurityGroupRuleListToHclSet(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		setSecurityGroupRules = common.ApplyFilter(DatasourceSecurityGroupRules().Schema, f.(*schema.Set), setSecurityGroupRules)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", setSecurityGroupRules)
	rd.Set("total_count", len(setSecurityGroupRules))

	return nil
}

func datasourceSecurityGroupRulesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("IsAllService"):  {Type: schema.TypeBool, Computed: true, Description: "Is all Service"},
			common.ToSnakeCase("RuleAction"):    {Type: schema.TypeString, Computed: true, Description: "Rule action"},
			common.ToSnakeCase("RuleDirection"): {Type: schema.TypeString, Computed: true, Description: "Rule direction"},
			common.ToSnakeCase("RuleId"):        {Type: schema.TypeString, Computed: true, Description: "Rule ID"},
			common.ToSnakeCase("RuleOwnerId"):   {Type: schema.TypeString, Computed: true, Description: "Rule Owner ID"},
			common.ToSnakeCase("RuleOwnerType"): {Type: schema.TypeString, Computed: true, Description: "Rule Owner type"},
			common.ToSnakeCase("RuleState"):     {Type: schema.TypeString, Computed: true, Description: "Rule state"},
			common.ToSnakeCase("TargetNetworks"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Target networks",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("TcpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of TCP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("UdpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of UDP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("IcmpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ICMP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("VendorObjectId"):  {Type: schema.TypeString, Computed: true, Description: "Vendor Object ID"},
			common.ToSnakeCase("VendorRuleId"):    {Type: schema.TypeInt, Computed: true, Description: "Vendor Rule ID"},
			common.ToSnakeCase("RuleDescription"): {Type: schema.TypeString, Computed: true, Description: "Rule description"},
		},
	}
}

func convertSecurityGroupRuleListToHclSet(securityGroupRules []securitygroup2.SecurityGroupRuleResponse) common.HclSetObject {
	var securityGroupRuleList common.HclSetObject
	for _, rule := range securityGroupRules {
		if len(rule.RuleId) == 0 {
			continue
		}
		kv := common.HclKeyValueObject{
			common.ToSnakeCase("RuleId"):          rule.RuleId,
			common.ToSnakeCase("IcmpServices"):    rule.IcmpServices,
			common.ToSnakeCase("IsAllService"):    rule.IsAllService,
			common.ToSnakeCase("RuleAction"):      rule.RuleAction,
			common.ToSnakeCase("RuleDirection"):   rule.RuleDirection,
			common.ToSnakeCase("RuleOwnerId"):     rule.RuleOwnerId,
			common.ToSnakeCase("RuleOwnerType"):   rule.RuleOwnerType,
			common.ToSnakeCase("RuleState"):       rule.RuleState,
			common.ToSnakeCase("TargetNetworks"):  rule.TargetNetworks,
			common.ToSnakeCase("TcpServices"):     rule.TcpServices,
			common.ToSnakeCase("UdpServices"):     rule.UdpServices,
			common.ToSnakeCase("VendorObjectId"):  rule.VendorObjectId,
			common.ToSnakeCase("VendorRuleId"):    rule.VendorRuleId,
			common.ToSnakeCase("RuleDescription"): rule.RuleDescription,
		}
		securityGroupRuleList = append(securityGroupRuleList, kv)
	}
	return securityGroupRuleList
}

func DatasourceSecurityGroupRuleInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupRuleInfoRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security Group ID",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule ID",
			},
			common.ToSnakeCase("IsAllService"):  {Type: schema.TypeBool, Computed: true, Description: "Is all Service"},
			common.ToSnakeCase("RuleAction"):    {Type: schema.TypeString, Computed: true, Description: "Rule action"},
			common.ToSnakeCase("RuleDirection"): {Type: schema.TypeString, Computed: true, Description: "Rule direction"},
			common.ToSnakeCase("RuleOwnerId"):   {Type: schema.TypeString, Computed: true, Description: "Rule Owner ID"},
			common.ToSnakeCase("RuleOwnerType"): {Type: schema.TypeString, Computed: true, Description: "Rule Owner type"},
			common.ToSnakeCase("RuleState"):     {Type: schema.TypeString, Computed: true, Description: "Rule state"},
			common.ToSnakeCase("TargetNetworks"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Target networks",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("TcpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of TCP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("UdpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of UDP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("IcmpServices"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ICMP Services",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("VendorObjectId"):  {Type: schema.TypeString, Computed: true, Description: "Vendor Object ID"},
			common.ToSnakeCase("VendorRuleId"):    {Type: schema.TypeInt, Computed: true, Description: "Vendor Rule ID"},
			common.ToSnakeCase("RuleDescription"): {Type: schema.TypeString, Computed: true, Description: "Rule description"},
		},
		Description: "Provides list of Security Group Rule Info",
	}
}

func datasourceSecurityGroupRuleInfoRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	SecurityGroupId := rd.Get("security_group_id").(string)
	ruleId := rd.Get("rule_id").(string)

	info, _, err := inst.Client.SecurityGroup.GetSecurityGroupRule(ctx, ruleId, SecurityGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.Set(common.ToSnakeCase("IcmpServices"), info.IcmpServices)
	rd.Set(common.ToSnakeCase("IsAllService"), info.IsAllService)
	rd.Set(common.ToSnakeCase("RuleAction"), info.RuleAction)
	rd.Set(common.ToSnakeCase("RuleDirection"), info.RuleDirection)
	rd.Set(common.ToSnakeCase("RuleId"), info.RuleId)
	rd.Set(common.ToSnakeCase("RuleOwnerId"), info.RuleOwnerId)
	rd.Set(common.ToSnakeCase("RuleOwnerType"), info.RuleOwnerType)
	rd.Set(common.ToSnakeCase("RuleState"), info.RuleState)
	rd.Set(common.ToSnakeCase("TargetNetworks"), info.TargetNetworks)
	rd.Set(common.ToSnakeCase("TcpServices"), info.TcpServices)
	rd.Set(common.ToSnakeCase("UdpServices"), info.UdpServices)
	rd.Set(common.ToSnakeCase("VendorObjectId"), info.VendorObjectId)
	rd.Set(common.ToSnakeCase("VendorRuleId"), info.VendorRuleId)
	rd.Set(common.ToSnakeCase("RuleDescription"), info.RuleDescription)

	rd.SetId(uuid.NewV4().String())

	return nil
}

func DatasourceSecurityGroupLogStorageInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupLogStorageInfoRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID",
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bucket ID for saving logs",
			},
			common.ToSnakeCase("LogStorageId"):   {Type: schema.TypeString, Computed: true, Description: "Log Storage ID"},
			common.ToSnakeCase("LogStorageType"): {Type: schema.TypeString, Computed: true, Description: "Log Storage Type"},
			common.ToSnakeCase("ObsBucketId"):    {Type: schema.TypeString, Computed: true, Description: "Bucket ID for saving logs"},
			common.ToSnakeCase("CreatedBy"):      {Type: schema.TypeString, Computed: true, Description: "creator"},
			common.ToSnakeCase("CreatedDt"):      {Type: schema.TypeString, Computed: true, Description: "created datetime"},
			common.ToSnakeCase("ModifiedBy"):     {Type: schema.TypeString, Computed: true, Description: "last modified user"},
			common.ToSnakeCase("ModifiedDt"):     {Type: schema.TypeString, Computed: true, Description: "Resource modified datetime"},
			common.ToSnakeCase("ProjectId"):      {Type: schema.TypeString, Computed: true, Description: "Project ID"},
		},
		Description: "Provides Security Group Log Storage Info",
	}
}

func datasourceSecurityGroupLogStorageInfoRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	VpcId := rd.Get("vpc_id").(string)
	ObsBucketId := rd.Get("obs_bucket_id").(string)

	logStorages, _, err := inst.Client.SecurityGroup.ListSecurityGroupLogStorages(ctx, VpcId, ObsBucketId)
	if err != nil {
		return diag.FromErr(err)
	}

	setLogStorages := common.ConvertStructToMaps(logStorages.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		setLogStorages = common.ApplyFilter(DatasourceSecurityGroupLogStorageInfo().Schema, f.(*schema.Set), setLogStorages)
	}

	if len(setLogStorages) == 0 {
		return diag.Errorf("no matching security group log storage found")
	}

	for k, v := range setLogStorages[0] {
		rd.Set(k, v)
	}

	rd.SetId(uuid.NewV4().String())

	return nil
}

func DatasourceSecurityGroupLogStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupLogStoragesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC id",
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OBS Bucket ID",
			},
			"log_storages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security Group log storage list", Elem: common.GetDatasourceItemsSchema(DatasourceSecurityGroupLogStorageInfo()),
			},
		},
		Description: "Provides Security Group Log Storages",
	}
}

func datasourceSecurityGroupLogStoragesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	VpcId := rd.Get("vpc_id").(string)
	ObsBucketId := ""

	if _, ok := rd.GetOk("obs_bucket_id"); ok {
		ObsBucketId = rd.Get("obs_bucket_id").(string)
	}

	logStorages, _, err := inst.Client.SecurityGroup.ListSecurityGroupLogStorages(ctx, VpcId, ObsBucketId)
	if err != nil {
		return diag.FromErr(err)
	}

	setLogStorages := common.ConvertStructToMaps(logStorages.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		setLogStorages = common.ApplyFilter(DatasourceSecurityGroupLogStorageInfo().Schema, f.(*schema.Set), setLogStorages)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("log_storages", setLogStorages)

	return nil
}

func DatasourceSecurityGroupUserIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupUserIpsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security Group ID",
			},
			"contents": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User IP List Attached to Security Group", Elem: datasourceSecurityGroupUserIpsElem(),
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total list size",
			},
		},
		Description: "Provides list of User IPs Attached to Security Group",
	}
}

func datasourceSecurityGroupUserIpsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	SecurityGroupId := rd.Get("security_group_id").(string)

	responses, err := inst.Client.SecurityGroup.ListUserIpsBySecurityGroupId(ctx, SecurityGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceSecurityGroups().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func datasourceSecurityGroupUserIpsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("UserIpType"):        {Type: schema.TypeString, Computed: true, Description: "Type of Directly Attached IP"},
			common.ToSnakeCase("UserIpAddress"):     {Type: schema.TypeString, Computed: true, Description: "Address of Directly Attached IP"},
			common.ToSnakeCase("UserIpDescription"): {Type: schema.TypeString, Computed: true, Description: "Description of Directly Attached IP"},
			common.ToSnakeCase("State"):             {Type: schema.TypeString, Computed: true, Description: "IP attach state"},
		},
	}
}
