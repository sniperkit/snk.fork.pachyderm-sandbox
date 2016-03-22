package data

import(
	"fmt"
	"io"

	"github.com/pachyderm/pachyderm/src/client"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"
	pfs_server "github.com/pachyderm/pachyderm/src/server/pfs"

	"github.com/pachyderm/sandbox/src/asset"
)

type SandboxRepo struct {
	DisplayName string
	pfs_server.Repo
}

type SandboxExample struct {
	Name string

	// Util
	client *client.APIClient

	// Data Pane
	Repo *SandboxRepo
	Files map[string]string
	Commits []pfs_server.Commit
	rawFiles *AssetHandler

	// Code Pane
	Code string
}

func NewExample(name string, APIClient *client.APIClient, assetHandler *AssetHandler) error {
	repo, err := e.createUniqueRepo()

	if err != nil {
		return err
	}

	code, err := assetHandler.FindOrPopulate(fmt.Sprintf("/examples/%v/code.go", name))

	if err != nil {
		return err
	}

	ex := &SandboxExample{
		Name: name,
		client: APIClient,
		Repo: repo,
		Files: make(map[string]string),
		rawFiles: assetHandler,
		Code: code,
	}

	ex.initialize()

	return ex
}

func (e *Example) initialize() error {

	files := e.populateRepo(e.Repo)

	if err != nil {
		return err
	}

	e.Files = files

	return nil
}

func (e *Example) populateRepo() error {

	if e.Repo == nil {
		return fmt.Errorf("No repo initialized")
	}

	var commits []pfs_server.Commit

	files := []string{
		fmt.Sprintf("/examples/%v/data/set1.txt", e.Name), 
		fmt.Sprintf("/examples/%v/data/set2.txt", e.Name),
	}

	for file := range(files) {
		commit, err := pfs_client.StartCommit(e.client, e.Repo.Name, "", "")

		if err != nil {
			return err
		}

		content, err := e.rawFiles.FindOrPopulate(file)

		if err != nil {
			return err
		}

		e.Files["sales"] = string(content)

		_, err := pfs_client.PutFile(e.client, e.Repo.Name, commit.ID, "sales", 0, NewCacheReader(content))

		if err != nil {
			return err
		}

		err := pfs_client.FinishCommit(e.client, e.Repo.Name, "", "")

		if err != nil {
			return err
		}

		commits = append(commits, commit)
	}

	e.Commits = commits

	return nil
}

type CacheReader struct {
	content []byte
	index int
}

func NewCacheReader(content []byte) {
	return &CacheReader{
		content: content,
		index: 0,
	}
}

func (cr *CacheReader) Read(p []byte) (n int, err error) {
	if len(p) < len(cr.content) - index {
		p[0:len(p)-1] = cr.content[index:len(p-1)]
		cr.index = len(p)

		return len(p), nil
	}

	bufferSize := len(cr.content) - index
	p[0:bufferSize] = cr.content[index:-1]

	return bufferSize, io.EOF
}

func (e *Example) createUniqueRepo() (*SandboxRepo, error) {

	unique_suffix := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	unique_name := "sales" + "-" + unique_suffix[0:12]

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

