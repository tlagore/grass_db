package main

import (
	"encoding/json"
	"fmt"
	"grass_scraper/db_manager"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)


var dbManager db_manager.DBManager
var dataDirectory = "data"
var digitRegExp = regexp.MustCompile(`\d+[\.]?[\d]*-\d+[\.]?[\d]*`)
type floraMap map[string]string
type floraData []floraMap

func main() {
	dbManager.Initialize("grass_user", os.Getenv("MYSQL_PSW"), "localhost", "grass_db")

	jsonReader, err := readJsons(dataDirectory)

	//min, max := getDigits()
	//fmt.Println(min)
	//fmt.Println(max)
	//fmt.Println("Culms 1–3 cm long" == "Culms 1-3 cm long")
	//fmt.Println("Culms" == "Culms")
	//min, max := getDigits("Culms 1–3 cm long")
	//fmt.Println(min)
	//fmt.Println(max)
	//return

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
		// replace all new lines with spaces and replace odd unicode character with regular dash
		json[key] = strings.Replace(strings.Replace(val, "\n", " ", -1), "–", "-", -1)
	}
	parseHabit(json["HABIT"], &row)

	return row
}

func parseHabit(fieldData string, row *db_manager.GrassEntry) {
	fields := strings.Split(fieldData, ". ")
	for _, field := range fields {
		pieces := strings.Split(field, ";")
		if strings.HasPrefix(field, "Culms") {
			parseCulms(pieces, row)
		}
	}
}

func parseCulms(culmParts []string, row *db_manager.GrassEntry) {
	for _, part := range culmParts {
		if strings.Contains(part, "diam") {
			row.CulmDiameterMinMm, row.CulmDiameterMaxMm = getDigits(part)
		}else if strings.Contains(part, "long") {
			row.CulmLengthMinCm, row.CulmLengthMaxCm = getDigits(part)
		}
	}
}

func getDigits(str string) (int, int) {
	var min, max = 0, 0
	findStr := digitRegExp.FindAllString(str, 1)
	fmt.Println(digitRegExp.String())
	if len(findStr) == 1 {
		digits := strings.Split(findStr[0], "-")
		if len(digits) == 2 {
			min, _ = strconv.Atoi(digits[0])
			max, _ = strconv.Atoi(digits[1])
		}
	}

	return min, max
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