package params

type ConfigType struct {
	Config string `json:"config" form:"config" query:"config"`
}

type UpdateConfig struct {
	Name       string `json:"name" form:"name"`
	Value      string `json:"value" form:"value"`
	ConfigType string `json:"conftype" form:"conftype"`
}
