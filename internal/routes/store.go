package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/models"
	"log"
	"strconv"
	"strings"
)

type qParam struct {
	Name  string
	Value string
}

func GetStores(c *gin.Context) {
	var stores []models.Store
	var params []qParam
	var args []interface{}
	var qWhere string

	storeName := c.Query("name")
	storeCountry := c.Query("country")
	storeRegion := c.Query("region")

	if storeName != "" {
		params = append(params, qParam{"'name'", storeName})
	}
	if storeCountry != "" {
		params = append(params, qParam{"country", storeCountry})
	}
	if storeRegion != "" {
		params = append(params, qParam{"region", storeRegion})
	}

	if len(params) > 0 {
		args, qWhere = obtainQueryArgs(params)
	}

	rows, err := db.SQL.Query("SELECT * FROM store "+qWhere, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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

func obtainQueryArgs(params []qParam) ([]interface{}, string) {
	var args []interface{}
	var strs []string
	for i, p := range params {
		args = append(args, p.Value)
		strs = append(strs, p.Name+" = $"+strconv.Itoa(i+1))
	}

	qWhere := " WHERE " + strings.Join(strs, " AND ")

	return args, qWhere
}

func GetStore(c *gin.Context) {
	var s models.Store
	id := c.Param("id")

	err := db.SQL.QueryRow("SELECT * FROM store WHERE id = $1", id).Scan(&s.ID, &s.Name, &s.Url, &s.Country, &s.Region, &s.BadPingCount)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, s)
}

func AddStore(c *gin.Context) {
	// TODO: consider if name should be unique
	var s models.Store
	var newStore models.Store
	var params []qParam
	c.BindJSON(&s)

	if s.Name == "" || s.Url == "" {
		c.JSON(400, gin.H{"error": "Name and URL are required"})
		return
	}

	params = append(params, qParam{"name", s.Name})
	params = append(params, qParam{"url", s.Url})

	if s.Country != "" {
		params = append(params, qParam{"country", s.Country})
	}
	if s.Region != "" {
		params = append(params, qParam{"region", s.Region})
	}

	cols, nums, args := obtainInsertArgs(params)

	query := "INSERT INTO store (" + cols + ") VALUES (" + nums + ") RETURNING *"

	err := db.SQL.QueryRow(query, args...).Scan(&newStore.ID, &newStore.Name, &newStore.Url, &newStore.Country,
		&newStore.Region, &newStore.BadPingCount)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, newStore)
}

func obtainInsertArgs(params []qParam) (string, string, []interface{}) {
	var args []interface{}
	var str_c []string
	var str_n []string

	for i, p := range params {
		str_c = append(str_c, p.Name)
		str_n = append(str_n, "$"+strconv.Itoa(i+1))
		args = append(args, p.Value)
	}

	cols := strings.Join(str_c, ", ")
	nums := strings.Join(str_n, ", ")

	return cols, nums, args
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
	c.BindJSON(&s)

	if s.Name != "" {
		params = append(params, qParam{"name", s.Name})
	}
	if s.Url != "" {
		params = append(params, qParam{"url", s.Url})
	}
	if s.Country != "" {
		params = append(params, qParam{"country", s.Country})
	}
	if s.Region != "" {
		params = append(params, qParam{"region", s.Region})
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
		log.Fatal(err)
	}

	c.JSON(200, updStore)

}

func obtainUpdateArgs(params []qParam) (string, []interface{}) {
	var args []interface{}
	var str_c []string

	for i, p := range params {
		str_c = append(str_c, p.Name+" = $"+strconv.Itoa(i+1))
		args = append(args, p.Value)
	}

	cols := strings.Join(str_c, ", ")

	return cols, args
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
