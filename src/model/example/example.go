/*
Sniperkit-Bot
- Status: analyzed
*/

package example

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/pachyderm/pachyderm/src/client"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"

	"github.com/sniperkit/snk.fork.pachyderm-sandbox/src/asset"
	"github.com/sniperkit/snk.fork.pachyderm-sandbox/src/model/repo"
)

type Example struct {
	Name string

	// Util
	client   *client.APIClient
	rawFiles *asset.AssetHandler

	// Data Pane
	Repo *repo.SandboxRepo

	// Code Pane
	Code string

	// Content Pane
	Steps      []template.HTML
	StepNumber int
}

func New(name string, APIClient *client.APIClient, assetHandler *asset.AssetHandler) (*Example, error) {
	r, err := repo.New(APIClient, name)

	if err != nil {
		return nil, err
	}

	ex := &Example{
		Name:     name,
		client:   APIClient,
		Repo:     r,
		rawFiles: assetHandler,
	}

	err = ex.loadSteps()

	if err != nil {
		return nil, err
	}

	err = ex.populateRepo()

	if err != nil {
		return nil, err
	}

	code, err := ex.loadPipeline()

	if err != nil {
		return nil, err
	}

	ex.Code = code

	ex.Repo.LoadFileData()

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

	r, err := repo.Load(APIClient, unique_name)

	if err != nil {
		return nil, errors.New("Could not load repo: " + unique_name)
	}

	if err != nil {
		return nil, err
	}

	ex := &Example{
		Name:     example_name,
		client:   APIClient,
		Repo:     r,
		rawFiles: assetHandler,
	}

	err = ex.loadSteps()

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

	for _, file := range files {
		commit, err := pfs_client.StartCommit(e.client, e.Repo.Name, "", "")

		if err != nil {
			return err
		}

		content, err := e.rawFiles.FindOrPopulate(file)

		if err != nil {
			return err
		}

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

func (ex *Example) loadSteps() error {
	path := fmt.Sprintf("assets/examples/%v/steps", ex.Name)
	entries, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		html, err := ex.rawFiles.FindOrPopulate(path + "/" + entry.Name())

		if err != nil {
			return err
		}

		ex.Steps = append(ex.Steps, template.HTML(string(html)))
	}

	ex.StepNumber = 0

	return nil
}
