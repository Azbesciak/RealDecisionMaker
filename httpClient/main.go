package main

import (
	"../logic/choquet"
	"../logic/electreIII"
	"../logic/owa"
	"../logic/weighted-sum"
	"../model"
	"encoding/json"
	"github.com/go-errors/errors"
	"log"
	"net/http"
)

var weightedSumF = &weighted_sum.WeightedSumPreferenceFunc{}
var owaF = &owa.OWAPreferenceFunc{}
var eleF = &electreIII.ElectreIIIPreferenceFunc{}
var choquetF = &choquet.ChoquetIntegralPreferenceFunc{}
var funcs = model.PreferenceFunctions{Functions: []model.PreferenceFunction{weightedSumF, owaF, eleF, choquetF}}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dm model.DecisionMaker
	err := decoder.Decode(&dm)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		writeError(err, &dm, &w)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			writeError(e, &dm, &w)
		}
	}()
	decision := dm.MakeDecision(funcs)
	log.Printf("%#v", requestSuccess{dm, *decision})
	writeJSON(decision, &w)
}

func writeError(e interface{}, dm *model.DecisionMaker, w *http.ResponseWriter) {
	log.Println(errors.Wrap(e, 1).ErrorStack())
	switch v := e.(type) {
	case error:
		e = v.Error()
	}
	err := requestError{
		Error:   e,
		Request: dm,
	}
	(*w).WriteHeader(400)
	writeJSON(err, w)
}

func writeJSON(data interface{}, w *http.ResponseWriter) {
	bytes, _ := json.Marshal(data)
	toCamelCaseJSON(&bytes)
	(*w).Write(bytes)
}

type requestError struct {
	Error   interface{}
	Request interface{}
}

type requestSuccess struct {
	Request  interface{}
	Response interface{}
}

func toCamelCaseJSON(jsonBytes *[]byte) {
	length := len(*jsonBytes)
	if length < 3 {
		return
	}
	bracketObj, comma, quote := "{"[0], ","[0], "\""[0]
	twoBefore, oneBefore := (*jsonBytes)[0], (*jsonBytes)[1]
	for i := 2; i < length; i++ {
		current := (*jsonBytes)[i]
		if oneBefore == quote && (twoBefore == bracketObj || twoBefore == comma) && current <= 90 && current >= 64 {
			(*jsonBytes)[i] = current + 32
		}
		twoBefore = oneBefore
		oneBefore = current
	}
}

func main() {
	http.HandleFunc("/decide", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
