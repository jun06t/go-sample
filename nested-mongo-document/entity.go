package main

type User struct {
	ID      string  `bson:"_id"`
	Name    string  `bson:"name,omitempty"`
	Age     int     `bson:"age,omitempty"`
	Address Address `bson:"address,omitempty"`
}

type Address struct {
	Country string `bson:"country,omitempty"`
	State   string `bson:"state,omitempty"`
	City    string `bson:"city,omitempty"`
	Zipcode string `bson:"zipcode,omitempty"`
}

type UserWithInline struct {
	ID      string  `bson:"_id"`
	Name    string  `bson:"name,omitempty"`
	Age     int     `bson:"age,omitempty"`
	Address Address `bson:"address,inline,omitempty"`
}
