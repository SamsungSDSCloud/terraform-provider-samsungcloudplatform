package product

type ListCategoriesRequest struct {
	CategoryId    string
	CategoryState string
	ExposureScope string
	ProductId     string
	ProductState  string
	LanguageCode  string
}

type ListMenusRequest struct {
	CategoryId    string
	ExposureType  string
	ExposureScope string
	ProductId     string
	ZoneIds       string
}
