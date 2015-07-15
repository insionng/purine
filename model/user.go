package model

import (
	"fmt"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/utils"
	"github.com/fuxiaohei/purine/vars"
	"time"
)

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

func (u *User) CheckPassword(pwd string) bool {
	return u.Password == utils.Sha256String(pwd+u.Salt)
}

func GetUserBy(col string, v interface{}) (*User, error) {
	u := new(User)
	if _, err := vars.Db.Where(col+" = ?", v).Get(u); err != nil {
		log.Error("Db|GetUserBy|%s,%v|%s", col, v, err.Error())
		return nil, err
	}
	if u.Id == 0 {
		return nil, nil
	}
	return u, nil
}

type Token struct {
	Id         int64
	UserId     int64
	Token      string `xorm:"unique"`
	CreateTime int64  `xorm:"created"`
	ExpireTime int64
}

func CreateToken(user, expire int64) (*Token, error) {
	t := &Token{
		UserId:     user,
		ExpireTime: time.Now().Unix() + expire,
	}
	t.Token = utils.Md5String(fmt.Sprintf("%d,%d", t.UserId, t.ExpireTime))
	if _, err := vars.Db.Insert(t); err != nil {
		log.Error("Db|CreateToken|%v|%s", t, err.Error())
		return nil, err
	}
	return t, nil
}

func GetValidToken(token string) (*Token, error) {
	t := new(Token)
	if _, err := vars.Db.Where("token = ?", token).Get(t); err != nil {
		log.Error("Db|GetValidToken|%s|%s", token, err.Error())
		return nil, err
	}
	// wrong token
	if t.Token != token {
		return nil, nil
	}
	// expired
	if time.Now().Unix() > t.ExpireTime {
		return nil, nil
	}
	return t, nil
}
