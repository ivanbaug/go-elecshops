package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/dbdriver"
	"strconv"
	"strings"
)

var db *dbdriver.DB

type qParam struct {
	Name    string
	Value   string
	Precise bool
}

func SetupRouter(rdb *dbdriver.DB) *gin.Engine {
	// Set db variable to the one passed in
	db = rdb

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// stores
	r.GET("/get-stores", GetStores)
	r.GET("/get-store/:id", GetStore)
	r.POST("/add-store", AddStore)
	r.PATCH("/update-store/:id", UpdateStore)
	r.DELETE("/delete-store/:id", DeleteStore)

	// products
	r.GET("/get-products", GetProducts)
	r.GET("/get-product/:id", GetProduct)
	r.POST("/add-product", AddProduct)
	r.PATCH("/update-product/:id", UpdateProduct)
	r.DELETE("/delete-product/:id", DeleteProduct)

	// TODO: Add routes by grouping their functions
	// TODO: Authorization for protected routes

	return r
}

// Helper functions, types

func newQParam(name string, value string) qParam {
	return qParam{name, value, true}
}

func obtainQueryArgs(params []qParam) ([]interface{}, string) {
	var args []interface{}
	var strs []string
	for i, p := range params {
		args = append(args, p.Value)
		if p.Precise {
			strs = append(strs, p.Name+" = $"+strconv.Itoa(i+1))
		} else {
			strs = append(strs, p.Name+" LIKE '%' || $"+strconv.Itoa(i+1)+" || '%'")
		}
	}

	qWhere := " WHERE " + strings.Join(strs, " AND ")

	return args, qWhere
}

func obtainInsertArgs(params []qParam) (string, string, []interface{}) {
	var args []interface{}
	var strC []string
	var strN []string

	for i, p := range params {
		strC = append(strC, p.Name)
		strN = append(strN, "$"+strconv.Itoa(i+1))
		args = append(args, p.Value)
	}

	cols := strings.Join(strC, ", ")
	nums := strings.Join(strN, ", ")

	return cols, nums, args
}

func obtainUpdateArgs(params []qParam) (string, []interface{}) {
	var args []interface{}
	var strC []string

	for i, p := range params {
		strC = append(strC, p.Name+" = $"+strconv.Itoa(i+1))
		args = append(args, p.Value)
	}

	cols := strings.Join(strC, ", ")

	return cols, args
}
