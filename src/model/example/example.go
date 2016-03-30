package example

type Example struct {
	Name string

	// Util
	client *client.APIClient
	rawFiles *asset.AssetHandler

	// Data Pane
	Repo *SandboxRepo

	// Code Pane
	Code string
}


func New(name string, APIClient *client.APIClient, assetHandler *asset.AssetHandler) (*Example, error) {
	repo, err := createUniqueRepo(APIClient)

	if err != nil {
		return nil, err
	}

	code, err := assetHandler.FindOrPopulate(fmt.Sprintf("assets/examples/%v/code.go", name))

	if err != nil { 
		return nil, err
	}

	ex := &Example{
		Name: name,
		client: APIClient,
		Repo: repo,
		rawFiles: assetHandler,
		Code: string(code),
	}

	err = ex.populateRepo()

	if err != nil {
		return nil, err
	}

	return ex, nil
}


func LoadFromCookie(session sessions.Session, APIClient *client.APIClient, assetHandler *asset.AssetHandler) (*Example, error) {


	value := session.Get("example_name")

	if value == nil {
		return nil, errors.New("Could not find example_name in session")
	}
	example_name := value.(string)
	fmt.Printf("Got example name (%v)\n", example_name)

	value = session.Get("repo_name")

	if value == nil {
		return nil, errors.New("Could not find repo_name in session")
	}

	unique_name := value.(string)
	name := strings.Split(unique_name, "-")[0]
	fmt.Printf("Got repo names (%v) (%v)\n", name, unique_name)

	repo := &SandboxRepo{
		DisplayName: name,
		Files: make(map[string]map[string][][]string),
		Repo: &pfs_server.Repo{
			Name: unique_name,
		},
	}

	code, err := assetHandler.FindOrPopulate(fmt.Sprintf("assets/examples/%v/code.go", example_name))

	if err != nil { 
		return nil, err
	}

	ex := &Example{
		Name: example_name,
		client: APIClient,
		Repo: repo,
		rawFiles: assetHandler,
		Code: string(code),
	}
	err = ex.loadFileData()

	if err != nil {
		return nil, err
	}	

	return ex, nil
}


func (e *Example) populateRepo() error {

	if e.Repo == nil {
		return fmt.Errorf("No repo initialized")
	}

	files := []string{
		fmt.Sprintf("assets/examples/%v/data/set1.txt", e.Name), 
		fmt.Sprintf("assets/examples/%v/data/set2.txt", e.Name),
	}

	// For now hardcode:
	destinationFile := "sales"

	for _, file := range(files) {
		commit, err := pfs_client.StartCommit(e.client, e.Repo.Name, "", "")

		if err != nil {
			return err
		}

		content, err := e.rawFiles.FindOrPopulate(file)

		if err != nil {
			return err
		}

		commitToContentMap, ok := e.Repo.Files[destinationFile]

		// SJ: this feels weird ... 
		if !ok {
			e.Repo.Files[destinationFile] = make(map[string][][]string)
			commitToContentMap = e.Repo.Files[destinationFile]
		}
		
		commitToContentMap[commit.ID] = parseData(string(content))

		_, err = pfs_client.PutFile(
			e.client,
			e.Repo.Name,
			commit.ID,
			destinationFile,
			0,
			strings.NewReader(string(content)),
		)

		if err != nil {
			return err
		}

		err = pfs_client.FinishCommit(e.client, e.Repo.Name, commit.ID)

		if err != nil {
			return err
		}

	}

	return nil
}

func createUniqueRepo(APIClient *client.APIClient) (*SandboxRepo, error) {
	unique_name := generateUniqueName("sales")

	err := pfs_client.CreateRepo(APIClient, unique_name)

	if err != nil {
		return nil, err
	}

	repo := &SandboxRepo{
		DisplayName: "Sales",
		Files: make(map[string]map[string][][]string),
		Repo: &pfs_server.Repo{
			Name: unique_name,
		},
	}
	return repo, nil
}
