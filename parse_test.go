package jia

import (
	"strings"
	"testing"
)

const testFile = `
package main

import (
	"time"
	"somepkg/models"
)

type Phone string

type User struct {
	Id    int
	Name  string
	Phone Phone
}

type Customer struct {
	User
	Id        int
	createdAt time.Time
	Order     []models.Device
}

func (c *Customer) SayHello() string {
	return ""
}

func removeCustomer(customerId int) error {
	// DO some
	return nil
}

func FindCustomer(userId, customerId int) (*Customer, error) {
	return nil, nil
}

func CreateOrder(customerId int,  order models.Device) (*models.Device, error) {
	return nil, nil
}
`

func TestParse(t *testing.T) {
	r := strings.NewReader(testFile)
	f, err := Parse("user.go", r)
	if err != nil {
		t.Error(err)
	}
	for _, fc := range f.Funcs {
		for _, p := range fc.Params {
			t.Log(p)
			t.Logf("%s", p.Type)
			t.Log(p.Type)
			t.Log(p.TypeKind())
		}
	}
}
