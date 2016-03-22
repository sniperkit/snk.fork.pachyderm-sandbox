package data

import(
	"fmt"
	"github.com/pachyderm/pachyderm/src/client"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"
	pfs_server "github.com/pachyderm/pachyderm/src/server/pfs"
)

type SandboxRepo struct {
	DisplayName string
	pfs_server.Repo
}

type SandboxExample struct {
	// Data Pane
	Repo *SandboxRepo
	Files map[string]pfs_server.File

	// Code Pane
	Code string
}

func initializeExample(APIClient *client.APIClient) (*SandboxExample, error) {
	repo, err := createUniqueRepo(APIClient)

	if err != nil {
		return nil, err
	}

	

	return &SandboxExample{
		Repo: repo,
		Files: files,
		Code: code,
	}, nil
}

func createUniqueRepo(APIClient *client.APIClient) (*SandboxRepo, error) {

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

