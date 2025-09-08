package mapper

type APIMapper struct {
	baseURL string
}

func NewAPIMapper(baseURL string) *APIMapper {
	return &APIMapper{
		baseURL: baseURL,
	}
}
