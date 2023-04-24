package base_proxy_vo

type HttpResult[T any] struct {
	Data    T      `json:"data"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
