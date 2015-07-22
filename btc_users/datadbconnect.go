package main

import (
	// "database/sql"
	// _ "mysql"
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	// "strconv"
	"gopkg.in/mgo.v2/bson"
)

func GetPredictionData(start int, end int) []PriceEntry {
	// connect to db
	sess, err := mgo.Dial("127.0.0.1")
	// try to connect to mongodb, return empty Prices array if we fail
	if err != nil {
		fmt.Printf("\nCould not connect to mongodb")
		return []PriceEntry{}
	}
	defer sess.Close()
	col := sess.DB("data").C("predictions")
	// get data
	var results []PriceEntry

	err = col.Find(bson.M{"timestamp":bson.M{"$lt":end, "$gt":start}}).Sort("timestamp").All(&results)
	return results
}

func GetPriceData(start int, end int) []Price {
	sess, err := mgo.Dial("127.0.0.1")
	// try to connect to mongodb, return empty Prices array if we fail
	if err != nil {
		fmt.Printf("\nCould not connect to mongodb")
		return []Price{}
	}
	defer sess.Close()
	col := sess.DB("data").C("prices")
	// get data
	var results []Price
	// query database, check for errors
	err = col.Find(bson.M{"timestamp":bson.M{"$lt":end, "$gt":start}}).Sort("timestamp").All(&results)
	if err != nil {
		fmt.Println("Error querying database")
	}
	return results
}

func AddPriceData(price Price) bool {
	sess, err := mgo.Dial("127.0.0.1")
	// try to connect to mongodb, return empty Prices array if we fail
	if err != nil {
		fmt.Printf("\nCould not connect to mongodb")
		return false
	}
	defer sess.Close()
	col := sess.DB("data").C("prices")
	// loop through preds int array, make PredictionEntry object
	err = col.Insert(price)
	if err != nil {
		fmt.Printf("\nFailed to add document to collection in database")
		return false
	} else {
		return true
	}
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