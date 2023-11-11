package dto

type AddSinglePersonAndMatchReq struct {
	Name       string `json:"name"`
	Height     int    `json:"height"`
	Gender     int    `json:"gender"`
	DateCounts int    `json:"date_counts"`
}
