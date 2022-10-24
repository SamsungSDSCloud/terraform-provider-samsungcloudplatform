package kubernetesapps

type ListStandardImageRequest struct {
	Category         string
	ImageId          string
	ImageName        string
	IsCarepack       string
	IsNew            string
	IsRecommended    string
	PricePolicy      string
	ProductGroupName string
	Page             int32
	Size             int32
	Sort             string
}
