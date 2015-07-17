package main

// heirarchical prediction object stored in database
type PriceEntry struct {
	Timestamp		int
	StockEntry		[]StockEntry
}

type StockEntry struct {
	StockId			string
	PredEntry		[]PredictionEntry
}

type PredictionEntry struct {
	PredType		string
	Predictions		[]float64
}
///////////////////////////////////////

// object returned in prediction API
type Prediction struct {
	Timestamp		int
	StockId			string
	PredictionType	string
	Predictions		[]float64
}

// heirarchical price object stored in database
type Price struct {
	Timestamp		int
	StockId			[]StockPrice
}

type StockPrice struct {
	Name			string
	Price			float64
	Volume			float64
}
///////////////////////////////////////

// object returned in price API
type SinglePrice struct {
	Timestamp		int
	Price			float64
	Volume			float64
}