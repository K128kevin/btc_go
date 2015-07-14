package main

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
	Predictions		[]int
}

type Prediction struct {
	Timestamp		int
	StockId			string
	PredictionType	string
	Predictions		[]int
}

type Price struct {
	Timestamp		int
	Prices			[]int
}