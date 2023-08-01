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

type UpdateEngineRequest struct {
	K8sVersion         string
	PublicAclIpAddress string
}

type UpdateEngineLoggingRequest struct {
	CloudLoggingEnabled bool
}

type UpdateEngineLoadBalancerRequest struct {
	LoadBalancerEnabled bool
	LbId                string
}

type UpdateEngineCifsVolumeRequest struct {
	CifsVolumeIdEnabled bool
	CifsVolumeId        string
}

type UpdateEnginePrivateAclRequest struct {
	PrivateAclResources []PrivateAclResourcesRequest
}

type CreateNodePoolRequest struct {
	AvailabilityZoneName string
	AutoRecovery         bool
	AutoScale            bool
	ContractId           string
	DesiredNodeCount     int32
	ImageId              string
	MaxNodeCount         int32
	MinNodeCount         int32
	NodePoolName         string
	ProductGroupId       string
	ScaleId              string
	ServiceLevelId       string
	StorageId            string
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
