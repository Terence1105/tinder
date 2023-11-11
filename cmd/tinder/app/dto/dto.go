package dto

type AddSinglePersonAndMatchReq struct {
	Name       string  `json:"name" binding:"required"`
	Height     float64 `json:"height" binding:"required"`
	Gender     int     `json:"gender" binding:"required"`
	DateCounts int     `json:"date_counts" binding:"required"`
}

type RemoveSinglePersonReq struct {
	Name   string `json:"name" binding:"required"`
	Gender int    `json:"gender" binding:"required"`
}

type QuerySinglePeopleReq struct {
	Counts int `json:"counts"`
}

type QuerySinglePeopleResp struct {
	Matches []Match `json:"matches"`
}
type Match struct {
	Boy  string
	Girl string
}
