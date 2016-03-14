package main

import(
	"net/http"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"

	"github.com/pachyderm/sandbox/src/asset"
	"github.com/pachyderm/pachyderm"
	"github.com/pachyderm/pachyderm/vendor/google.golang.org/grpc"

)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()

func getPachClient() (*pachyderm.APIClient, error) {
	pachAddr := os.Getenv("PACHD_PORT_650_TCP_ADDR")
	if pachAddr == "" {
		return nil, fmt.Errorf("PACHD_PORT_650_TCP_ADDR not set")
	}
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:650", pachAddr), grpc.WithInsecure())
	return pachyderm.NewAPIClient(clientConn), nil
}

func main() {

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.HTMLRender = loadTemplates()

	router.GET("/", handle("main"))

	router.Run(":9080")
	
}

func handle(page string) ( func (*gin.Context) ){
	return func(c *gin.Context) {
		pachd, err := getPachClient()
		var errors []error

		if err != nil {
			errors = append(errors, err)
		}

		if gin.Mode() == "debug" {
			router.HTMLRender = loadTemplates()
		}
		

		c.HTML(http.StatusOK, page, gin.H{
			"title" : "THE FINAL COUNTDOWN",
			"repos" : "",
			"errors": errors,
		})
	}
}


func loadTemplates() multitemplate.Render {
	templates := multitemplate.New()
	templates.AddFromFiles(
		"main",
		"views/base.html",
		"views/main.html",
		"views/data.html",
		"views/copy.html",
		"views/code.html",
	)

	return templates
}
