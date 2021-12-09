package models

type AllInfo struct {
	TableID          string
	Id               int    `json:"id"`
	ShopId           string `json:"app_id"`
	UserId           string `json:"user_id"`
	ResourceType     int    `json:"resource_type"`
	ResourceId       string `json:"resource_id"`
	LearnProgress    int    `json:"learn_progress"`
	MaxLearnProgress int    `json:"max_learn_progress"`
	OrgLearnProgress string `json:"org_learn_progress"`
	AgentType        int    `json:"agent_type"`
	IsFinish         int    `json:"is_finish"`
	Fa               string `json:"finished_at"`
	St               int    `json:"stay_time"`
	Spt              int    `json:"spend_time"`
	Llt              string `json:"last_learn_time"`
	Ca               string `json:"created_at"`
	Ua               string `json:"updated_at"`
	Cappid           string `json:"content_app_id"`
	Dstate           int    `json:"display_state"`
	ProductID        string `json:"product_id"`
	State            int    `json:"state"`
	ResLen           int
}

type ResInfo struct {
	Id  string `json:"id"`
	Len string `json:"video_length"`
}

type TraningInfo struct {
	Id         int    `json:"id"`
	ShopId     string `json:"app_id"`
	UserId     string `json:"user_id"`
	ResourceId string
	TermId     string
}
