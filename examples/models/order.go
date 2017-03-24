//go:generate jia -format -f $GOFILE -o ../routes/  -t ../route.tmpl
package models

type Order struct {
	Id    int
	Title string
	Price float64
}

func FindOrderById(userid, orderid int) (*Order, error) {
	// ...
	return &Order{}, nil
}

func CreateOrder(userid int, title string, price float64) (*Order, error) {
	// ...
	return &Order{}, nil
}
