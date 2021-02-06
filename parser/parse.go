package main

import (
	"encoding/json"
	"grass_scraper/db_manager"
	"io/ioutil"
	"os"
	"strings"
)


var dbManager db_manager.DBManager
var dataDirectory = "data"
type floraMap map[string]string
type floraData []floraMap

func main() {
	dbManager.Initialize("grass_user", os.Getenv("MYSQL_PSW"), "localhost", "grass_db")

	jsonReader, err := readJsons(dataDirectory)

	if err != nil {
		panic(err)
	}

	for json := range jsonReader {
		// idx, value. Value = a
		row := parseFloraJson(json)
		return
		dbManager.InsertRow(row)
	}
}

func parseFloraJson(json floraMap) db_manager.GrassEntry {
	var row db_manager.GrassEntry

	/*
	var floraDescriptors = []string {
		"DISTRIBUTION",
		"FERTILE SPIKELETS",
		"FLORETS",
		"FLOWER",
		"FRUIT",
		"GLUMES",
		"HABIT",
		"INFLORESCENCE",
		"NOTES",
	}
	*/

	row.GenusSpecies = json["Name"]
	for key, val := range json {
		json[key] = strings.Replace(val, "\n", " ", -1)
	}
	parseHabit(json["HABIT"], &row)

	return row
}

func parseHabit(fieldData string, row *db_manager.GrassEntry) {
	fields := strings.Split(fieldData, ". ")
	for field := range fields {
		print(field)
	}
}

func readJsons(dir string) (<- chan floraMap, error) {
	files := GetFiles(dir)
	channel := make(chan floraMap)

	for _, file := range files {
		fileName := file.Name()
		fullPath := strings.Join([]string{dir, fileName}, "\\")
		data := GetFileData(fullPath)
		var parsedData floraData
		err := json.Unmarshal(data, &parsedData)

		if err != nil {
			panic(err)
		}

		// remove new line from all the fields
		// idx, v. Value = A single floraMap
		go func() {
			for _, v := range parsedData {
				channel <- v
			}
			close(channel)
		}()
	}

	return channel, nil
}

func GetFiles(directory string) []os.FileInfo {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		panic(err)
	}

	return files
}

func GetFileData(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return data
}