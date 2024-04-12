package main

import (
    "log"
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "postgresql://localhost/url_short?user=me&password=me"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("An error occured while connecting to psql: %v", err)
        return
    }

    db.Exec("CREATE TABLE users(name TEXT PRIMARY KEY);")
    db.Exec("INSERT INTO users (name) VALUES ($1)", "a")
    db.Exec("CREATE TABLE urls(uuid PRIMARY KEY DEFAULT gen_random_uuid(), identifier INT, owner TEXT references users (name), url TEXT)")
}

// TODO: MIGRATOR, for now we just will init a default table
