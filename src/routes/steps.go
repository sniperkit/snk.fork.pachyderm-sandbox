package routes

import(
	"fmt"

	"github.com/pachyderm/sandbox/src/example"	
)

func step1(c *gin.Context) (ex *example.Example, errors []error){

	ex, err := example.New("fruit-stand", APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	} else {
		// Silly ... but go compiler doesn't know I'm using it in a view
		fmt.Printf("Loaded %v\n", ex.Name)			
	}

	repos, err := pfs_client.ListRepo(APIClient)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	return ex, errors
}

func extractCookie(c *gin.Context, name string) (string, error) {
	cookies := strings.Split(c.Request.Header["Cookie"], ";")

	for _, cookie := range(cookies) {
		tokens := strings.Split(cookie, "=")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Invalid cookie value")
		}

		if name == tokens[0] {
			return value, nil
		}
	}

	return nil, nil
}

func step1submit(c *gin.Context) (ex *example.Example, errors []error) {
	cookie := extractCookie(c, "example")

	ex, err := example.LoadFromCookie(cookie, APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	return ex, errors
}
