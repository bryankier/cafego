// NOTE: You may need to import go-sqlite3 explicitly.
// Go won't know how to import this on its own
package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func initDB() {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	database = db
	queries := []string{
		"CREATE TABLE IF NOT EXISTS cgo_user (username TEXT, password TEXT)",
		"CREATE TABLE IF NOT EXISTS cgo_product (name TEXT, price INTEGER, description TEXT)",
		"CREATE TABLE IF NOT EXISTS cgo_session (token TEXT, user_id INTEGER)",
		"CREATE TABLE IF NOT EXISTS cgo_cart_item (product_id INTEGER, quantity INTEGER, user_id INTEGER)",
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Seed data
	var q string
	var count int
	// cgo_user
	q = "SELECT COUNT(*) FROM cgo_user"
	err = db.QueryRow(q).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		q = "INSERT INTO cgo_user (username, password) VALUES (?, ?)"
		for _, u := range seedUsers {
			_, err = db.Exec(q, u.Username, u.Password)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	// cgo_product
	q = "SELECT COUNT(*) FROM cgo_product"
	err = db.QueryRow(q).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		q = "INSERT INTO cgo_product (name, price, description) VALUES (?, ?, ?)"
		for _, p := range seedProducts {
			_, err = db.Exec(q, p.Name, p.Price, p.Description)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

type Product struct {
	Id          int
	Name        string
	Price       int
	Description string
}

type User struct {
	Id       int
	Username string
	Password string
}

type Session struct {
	Token  string
	UserId int
}

var sessions = []Session{}

var seedUsers = []User{
	{
		Id:       1,
		Username: "zagreus",
		Password: "cerberus",
	},
	{
		Id:       2,
		Username: "melinoe",
		Password: "b4d3ec1",
	},
}
var seedProducts = []Product{
	{Id: 1, Name: "Americano", Price: 100, Description: "Espresso, diluted for a lighter experience"},
	{Id: 2, Name: "Cappuccino", Price: 110, Description: "Espresso with steamed milk"},
	{Id: 3, Name: "Espresso", Price: 90, Description: "A strong shot of coffee"},
	{Id: 4, Name: "Macchiato", Price: 120, Description: "Espresso with a small amount of milk"},
}

func getProducts() []Product {
	var result []Product
	q := "SELECT rowid, name, price, description FROM cgo_product"
	rows, err := database.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, product)
	}
	return result
}

func getUsers() []User {
	var result []User
	q := "SELECT rowid, username, password FROM cgo_user"
	rows, err := database.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, user)
	}
	return result
}

func getSessions() []Session {
	return sessions
}

func setSession(token string, user User) {
	q := "INSERT INTO cgo_session (token, user_id) VALUES (?, ?)"
	_, err := database.Exec(q, token, user.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserFromSessionToken(token string) User {
	q := `
	SELECT
		cgo_session.user_id,
		cgo_user.username,
		cgo_user.password
	FROM cgo_session
	INNER JOIN cgo_user
	ON cgo_session.user_id = cgo_user.rowid
	WHERE cgo_session.token = ?
	LIMIT 1;
	`
	var user User
	err := database.QueryRow(q, token).Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return User{}
	} else if err != nil {
		log.Fatal(err)
	}
	return user
}
func createCartItem(userId int, productId int, quantity int) {
	q := "INSERT INTO cgo_cart_item (user_id, product_id, quantity) VALUES (?, ?, ?)"
	_, err := database.Exec(q, userId, productId, quantity)
	if err != nil {
		log.Fatal(err)
	}
}
