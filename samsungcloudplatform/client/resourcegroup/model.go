package resourcegroup

type ResourceGroupRequest struct {
	ResourceGroupName        string
	TargetResourceTags       []interface{}
	TargetResourceTypes      []interface{}
	ResourceGroupDescription string
}

type ListResourceGroupRequest struct {
	CreatedById       string
	ModifiedByEmail   string
	ModifiedById      string
	ResourceGroupName string
}

type DeleteResourceGroupRequest struct {
	ResourceGroupIds []string
}

type ListResourceGroupResourcesRequest struct {
	CreatedById  string
	ModifiedById string
	ResourceId   string
	ResourceName string
}

type ListResourceRequest struct {
	CreatedById         string
	DisplayServiceNames []interface{}
	FromCreatedAt       string
	IncludeDeleted      string
	Location            string
	ModifiedById        string
	MyCreate            string
	Partitions          []interface{}
	Regions             []interface{}
	ResourceId          string
	ResourceName        string
	ResourceTypes       []interface{}
	ServiceTypes        []interface{}
	ServiceZones        []interface{}
	Tags                []interface{}
	ToCreatedAt         string
}
