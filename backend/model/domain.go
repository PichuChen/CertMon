package model

type Domain struct {
	ID       int64  `json:"id"`
	Domain   string `json:"domain"`
	ValidTo  string `json:"valid_to"`
	DaysLeft int    `json:"days_left"`
	Status   string `json:"status"`
}
