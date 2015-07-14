package model

const (
	USER_ROLE_ADMIN  = "admin"
	USER_ROLE_WRITER = "writer"
	USER_ROLE_READER = "reader"

	USER_STATUS_ACTIVE = "active"
	USER_STATUS_BLOCK  = "block"
	USER_STATUS_DELETE = "delete"
)

type User struct {
	Id       int64
	Name     string `xorm:"unique"`
	Password string
	Salt     string

	Email      string `xorm:"unique"`
	Url        string
	AvatarUrl  string
	Profile    string
	CreateTime int64 `xorm:"created"`

	Role   string `xorm:"index(role)"`
	Status string
}
