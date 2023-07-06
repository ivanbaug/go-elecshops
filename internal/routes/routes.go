package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-elecshops/internal/dbdriver"
	"os"
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

	// Stores
	rStores := r.Group("/stores")
	{
		rStores.Use(CorsMiddlewareGet())
		rStores.GET("/", GetStores)
		rStores.GET("/:id", GetStore)
	}

	// products
	rProducts := r.Group("/products")
	{
		rProducts.Use(CorsMiddlewareGet())
		rProducts.GET("/", GetProductSort)
		rProducts.GET("/:id", GetProduct)
	}

	// protected
	rProtected := r.Group("/p")
	rProtected.Use(simpleAuthMiddleware())
	{
		pgStores := rProtected.Group("/stores")
		{
			pgStores.POST("/", AddStore)
			pgStores.PATCH("/:id", UpdateStore)
			pgStores.DELETE("/:id", DeleteStore)
		}
		pgProducts := rProtected.Group("/products")
		{
			pgProducts.POST("/", AddProduct)
			pgProducts.POST("/upsert/", UpsertProduct)
			pgProducts.PATCH("/:id", UpdateProduct)
			pgProducts.DELETE("/:id", DeleteProduct)
		}
	}

	return r
}

func CorsMiddlewareGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		//c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Next()
	}
}

// Helper functions, types

func newQParam(name string, value string) qParam {
	return qParam{name, value, true}
}

func obtainQueryArgs(params []qParam) ([]interface{}, string) {
	var args []interface{}
	var strs []string
	for i, p := range params {
		args = append(args, strings.ToLower(p.Value))
		if p.Precise {
			strs = append(strs, " LOWER("+p.Name+") = $"+strconv.Itoa(i+1))
		} else {
			strs = append(strs, " LOWER("+p.Name+") LIKE '%' || $"+strconv.Itoa(i+1)+" || '%'")
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

// unsafe, probably
func simpleTokenValidator(c *gin.Context) bool {
	bearerToken := c.Request.Header.Get("Authorization")
	tkn := ""
	if len(strings.Split(bearerToken, " ")) == 2 {
		tkn = strings.Split(bearerToken, " ")[1]
	}
	validTkn := os.Getenv("API_KEY")
	if tkn == validTkn {
		return true
	}
	return false
}

func simpleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !simpleTokenValidator(c) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
