package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

const numRecords = 5000

func connectTestDB(b *testing.B) *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func BenchmarkSelectSingleWithQuery(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()
	email := "user100@example.com"

	b.ResetTimer()
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

func BenchmarkSelectVariousWithQuery(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// email を変化させる
		email := fmt.Sprintf("user%d@example.com", i%numRecords)

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

func BenchmarkSelectSingleWithPrepared(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()

	stmt, err := db.Prepare(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`)
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	email := "user100@example.com"
	b.ResetTimer()
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

func BenchmarkSelectVariousWithPrepared(b *testing.B) {
	db := connectTestDB(b)
	defer db.Close()

	stmt, err := db.Prepare(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`)
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// email を変化させる
		email := fmt.Sprintf("user%d@example.com", i%numRecords)

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
