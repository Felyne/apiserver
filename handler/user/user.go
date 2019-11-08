package user

import (
	"apiserver/model"
)

// POST 请求的body，swap使用演示
type CreateRequest struct {
	Username string `json:"username" example:"小明" format:"string" binding:"required" minLength:"1" maxLength:"16"` // 用户名
	Password string `json:"password" example:"123456" format:"string" binding:"required" default:"123456"` // 密码
	//Age int `json:"age" example:"10" binding:"required" minimum:"1" maximum:"150" default:"10"`   // 年龄
	//Height float64 `json:"height" example:"1.80" binding:"required" minimum:"0.0" maximum:"9.99"` // 身高,单位米
	//Status string `json:"status" enums:"healthy,ill"` // 状态
}

type CreateResponse struct {
	Username string `json:"username" example:"小明" format:"string" binding:"required"` // 用户名
}

// GET请求的body,swap不起作用?
type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset" example:"0" format:"int64" binding:"required" minimum:"0" maximum:"1000"`
	Limit    int    `json:"limit" example:"10" format:"int64" minimum:"1" maximum:"1000"`
}

type ListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []*model.UserInfo `json:"userList" `
}

type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount" format:"int64" validate:"required"`
	UserList   []model.UserInfo `json:"userList" binding:"required"`
}
