package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var fileLoc = "..\\data"
type floraMap map[string]string
type floraData []floraMap

func main() {
	files := getFiles(fileLoc)
	for _, file := range files {
		fileName := file.Name()
		fullPath := strings.Join([]string{fileLoc, fileName}, "\\")
		fmt.Println(fullPath)
		data := getFileData(fullPath)
		var parsedData floraData
		err := json.Unmarshal(data, &parsedData)

		if err != nil {
			panic(err)
		}

		// remove new line from all the fields
		// idx, value. Value = A single floraMap
		for idx, v := range parsedData {
			// idx, value. Value = a
			for key, val := range v {
				parsedData[idx][key] = strings.Replace(val, "\n", "", -1)
				fmt.Println(parsedData[idx][key])
			}
		}
	}
}

func getFiles(directory string) []os.FileInfo {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		panic(err)
	}

	return files
}

func getFileData(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return data
}