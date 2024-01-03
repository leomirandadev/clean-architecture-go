package models

type GoogleCallbackReq struct {
	State string `json:"state"`
	Code  string `json:"code"`
}
