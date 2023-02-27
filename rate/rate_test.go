package rate

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRate_LoadToDB(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath := "./test.db"
	defer os.Remove(dbPath)

	rate := NewRate(xmlPath, dbPath)
	err := rate.LoadToDB()
	assert.NoError(t, err)
	assert.FileExists(t, dbPath)

	var count int64
	rate.GetDB().QueryRow("select count(*) from hist").Scan(&count)
	t.Log(count)
	assert.True(t, count > 0)
}

func TestRate_GetByTime(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath := "./test.db"
	defer os.Remove(dbPath)

	rate := NewRate(xmlPath, dbPath)
	loadToTestDB(rate.GetDB())
	//t.Log(rate)

	time := "2022-07-06"
	rates, err := rate.GetByTime(time)
	assert.NoError(t, err)
	expected := map[string]float64{
		"AUD": 1.4961,
		"BGN": 1.9558,
		"USD": 1.0177,
		"ZAR": 17.0246,
	}
	assert.Equal(t, expected, rates)

	time = "2022-07-05"
	rates, err = rate.GetByTime(time)
	assert.NoError(t, err)
	expected = map[string]float64{
		"AUD": 1.518,
		"BGN": 1.9558,
		"USD": 1.029,
		"ZAR": 16.9143,
	}
	assert.Equal(t, expected, rates)

	time = "2022-07-04"
	rates, _ = rate.GetByTime(time)
	assert.Empty(t, rates)
}

func TestRate_GetLatest(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath := "./test.db"
	defer os.Remove(dbPath)
	rate := NewRate(xmlPath, dbPath)
	loadToTestDB(rate.GetDB())

	rates, err := rate.GetLatest()
	assert.NoError(t, err)
	//t.Log(rate)

	expected := map[string]float64{
		"AUD": 1.4961,
		"BGN": 1.9558,
		"USD": 1.0177,
		"ZAR": 17.0246,
	}
	assert.Equal(t, expected, rates)
}

func TestRate_Analyze(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath := "./test.db"
	defer os.Remove(dbPath)
	rate := NewRate(xmlPath, dbPath)
	loadToTestDB(rate.GetDB())

	rates, err := rate.Analyze()
	assert.NoError(t, err)
	//t.Log(rate)

	expected := map[string]map[string]float64{
		"AUD": {
			"min": 1.4961,
			"max": 1.518,
			"avg": (tests[3].rate + tests[7].rate) / 2,
		},
		"BGN": {
			"min": 1.9558,
			"max": 1.9558,
			"avg": (tests[2].rate + tests[6].rate) / 2,
		},
		"USD": {
			"min": 1.0177,
			"max": 1.029,
			"avg": (tests[1].rate + tests[5].rate) / 2,
		},
		"ZAR": {
			"min": 16.9143,
			"max": 17.0246,
			"avg": (tests[0].rate + tests[4].rate) / 2,
		},
	}
	assert.Equal(t, expected, rates)
}

var tests = []struct {
	time     string
	currency string
	rate     float64
}{
	{
		time:     "2022-07-06",
		currency: "ZAR",
		rate:     17.0246,
	},
	{
		time:     "2022-07-06",
		currency: "USD",
		rate:     1.0177,
	},
	{
		time:     "2022-07-06",
		currency: "BGN",
		rate:     1.9558,
	},
	{
		time:     "2022-07-06",
		currency: "AUD",
		rate:     1.4961,
	},
	{
		time:     "2022-07-05",
		currency: "ZAR",
		rate:     16.9143,
	},
	{
		time:     "2022-07-05",
		currency: "USD",
		rate:     1.029,
	},
	{
		time:     "2022-07-05",
		currency: "BGN",
		rate:     1.9558,
	},
	{
		time:     "2022-07-05",
		currency: "AUD",
		rate:     1.518,
	},
}

func loadToTestDB(db *sql.DB) {
	db.Exec(ddl)
	db.Exec("delete from hist")
	for _, test := range tests {
		db.Exec("insert into hist(time, currency, rate) values (?, ?, ?)", test.time, test.currency, test.rate)
	}
}
