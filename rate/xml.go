package rate

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type RootNode struct {
	Cube *Cube `xml:"Cube"`
}

type Cube struct {
	CubeTime []*CubeTime `xml:"Cube"`
}

type CubeTime struct {
	CubeCurrency []*CubeCurrency `xml:"Cube"`
	Time         string          `xml:"time,attr"`
}

type CubeCurrency struct {
	Currency string  `xml:"currency,attr"`
	Rate     float32 `xml:"rate,attr"`
}

type XMLReader struct {
	url string
}

func NewXMLReader(url string) *XMLReader {
	return &XMLReader{
		url: url,
	}
}

func (r *XMLReader) GetCube() (*Cube, error) {
	b, err := r.readXML()
	if err != nil {
		return nil, err
	}
	root := RootNode{}
	err = xml.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}
	return root.Cube, nil
}

func (r *XMLReader) readXML() ([]byte, error) {
	resp, err := http.Get(r.url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	return body, nil
}
