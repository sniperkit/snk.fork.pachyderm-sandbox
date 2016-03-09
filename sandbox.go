package main

import(
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		fmt.Printf("wow okzz")
		name := c.Query("name")
		c.String(http.StatusOK, "Hello [%s]", name)
	})

	router.Run(":5678")
}
