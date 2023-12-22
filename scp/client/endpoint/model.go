package endpoint

type CreateEndpointRequest struct {
	EndpointIpAddress   string
	EndpointName        string
	EndpointType        string
	ObjectId            string
	VpcId               string
	EndpointDescription string
	ServiceZoneId       string
	//Tags                []TagRequest
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
