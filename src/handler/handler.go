package handler

import(
	"io/ioutil"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	files map[string][]byte
}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{
		files: make(map[string][]byte),
	}
}

func (a *AssetHandler) Serve(c *gin.Context) {
	
	path := fmt.Sprintf(".%v", c.Request.URL.Path)
	content, ok := a.files[path]

	if !ok || gin.Mode() == "debug" {
		fmt.Println(path)
		content, err := ioutil.ReadFile(path)

		fmt.Printf("Found content: %v\n", string(content))

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "Asset not found")
			return
		}

		a.files[path] = content
	}

	c.String(http.StatusOK, string(content) )
}


