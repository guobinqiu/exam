package rate

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type RateController struct {
	rate *Rate //rate service
}

func NewRateController(xmlPath, dbPath string) *RateController {
	rate := NewRate(xmlPath, dbPath)
	err := rate.LoadToDB()
	if err != nil {
		log.Fatal(err)
	}
	return &RateController{rate}
}

func (s *RateController) GetLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("unsupported http method")
		fail(w, []byte("unsupported http method"))
		return
	}

	rates, err := s.rate.GetLatest()
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
		return
	}

	resp := make(map[string]interface{})
	resp["base"] = "EUR"
	resp["rates"] = rates

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (s *RateController) GetByTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("unsupported http method")
		fail(w, []byte("unsupported http method"))
		return
	}

	keys := strings.Split(r.URL.Path, "/")
	if len(keys) != 3 {
		log.Println("invalid url format")
		fail(w, []byte("invalid url format"))
		return
	}

	_, err := time.Parse("2006-01-02", keys[2])
	if err != nil {
		log.Println("invalid url format")
		fail(w, []byte("invalid url format"))
		return
	}

	rates, err := s.rate.GetByTime(keys[2])
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
		return
	}

	resp := make(map[string]interface{})
	resp["base"] = "EUR"
	resp["rates"] = rates

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
		return
	}

	success(w, jsonResp)
}

func (s *RateController) Analyze(w http.ResponseWriter, r *http.Request) {
	rates, err := s.rate.Analyze()
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
	}

	resp := make(map[string]interface{})
	resp["base"] = "EUR"
	resp["rates_analyze"] = rates
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		fail(w, []byte("ask your programmer for help"), http.StatusInternalServerError)
		return
	}

	success(w, jsonResp)
}

//helper methods
func success(w http.ResponseWriter, msg []byte, statusCode ...int) {
	if len(statusCode) == 0 {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}

func fail(w http.ResponseWriter, msg []byte, statusCode ...int) {
	if len(statusCode) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}
