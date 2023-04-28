package api_manage_vo

type ApiProducerInfo struct {
	ProducerId  int    `json:"producerId"`
	ServiceCode string `json:"serviceCode"`
	CategoryId  string `json:"categoryId"`
	ProjectCode string `json:"projectCode"`
}
