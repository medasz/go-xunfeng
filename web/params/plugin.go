package params

type Plugin struct {
	Type   string `json:"type" form:"type"`
	Risk   string `json:"risk" form:"risk"`
	Search string `json:"search" form:"search"`
}

type PluginList struct {
	Page int `json:"page" form:"page,default=1"`
}
