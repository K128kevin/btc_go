package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"io"
	"io/ioutil"
	"strings"
)

var MAX_DATA = 500
var MIN_INTERVAL = 300
var MAX_INTERVAL = 3600
var EARLIEST_TS = 1400000000

func ValidateDataParams(w http.ResponseWriter, r *http.Request) (bool, int, int, int) {

	// get data from URL
    start := r.FormValue("start")
    end := r.FormValue("end")
    interval := r.FormValue("interval")

	// make sure all parameters are included
	returnVal := true
	if start == "" {
		fmt.Fprintln(w, "Please provide a start timestamp (&start=1234567890)")
		returnVal = false
	}
	if end == "" {
		fmt.Fprintln(w, "Please provide an end timestamp (&end=1234567890)")
		returnVal = false
	}
	if interval == "" {
		fmt.Fprintln(w, "Please provide an interval size in seconds (&interval=900)")
		returnVal = false
	}
	if returnVal == false {
		return returnVal, 0, 0, 0
	}

	// fmt.Printf("\nStock ID: %s\nStart timestamp: %s\nEnd timestamp: %s\nInterval: %s\n", stockId, start, end, interval)

    // convert start, end, and interval strings to ints
    startInt, err := strconv.Atoi(start)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }
    endInt, err := strconv.Atoi(end)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }
    intervalInt, err := strconv.Atoi(interval)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }
    // -----------------------------------------------

    // ERROR CHECKING
    t := time.Now()
	currentTs := t.Format("20060102150405")
	tsInt, err := strconv.Atoi(currentTs)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }
	if (endInt - startInt) / intervalInt > MAX_DATA {
		fmt.Fprintf(w, "Combination of range size and interval size yields too much data. Please decrease one or both of these values")
		returnVal = false
	}
	if intervalInt < MIN_INTERVAL || intervalInt > MAX_INTERVAL {
		fmt.Fprintf(w, "Interval size must be between %d and %d", MIN_INTERVAL, MAX_INTERVAL)
		returnVal = false
	}
	if startInt < EARLIEST_TS {
		fmt.Fprintf(w, "Start date is too early - earliest available data is at timestamp %d", EARLIEST_TS)
		returnVal = false
	}
	if endInt > tsInt + (intervalInt * MAX_DATA) {
		fmt.Fprintf(w, "End date is too far into the future (maximum end date is %d)", tsInt + (intervalInt * MAX_DATA))
		returnVal = false
	}
    // -----------------------------------------------
    return returnVal, startInt, endInt, intervalInt

}

func ValidatePredictionParams(w http.ResponseWriter, r *http.Request) (bool, int, int, string, string) {
	// get data from URL
    start := r.FormValue("start")
    end := r.FormValue("end")
	stockId := r.FormValue("stockId")
	predType := r.FormValue("predType")

	// make sure all parameters are included
	returnVal := true
	if start == "" {
		fmt.Fprintln(w, "Please provide a start timestamp (&start=1234567890)")
		returnVal = false
	}
	if end == "" {
		fmt.Fprintln(w, "Please provide an end timestamp (&end=1234567890)")
		returnVal = false
	}
	if stockId == "" {
		fmt.Fprintln(w, "Please provide a stock id (&stockId=MSFT)")
		returnVal = false
	}
	if predType == "" {
		fmt.Fprintln(w, "Please provide a prediction type (&predType=prediction1)")
		returnVal = false
	}
	if returnVal == false {
		return false, 0, 0, "", ""
	}



	// convert start, end, and interval strings to ints
    startInt, err := strconv.Atoi(start)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }
    endInt, err := strconv.Atoi(end)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }

	if startInt > endInt {
		fmt.Fprintln(w, "Start timestmap must be before (less than) end timestamp")
		return false, 0, 0, "", ""
	}

	t := time.Now()
	currentTs := t.Format("20060102150405")
	tsInt, err := strconv.Atoi(currentTs)
    if err != nil {
        fmt.Println(err)
        returnVal = false
    }

    if endInt > tsInt || startInt < EARLIEST_TS {
		fmt.Fprintf(w, "Either start date is too early or end date is too late")
		fmt.Fprintf(w, "\nEarliest available data is at timestamp %d", EARLIEST_TS)
		fmt.Fprintf(w, "\nLatest available data is at timestamp %d (current timestamp)", tsInt)
		return false, 0, 0, "", ""
    }

	return true, startInt, endInt, stockId, predType
}

// Creates array of PriceEntry structs and converts it to json, which is then returned as string (Fprintf)
// struct is populated by querying database (mongodb) using params in request URL
func PredictionGet(w http.ResponseWriter, r *http.Request) {
	// first, check for a valid auth header
	headerToken := r.Header.Get(authTokenKey)
	_, ok := sessions[headerToken]
	if headerToken == "" || !ok {
		var resp LoginResponse
		resp.Error = true
		resp.Message = "Missing or invalid auth token header"
		fmt.Printf("\n%s", headerToken)
		retVal, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, string(retVal))
		fmt.Printf("\n%s", string(retVal))
		return
	}
	err, startInt, endInt, stockId, predType := ValidatePredictionParams(w, r)

    w.Header().Set("Access-Control-Allow-Origin", "*") // cors
    if !err {
    	fmt.Fprintln(w, "\nFailed to validate parameters provided")
    	return
    }
    // if we get here that means we have a valid start/end timestamp range with a valid interval
    // now time to query the database for that data and return it in json format
    var pricey []PriceEntry
    var preddy []Prediction
    pricey = GetPredictionData(startInt, endInt) // connect to database and get all PriceEntry structs in timestamp range
    preddy = PriceEntryToPrediction(pricey, stockId, predType) // convert this array to an array of Prediction structs
    retVal, _ := json.Marshal(preddy) // convert to json string and return this in response
    fmt.Fprintf(w, "%s", retVal)
}

func PriceEntryToPrediction(pe []PriceEntry, sid string, pt string) []Prediction {
	// find the right entry and return it
	var i int
	var preddy []Prediction // array to be returned
	for i = 0; i < len(pe); i++ { // loop through every PriceEntry to convert it to a Prediction based on StockId and PredType
		// create temporary new Prediction object, which will be added to array of Predictions that will be returned
		var tempPred Prediction
		tempPred.Timestamp = pe[i].Timestamp
		tempPred.StockId = sid
		tempPred.PredictionType = pt
		// now we have to populate its predictions
		// so we have to find the predictions for the given StockId and PredType in the PriceEntry
		for _, se := range pe[i].StockEntry { // loop through each StockEntry within each PredEntry
			if se.StockId == sid { // only continue when StockId matches given sid
				for _, pre := range se.PredEntry { // loop through each Pred Entry where StockId was found to be a match
					if pre.PredType == pt { // only continue when PredType matches given pt
						tempPred.Predictions = pre.Predictions
					}
				}
			}
		}
		if len(tempPred.Predictions) > 0 {
			preddy = append(preddy, tempPred)
		}
	}
	return preddy // CompressPredictionArray(preddy, interval)
}

func PriceGet(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nStarting PriceGet...")
	err, startInt, endInt, intervalInt := ValidateDataParams(w, r)
	stockId := r.FormValue("stockId")
	if stockId == "" {
		fmt.Fprintln(w, "Please provide a stockId (&stockId=MSFT)")
		return
	}

    w.Header().Set("Access-Control-Allow-Origin", "*") // cors
    if !err {
    	fmt.Fprintln(w, "Failed to validate parameters provided")
    }

    // if we get here that means we have a valid start/end timestamp range with a valid interval
    // now time to query the database for that data and return it in json format
    
    // first query db and get Price array
    var prices []Price
	fmt.Printf("\nAbout to get price data...")
    prices = GetPriceData(startInt, endInt)
	fmt.Printf("\nGot price data! Timestamp: %d", prices[0].Timestamp);
    // now compress prices based on interval
    var smallerPrices []SinglePrice
	fmt.Printf("\nAbout to compress prices...")
    smallerPrices = CompressPrices(prices, stockId, intervalInt)
	fmt.Printf("\nPrices compressed!")
	var i int
	for i = 0; i < len(smallerPrices); i++ {
		fmt.Printf("\nPrice: %f", smallerPrices[i].Price)
	}
    // finally, jsonify it and print to responsewriter
    retVal, _ := json.Marshal(smallerPrices)
    fmt.Fprintf(w, "%s", retVal)
}

func CompressPrices(prices []Price, stockId string, interval int) []SinglePrice {
	if len(prices) < 1 {
		fmt.Printf("CompressPrices - no prices given, returning empty SinglePrice array")
		return []SinglePrice{}
	}
	count := 1.0
	priceSum := 0.0
	volSum := 0.0
	start := prices[0].Timestamp
	var newPrices []SinglePrice
	var tempPrice SinglePrice
	var i int
	for i = 0; i < len(prices); i++ {
		// find stockId in Price object, add its price to sum
		inner: for _, val := range prices[i].StockId {
			if val.Name == stockId {
				priceSum += val.Price
				volSum += val.Volume
				break inner
			}
		}
		if prices[i].Timestamp - start >= interval || i == len(prices) - 1 {
			start = prices[i].Timestamp // set new start point to current TS
			priceSum /= count // actually calculate average
			volSum /= count
			tempPrice.Timestamp = start + interval // set timestamp for this price average
			tempPrice.Volume = volSum
			tempPrice.Price = priceSum
			// add price to return array
			newPrices = append(newPrices, tempPrice)
			// reset sum and count so we can start over
			priceSum = 0
			volSum = 0
			count = 0
		} 
		count++
	}
	return newPrices
}

func PriceAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // cors
	var price Price
	// parse request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	bodyString := string(body)
	fmt.Printf("\nData to be parsed and added:\n\n%s\n\n", bodyString)

	if strings.Contains(bodyString, "\\") {
		bodyString = strings.Replace(bodyString, "\\", "", -1)
		bodyString = strings.TrimRight(bodyString, "\"")
		bodyString = strings.TrimLeft(bodyString, "\"")
		fmt.Printf("Removed excape chars, should be parseable json now: \n%s", bodyString)
	}

	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	if err := json.Unmarshal([]byte(bodyString), &price); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	fmt.Printf("\nAdding to database...\n")
	fmt.Printf("\nNumber of stocks: %d", len(price.StockId))
	// add price object to database
	if AddPriceData(price) {
		fmt.Fprintf(w, "Data added successfully! :)")
		fmt.Printf("\nData added successfully")
	} else {
		fmt.Fprintf(w, "Failed to add data :(")
		fmt.Printf("Failed to add data")
	}
}

// add json data in body to mongodb list of data points
func PredictionAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // cors
	var prediction PriceEntry
	// parse request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	if err := json.Unmarshal(body, &prediction); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	// create a heirarchical prediction object and enter it into the database
	if AddPredictionData(prediction) {
		fmt.Fprintf(w, "Data added successfully! :)")
	} else {
		fmt.Fprintf(w, "Failed to add data :(")
	}
}
