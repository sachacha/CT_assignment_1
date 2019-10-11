package APIs

import (
	"encoding/json"
	"net/http"
	"time"
)

var START_TIME = time.Now()

type DiagAnswer struct {
	Gbif          int     `json:"gbif"`
	Restcountries int     `json:"restcountries"`
	Version       string  `json:"verion"`
	Uptime        float64 `json:"uptime"`
}

func HandlerDiag(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	respGBIF, errGBIF := http.Get("http://api.gbif.org/v1/occurrence/search")

	if errGBIF != nil {
		http.Error(w, errGBIF.Error(), http.StatusBadRequest)
		return
	}

	respCountry, errCountry := http.Get("https://restcountries.eu/rest/v2/all")

	if errCountry != nil {
		http.Error(w, errCountry.Error(), http.StatusBadRequest)
		return
	}

	uptime := time.Since(START_TIME).Seconds()

	diagAnswer := DiagAnswer{respGBIF.StatusCode, respCountry.StatusCode, "v1", uptime}

	json.NewEncoder(w).Encode(diagAnswer)
}
