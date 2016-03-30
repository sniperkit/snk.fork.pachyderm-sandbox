package repo


type SandboxRepo struct {
	DisplayName string
	Files map[string]map[string][][]string //name -> commit -> 2D data
	client *pfs_client

	*pfs_server.Repo
}

func (r *SandboxRepo) LoadFileData() error {
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
			var buffer bytes.Buffer
			err = pfs_client.GetFile(
				e.client, 
				e.Repo.Name, 
				commitID, 
				"sales", 
				0, 
				0,
				"", 
				nil, 
				&buffer)
			
			if err != nil {
				return err
			}
			
			_, ok := e.Repo.Files[fileInfo.File.Path]

			if !ok {
				e.Repo.Files[fileInfo.File.Path] = make(map[string][][]string)
			}

			e.Repo.Files[fileInfo.File.Path][commitID] = parseData(buffer.String())
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
