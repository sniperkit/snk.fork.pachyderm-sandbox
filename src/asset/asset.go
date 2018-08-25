/*
Sniperkit-Bot
- Status: analyzed
*/

package asset

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	content, err := a.FindOrPopulate(path)

	if err != nil {
		c.String(http.StatusNotFound, "Asset not found")
	}

	c.String(http.StatusOK, string(content))
	setMIMEType(c, path)

}

func (a *AssetHandler) FindOrPopulate(path string) (content []byte, err error) {
	_, ok := a.files[path]

	if !ok || gin.Mode() == "debug" {
		content, err := ioutil.ReadFile(path)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		a.files[path] = content
	}

	return a.files[path], nil
}

func setMIMEType(c *gin.Context, path string) {
	tokens := strings.Split(path, ".")
	suffix := tokens[len(tokens)-1]
	contentType := fmt.Sprintf("text/%v", suffix)
	c.Request.Header.Set("Content-Type", contentType)
}
