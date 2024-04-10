package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSON map[string]interface{}

func main() {
    // TODO: ENV VARS
    port := "8080"

    r := gin.Default()

    r.GET("/ping", ping)
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

func ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "pong": "Im healthy",
    })
}

/* POST /shorten
 * Body Shape ={"url": "urlToShorten"(string) }
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

    c.JSON(http.StatusOK, gin.H{
        "valid": "body",
    })
}
