package apidocmd

import (
	"encoding/json"
	"os"
)

// readProject read api-project json file
func readProject(path string) (*apiProject, error) {
	file, err := os.Open(path + "/api_project.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data apiProject
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// readProject read api-data json file
func readData(path string) (*groupDataList, error) {
	file, err := os.Open(path + "/api_data.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var list apiDataList
	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return nil, err
	}

	glist := list.Grouping()

	return &glist, nil
}
