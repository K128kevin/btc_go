package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"time"
)

var MAX_DATA = 500
var MIN_INTERVAL = 300
var MAX_INTERVAL = 1800
var EARLIEST_TS = 1400000000

func ValidateDataParams(w http.ResponseWriter, r *http.Request) (bool, int, int, int) {

	// get data from URL
    start := r.FormValue("start")
    end := r.FormValue("end")
    interval := r.FormValue("interval")
	vars := mux.Vars(r)
	stockId := vars["stockId"]

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

	fmt.Printf("\nStock ID: %s\nStart timestamp: %s\nEnd timestamp: %s\nInterval: %s\n", stockId, start, end, interval)

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

// handle requests to root ("/")
func PredictionGet(w http.ResponseWriter, r *http.Request) {
	err, startInt, endInt, intervalInt := ValidateDataParams(w, r)

    w.Header().Set("Access-Control-Allow-Origin", "*") // cors
    if !err {
    	fmt.Fprintln(w, "Failed to validate parameters provided")
    }
    // if we get here that means we have a valid start/end timestamp range with a valid interval
    // now time to query the database for that data and return it in json format
    if !GetPredictionData(startInt, endInt, intervalInt) {
    	fmt.Fprintf(w, "A valid start, end, and interval value were provided, but there was an unknown error querying the database :(")
    	return
    }

}

// handle requests to root ("/")
func PriceGet(w http.ResponseWriter, r *http.Request) {
	err, startInt, endInt, intervalInt := ValidateDataParams(w, r)

    w.Header().Set("Access-Control-Allow-Origin", "*") // cors
    if !err {
    	fmt.Fprintln(w, "Failed to validate parameters provided")
    }
    // if we get here that means we have a valid start/end timestamp range with a valid interval
    // now time to query the database for that data and return it in json format
    if !GetPriceData(startInt, endInt, intervalInt) {
    	fmt.Fprintf(w, "A valid start, end, and interval value were provided, but there was an unknown error querying the database :(")
    	return
    }

}


