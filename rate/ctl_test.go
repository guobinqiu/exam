package rate

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	xmlPath = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath  = "./rate.db"
)

func TestGetLatest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rates/latest", nil)
	w := httptest.NewRecorder()
	ctl := NewRateController(xmlPath, dbPath)
	ctl.GetLatest(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resp))
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestGetByTime(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rates/2022-07-06", nil)
	w := httptest.NewRecorder()
	ctl := NewRateController(xmlPath, dbPath)
	ctl.GetByTime(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resp))
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestGetByTime_invalidTime(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rates/2022-07", nil)
	w := httptest.NewRecorder()
	ctl := NewRateController(xmlPath, dbPath)
	ctl.GetByTime(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resp))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
	assert.Equal(t, string(resp), "invalid url format")
}

func TestGetByTime_invalidHttpMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/rates/2022-07-06", nil)
	w := httptest.NewRecorder()
	ctl := NewRateController(xmlPath, dbPath)
	ctl.GetByTime(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resp))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
	assert.Equal(t, string(resp), "unsupported http method")
}

func TestAnalyze(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rates/analyze", nil)
	w := httptest.NewRecorder()
	ctl := NewRateController(xmlPath, dbPath)
	ctl.Analyze(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resp))
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
