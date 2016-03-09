package main

import(
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pachyderm/sandbox/src/asset_handler"
)

var assetHandler = asset_handler.NewAssetHandler()

func main() {
	router := gin.Default()

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.GET("/", func(c *gin.Context) {
		fmt.Printf("wow okzz")
		name := c.Query("name")
		c.String(http.StatusOK, "Hello [%s]", name)
	})

	router.Run(":5678")
}
