package rate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadXML(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	xml := NewXMLReader(xmlPath)
	b, _ := xml.readXML()
	//t.Log(string(b))
	assert.NotEmpty(t, b)
	assert.Contains(t, string(b), "Cube")
}

func TestGetCube(t *testing.T) {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	xml := NewXMLReader(xmlPath)
	cube, _ := xml.GetCube()
	//t.Log(cube)

	theCubeTime := &CubeTime{}
	for _, cubeTime := range cube.CubeTime {
		if cubeTime.Time == "2022-07-06" {
			theCubeTime = cubeTime
			break
		}
	}

	assert.Equal(t, theCubeTime.CubeCurrency[0].Currency, "USD")
	assert.Equal(t, theCubeTime.CubeCurrency[0].Rate, float32(1.0177))
	assert.Equal(t, theCubeTime.CubeCurrency[1].Currency, "JPY")
	assert.Equal(t, theCubeTime.CubeCurrency[1].Rate, float32(137.71))
}
