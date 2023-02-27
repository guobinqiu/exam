package rate

import (
	"database/sql"
)

const ddl = `
		CREATE TABLE if NOT EXISTS hist (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "time" TEXT,
        "currency" TEXT,
        "rate" REAL);
	`

type Rate struct {
	db      *MyDB
	xmlPath string
}

func NewRate(xmlPath string, dbPath string) *Rate {
	return &Rate{
		db:      NewMyDB(dbPath),
		xmlPath: xmlPath,
	}
}

func (r *Rate) LoadToDB() error {
	xmlReader := NewXMLReader(r.xmlPath)
	cube, err := xmlReader.GetCube()
	if err != nil {
		return err
	}

	if _, err := r.GetDB().Exec(ddl); err != nil {
		return err
	}

	//not sure if the xml content is incremental by day, so temporarily run delete for now
	if _, err := r.GetDB().Exec("DELETE FROM hist"); err != nil {
		return err
	}

	for _, cubeTime := range cube.CubeTime {
		for _, cubeCurrency := range cubeTime.CubeCurrency {
			//fmt.Println(cubeTime.Time, cubeCurrency.Currency, cubeCurrency.Rate)
			_, err := r.GetDB().Exec("INSERT INTO hist(time, currency, rate) VALUES (?, ?, ?)", cubeTime.Time, cubeCurrency.Currency, cubeCurrency.Rate)
			if err != nil {
				return err
			}
		}
	}

	r.db.Close()

	return nil
}

func (r *Rate) GetLatest() (map[string]float64, error) {
	time, err := r.getLatestTime()
	if err != nil {
		return nil, err
	}
	return r.GetByTime(time)
}

func (r *Rate) GetByTime(time string) (map[string]float64, error) {
	rows, err := r.GetDB().Query("SELECT currency, rate FROM hist WHERE time = ? ORDER BY currency", time)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var currency string
	var rate float64

	rates := make(map[string]float64)

	for rows.Next() {
		if err := rows.Scan(&currency, &rate); err != nil {
			return nil, err
		}
		rates[currency] = rate
	}

	return rates, nil
}

func (r *Rate) Analyze() (map[string]map[string]float64, error) {
	rows, err := r.GetDB().Query("SELECT currency, min(rate), max(rate), avg(rate) FROM hist GROUP BY currency")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var currency string
	var min float64
	var max float64
	var avg float64

	rates := make(map[string]map[string]float64)

	for rows.Next() {
		if err := rows.Scan(&currency, &min, &max, &avg); err != nil {
			return nil, err
		}
		aggr := make(map[string]float64, 3)
		aggr["min"] = min
		aggr["max"] = max
		aggr["avg"] = avg
		rates[currency] = aggr
	}

	return rates, nil
}

func (r *Rate) getLatestTime() (string, error) {
	var time string
	row := r.GetDB().QueryRow("SELECT time FROM hist ORDER BY time DESC LIMIT 1")
	if err := row.Scan(&time); err != nil {
		return "", err
	}
	return time, nil
}

func (r *Rate) GetDB() *sql.DB {
	return r.db.GetDB()
}
