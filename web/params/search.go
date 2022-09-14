package params

type Search struct {
	Page int    `bson:"page" form:"page,default=1"`
	Q    string `json:"q" form:"q"`
}
