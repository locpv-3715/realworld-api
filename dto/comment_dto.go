package dto

import "time"

type CommentReq struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment"`
}

type CommentRes struct {
	Comment struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Body      string    `json:"body"`
		Author    struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"comment"`
}

type MultipleCommentsRes struct {
	Comments []interface{} `json:"comments"`
}
