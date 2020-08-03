package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

func TestReferenceCriteriaManager_ForParams_onMissing(t *testing.T) {
	manager := ReferenceCriteriaManager{
		factories: []ReferenceCriterionFactory{
			&RandomWeightedReferenceCriterionManager{},
			&RandomUniformReferenceCriterionManager{},
			&ImportanceRatioReferenceCriterionManager{},
		},
	}
	var m interface{} = utils.Map{}
	result := manager.ForParams(&m)
	if _, ok := result.(*RandomWeightedReferenceCriterionProvider); !ok {
		t.Errorf("expected RandomWeightedReferenceCriterionProvider (the first one) when nothing provided")
	}
	manager = ReferenceCriteriaManager{
		factories: []ReferenceCriterionFactory{
			&ImportanceRatioReferenceCriterionManager{},
			&RandomUniformReferenceCriterionManager{},
			&RandomWeightedReferenceCriterionManager{},
		},
	}
	result = manager.ForParams(&m)
	if _, ok := result.(*ImportanceRatioReferenceCriterionProvider); !ok {
		t.Errorf("expected ImportanceRatioReferenceCriterionManager (the first one) when nothing provided")
	}
}

func TestReferenceCriteriaManager_ForParams_noManagerDeclared(t *testing.T) {
	manager := ReferenceCriteriaManager{
		factories: []ReferenceCriterionFactory{},
	}
	defer utils.ExpectError(t, "no ReferenceCriterionFactory has been declared")()
	var params interface{} = utils.Map{"referenceCriterionType": "importanceRatio"}
	manager.ForParams(&params)
}
func TestReferenceCriteriaManager_ForParams_onInvalid(t *testing.T) {
	manager := ReferenceCriteriaManager{
		factories: []ReferenceCriterionFactory{
			&ImportanceRatioReferenceCriterionManager{},
			&RandomUniformReferenceCriterionManager{},
			&RandomWeightedReferenceCriterionManager{},
		},
	}
	defer utils.ExpectError(t, "no reference criterion factory found for 'abcd' in [importanceRatio randomUniform randomWeighted]")()
	var params interface{} = utils.Map{"referenceCriterionType": "abcd"}
	manager.ForParams(&params)
}

func TestReferenceCriteriaManager_ForParams_validFetching(t *testing.T) {
	factories := []ReferenceCriterionFactory{
		&RandomWeightedReferenceCriterionManager{},
		&RandomUniformReferenceCriterionManager{},
		&ImportanceRatioReferenceCriterionManager{},
	}
	manager := ReferenceCriteriaManager{factories: factories}
	for _, f := range factories {
		var params interface{} = utils.Map{"referenceCriterionType": f.Identifier()}
		result := manager.ForParams(&params)
		actualType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf(f.NewProvider())
		if actualType != expectedType {
			t.Errorf("expected '%s' type of provider, got '%s'", expectedType, actualType)
		}
	}
}
