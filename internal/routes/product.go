package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/models"
)

func AddProduct(c *gin.Context) {
	var p = models.Product{
		Sku:         "k94940505",
		Description: "TestDescription",
		Vendor:      "SomeVendor",
		Stock:       100,
		Price:       101,
	}

	db.SQL.Query("INSERT INTO product (sku, description, vendor, stock, price) VALUES (?, ?, ?, ?, ?)",
		p.Sku, p.Description, p.Vendor, p.Stock, p.Price)
	c.String(200, "This should appear")
}

func GetProducts(c *gin.Context) {
	// TODO: https://go.dev/doc/database/querying
	var products []models.Product

	rows, err := db.SQL.Query("SELECT * FROM product")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Sku, &p.Description, &p.Vendor, &p.Stock, &p.Price, &p.TimesClickedUpdate, &p.IdStore, &p.LastUpdate, &p.FirstUpdate, &p.NumUpdates)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)
	}

	c.JSON(200, products)
}
