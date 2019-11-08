package model

import (
	"sync"
	"time"
)

type BaseModel struct {
	Id        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

type UserInfo struct {
	Id        uint64 `json:"id" binding:"required"` //用户id
	Username  string `json:"username" example:"小明" format:"string" binding:"required"` //用户名
	SayHello  string `json:"sayHello" binding:"required"` //测试
	Password  string `json:"password" binding:"required"` //密码
	CreatedAt string `json:"createdAt" binding:"required"` //创建时间
	UpdatedAt string `json:"updatedAt" binding:"required"` //更新时间
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}
