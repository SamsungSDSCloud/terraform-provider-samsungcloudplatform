package gslb

type CreateGslbRequest struct {
	// GSLB Algorithm
	GslbAlgorithm string
	// GSLB Environment Usage
	GslbEnvUsage string
	// GSLB Health Check
	GslbHealthCheck GslbHealthCheckRequest
	// GSLB Name(Domain Name : omitted when creating)
	GslbName string
	// GSLB Resource
	GslbResources []GslbResourceRequest
	Tags          map[string]interface{}
}

type GslbHealthCheckRequest struct {
	// GSLB Health Check Interval
	GslbHealthCheckInterval int32
	// GSLB Health Check Timeout
	GslbHealthCheckTimeout int32
	// GSLB Health Check User ID
	GslbHealthCheckUserId string
	// GSLB Health Check User Password
	GslbHealthCheckUserPassword string
	// GSLB Health Check Response String
	GslbResponseString string
	// GSLB Health Check Send String
	GslbSendString string
	// GSLB Health Check Probe Timeout
	ProbeTimeout int32
	// GSLB Health Check Protocol
	Protocol string
	// GSLB Health Check Service Port
	ServicePort int32
}

type GslbResourceRequest struct {
	// GSLB Resource Destination
	GslbDestination string
	// GSLB Resource Region
	GslbRegion string
	// GSLB Resource Weight
	GslbResourceWeight int32
	// GSLB Resource Description
	GslbResourceDescription string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
