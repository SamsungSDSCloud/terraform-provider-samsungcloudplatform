package placementgroup

type CreateRequest struct {
	// Availability zone name
	AvailabilityZoneName string
	// Placement Group name
	PlacementGroupName string
	// Zone ID serviceZoneId is obtained through @[View Project Details]
	ServiceZoneId string
	Tags          []TagRequest
	// Virtual server type
	VirtualServerType string
	// Placement Group description
	PlacementGroupDescription string
}

type ListPlacementGroupsRequestParam struct {
	PlacementGroupName      string
	PlacementGroupStateList []string
	ServiceZoneId           string
	VirtualServerType       string
	CreatedBy               string
	Page                    int32
	Size                    int32
	Sort                    string
}

type UpdatePlacementGroupDescriptionRequest struct {
	PlacementGroupDescription string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
