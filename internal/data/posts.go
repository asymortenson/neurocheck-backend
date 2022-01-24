package data

type PostItem struct {
	Id      uint32 `json:"id"`
	OwnerId int    `json:"owner_id"`
	Date    uint32 `json:"date"`
	Text    string `json:"text"`
}

type Post struct {
	Count int32      `json:"count"`
	Items []PostItem `json:"items"`
}

type Response struct {
	Response *Post `json:"response"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RequestData struct {
	Offsets     int32   `json:"offsets"`
	AccessToken string  `json:"access_token"`
	Fields      []Field `json:"fields"`
}

type Comment struct {
	PostId   int32         `json:"post_id"`
	Comments []interface{} `json:"comments"`
}