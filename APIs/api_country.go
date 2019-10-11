package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func floatInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type CountryAnswer struct {
	Code        string    `json:"code"`
	CountryName string    `json:"countryname"`
	CountryFlag string    `json:"countryflag"`
	Species     []string  `json:"species"`
	SpeciesKey  []float64 `json:"speciesKey"`
}

var countryAnswer = CountryAnswer{}

func HandlerCountry(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 5 || parts[3] != "country" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	limitRequest := 100

	// get the limit number if it's not the default one
	limit, ok := r.URL.Query()["limit"]

	if ok {
		limit, err := strconv.Atoi(limit[0])

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		limitRequest = limit
	}

	// get the species json
	var getArgument = fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?country=%s", parts[4])

	resp, err := http.Get(getArgument)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var specyJson interface{}

	err = json.Unmarshal(body, &specyJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var speciesMap = specyJson.(map[string]interface{})

	// define the final limit
	if limitJson, ok := speciesMap["limit"].(float64); ok {
		if limitRequest > int(limitJson) {
			limitRequest = int(limitJson)
		}
	}

	// create a map for each specy that we store in an array
	var speciesResult = speciesMap["results"].([]interface{})
	species := make(map[int]map[string]interface{})

	for i := 0; i < limitRequest; i++ {
		species[i] = speciesResult[i].(map[string]interface{})
	}

	// append species name and key to our answer
	for i := 0; i < limitRequest; i++ {
		if specyKey, ok0 := species[i]["speciesKey"].(float64); ok0 {
			if !(floatInSlice(specyKey, countryAnswer.SpeciesKey)) {
				if specyName, ok1 := species[i]["genericName"].(string); ok1 {
					countryAnswer.Species = append(countryAnswer.Species, specyName)
					countryAnswer.SpeciesKey = append(countryAnswer.SpeciesKey, specyKey)
				}
			}
		}
	}

	// get the country json
	getArgument = fmt.Sprintf("https://restcountries.eu/rest/v2/alpha/%s", parts[4])

	resp, err = http.Get(getArgument)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var countryJson interface{}

	err = json.Unmarshal(body, &countryJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var countryMap = countryJson.(map[string]interface{})

	// put country name and flag in our answer
	if name, okk := countryMap["name"].(string); okk {
		countryAnswer.CountryName = name
	}

	if flag, okkk := countryMap["flag"].(string); okkk {
		countryAnswer.CountryFlag = flag
	}

	// encoding the answer
	countryAnswer.Code = parts[4]
	json.NewEncoder(w).Encode(countryAnswer)
}
