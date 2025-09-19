package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func connectTestDB(b *testing.B) *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func BenchmarkSelectWithQuery(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()

	email := "user100@example.com"
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`, email)
		if err != nil {
			b.Fatal(err)
		}
		var id int
		var name string
		for rows.Next() {
			rows.Scan(&id, &name)
		}
		rows.Close()
	}
}

func BenchmarkSelectWithPrepared(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()

	stmt, err := db.Prepare(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`)
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	email := "user100@example.com"
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query(email)
		if err != nil {
			b.Fatal(err)
		}
		var id int
		var name string
		for rows.Next() {
			rows.Scan(&id, &name)
		}
		rows.Close()
	}
}
