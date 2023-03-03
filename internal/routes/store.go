package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/models"
	"strconv"
)

func GetStores(c *gin.Context) {
	var stores []models.Store
	var params []qParam
	var args []interface{}
	var qWhere string

	storeName := c.Query("name")
	storeCountry := c.Query("country")
	storeRegion := c.Query("region")

	if storeName != "" {
		p := newQParam("name", storeName)
		p.Precise = false
		params = append(params, p)
	}
	if storeCountry != "" {
		p := newQParam("country", storeCountry)
		p.Precise = false
		params = append(params, p)
	}
	if storeRegion != "" {
		p := newQParam("region", storeRegion)
		p.Precise = false
		params = append(params, p)
	}

	if len(params) > 0 {
		args, qWhere = obtainQueryArgs(params)
	}

	rows, err := db.SQL.Query("SELECT * FROM store "+qWhere, args...)
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
		var s models.Store
		err := rows.Scan(&s.ID, &s.Name, &s.Url, &s.Country, &s.Region, &s.BadPingCount)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error(), "row_id": s.ID})
			return
		}
		stores = append(stores, s)
	}

	c.JSON(200, stores)
}

func GetStore(c *gin.Context) {
	var s models.Store
	id := c.Param("id")

	err := db.SQL.QueryRow("SELECT * FROM store WHERE id = $1", id).Scan(&s.ID, &s.Name, &s.Url, &s.Country, &s.Region, &s.BadPingCount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, s)
}

func AddStore(c *gin.Context) {
	// TODO: consider if name should be unique
	var s models.Store
	var newStore models.Store
	var params []qParam
	err := c.BindJSON(&s)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if s.Name == "" || s.Url == "" {
		c.JSON(400, gin.H{"error": "Name and URL are required"})
		return
	}

	params = append(params, newQParam("name", s.Name))
	params = append(params, newQParam("url", s.Url))

	if s.Country != "" {
		params = append(params, newQParam("country", s.Country))
	}
	if s.Region != "" {
		params = append(params, newQParam("region", s.Region))
	}

	cols, nums, args := obtainInsertArgs(params)

	query := "INSERT INTO store (" + cols + ") VALUES (" + nums + ") RETURNING *"

	err = db.SQL.QueryRow(query, args...).Scan(&newStore.ID, &newStore.Name, &newStore.Url, &newStore.Country,
		&newStore.Region, &newStore.BadPingCount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, newStore)
}

func UpdateStore(c *gin.Context) {
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

	var s models.Store
	var updStore models.Store
	var params []qParam
	err = c.BindJSON(&s)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if s.Name != "" {
		params = append(params, newQParam("name", s.Name))
	}
	if s.Url != "" {
		params = append(params, newQParam("url", s.Url))
	}
	if s.Country != "" {
		params = append(params, newQParam("country", s.Country))
	}
	if s.Region != "" {
		params = append(params, newQParam("region", s.Region))
	}

	if len(params) == 0 {
		c.JSON(400, gin.H{"error": "Nothing to update"})
		return
	}

	cols, args := obtainUpdateArgs(params)
	args = append(args, id)

	query := "UPDATE store SET " + cols + " WHERE id = $" + strconv.Itoa(len(params)+1) + " RETURNING *"

	err = db.SQL.QueryRow(query, args...).Scan(&updStore.ID, &updStore.Name, &updStore.Url, &updStore.Country,
		&updStore.Region, &updStore.BadPingCount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, updStore)

}

func DeleteStore(c *gin.Context) {
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
	_, err = db.SQL.Exec("DELETE FROM store WHERE id = $1", id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Store with ID:" + sid + " deleted"})
}
