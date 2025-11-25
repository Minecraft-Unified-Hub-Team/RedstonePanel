package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"net/http"
)

type Version struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
	Loaders      []string `json:"loaders"`
	Id       string `json:"id"`
}

type File struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Primary  bool   `json:"primary"`
	Size     int    `json:"size"`
}

func contains(list []string, example string) (bool) {
	for _, str := range list {
		if str == example {
			return true
		}
	}
	return false
}

func GetProjectFiles(projectName, versionName, loadersName string) ([]File, string, error) {
	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", projectName)
	loadersName = strings.ToLower(loadersName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	
	var versions []Version
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, "", err
	}
	
	for _, version := range versions {
		if (loadersName == "" || contains(version.Loaders, loadersName)) {
			return version.Files, version.Id, nil
		}
	}
	if (len(versions) == 0) {
		return nil, "", fmt.Errorf("version is not found")
	}
	return versions[0].Files, "", nil
}

func main() {

	files, mod_id, err := GetProjectFiles("iris", "0.91.0+1.20.1", "quilt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Found files:")
	for _, file := range files {
		fmt.Printf("- %s (%d bytes)\n", file.Filename, file.Size)
		if file.Primary {
			fmt.Printf("URL: %s id: %s\n", file.URL, mod_id)
		}
	}
}
