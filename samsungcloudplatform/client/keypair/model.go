package keypair

type CreateRequest struct {
	KeyPairName string
	Tags        map[string]interface{}
}

type TagRequest struct {
	TagKey   string
	TagValue string
}

type ListKeyPairsRequestParam struct {
	KeyPairName string
	CreatedBy   string
	Page        int32
	Size        int32
	Sort        string
}
