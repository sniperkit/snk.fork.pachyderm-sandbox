package example

import(
	"fmt"
	"errors"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/pachyderm/pachyderm/src/client"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"
	pfs_server "github.com/pachyderm/pachyderm/src/server/pfs"

	"github.com/pachyderm/sandbox/src/asset"
)

type Example struct {
	Name string

	// Util
	client *client.APIClient

	// Data Pane
	Repo *SandboxRepo
	Files map[string]map[string]string //name -> commit -> content
	rawFiles *asset.AssetHandler

	// Code Pane
	Code string
}

type SandboxRepo struct {
	DisplayName string
	*pfs_server.Repo
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
		Repo: &pfs_server.Repo{
			Name: unique_name,
		},
	}

	fmt.Printf("Created repo")
	code, err := assetHandler.FindOrPopulate(fmt.Sprintf("assets/examples/%v/code.go", example_name))

	if err != nil { 
		return nil, err
	}

	fmt.Printf("Loaded code")
	ex := &Example{
		Name: example_name,
		client: APIClient,
		Repo: repo,
		Files: make(map[string]map[string]string), // Initialize filename -> commitID[content] map
		rawFiles: assetHandler,
		Code: string(code),
	}
	fmt.Printf("Created example")
	err = ex.loadFileData()
	fmt.Printf("Loaded example data")

	if err != nil {
		return nil, err
	}	

	fmt.Printf("succeeded in load")
	return ex, nil
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
		Files: make(map[string]map[string]string), // Initialize filename -> commitID[content] map
		rawFiles: assetHandler,
		Code: string(code),
	}

	err = ex.populateRepo()

	if err != nil {
		return nil, err
	}

	return ex, nil
}

func (e *Example) loadFileData() error {
	commitInfos, err := pfs_client.ListCommit(e.client, []string{ e.Repo.Name })

	if err != nil {
		return err
	}

	for _, commitInfo := range(commitInfos) {
		commitID := commitInfo.Commit.ID
		fmt.Printf("Looking at commit %v\n", commitID)

		fileInfos, err := pfs_client.ListFile(e.client, e.Repo.Name, commitID, "", "", nil)
		if err != nil {
			fmt.Printf("err listing file")
			return err
		}		

		for _, fileInfo := range(fileInfos) {
			writer := NewBufferWriter()
			fmt.Printf("Getting file %v\n", fileInfo.File.Path)

			err = pfs_client.GetFile(
				e.client, 
				e.Repo.Name, 
				commitID, 
				"", 
				0, 
				int64(fileInfo.SizeBytes), 
				"", 
				nil, 
				writer)
			
			if err != nil {
				fmt.Printf("error getting file: %v\n", err)
//				return err
			}
			
			e.Files[fileInfo.File.Path][commitID] = string(writer.Content)
		}

	}

	return nil
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

		commitToContentMap, ok := e.Files[destinationFile]

		// SJ: this feels weird ... 
		if !ok {
			e.Files[destinationFile] = make(map[string]string)
			commitToContentMap = e.Files[destinationFile]
		}
		
		commitToContentMap[commit.ID] = string(content)

		contentReader := NewCacheReader(content)

		_, err = pfs_client.PutFile(e.client, e.Repo.Name, commit.ID, destinationFile, 0, contentReader)

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
		Repo: &pfs_server.Repo{
			Name: unique_name,
		},
	}
	return repo, nil
}

