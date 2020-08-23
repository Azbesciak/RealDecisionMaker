package main

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/anchoring"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/criteria-concealment"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/criteria-mixing"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/criteria-omission"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/fatigue"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/biases/preference-reversal"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/aspect-elimination"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/majority"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/choquet"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/electreIII"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/owa"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/weighted-sum"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-ordering"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"log"
	"net/http"
	"sync"
)

//go:generate easytags $GOFILE json:camel

var increasingSatisfactionLevels = []satisfaction_levels.SatisfactionLevelsSource{
	&satisfaction_levels.IdealIncreasingMulCoefficientSatisfaction,
	&satisfaction_levels.IdealAdditiveCoefficientSatisfaction,
	&satisfaction_levels.IncreasingThresholds,
}

var decreasingSatisfactionLevels = []satisfaction_levels.SatisfactionLevelsSource{
	&satisfaction_levels.IdealDecreasingMulCoefficientSatisfaction,
	&satisfaction_levels.IdealSubtrCoefficientSatisfaction,
	&satisfaction_levels.DecreasingThresholds,
}

var decreasingSatisfactionLevelsUpdates = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.Thresholds:         &satisfaction_levels.DecreasingThresholds,
		satisfaction_levels.IdealDecreasingMul: &satisfaction_levels.IdealDecreasingMulCoefficientSatisfaction,
		satisfaction_levels.IdealSubtractive:   &satisfaction_levels.IdealSubtrCoefficientSatisfaction,
	},
}

var increasingSatisfactionLevelsUpdates = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.Thresholds:         &satisfaction_levels.IncreasingThresholds,
		satisfaction_levels.IdealIncreasingMul: &satisfaction_levels.IdealIncreasingMulCoefficientSatisfaction,
		satisfaction_levels.IdealAdditive:      &satisfaction_levels.IdealAdditiveCoefficientSatisfaction,
	},
}

var funcs = model.PreferenceFunctions{
	Functions: []model.PreferenceFunction{
		&weighted_sum.WeightedSumPreferenceFunc{},
		&owa.OWAPreferenceFunc{},
		&electreIII.ElectreIIIPreferenceFunc{},
		&choquet.ChoquetIntegralPreferenceFunc{},
		aspect_elimination.NewAspectEliminationHeuristic(increasingSatisfactionLevels, utils.RandomBasedSeedValueGenerator),
		majority.NewMajority(utils.RandomBasedSeedValueGenerator, []majority.DrawResolver{
			&majority.DrawAllowedResolver{},
			&majority.CurrentIsWinnerDrawResolver{},
			&majority.NewerIsWinnerResolver{},
			&majority.RandomWinnerResolver{},
		}),
		satisfaction.NewSatisfaction(utils.RandomBasedSeedValueGenerator, decreasingSatisfactionLevels),
	},
}
var biasListeners = model.BiasListeners{
	Listeners: []model.BiasListener{
		&weighted_sum.WeightedSumBiasListener{},
		&owa.OwaBiasListener{},
		&electreIII.ElectreIIIBiasLIstener{},
		&choquet.ChoquetIntegralBiasListener{},
		aspect_elimination.NewAspectEliminationBiasListener(increasingSatisfactionLevelsUpdates),
		&majority.MajorityBiasListener{},
		satisfaction.NewSatisfactionBiasListener(decreasingSatisfactionLevelsUpdates),
	},
}

var referenceCriterionManager = *reference_criterion.NewReferenceCriteriaManager(
	[]reference_criterion.ReferenceCriterionFactory{
		&reference_criterion.ImportanceRatioReferenceCriterionManager{},
		&reference_criterion.RandomUniformReferenceCriterionManager{
			RandomFactory: utils.RandomBasedSeedValueGenerator,
		},
		&reference_criterion.RandomWeightedReferenceCriterionManager{
			RandomFactory: utils.RandomBasedSeedValueGenerator,
		},
	},
)

var criteriaOrdering = []criteria_ordering.CriteriaOrderingResolver{
	&criteria_ordering.WeakestCriteriaOrderingResolver{},
	&criteria_ordering.StrongestCriteriaOrderingResolver{},
	&criteria_ordering.RandomCriteriaOrderingResolver{
		Generator: utils.RandomBasedSeedValueGenerator,
	},
	&criteria_ordering.WeakestByProbabilityCriteriaOrderingResolver{
		Generator: utils.RandomBasedSeedValueGenerator,
	},
	&criteria_ordering.StrongestByProbabilityCriteriaOrderingResolver{
		WeakestByProbability: &criteria_ordering.WeakestByProbabilityCriteriaOrderingResolver{
			Generator: utils.RandomBasedSeedValueGenerator,
		},
	},
}

var biases = model.BiasMap{
	anchoring.BiasName: anchoring.NewAnchoring(
		[]anchoring.AnchoringEvaluator{
			&anchoring.LinearAnchoringEvaluator{},
			&anchoring.ExpFromZeroAnchoringEvaluator{},
		},
		[]anchoring.ReferencePointsEvaluator{
			&anchoring.IdealReferenceAlternativeEvaluator{},
			&anchoring.NadirReferenceAlternativeEvaluator{},
		},
		[]anchoring.AnchoringApplier{
			&anchoring.InlineAnchoringApplier{},
		},
	),
	criteria_concealment.BiasName: criteria_concealment.NewCriteriaConcealment(
		utils.RandomBasedSeedValueGenerator,
		referenceCriterionManager,
	),
	criteria_mixing.BiasName: criteria_mixing.NewCriteriaMixing(
		utils.RandomBasedSeedValueGenerator,
		referenceCriterionManager,
	),
	preference_reversal.BiasName: preference_reversal.NewPreferenceReversal(criteriaOrdering),
	criteria_omission.BiasName:   criteria_omission.NewCriteriaOmission(criteriaOrdering),
	fatigue.BiasName: fatigue.NewFatigue(
		utils.RandomBasedSeedValueGenerator,
		utils.RandomBasedSeedValueGenerator,
		[]fatigue.FatigueFunction{
			&fatigue.ExponentialFromZeroFatigue{},
			&fatigue.ConstFatigueFunction{},
		},
	),
}

type LazyFunctions func() *utils.Map

func Make(f LazyFunctions) LazyFunctions {
	var v *utils.Map
	var once sync.Once
	return func() *utils.Map {
		once.Do(func() {
			v = f()
			f = nil
		})
		return v
	}
}

var funcRequirements = Make(func() *utils.Map {
	return funcs.FetchParameters()
})

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
	decision := dm.MakeDecision(funcs, biasListeners, &biases, utils.RandomBasedSeedValueGenerator)
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

func functionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, funcRequirements())
}

func main() {
	r := gin.Default()
	// Dont worry about this line just yet, it will make sense in the Dockerise bit!
	r.Use(static.Serve("/", static.LocalFile("./web", true)))
	r.Use(cors.Default())
	api := r.Group("/api")
	api.POST("/decide", decideHandler)
	api.GET("/preferenceFunctions", functionsHandler)
	log.Fatal(r.Run())
}
