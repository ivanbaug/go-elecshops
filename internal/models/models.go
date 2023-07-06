package models

import "time"

type Product struct {
	ID                 int64     `json:"id"`
	Sku                string    `json:"sku"`
	Description        string    `json:"description"`
	Vendor             string    `json:"vendor"`
	Stock              int64     `json:"stock"`
	Price              float64   `json:"price"`
	TimesClickedUpdate int64     `json:"times_clicked_update"`
	IdStore            int64     `json:"id_store"`
	LastUpdate         time.Time `json:"last_update"`
	FirstUpdate        time.Time `json:"first_update"`
	NumUpdates         int64     `json:"num_updates"`
	Url                string    `json:"url"`
}

type Store struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	Country      string    `json:"country"`
	Region       string    `json:"region"`
	BadPingCount int64     `json:"bad_ping_count"`
	LastUpdate   time.Time `json:"last_update"`
}

type VwProduct struct {
	ID                 int64     `json:"id"`
	Sku                string    `json:"sku"`
	Description        string    `json:"description"`
	Vendor             string    `json:"vendor"`
	Stock              int64     `json:"stock"`
	Price              float64   `json:"price"`
	TimesClickedUpdate int64     `json:"times_clicked_update"`
	IdStore            int64     `json:"id_store"`
	LastUpdate         time.Time `json:"last_update"`
	FirstUpdate        time.Time `json:"first_update"`
	NumUpdates         int64     `json:"num_updates"`
	Url                string    `json:"url"`
	StoreName          string    `json:"store_name"`
	Country            string    `json:"country"`
}
