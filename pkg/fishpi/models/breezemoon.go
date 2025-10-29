package models

// Breezemoon 清风明月数据结构
type Breezemoon struct {
	OID                            string `json:"oId"`
	BreezemoonAuthorName           string `json:"breezemoonAuthorName"`
	BreezemoonAuthorThumbnailURL48 string `json:"breezemoonAuthorThumbnailURL48"`
	BreezemoonContent              string `json:"breezemoonContent"`
	BreezemoonCity                 string `json:"breezemoonCity"`
	BreezemoonCreated              int64  `json:"breezemoonCreated"`
	BreezemoonUpdated              int64  `json:"breezemoonUpdated"`
	BreezemoonCreateTime           string `json:"breezemoonCreateTime"`
	TimeAgo                        string `json:"timeAgo"`
}

// BreezemoonListResponse 清风明月列表响应
type BreezemoonListResponse struct {
	Code        int          `json:"code"`
	Msg         string       `json:"msg"`
	Breezemoons []Breezemoon `json:"breezemoons"`
}

// BreezemoonPostResponse 发布清风明月响应
type BreezemoonPostResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
