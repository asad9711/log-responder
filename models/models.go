package models

type ReqPayload struct {
	NumberOfLines string `json:"numberOfLines"`
	FileName      string `json:"fileName"`
}
