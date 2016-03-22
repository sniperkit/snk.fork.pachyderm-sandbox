package main

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

func step1submit(c *gin.Context) (ex *example.Example, errors []error) {
	ex, err := LoadFromCookie(cookie, APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	return ex, errors
}
