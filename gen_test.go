package main

import (
	"fmt"
	"strings"
	"testing"
)

var (
	src string = `
package subtitle

import (
	"context"
	"github.com/relaxgo/subtitle-server/db"
)

// some des
func FindProjects(c context.Context, where, sorts string, skip int, limit int) ([]Project, error) {
	list := make([]Project, 0)

	query := db.DB.
		Preload("Ower").
		Offset(skip).
		Limit(limit).
		Find(&list)

	return list, query.Error
}

/*
* param
*/
func AddProject(c context.Context, project *Project) (*Project, error) {
	query := db.DB.Create(project)
	return project, query.Error
}
	`
	src2 = `
package models

func Login(phone, password string) (*Token, error) {
	user := &User{}
	if db.DB == nil {
		return nil, errs.CMDBFaild
	}
	query := db.DB.Where("username = ?", phone).First(user)
	if query.Error != nil {
		return nil, errs.LoginErrorPhonePassword
	}
	md5Pw := fmt.Sprintf("%X", md5.Sum([]byte(password)))
	fmt.Println(user.Password, md5Pw)
	if user.Password != md5Pw {
		return nil, errs.LoginErrorPhonePassword
	}
	return &Token{"code", user.Id}, nil
}

	`
)

func TestGenRouter(t *testing.T) {
	r := strings.NewReader(src2)
	data := Gen(r)
	fmt.Println(string(data))
}
