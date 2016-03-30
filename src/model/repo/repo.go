package repo

import(
	"strings"
	"bytes"

	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"
	pfs_server "github.com/pachyderm/pachyderm/src/server/pfs"
	"github.com/pachyderm/pachyderm/src/client"

	"github.com/pachyderm/sandbox/src/util"
)

type SandboxRepo struct {
	DisplayName string
	Files map[string]map[string][][]string //name -> commit -> 2D data
	client *client.APIClient

	*pfs_server.Repo
}

func New(APIClient *client.APIClient, name string) (*SandboxRepo, error) {

	unique_name := util.GenerateUniqueName(name)

	err := pfs_client.CreateRepo(APIClient, unique_name)

	if err != nil {
		return nil, err
	}

	r := createRepo(APIClient, name, unique_name)

	return r, nil
}

func createRepo(APIClient *client.APIClient, name string, unique_name string) *SandboxRepo {
	return &SandboxRepo{
		DisplayName: name,
		client: APIClient,
		Files: make(map[string]map[string][][]string),
		Repo: &pfs_server.Repo{
			Name: unique_name,
		},
	}
}

func Load(APIClient *client.APIClient, unique_name string) (*SandboxRepo, error) {
	name := util.NameFromUniqueName(unique_name)

	r := createRepo(APIClient, name, unique_name)

	if err := r.LoadFileData(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SandboxRepo) LoadFileData() error {
	commitInfos, err := pfs_client.ListCommit(r.client, []string{ r.Name })

	if err != nil {
		return err
	}

	for _, commitInfo := range(commitInfos) {
		commitID := commitInfo.Commit.ID

		fileInfos, err := pfs_client.ListFile(r.client, r.Name, commitID, "", "", nil)
		if err != nil {
			return err
		}		

		for _, fileInfo := range(fileInfos) {
			var buffer bytes.Buffer
			err = pfs_client.GetFile(
				r.client, 
				r.Name, 
				commitID, 
				fileInfo.File.Path, 
				0, 
				0,
				"", 
				nil, 
				&buffer)
			
			if err != nil {
				return err
			}
			
			_, ok := r.Files[fileInfo.File.Path]

			if !ok {
				r.Files[fileInfo.File.Path] = make(map[string][][]string)
			}

			r.Files[fileInfo.File.Path][commitID] = parseData(buffer.String())
		}

	}

	return nil
}

func parseData(raw string) (data [][]string) {

	lines := strings.Split(raw,"\n")

	for _, line := range(lines) {
		tokens := strings.Fields(line)

		datum := make([]string, 0)

		for _, token := range(tokens) {
			datum = append(datum, token)
		}
		data = append(data, datum)
	}

	return data
}
