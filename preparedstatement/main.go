package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// DB接続
func connectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	must(err)
	return db
}

// データ準備
func setupTable(db *sql.DB, n int) {
	_, _ = db.Exec(`DROP TABLE IF EXISTS users`)
	_, _ = db.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT, email TEXT, status TEXT)`)

	stmt, _ := db.Prepare(`INSERT INTO users (name, email, status) VALUES ($1, $2, $3)`)
	defer stmt.Close()

	for i := 0; i < n; i++ {
		_, _ = stmt.Exec(
			fmt.Sprintf("User%d", i),
			fmt.Sprintf("user%d@example.com", i),
			[]string{"active", "disabled"}[i%2],
		)
	}
}

// 普通の Query
func selectWithQuery(db *sql.DB, emails []string) time.Duration {
	start := time.Now()
	for _, email := range emails {
		rows, err := db.Query(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`, email)
		must(err)
		var id int
		var name string
		for rows.Next() {
			rows.Scan(&id, &name)
		}
		rows.Close()
	}
	return time.Since(start)
}

// PreparedStatement を利用
func selectWithPrepared(db *sql.DB, emails []string) time.Duration {
	start := time.Now()
	stmt, err := db.Prepare(`SELECT id, name FROM users WHERE email = $1 AND status = 'active'`)
	must(err)
	defer stmt.Close()

	for _, email := range emails {
		rows, err := stmt.Query(email)
		must(err)
		var id int
		var name string
		for rows.Next() {
			rows.Scan(&id, &name)
		}
		rows.Close()
	}
	return time.Since(start)
}

func main() {
	db := connectDB()
	defer db.Close()

	// データ5万件
	setupTable(db, 50000)

	// テスト対象のクエリ
	emails := []string{}
	for i := 0; i < 5000; i++ {
		emails = append(emails, fmt.Sprintf("user%d@example.com", i))
	}

	// 通常の Query
	d1 := selectWithQuery(db, emails)
	fmt.Printf("普通の Query: %v\n", d1)

	// PreparedStatement
	d2 := selectWithPrepared(db, emails)
	fmt.Printf("PreparedStatement: %v\n", d2)
}
