package main

import (
	// "database/sql"
	// _ "mysql"
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetPredictionData(timestamp int, stockId string, predType string) []int {
	// connect to db
	sess, err := mgo.Dial("127.0.0.1")
	// try to connect to mongodb, return empty Prices array if we fail
	if err != nil {
		fmt.Printf("\nCould not connect to mongodb")
		return []int{}
	}
	defer sess.Close()
	col := sess.DB("data").C("predictions")
	// get data
	var results []PriceEntry
	err = col.Find(bson.M{"timestamp":timestamp}).All(&results)

	// find the right entry and return it
	for _, result := range results {
		for _, stock := range result.StockEntry {
			if stock.StockId == stockId {
				for _, pred := range stock.PredEntry {
					if pred.PredType == predType {
						return pred.Predictions
					}
				}
			}
		}
	}
	return []int{}
}

func GetPriceData(start int, end int, interval int) bool {
	// get data

	// parse data, save in int array

	// return parsed data
	return true
}

func AddPriceData(ts int, price int) bool {
	// validate session
	return true
}

func AddPredictionData(prediction PriceEntry) bool {
	sess, err := mgo.Dial("127.0.0.1")
	// try to connect to mongodb, return empty Prices array if we fail
	if err != nil {
		fmt.Printf("\nCould not connect to mongodb")
		return false
	}
	defer sess.Close()
	col := sess.DB("data").C("predictions")
	// loop through preds int array, make PredictionEntry object
	err = col.Insert(prediction)
	if err != nil {
		fmt.Printf("\nFailed to add document to collection in database")
		return false
	} else {
		return true
	}
}