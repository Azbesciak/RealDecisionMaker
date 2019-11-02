package main

import (
	"github.com/Azbesciak/RealDecisionMaker/logic/choquet"
	"github.com/Azbesciak/RealDecisionMaker/logic/electreIII"
	"github.com/Azbesciak/RealDecisionMaker/logic/owa"
	"github.com/Azbesciak/RealDecisionMaker/logic/weighted-sum"
	"github.com/Azbesciak/RealDecisionMaker/model"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"log"
	"net/http"
)

//go:generate easytags $GOFILE json:camel

var weightedSumF = &weighted_sum.WeightedSumPreferenceFunc{}
var owaF = &owa.OWAPreferenceFunc{}
var eleF = &electreIII.ElectreIIIPreferenceFunc{}
var choquetF = &choquet.ChoquetIntegralPreferenceFunc{}
var funcs = model.PreferenceFunctions{Functions: []model.PreferenceFunction{weightedSumF, owaF, eleF, choquetF}}

func decideHandler(c *gin.Context) {
	var dm model.DecisionMaker
	if err := c.ShouldBindJSON(&dm); err != nil {
		writeError(err, &dm, c)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			writeError(e, &dm, c)
		}
	}()
	decision := dm.MakeDecision(funcs)
	log.Printf("%#v", requestSuccess{dm, *decision})
	writeJSON(decision, c)
}

func writeError(e interface{}, dm *model.DecisionMaker, c *gin.Context) {
	log.Println(errors.Wrap(e, 1).ErrorStack())
	switch v := e.(type) {
	case error:
		e = v.Error()
	}
	err := requestError{
		Error:   e,
		Request: dm,
	}
	c.JSON(http.StatusBadRequest, err)
}

func writeJSON(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, data)
}

type requestError struct {
	Error   interface{} `json:"error"`
	Request interface{} `json:"request"`
}

type requestSuccess struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}

func main() {
	r := gin.Default()
	r.POST("/decide", decideHandler)
	log.Fatal(r.Run())
}
