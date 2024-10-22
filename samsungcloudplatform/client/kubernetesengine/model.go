package kubernetesengine

type CreateEngineRequest struct {
	CloudLoggingEnabled  bool
	K8sVersion           string
	KubernetesEngineName string
	LbId                 string
	PrivateAclResources  []PrivateAclResourcesRequest
	PublicAclIpAddress   string
	SecurityGroupId      string
	SubnetId             string
	VolumeId             string
	CifsVolumeId         string
	VpcId                string
	ZoneId               string
	Tags                 map[string]interface{}
}

type ListEngineRequest struct {
	K8sVersion             []string
	KubernetesEngineName   string
	KubernetesEngineStatus []string
	Region                 []string
	CreatedBy              string
	Page                   int32
	Size                   int32
	Sort                   string
}

type UpdatePublicEndpointAccessControlRequest struct {
	PublicAclIpAddress string
}

type UpgradeRequest struct {
	K8sVersion string
}

type UpdateEngineLoggingRequest struct {
	CloudLoggingEnabled bool
}

type UpdateEngineLoadBalancerRequest struct {
	LbId string
}

type UpdateEngineCifsVolumeRequest struct {
	CifsVolumeId string
}

type UpdateEnginePrivateAclRequest struct {
	PrivateAclResourcesToUpdate []PrivateAclResourcesRequestToUpdate
}

type CreateNodePoolRequest struct {
	AvailabilityZoneName string
	AutoRecovery         bool
	AutoScale            bool
	ContractName         string
	DesiredNodeCount     int32
	EncryptEnabled       bool
	ImageId              string
	MaxNodeCount         int32
	MinNodeCount         int32
	NodePoolName         string
	ScaleName            string
	ServerType           string
	ServiceLevelName     string
	StorageName          string
	StorageSize          string
}

type NodePoolUpdateRequest struct {
	AutoRecovery     bool
	AutoScale        bool
	ContractId       string
	DesiredNodeCount int32
	ImageId          string
	MaxNodeCount     int32
	MinNodeCount     int32
	NodePoolName     string
	ProductGroupId   string
	ScaleId          string
	ServiceLevelId   string
	StorageId        string
	StorageSize      string
}

type PrivateAclResourcesRequest struct {
	Id    string
	Type  string
	Value string
}

type PrivateAclResourcesRequestToUpdate struct {
	ResourceId    string
	ResourceType  string
	ResourceValue string
}
