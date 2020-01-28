package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type Trade struct {
	Security   string  `csv:"Security Type"`
	Symbol     string  `csv:"Symbol"`
	Time       string  `csv:"Time"`
	Date       string  `csv:"Date"`
	Quantity   int64   `csv:"Quantity"`
	Price      float64 `csv:"Price"`
	Action     string  `csv:"Action"`
	Commission float64 `csv:"Commission"`
	Underlying string  `csv:"Underlying"`
	Account    string  `csv:"Account"`
}

type TradeOutput struct {
	TradeDate    string `csv:"TradeDate"`
	TradeTime    string `csv:"TradeTime"`
	BuySell      string `csv:"Buy/Sell"`
	AssetClass   string `csv:"AssetClass"`
	Symbol       string `csv:"Symbol"`
	Quantity     string `csv:"Quantity"`
	TradePrice   string `csv:"TradePrice"`
	IBCommission string `csv:"IBCommission"`
	Underlying   string `csv:"UnderlyingSymbol"`
	NetCash      string `csv:"NetCash"`
}

type Multiplier struct {
	Underlying string  `json:"underlying"`
	Multiplier float64 `json:"multiplier"`
}

type Configuration struct {
	InputPath    string       `json:"input_path"`
	OutputPath   string       `json:"output_path"`
	OutputPrefix string       `json:"output_prefix"`
	Multipliers  []Multiplier `json:"multipliers"`
}

var multiplierMap map[string]float64
var configuration Configuration
var flagDate string
var inputFileName string

func init() {
	flag.StringVar(&flagDate, "date", "", "tradelog export file date YYYYMMdd, default today")
	flag.StringVar(&flagDate, "d", "", "tradelog export file date YYYYMMdd, default today")
	flag.StringVar(&inputFileName, "file", "", "tws export file name, default is trades.YYYYMMdd.csv")
	flag.StringVar(&inputFileName, "f", "", "tws export file name, default is trades.YYYYMMdd.csv")
}

func main() {
	flag.Parse()

	configFile, _ := os.Open("config.json")
	configDecoder := json.NewDecoder(configFile)
	err := configDecoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	// initialize multipliers
	multiplierMap = make(map[string]float64)

	for _, m := range configuration.Multipliers {
		multiplierMap[m.Underlying] = m.Multiplier
	}

	// expect to use default tws export name
	t := time.Now()
	fileDate := t.Format("20060102")

	if len(flagDate) > 0 {
		fileDate = flagDate
	}

	if len(inputFileName) == 0 {
		inputFileName = "trades." + fileDate + ".csv"
	}
	inputFile := filepath.FromSlash(configuration.InputPath+"/") + inputFileName

	tradeFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer tradeFile.Close()

	trades := []*Trade{}

	if err := gocsv.UnmarshalFile(tradeFile, &trades); err != nil {
		panic(err)
	}

	TradeOutputMap := make(map[string][]TradeOutput)

	for _, t := range trades {

		if t.Security == "BAG" {
			continue
		}

		m := multiplierMap[t.Underlying]
		if m == 0.0 {
			m = 100.0
		}

		var qu int64
		var bs string
		var nc float64

		price := t.Price
		if t.Underlying == "EUR" {
			price *= 100
		}

		if t.Action == "SLD" {
			qu = -1
			bs = "SELL"
			nc = 1
		} else {
			qu = 1
			bs = "BUY"
			nc = -1
		}

		o := TradeOutput{}

		layout := "20060102 15:04:05 MST"

		timeString := t.Date + " " + t.Time + " UTC"

		pt, _ := time.Parse(layout, timeString)

		loc, _ := time.LoadLocation("America/Los_Angeles")
		pt = pt.In(loc)

		o.TradeDate = pt.Format("20060102")
		o.TradeTime = pt.Format("150405")

		o.BuySell = bs
		o.AssetClass = t.Security
		o.Symbol = t.Symbol
		o.Quantity = strconv.FormatInt(qu*t.Quantity, 10)
		o.TradePrice = strconv.FormatFloat(price, 'f', 2, 64)
		o.IBCommission = strconv.FormatFloat(t.Commission*-1, 'f', 2, 64)
		o.NetCash = strconv.FormatFloat((float64(t.Quantity)*m*t.Price*nc)-t.Commission, 'f', 2, 64)
		o.Underlying = t.Underlying

		TradeOutputMap[t.Account] = append(TradeOutputMap[t.Account], o)

	}

	for k, v := range TradeOutputMap {
		outputFile := filepath.FromSlash(configuration.OutputPath+"/") +
			configuration.OutputPrefix + k + "_" + inputFileName

		out, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		err = gocsv.MarshalFile(&v, out)
		if err != nil {
			panic(err)
		}
	}
}
