package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/models"
	"strconv"
	"strings"
)

func GetProducts(c *gin.Context) {
	var products []models.Product
	var params []qParam
	var args []interface{}
	var qWhere string

	productSku := c.Query("sku")
	productDescription := c.Query("description")
	productVendor := c.Query("vendor")
	productUrl := c.Query("url")
	productIdStore := c.Query("id_store")

	if productSku != "" {
		p := newQParam("sku", productSku)
		p.Precise = false
		params = append(params, p)
	}
	if productDescription != "" {
		p := newQParam("description", productDescription)
		p.Precise = false
		params = append(params, p)
	}
	if productVendor != "" {
		p := newQParam("vendor", productVendor)
		p.Precise = false
		params = append(params, p)
	}
	if productUrl != "" {
		p := newQParam("url", productUrl)
		params = append(params, p)
	}
	if productIdStore != "" {
		p := newQParam("id_store", productIdStore)
		params = append(params, p)
	}

	if len(params) > 0 {
		args, qWhere = obtainQueryArgs(params)
	}

	rows, err := db.SQL.Query("SELECT * FROM product "+qWhere, args...)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
	}(rows)

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Sku, &p.Description, &p.Vendor, &p.Stock, &p.Price, &p.TimesClickedUpdate,
			&p.IdStore, &p.LastUpdate, &p.FirstUpdate, &p.NumUpdates, &p.Url)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error(), "row_id": p.ID})
			return
		}
		products = append(products, p)
	}

	c.JSON(200, products)
}

func GetProduct(c *gin.Context) {
	var p models.Product
	id := c.Param("id")

	err := db.SQL.QueryRow("SELECT * FROM product WHERE id = $1", id).Scan(&p.ID, &p.Sku, &p.Description,
		&p.Vendor, &p.Stock, &p.Price, &p.TimesClickedUpdate, &p.IdStore, &p.LastUpdate, &p.FirstUpdate,
		&p.NumUpdates, &p.Url)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, p)
}

func AddProduct(c *gin.Context) {
	var p models.Product
	var newProduct models.Product
	var params []qParam
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if p.Sku == "" || p.IdStore <= 0 {
		c.JSON(400, gin.H{"error": "sku and id_store are required"})
		return
	}

	params = append(params, newQParam("sku", p.Sku))
	params = append(params, newQParam("id_store", strconv.FormatInt(p.IdStore, 10)))

	if p.Description != "" {
		params = append(params, newQParam("description", p.Description))
	}
	if p.Vendor != "" {
		params = append(params, newQParam("vendor", p.Vendor))
	}
	if p.Stock > 0 {
		params = append(params, newQParam("stock", strconv.FormatInt(p.Stock, 10)))
	}
	if p.Price > 0 {
		params = append(params, newQParam("price", strconv.FormatFloat(p.Price, 'f', 2, 64)))
	}
	if p.Url != "" {
		params = append(params, newQParam("url", p.Url))
	}

	cols, nums, args := obtainInsertArgs(params)

	query := "INSERT INTO product (" + cols + ") VALUES (" + nums + ") RETURNING *"

	err = db.SQL.QueryRow(query, args...).Scan(&newProduct.ID, &newProduct.Sku, &newProduct.Description,
		&newProduct.Vendor, &newProduct.Stock, &newProduct.Price, &newProduct.TimesClickedUpdate,
		&newProduct.IdStore, &newProduct.LastUpdate, &newProduct.FirstUpdate, &newProduct.NumUpdates, &newProduct.Url)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, newProduct)
}

func UpsertProduct(c *gin.Context) {
	var p models.Product
	var newProduct models.Product
	var params []qParam
	var upsrtParams []string
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if p.Sku == "" || p.IdStore <= 0 {
		c.JSON(400, gin.H{"error": "sku and id_store are required"})
		return
	}

	params = append(params, newQParam("sku", p.Sku))
	params = append(params, newQParam("id_store", strconv.FormatInt(p.IdStore, 10)))

	if p.Description != "" {
		params = append(params, newQParam("description", p.Description))
		upsrtParams = append(upsrtParams, "description = EXCLUDED.description")

	}
	if p.Vendor != "" {
		params = append(params, newQParam("vendor", p.Vendor))
		upsrtParams = append(upsrtParams, "vendor = EXCLUDED.vendor")
	}
	if p.Stock > 0 {
		params = append(params, newQParam("stock", strconv.FormatInt(p.Stock, 10)))
		upsrtParams = append(upsrtParams, "stock = EXCLUDED.stock")
	}
	if p.Price > 0 {
		params = append(params, newQParam("price", strconv.FormatFloat(p.Price, 'f', 2, 64)))
		upsrtParams = append(upsrtParams, "price = EXCLUDED.price")
	}
	if p.Url != "" {
		params = append(params, newQParam("url", p.Url))
		upsrtParams = append(upsrtParams, "url = EXCLUDED.url")
	}

	cols, nums, args := obtainInsertArgs(params)
	updStr := strings.Join(upsrtParams, ", ")

	query := "INSERT INTO product (" + cols + ") VALUES (" + nums + ") RETURNING *"

	if len(upsrtParams) > 0 {
		query = "INSERT INTO product (" + cols + ") VALUES (" + nums + ")" +
			" ON CONFLICT (sku,id_store) DO UPDATE SET " + updStr +
			", last_update = CURRENT_TIMESTAMP, num_updates = product.num_updates + 1 " +
			" RETURNING *"
	}

	err = db.SQL.QueryRow(query, args...).Scan(&newProduct.ID, &newProduct.Sku, &newProduct.Description,
		&newProduct.Vendor, &newProduct.Stock, &newProduct.Price, &newProduct.TimesClickedUpdate,
		&newProduct.IdStore, &newProduct.LastUpdate, &newProduct.FirstUpdate, &newProduct.NumUpdates, &newProduct.Url)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, newProduct)
}

func UpdateProduct(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(400, gin.H{"error": "Received ID is not a number"})
		return
	}
	if id <= 0 {
		c.JSON(400, gin.H{"error": "Wrong ID"})
		return
	}

	var p models.Product
	var updProduct models.Product
	var params []qParam
	err = c.BindJSON(&p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if p.Sku != "" {
		params = append(params, newQParam("sku", p.Sku))
	}
	if p.Description != "" {
		params = append(params, newQParam("description", p.Description))
	}
	if p.Vendor != "" {
		params = append(params, newQParam("vendor", p.Vendor))
	}
	if p.Stock >= 0 {
		params = append(params, newQParam("stock", strconv.FormatInt(p.Stock, 10)))
	}
	if p.Price > 0 {
		params = append(params, newQParam("price", strconv.FormatFloat(p.Price, 'f', 2, 64)))
	}
	if p.Url != "" {
		params = append(params, newQParam("url", p.Url))
	}

	if len(params) == 0 {
		c.JSON(400, gin.H{"error": "Nothing to update"})
		return
	}

	cols, args := obtainUpdateArgs(params)
	args = append(args, id)

	query := "UPDATE product SET " + cols + ", last_update=NOW() WHERE id = $" + strconv.Itoa(len(args)) +
		" RETURNING *"

	err = db.SQL.QueryRow(query, args...).Scan(&updProduct.ID, &updProduct.Sku, &updProduct.Description,
		&updProduct.Vendor, &updProduct.Stock, &updProduct.Price, &updProduct.TimesClickedUpdate,
		&updProduct.IdStore, &updProduct.LastUpdate, &updProduct.FirstUpdate, &updProduct.NumUpdates, &updProduct.Url)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, updProduct)
}

func DeleteProduct(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(400, gin.H{"error": "Received ID is not a number"})
		return
	}
	if id <= 0 {
		c.JSON(400, gin.H{"error": "Wrong ID"})
		return
	}

	_, err = db.SQL.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product with ID: " + sid + " deleted"})
}
