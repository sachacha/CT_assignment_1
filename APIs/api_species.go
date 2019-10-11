package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SpecyAnswer struct {
	Key            float64 `json:"key"`
	Kingdom        string  `json:"kingdom"`
	Phylum         string  `json:"phylum"`
	Order          string  `json:"order"`
	Family         string  `json:"family"`
	Genus          string  `json:"genus"`
	ScientificName string  `json:"scientificName"`
	CanonicalName  string  `json:"canonicalName"`
	Year           float64 `json:"year"`
}

func HandlerSpecies(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 5 || parts[3] != "species" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	// get the specy json
	var getArgument = fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?speciesKey=%s", parts[4])

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

	// define the limit
	limitRequest := 100
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

	// append speciy details to our answer
	speciesAnswer := make(map[int]*SpecyAnswer)

	for i := 0; i < limitRequest; i++ {
		if specyKey, ok0 := species[i]["speciesKey"].(float64); ok0 {
			if specyKingdom, ok1 := species[i]["kingdom"].(string); ok1 {
				if specyPhylum, ok2 := species[i]["phylum"].(string); ok2 {
					if specyOrder, ok3 := species[i]["order"].(string); ok3 {
						if specyFamily, ok4 := species[i]["family"].(string); ok4 {
							if specyGenus, ok5 := species[i]["genus"].(string); ok5 {
								if specyScientificName, ok6 := species[i]["scientificName"].(string); ok6 {
									if specyCanonicalName, ok7 := species[i]["genericName"].(string); ok7 {
										if specyYear, ok8 := species[i]["year"].(float64); ok8 {
											speciesAnswer[i] = &SpecyAnswer{specyKey, specyKingdom, specyPhylum, specyOrder, specyFamily, specyGenus, specyScientificName, specyCanonicalName, specyYear}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// encoding the answer
	json.NewEncoder(w).Encode(speciesAnswer)
}
