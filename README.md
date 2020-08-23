# RealDecisionMaker
Simulator of the real human-like decision maker.

# Why and what
In a real life we, as a Decision Makers (DM), are often trying to make the best possible decision - at least we say so.
However, world is not so simple; we are lazy, often tired (or think so, which is not always equal)
 or not able to make reliable comparison because of alternatives and criteria complexity.

There are multiple methods for decision-making process support, but most of them are suited to give clear and mathematically valid results.
 The result can be used as some model or reference, but as we know:
> models are good, some are useful

So basically we can assume this as some reference of recommendation for the Decision Maker.
This fails when we are trying to predict real human decisions, because these methods does not take 
into account mentioned earlier problems in decision-making, which are just a top of the ice berg.
The issue can be also that we are not looking for the best solution, but for the first satisfying one,
 or for example our preference changes during the process because of comparision order.

## Valuable resources
- [Thinking, fast and slow. Daniel Kahneman (2011)](https://www.amazon.com/Thinking-Fast-Slow-Daniel-Kahneman/dp/0374533555)
- [Thinking fast and slow - goto; conference 2019. Linda Rising](https://youtu.be/XjbTLIqnq-o)
- [Modeling Behavior-Realistic Artificial Decision-Makers to Test Preference-Based Multiple Objective Optimization Methods](https://www.research.manchester.ac.uk/portal/en/publications/modeling-behaviorrealistic-artificial-decisionmakers-to-test-preferencebased-multiple-objective-optimization-methods(61d4687a-307b-454b-85cd-a8f6a72b26da).html)
- [Decision Making: Factors that Influence Decision Making, Heuristics Used, and Decision Outcomes](http://www.inquiriesjournal.com/articles/180/decision-making-factors-that-influence-decision-making-heuristics-used-and-decision-outcomes)
- [Decision-Making and Cognitive Biases](https://www.researchgate.net/publication/301662722_Decision-Making_and_Cognitive_Biases)
- [Naturalistic Decision Making](http://www.decisionskills.com/uploads/5/1/6/0/5160560/klein_2008.pdf)
- [Human Factors Influencing Decision Making](https://apps.dtic.mil/dtic/tr/fulltext/u2/a351910.pdf)
- [Modeling Human and Organizational Behavior: Application to Military Simulations (1998)](https://www.nap.edu/read/6173/chapter/8)

# Solution architecture
DM needs to have some weighting function. 
These are located in [`lib/logic/preference-func`](lib/logic/preference-func); available are:

- [OWA](#owa)
- [Weighted sum](#weighted-sum)
- [Choquet integral](#choquet-integral)
- [Electre III](#electre-iii)

However, we often decide in other ways, making some assumptions or fulfilling previously set goals.
Therefore, it is also possible to chose not fully rational weighing - *heuristic*, which are:

- [aspect elimination heuristic](#aspect-elimination-heuristic)
- [majority heuristic](#majority-heuristic)
- [satisfaction heuristic](#satisfaction-heuristic)

You can find their implementation in [`lib/logic/limited-rationality`](lib/logic/limited-rationality).

Finally, our decisions are often not clear or rational.
 We sometime just don't want to decide, or we change our mind, depending on the current knowledge,
  which can be in the decision-making context, but does not influence considered alternatives.
To simulate these behaviors we have following biases, which implementation can be found in [`lib/logic/biases`](lib/logic/biases):

- [criteria omission](#criteria-omission)
- [criteria concealment](#criteria-concealment)
- [criteria mixing](#criteria-mixing)
- [preference reversal](#preference-reversal)
- [fatigue](#fatigue)
- [anchoring](#anchoring)

They are mostly based on the paper *Modeling Behavior-Realistic Artificial Decision-Makers to
 Test Preference-Based Multiple Objective Optimization Methods* (mentioned earlier)

# Methods descriptions and usage
This library can be used as a go module, but there is also a http client, which can be used via JSON interface.
> all fields are named in the `PascalCase` fashion because of the go naming convention and public visibility;
> in JSON these are in `camelCase` fashion, so for example in go field `PreferenceFunction` is named `preferenceFunction` in JSON input.

### Input parameters
The highest level abstraction is [`DecisionMaker`](lib/model/decision-maker.go), which is the model of our Decision Maker.
He has following properties:

One important assumption was to make the code stateless, so all necessary information is passed in the input, and all possibly necessary later - as an output.
Other was to make the process clear and repeatable, so every change in the parameters is also reported for the given method.
When the random value can occur, each occurrence can be configured with the `randomSeed` separately (default is 0). 

- `PreferenceFunction`  -  `string`,  name of used [preference function](#preference-methods) of [limited rationality heuristic](#limited-rationality).          
- `Biases`              -  [`BiasesParams`](lib/model/bias.go), definitions of used biases (`name`) and their params (`props`) with optional `disabled` state and `applyProbability` (default = 1, applied for sure).
 One bias can be used multiple times. *Order matters.*       
- `BiasApplyRandomSeed` -  `int`, seed used for random when calculating bias activation probability
- `KnownAlternatives`   -  [`[]AlternativeWithCriteria`](lib/model/alternative.go), definition of all known
 to decision maker alternatives (so also these which are not considered in current choice, but somehow can influence on the final result)
- `ChoseToMake`         -  [`[]Alternative`](lib/model/alternative.go), names of alternatives which are considered during this iteration             
- `Criteria`            -  [`Criteria`](lib/model/criterion.go), definitions for the criteria;
name and the type (`gain` [default] or `cost` [otherwise, if not gain]). Weights are defined specific for each preference function/heuristic.
 There is also optional `valuesRange` used in biases calculation (otherwise it is calculated based on criteria values)      
- `MethodParameters`    -  [`RawMethodParameters`](lib/model/decision-maker.go), parameters for the [preference function](#preference-methods) of [limited rationality heuristic](#limited-rationality).          

Each method or bias are called in `camelCase` fashion, so when you want to call OWA (as a preference function),
 you use `preferenceFunction: "owa"`, when Weighted sum = `weightedSum` etc.

### Output
As an output you receive 
- `Result` - [`AlternativesRanking`](lib/model/alternative.go), which is the final ranking of the preference.
 Order may matter, but it is not said, because some methods may cause incomparability.
 Therefore some kind of [Hasse Diagram](https://en.wikipedia.org/wiki/Hasse_diagram) is used, 
 with the modification that each alternative have the list of worse or equal next in order alternatives - outgoing connections in graph. 
- `Biases` - [`BiasesParams`](lib/model/bias.go), which are biases changes in the original data.
 Notice that each bias passes further modified data, so each bias might work on the changed dataset.        

### Examples
Examples are located in [`httpClient/examples`](./httpClient/examples).
Directory name describes used preference/heuristic and optionally biases.
Each dictionary contains request and response.

## Preference methods
Each preference method has their own parameters in `MethodParameters` field in `DecisionMaker`.

### OWA
- [method description](https://en.wikipedia.org/wiki/Type-1_OWA_operators)
- [implementation](lib/logic/preference-func/owa)

The simplest preference function. Here the simplest variant is implemented:
 - weights are not fixed to the criteria, nevertheless *each weight should have valid criterion name assigned*
 - criteria values for each alternative are sorted and zipped with sorted criteria weights, then multiplied
 - *only gain criteria are allowed*
 
##### Input parameters:
Method takes only weights for each criterion:
```json
{
    "weights": {
      "c1": 12,
      "c2": 1.1
    }
}
```
  
### Weighted sum
- [method description](https://en.wikipedia.org/wiki/Weighted_sum_model)
- [implementation](lib/logic/preference-func/weighted-sum)

Criteria values are not scaled in any way. Any scalling must be done by a proper weighting *on the client side*.
Method, same as [OWA](#owa), method takes only weights for each criterion.

### Choquet integral
- [method description](https://en.wikipedia.org/wiki/Choquet_integral)
- [implementation](lib/logic/preference-func/choquet)

In short, this preference method allows to consider criteria interaction.
There are two requirements:
- for each criteria combination weights are required (don't have to be sorted) - combinations are joined with `,` (comma).
Their combinations can be described as a [Power set](https://en.wikipedia.org/wiki/Power_set).
> For example, for criteria `1` and `2` we need weight for `1`, `2` and `1,2`.
- weights must be in range [0,1].
Again, input format is the same as in [OWA](#owa) and [Weighted sum](#weighted-sum)

### Electre III
- [method description](https://www.sciencedirect.com/science/article/abs/pii/S0926580510002050)
- [implementation](lib/logic/preference-func/electreIII)

The aim for this method was to provide some way to model incomparability between alternatives.

##### Method requires following parameters:
- `criteria` which contains thresholds for each criterion:
    - `q` - indistinguishability
    - `p` - preference
    - `v` - veto
    and `k` which is given criterion weight (voting power).
- `distilationFunction` - linear function with `a` and `b` params, as a default `a = -0.15` and `b = 0.3`, used in distillation process.

for example:
```json
{
  "criteria": {
    "c1": {"q": 1, "p": 2, "v": 4},
    "c2": {"q": 20.4, "p": 50, "v": 100} 
  },
  "distillationFunction": {"a": -0.1, "b": 0.2}
}
```

## Limited rationality
These methods are used in mutal exclusion with preference functions; their parameters are also passed via `methodParameters`.

### Majority heuristic
- [implementation](lib/logic/limited-rationality/majority) 

- 2 alternatives are compared at given time on each criterion.
- Alternative receives points equal to criterion weight, when its criterion value is better that for other alternative
- The one which had more points passes further and is compared with the next one
- The final ranking is created by reversing drop out order:
    - last (which recently won) alternative is the best,
    - the one which was worse in the last comparison is the second, etc.

##### Input parameters:
- `Weights` - Weights for each criterion
- `CurrentChoice`  - optional, id of currently chosen alternative. It does not have to be in `consideredAlternatives`,
 but need to be known alternative. This alternative will be the first in comparison, also will occur in the final result.
- `RandomAlternativesOrdering` - whether alternatives should be shuffled before comparison (`CurrentChoise` will be still the first one)
- `RandomSeed` - seed for random, useful when alternatives are shuffled.
- `DrawResolution` - which alternative should win in case of draw, possible values are:
    - `allow` (draws are allowed)
    - `current` (earlier in comparison or currently winning one)
    - `newer`
    - `random`

For each alternative the method is returning recent evaluation result, which is
- `Value` - last comparison value
- `ComparedWith` - name of the alternative which was better or equal in the last comparison (not provided for the best/first in ranking)
- `ComparedAlternativeValue` - value of the better alternative in the recent comparison (not provided for the best/first in ranking)

### Aspect elimination heuristic
- [implementation](lib/logic/limited-rationality/aspect-elimination)

For each criterion
- Each alternative is checked whether it meets given threshold value
    - if not, this alternative is removed
    - check is made till one alternative left
- When all criteria were checked and there are more than 2 alternatives, next criteria threshold is checked
    - thresholds rather should have *increasing* character - we have higher expectations for each of them
- When there is no more criteria thresholds and more than 1 alternative those are copied to the final ranking to the top in reverse order,
so when alternative `1` and `2` left (in given order), the final ranking will be `2` as first, and `1` as the second.
        
##### Input parameters:
- `Function` - name of the thresholds type function. Possible are:
    - `thresholds` - these thresholds will be checked for each alternative.
     Takes only object with `thresholds` field, which contains list of weights for *every* criterion:
    ```json
    {
      "thresholds": [{"c1": 1, "c2": 2},{"c1":  2, "c2": 3}]
    }
    ```
    - `idealMultipliedCoefficient` - thresholds are generated by multiplication based on criteria values range. Parameters:
        - `minValue` - minimum allowed value in range [0, 1]
        - `maxValue` - maximum allowed value in range [0, 1]
        - `coefficient` - value in range (0, 1).
        > We assume the increasing tendency, so under the hood coefficient is in reality 1 + x (0.2 -> 1.2).
        > Each threshold is generated by equation `min((1 + minValue) * (1 + coefficient)^i - 1, 1)`
        > with upperbound 1 (`i` is iteration, starting from 0).
        > However, the higher threshold cannot be higher than `maxValue`
    - `idealAdditiveCoefficient` - thresholds are generated by adding constant values based on criteria values range to the current value.
     Parameters are the same as for `idealMultipliedCoefficient`, with 2 differences:
        - `coefficient` is raw (no `1` addition) 
        - thresholds are calculated with equation `min(minValue + coefficient * i, 1)` 
    
- `Params` - parameters for given function, described for `Function`
- `RandomAlternativesOrdering` - whether to shuffle alternatives ordering
- `Weights` - weights for criteria, the higher the earlier criterion will be checked
- `RandomSeed` - seed for ordering when `RandomAlternativesOrdering` is truly or when criteria have equal weights.

### Satisfaction Heuristic
- [implementation](lib/logic/limited-rationality/satisfaction)

Works similarly to [aspect elimination](#aspect-elimination-heuristic), with 3 differences
    - search is alternative-wise - each alternative is checked whether it meets our expectations
    - the earlier alternative will meet threshold, the higher it is in the final ranking
    - therefore, thresholds probably should be *decreasing*; when nothing meets our expectations, we should decrease them
    
##### Input parameters:
- `Function` - threshold function type, allowed are:
    - `thresholds` - given criteria thresholds, `Params` are - same as in aspect elimination,
     just an object with `thresholds` field which contains list of weights for each criterion
    - `idealMultipliedCoefficient` - params are same as for aspect elimination (`coefficient`, `minValue`, `maxValue`) with three differences:
        - thresholds are calculated via `maxValue * coefficient^i` (`i` is iteration, starting from 0),
         so probably coefficient should have value greater than 0.5 (*1 is not added to it*).
        - `minValue` and `maxValue` are limited to (0,1] because of current threshold value calculation
         (0 is reachable only when `coefficient == 0`)
        - value is limited to 0 (the least threshold won't have lower value than `minValue`)
    - `idealSubtractiveCoefficient` - same params as `idealSubtractiveCoefficient`, but current value is calculated with equation
    `max(maxValue - coefficient*i, 0)` (contrary to `idealAdditiveCoefficient`).
- `Params` - parameters for `Function`
- `CurrentChoice` - optional, currently *possessed* alternative, will be the first in search order,
 will occur in the final result even when not present in `ChooseToMake`
- `RandomAlternativesOrdering` - whether to shuffle alternatives order (if `CurrentChoise` is passed, it will be the first one anyway)
- `RandomSeed` - used for alternatives order shuffle and 

## Biases
As humans, we often tends to see things differently, forget about something, conceal certain things etc. These ideas tries to address mentioned problems.

##### Apply probability
Each bias can be applied with some probability.

##### Criteria ordering
For biases which operates on multiple criteria in purpose to modify them there is possibility to describe
- `ordering` in which these are used,
- `ratio` which depends on total number of criteria, so value belongs to [0, 1]. Final value is a floor of, for example for 3 criteria and `ratio` = 0.3, no criterion is considered, whereas for `ratio` = 0.34 - only one (if min >= 1)
- `min` and `max` is a criteria number which is going to be processed, considering `ratio` also (hard bounds). `max` cannot be lower than `min`.

`ordering` can have following values:
- `weakest` - default, criteria are sorted ascending by the influence on the final result
- `weakestByProbability` - each criterion is processed with a probability inversely proportional to its influence (the more important, the lower probability to be took as a next)
- `strongest` - contrary to `weakest`
- `strongestByProbability` - contrary to `weakestByProbability`
- `random` - criteria are processed randomly (configured via `RandomSeed` on the same level)

##### Reference criterion
In some biases there is a need to create a new criterion. However, because of multiple preference function
 or heuristics it is hard to compute weights for each method and make it useful.
Because of it existing criterion is took as a reference - we call it `referenceCriterion`.
This criterion strategy choosing can be configured by setting `ReferenceCriterionType` property on the bias level.
Allowed values are:
- `importanceRatio` - default one, has parameter `newCriterionImportance` which take a value [0, 1] (default 0); 
reference criterion will be chosen based on its influence on the final result.
- `randomUniform` - has parameter `newCriterionRandomSeed`, each criterion has the same possibility to be the reference one
- `randomWeighted` - has parameter `newCriterionRandomSeed`, criterion weight is calculated by equation `weight = minWeight/value`,
 so the most influencing one have the lowest change to be the reference criterion.

### Criteria omission
- [implementation](lib/logic/biases/criteria-omission)

Sometimes we say that something is important for us, but in reality we don't think,
 criteria number is overwhelming or other criteria just take our attention.
 Therefore, we are omitting certain criteria.
 
##### Input parameters:
Same as for [criteria ordering](#criteria-ordering)

As a result method is returning:
- `OmittedCriteria` - list of criteria names which were omitted

### Criteria concealment
- [implementation](lib/logic/biases/criteria-concealment)

Sometimes we consider some criteria, but don't reveal this fact to the audience/coordinator. These criteria are concealed.

##### Input parameters:
- `RandomSeed` - seed for criteria values generation
- `NewCriterionScaling` - `float`, option to equally shrink or widen newly created criterion values range. *Cannot be zero*

Also, new criterion parametrization is allowed like described in [reference criterion](#reference-criterion).

The method returns:
- `AddedCriteria` - list of added criteria (empty or 1). Each item is an object with:
    - `Id` - id of created criterion
    - `Type` - gain or cost,
    - `AlternativesValues` - values for each alternative,
    - `MethodParameters` - parameters for given preference function or heuristic (name will occur there)
    - `ValuesRange` - possible criterion values range

### Criteria mixing
- [implementation](lib/logic/biases/criteria-mixing)

*Even if the criteria internally considered by the DM are preferentially independent, 
they may have been inadvertently corrupted when modeling the problem by
 mixing them in such a way that violates preferential independence*

##### Input parameters:
- `RandomSeed` - seed for criteria to mix choose
- `MixingRatio` - [0,1], ratio of each component impact on the final criterion, default is `0.5`

Reference criterion parameters can be configured as described in [reference criterion](#reference-criterion)

As an output following information is returned:
- `Component1` - first criterion which was used to mix (object with criterion name (`id`), `type` and `scaledValues` for each alternative)
- `Component2` - second criterion used in mix
- `NewCriterion` - created criterion
- `Params` - added params for the preference function/heuristic for the new criterion

### Preference reversal
- [implementation](lib/logic/biases/preference-reversal)

Sometimes we change our preference during the process, like we declared price as cost criterion,
 but later we noticed that the lower price is worse for us because of others perception...

##### Input parameters:
Same as for [criteria ordering](#criteria-ordering)

The method returns:
- `ReversedPreferenceCriteria` - list of criteria which preference was reversed.
 Because of other methods and possible conditions violations preference is changed only withing given criterion 
 values range (`min` becomes `max`, `max` becomes `min`, rest is `max - val + min).
 Each item is an object with:
    - `Id` - id of criterion, which preference was changed,
    - `Type` - gain or cost, same as for original criterion,
    - `AlternativesValues` - values for each alternative,
    - `ValuesRange` - criterion values boundaries

### Fatigue
- [implementation](lib/logic/biases/fatigue)

When we are tired our mental abilities get weaker over time. It also influences on our judgements and evaluation.

##### Input parameters:
- `Function` - name of the fatigue coefficient function, allowed are:
    - `expFromZero` - exponential function with ground = 0, expressed by equation `multiplier * e^(alpha * queryNumber) - multiplier`
    - `const` - given `value` will be the same considered as a function output
- `Params` - parameters for `Function`
- `RandomSeed` - used for criterion value blur generation

Output from the fatigue function is used as a coefficient in criterion value blur, which depends also on the criterion value range.
Blur is added or subtracted from the current criterion value depending on the generated sign.
*The result may be also negative!*

The method returns:
- `EffectiveFatigueRatio` - output of the fatigue function
- `ConsideredAlternatives` - values for criteria for each considered alternatives after fatigue application
- `NotConsideredAlternatives` - same as for `ConsideredAlternatives`, but for those not considered

### Anchoring
- [implementation](lib/logic/biases/anchoring)

We often don't have our own opinion in given subject. We are not an expert in every area.
However, when someone gives us some reference, even not explicit, we trust in it or shape our point of view
 taking it into account, even when our intuition tells us that this is not relevant.
  
This is called *anchoring*.

The idea is like follows:
- we need some anchoring, there can be many of them
- those anchoring can be reduced to certain reference points. There can be configured.
- for each alternative we have also differences for every reference point. 
    - Given alternative might be better or worse on certain criterion. 
    - we have different approach in case of gain and loss.
- there can be several ways how to apply anchoring result, but for sure it influences our final evaluation.

##### Input parameters:
> Each `ReferencePoints`, `Loss`, `Gain` and `Applier` are object with fields `Function` and `Params`, where 
> `Function` describes selected strategy and `Params` are dedicated parameters.

- `AnchoringAlternatives` - list of alternatives which creates an anchoring point(s).
 Each object in that list contains `alternative` name and its `coefficient` in anchoring.
  The higher `coefficient`, the more important alternative is. Low coefficient can be compared to forgetting.
- `ReferencePoints` - object describing reference (anchoring) points evaluation strategy. Possible strategies are
    - `ideal` - `AnchoringAlternatives` are reduced to the one containing *the best values* of criteria,
     taking into account coefficients of each alternative and its criteria   
    - `nadir` - `AnchoringAlternatives` are reduced to the one containing *the worst values*.
     Higher coefficient for greater criterion value makes it worse than when criterion value is lower.  
- `Loss` - definition for loss function (when given alternative has worse evaluation than the reference one).
    Function takes as an argument `dif`, which is always in range [0, 1] - difference is scaled according to criteria values range.
    The final result will have different oposite sign (so when positive value is returned from function, it will be ultimately negative).
    
    Possible functions are:
    - `expFromZero` - calculated as `multiplier * e^(alpha * dif) - multiplier`.`multiplier` and `alpha` are the parameters.
    - `linear` - calculated as `a * dif + b`. `a` and `b` are the parameters.
- `Gain` - same as loss, but for cases when given alternative has better criterion value than the reference one.
- `Applier` - defines strategy for applying reference points difference at the end of processing. Possible are:
    - `inline` - for every alternative applies averaged reference points differences directly on the alternative criteria.
     Processed differences are rescaled to the original criteria range. By default, results are trimmed to 
     the original criteria values range, but it is also possible to disable this behavior.
     When parameter `unbounded` is set to `true`, resulting criterion value can be negative or much above normal
      criterion values range, depending on the evaluation.
      
Method returns following data:
- `ReferencePoints` - evaluated reference points as alternatives with their criteria values 
- `CriteriaScaling` - scaling for each criterion used during computations
- `PerReferencePointsDifferences` - each alternative differences for each reference points with values already computed by the functions.