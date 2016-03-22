package example

import(
	"fmt"
	"io"
	"strings"

	"github.com/satori/go.uuid"
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

func LoadFromCookie(cookie string, APIClient *client.APIClient, assetHandler *asset.AssetHandler) (*Example, error) {

	// cookie is of type example=name/customtoken

	tokens := strings.Split(cookie, "/")

	if len(tokens) != 2 {
		return nil, fmt.Errorf("Error extracting cookie information")
	}

	name := tokens[0]
	id := tokens[1]

	repo := &SandboxRepo{
		DisplayName: name,
		Repo: &pfs_server.Repo{
			Name: uniqueNameFromToken(name,token),
		},
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

	err = ex.loadFileData()

	if err != nil {
		return nil, err
	}	

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
		
		fileInfos, err := pfs_client.ListFile(e.client, e.Repo.Name, commitID, "", "", nil)
		if err != nil {
			return err
		}		

		for _, fileInfo := range(fileInfos) {
			writer := NewBufferWriter()

			err = pfs_client.GetFile(e.client, e.Repo.Name, commitID, "", 0, fileInfo.size, "", nil, writer)
			if err != nil {
				return err
			}
			
			e.Files[fileInfo.File.Path][commitID] = writer.content
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

