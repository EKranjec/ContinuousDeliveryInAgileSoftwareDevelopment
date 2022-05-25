package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type order struct {
	ID        int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"qty"`
}

func getOrders(db *sql.DB, orderid int) ([]order, error) {
	rows, err := db.Query("SELECT product_id, qty FROM orders WHERE order_id=$1", orderid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []order{}

	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ProductID, &o.Quantity); err != nil {
			fmt.Printf("Err: %s", err)
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (o *order) createOrder(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO orders(product_id, qty) VALUES($1, $2) RETURNING order_id",
		o.ProductID, o.Quantity).Scan(&o.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)
	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
