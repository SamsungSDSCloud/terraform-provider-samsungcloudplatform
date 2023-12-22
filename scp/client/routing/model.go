package routing

type ListVpcRoutingTableRequest struct {
	RoutingTableId   string
	RoutingTableName string
	VpcId            string
	CreatedBy        string
	Sort             string
}

type ListVpcRoutingRulesRequest struct {
	DestinationNetworkCidr   string
	Editable                 string
	RoutingRuleId            string
	SourceServiceInterfaceId string
	Sort                     string
}

type CreateRoutingRulesRequest struct {
	// Routing Rules
	RoutingRules []RoutingRule
}

type RoutingRule struct {
	// Destination Network
	DestinationNetworkCidr string
	// Source Service Interface ID
	SourceServiceInterfaceId string
	// Source Service Interface Name
	SourceServiceInterfaceName string
}

type DeleteRoutingRulesRequest struct {
	// Routing Rule IDs
	RoutingRuleIds []string
}

type ListTgwRoutingTableRequest struct {
	RoutingTableId             string
	RoutingTableName           string
	TransitGatewayConnectionId string
	CreatedBy                  string
	Sort                       string
}

type ListTgwRoutingRuleRequest struct {
	DestinationNetworkCidr   string
	Editable                 string
	RoutingRuleId            string
	SourceServiceInterfaceId string
	Sort                     string
}
