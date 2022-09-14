package params

type Task struct {
	Title       string `json:"title" form:"title"`
	Condition   string `json:"condition" form:"condition"`
	Plugin      string `json:"plugin" form:"plugin"`
	Ids         string `json:"ids" form:"ids"`
	Plan        int    `json:"plan" form:"plan"`
	IsUpdate    bool   `json:"isupdate" form:"isupdate"`
	ResultCheck bool   `json:"resultcheck" form:"resultcheck"`
}

type TaskList struct {
	Page int `json:"page" form:"page,default=1"`
}

type TaskId struct {
	Oid string `json:"oid" form:"oid"`
}

type TaskDownloadXls struct {
	TaskId   string `json:"taskid" form:"taskid"`
	TaskDate string `json:"taskdate" form:"taskdate"`
}
