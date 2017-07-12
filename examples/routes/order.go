// generate by jia
package routes

import (
	"net/http"

	"github.com/relaxgo/tangram/param"
	"somepkg/models"
)

func FindOrderById(w http.ResponseWriter, r *http.Request) {
	p := NewRequestValue(r)

	userid := param.Int(p, "userid")
	orderid := param.Int(p, "orderid")

	v, err := models.FindOrderById(userid, orderid)
	Respond(w, r, v, err)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	p := NewRequestValue(r)

	userid := param.Int(p, "userid")
	title := param.String(p, "title")
	price := param.Float64(p, "price")

	v, err := models.CreateOrder(userid, title, price)
	Respond(w, r, v, err)
}
