package main

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

func getProducts() []Product {
	return []Product{
		{Id: 1, Name: "Americano", Price: 100, Description: "Espresso, diluted for a lighter experience"},
		{Id: 2, Name: "Cappuccino", Price: 110, Description: "Espresso with steamed milk"},
		{Id: 3, Name: "Espresso", Price: 90, Description: "A strong shot of coffee"},
		{Id: 4, Name: "Macchiato", Price: 120, Description: "Espresso topped with a thin layer of foamed milk"},
	}
}

func getUsers() []User {
	return []User{
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
}

func getSessions() []Session {
	return sessions
}

// In database.go

func setSession(token string, user User) {
	sessions = append(sessions, Session{Token: token, UserId: user.Id})
}

func getUserFromSessionToken(token string) User {
	var userId int
	for _, session := range sessions {
		if session.Token == token {
			userId = session.UserId
		}
	}
	for _, u := range getUsers() {
		if u.Id == userId {
			return u
		}
	}
	return User{}
}
