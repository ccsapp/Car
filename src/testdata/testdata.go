package testdata

import (
	_ "embed"
	"strings"
)

//go:embed exampleCar.json
var ExampleCar string

//go:embed exampleCar2.json
var ExampleCar2 string

//go:embed exampleCarDuplicate.json
var ExampleCarDuplicate string

//go:embed exampleCarWithDynamicData.json
var ExampleCarWithDynamicData string

const ExampleCarVinString = "WVWAA71K08W201030"
const ExampleCar2VinString = "WVWAA71K08W201031"

var ExampleCarVin = QuoteString(ExampleCarVinString)
var ExampleCarVinArray = ArrayString(ExampleCarVin)
var ExampleCar2Vin = QuoteString(ExampleCar2VinString)

//go:embed exampleNoJson.txt
var ExampleNoJson string

//go:embed exampleNoCar.json
var ExampleNoCar string

//go:embed exampleCarWrongEnum.json
var ExampleCarWrongEnum string

func QuoteString(unquoted string) string {
	return "\"" + unquoted + "\""
}

func ArrayString(values ...string) string {
	return "[" + strings.Join(values, ", ") + "]"
}
