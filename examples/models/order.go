//go:generate jia -format=go -file $GOFILE -out ../routes/  -tpl ../route.tmpl
package models

type Order struct {
	Id    int
	Title string
	Price float64
	User  User
}

func FindOrderById(userid, orderid int) (*Order, error) {
	// ...
	return &Order{}, nil
}

func CreateOrder(userid int, title string, price float64) (*Order, error) {
	// ...
	return &Order{}, nil
}
