package model

import (
	"fmt"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/utils"
	"github.com/fuxiaohei/purine/src/vars"
	"time"
)

const (
	USER_ROLE_ADMIN  = "admin"
	USER_ROLE_WRITER = "writer"
	USER_ROLE_READER = "reader"

	USER_STATUS_ACTIVE = "active"
	USER_STATUS_BLOCK  = "block"
	USER_STATUS_DELETE = "deleted"
)

// user struct
type User struct {
	Id       int64
	Name     string `xorm:"unique"`
	Password string
	Salt     string

	Nick       string
	Email      string `xorm:"unique"`
	Url        string
	AvatarUrl  string
	Profile    string
	CreateTime int64 `xorm:"created"`

	Role   string `xorm:"index(role)"`
	Status string
}

// check user's password
func (u *User) CheckPassword(pwd string) bool {
	return u.Password == utils.Sha256String(pwd+u.Salt)
}

// get user by column and value
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

// update user profile columns
func UpdateUser(u *User) error {
	if _, err := vars.Db.Cols("name,nick,email,url,profile,avatar_url").
		Where("id = ?", u.Id).Update(u); err != nil {
		log.Error("Db|UpdateUser|%d|%s", u.Id, err.Error())
		return err
	}
	return nil
}

// update user password with user id and new password
func UpdatePassword(id int64, newPassword string) error {
	u := new(User)
	u.Salt = utils.Md5String(newPassword)[8:24]
	u.Password = utils.Sha256String(newPassword + u.Salt)
	if _, err := vars.Db.Cols("password,salt").Where("id = ?", id).Update(u); err != nil {
		log.Error("Db|UpdatePassword|%d|%s", id, err.Error())
		return err
	}
	return nil
}

// token struct
type Token struct {
	Id         int64
	UserId     int64
	Token      string `xorm:"unique"`
	CreateTime int64  `xorm:"created"`
	ExpireTime int64
}

// create new token with user id and expiration duration
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

// get valid token.
// it checks expiration.
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
