package main

import (
	"bufio"
	"encoding/json"
	"grassscraper.ty/db_manager"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var dbManager db_manager.DBManager
var dataDirectory = "data"
var bambooSpeciesFile = "bamboo_species.txt"
var bracketedDigit = regexp.MustCompile(`[\d]+`)
var notDigit = regexp.MustCompile(`[^\d]*`)
var digits = regexp.MustCompile(`\d+[.]?[\d]*`)
var digitRegExp = regexp.MustCompile(`\d+[\.]?[\d]*-\d+[\.]?[\d]*`)
type floraMap map[string]string
type floraData []floraMap


/*
 Parse grass db jsons into database
 */
func main() {

	runtime.GOMAXPROCS(30)
	dbManager.Initialize("grass_user", os.Getenv("MYSQL_PSW"), "tcp(localhost:3306)", "grass_db")
	if !dbManager.TableExists("grass_table") {
		panic("grass_table does not exist on given connection! Create it with create_schema.sql first.")
	}

	if !dbManager.TableExists("bamboo_species") {
		panic ("bamboo_table does not exist on given connection! Create it with create_schema.sql first.")
	}
	parseBambooFile(bambooSpeciesFile)

	jsonReader, err := readJsons(dataDirectory)

	if err != nil {
		panic(err)
	}

	// idx := 0
	for jsonData := range jsonReader {
		row := parseFloraJson(jsonData)

		insertError := dbManager.InsertGrassRow(&row)
		if insertError != nil {
			log.Printf("%s", insertError)
		}
	}
}

func parseBambooSpecies(line string) db_manager.BambooEntry {
	var row db_manager.BambooEntry

	if strings.Contains(line, "+") {
		row.IsInvasive = true
		line = strings.Replace(line, "+", "", 1)
	}

	if strings.Contains(line, "*") {
		findStr := bracketedDigit.FindString(line)
		if findStr != "" {
			row.NumIntroductions, _ = strconv.Atoi(findStr)
			line = strings.Replace(line, "("+findStr+")", "", 1)
		}

		line = strings.Replace(line, "*",  "",1)
	}

	if strings.Contains(line, "?") {
		row.DisputedNativeRange = true
		line = strings.Replace(line, "?", "", 1)
	}

	line = strings.Trim(line, " ")
	row.GenusSpecies = line
	return row
}

func parseBambooFile(fileName string) {
	readFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		row := parseBambooSpecies(fileScanner.Text())
		dbManager.InsertBambooRow(&row, "bamboo_species")
	}
}

func parseFloraJson(json floraMap) db_manager.GrassEntry {
	var row db_manager.GrassEntry

	for key, val := range json {
		// replace all new lines with spaces and replace odd unicode character with regular dash
		json[key] = strings.Replace(strings.Replace(val, "\n", " ", -1), "–", "-", -1)
	}
	row.GrassSpecies = json["Name"]
	parseHabit(strings.ToLower(json["HABIT"]), &row)
	parseLocation(strings.ToLower(json["DISTRIBUTION"]), &row)
	row.Notes = strings.Replace(json["NOTES"], "NOTES ", "", 1)

	return row
}

func parseLocation(fieldData string, row *db_manager.GrassEntry) {
	fieldData = strings.Replace(fieldData, "distribution", "", 1)

	locations := strings.Split(fieldData, ".")
	locationNarrow := ""
	locationBroad := ""
	for _, location := range locations {
		location = strings.Trim(location, " ")

		if location != "" {
			if strings.Contains(location, ":") {
				locParts := strings.Split(location, ":")
				if locationNarrow != "" {
					locationNarrow += "; "
				}

				if locationBroad != "" {
					locationBroad += "; "
				}

				locationNarrow += strings.Trim(locParts[0], " ")
				locationBroad += parseDefHelper(locParts[1:len(locParts)], 0)
			} else {
				if locationBroad != "" {
					locationBroad += "; "
				}
				locationBroad += fieldData
				locationBroad = fieldData
			}
		}
	}

	row.LocationBroad = locationBroad
	row.LocationNarrow = locationNarrow
}

func parseHabit(fieldData string, row *db_manager.GrassEntry) {
	fieldData = strings.Trim(strings.Replace(fieldData, "habit", "", 1), " ")

	fields := strings.Split(fieldData, ". ")
	for _, field := range fields {
		if strings.Contains(field, "perennial"){
			row.IsPerennial = true
			parseCulmDensity(field, row)
		}

		if strings.Contains(field, "annual") {
			row.IsAnnual = true
		}

		if strings.Contains(field, "woody") {
			row.IsWoody = true
		}

		if strings.HasPrefix(field, "rhizomes") {
			parseRootingCharacteristics(field, row)
		}

		if strings.HasPrefix(field, "culms") {
			parseCulmDef(field, row)
			parseCulms(field, row)
		}

		if strings.HasPrefix(field, "culm-internodes") {
			parseCulmInternodes(field, row)
		}
	}
}

func parseCulmDensity (field string, row *db_manager.GrassEntry) {
	if strings.Contains(field, ";") {
		parts := strings.Split(field, ";")
		culmDensity := parseDefHelper(parts[1:len(parts)], 0)

		if culmDensity != "" {
			row.CulmDensity = culmDensity
		}
	}
}

func parseDefHelper(parts []string, expectedStart int) string {
	def := ""

	for idx, val := range parts {
		val = strings.Trim(val, " ")
		if val != "" {
			def = concatHelper(def, val, idx, expectedStart)
		}
	}

	return def
}

/*
 parse the Culm-internodes definition. field form:
	Culm-internodes def1;def2;def3
 */
func parseCulmInternodes (field string, row *db_manager.GrassEntry) {
	field = strings.Replace(field, "culm-internodes ", "", 1)
	parts := strings.Split(field, ";")
	internodesDef := parseDefHelper(parts, 0)

	if internodesDef != "" {
		row.CulmInternode = internodesDef
	}
}

func concatHelper(s string, toConcat string, idx int, expectedStart int) string {
	if idx == expectedStart {
		return s + toConcat
	} else {
		return s + ", " + toConcat
	}
}

/*
	parse rooting characteristics
 */
func parseRootingCharacteristics(rootingChars string, row *db_manager.GrassEntry) {
	rootingChars = strings.Replace(rootingChars, "rhizomes ", "", 1)
	parts := strings.Split(rootingChars, ";")
	rootingCharDef := parseDefHelper(parts, 0)

	if rootingCharDef != "" {
		row.RootingCharactersitic = rootingCharDef
	}
}

/*
	Extract the length and diameter from the Culms definition

	Looks something like:
	Culms erect; 300–700 cm long; 10–30 mm diam.; woody.
 */
func parseCulms(field string, row *db_manager.GrassEntry) {
	culmParts := strings.Split(field, ";")

	for _, part := range culmParts {
		// Try to extract
		if strings.Contains(part, "diam") || strings.Contains(part, "wide") {
			row.CulmDiameterMinMm, row.CulmDiameterMaxMm = getDigits(part)
		}else if strings.Contains(part, "long") {
			row.CulmLengthMinCm, row.CulmLengthMaxCm = getDigits(part)
		}
	}
}

/*
	Culm definition is a part of the main definition of Culms and may not exist.

	If it exists it'll look like something Culms erect; other descritor; 40-30cm long. So we parse up until we see a digit
	with the notDigit regex
 */
func parseCulmDef(field string, row *db_manager.GrassEntry) {
	matchStr := notDigit.FindString(field)
	if matchStr != "" {
		matchStr = strings.Replace(matchStr, "culms ", "", 1)
		if matchStr != "" {
			culmDefParts := strings.Split(matchStr, ";")
			culmDef := parseDefHelper(culmDefParts, 0)

			if culmDef != "" {
				row.CulmGrowth = culmDef
			}
		}
	}
}

func getDigits(str string) (float64, float64) {
	var min, max = 0.0, 0.0

	if strings.Contains(str, "-") {
		findStr := digitRegExp.FindAllString(str, 1)
		if len(findStr) == 1 {
			digits := strings.Split(findStr[0], "-")
			if len(digits) == 2 {
				min, _ = strconv.ParseFloat(digits[0],64)
				max, _ = strconv.ParseFloat(digits[1],64)
			}
		}
	} else {
		findStr := digits.FindAllString(str, 1)
		if len(findStr) == 1 {
			min, _ = strconv.ParseFloat(findStr[0], 64)
			max = min
		}
	}

	return min, max
}

func readJsons(dir string) (<- chan floraMap, error) {
	files := GetFiles(dir)
	channel := make(chan floraMap)
	wg := &sync.WaitGroup{}

	for _, file := range files {
		fileName := file.Name()
		fullPath := strings.Join([]string{dir, fileName}, "\\")
		data := GetFileData(fullPath)
		var parsedData floraData
		err := json.Unmarshal(data, &parsedData)

		if err != nil {
			panic(err)
		}

		// WaitGroup to know when to close channel
		wg.Add(1)

		// v = A single floraMap
		go func() {
			defer wg.Done()
			for _, v := range parsedData {
				channel <- v
			}
		}()
	}

	// when all producers are done, close the channel
	go func() {
		wg.Wait()
		close(channel)
	}()

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