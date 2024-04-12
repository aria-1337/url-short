package main

import (
	"encoding/json"
	"fmt"
    "log"
	"io"
	"net/http"
    "database/sql"
    _ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type JSON map[string]interface{}

func main() {
    // TODO: ENV VARS
    port := "8080"

    r := gin.Default()
    psql := db()

    r.GET("/ping", func(c *gin.Context) {
        ping(c, psql)
    })
    r.POST("/shorten", shorten)

    r.Run(fmt.Sprintf(":%s", port))
}

// TODO: handle errors better
func decodeBody(c *gin.Context) JSON {
    j, err := io.ReadAll(c.Request.Body)
    var body JSON
    if err != nil {
        return JSON{ "error": "500" }
    }

    if err := json.Unmarshal(j, &body); err != nil {
        return JSON{ "error": "500" }
    }
    return body
}

func validateBody(keys []string, body JSON) bool {
    for _, key := range keys {
        if _, ok := body[key]; ok {
            continue
        }
        return false
    }
    return true
}

func ping(c *gin.Context, d *sql.DB) {
    val := query(d, "SELECT $1 as val", 1)
    c.JSON(http.StatusOK, gin.H{
        "pong": "Im healthy",
        "val": val,
    })
}

/* POST /shorten
* Body Shape ={"user": "username"(string|not required) "url": "urlToShorten"(string) }
 */ 
func shorten(c *gin.Context) {
    body := decodeBody(c)
    valid := validateBody([]string{"url"}, body);
    if !valid {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "message": "Bad request body; expected: { url: yourUrlToShorten(string) }", 
        })
        return
    }

    /*
    // Retrieve user or set anon
    user := "a"
    if v, ok := body["user"].(string); ok {
        user = string(v)
    }
    */
    // Get link id

    c.JSON(http.StatusOK, gin.H{
        "valid": "body",
    })
}

func db() *sql.DB {
    connStr := "postgresql://localhost/url_short?user=me&password=me"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("An error occured while connecting to psql: %v", err)
    }
    return db
}

func query(d *sql.DB, query string, args ...interface{}) []interface{} {
    var res interface{}
    var all []interface{}
    rows, err := d.Query(query, args...)
    defer rows.Close()
    if err != nil {
        log.Fatalf("An error occured while querying: %v", err)
    }
    for rows.Next() {
        rows.Scan(&res)
        all = append(all, res)
    }
    return all
}
