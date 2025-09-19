setup:
	psql -h localhost -p 5432 -U postgres -d testdb < setup.sql

bench:
	go test -bench . -benchtime 5s -benchmem
