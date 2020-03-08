import React from 'react';
import MethodsList from "./MethodsList";
import CriteriaContainerComponent from "./criteria/CriteriaContainerComponent";
import {AlternativesContainerComponent} from "./alternatives/AlternativesContainerComponent";
import {Criterion} from "./criteria/CriterionComponent";
import {Collection} from "./utils/ValuesContainerComponent";
import {Alternative} from "./alternatives/AlternativeComponent";
import {fromCollection, remapCollection} from "./utils/utils";
import {WeightedSumFactory} from "./methods/WeightedSum";
import {MethodFactory} from "./methods/declarations";
import AcceptButton from "./utils/AcceptButton";
import ErrorMessage from "./ErrorMessage";
import {ChoquetIntegralFactory} from "./methods/ChoquetIntegral";
import {OWAFactory} from "./methods/OWA";
import {ElectreIIIFactory} from "./methods/electre/ElectreIII";

interface DecisionError {
    error: string | null;
}

export interface NamedAlternative {
    id: string;
    criteria: Collection<number>;
}

interface AlternativeResult {
    alternative: NamedAlternative;
    value: number;
    betterThanOrSameAs: string[]
}

interface DecisionResult {
    result: AlternativeResult[];
}

interface DecisionMakerQuery {
    preferenceFunctions: { [key: string]: any };
    criteria: Collection<Criterion>;
    alternatives: Collection<Alternative>;
    selectedMethod: MethodFactory<any>;
    decision: Partial<DecisionError & DecisionResult>;
}

class QueryForm extends React.Component<any, DecisionMakerQuery> {
    private lastRequestId = 0;
    private readonly update = () => this.setState({});
    private functions = [
        new ChoquetIntegralFactory(this.update),
        new ElectreIIIFactory(this.update),
        new OWAFactory(this.update),
        new WeightedSumFactory(this.update),
    ];
    state: DecisionMakerQuery = {
        preferenceFunctions: {},
        criteria: {},
        alternatives: {},
        decision: {},
        selectedMethod: this.functions[0]
    };

    componentDidMount() {
        fetch(`${process.env.REACT_APP_API_ROOT}/preferenceFunctions`)
            .then(res => res.json())
            .then(results => this.setState({preferenceFunctions: results}));
    }

    onCriteriaUpdated = (criteria: Collection<Criterion>) => this.setState(s => {
        const alternatives = remapCollection(s.alternatives, (a, value) => this.applyCriteriaToAlternative(value, criteria));
        return {criteria, alternatives}
    });

    onAlternativesUpdated = (alts: Collection<Alternative>) => this.setState(s => (
        {alternatives: remapCollection(alts, (k, v) => s.alternatives[k] ? v : this.applyCriteriaToAlternative(v, s.criteria))}
    ));

    onMethodSelected = (method: string) => {
        const selectedMethod = this.functions.find(m => m.methodName === method) as MethodFactory<any>;
        this.setState({selectedMethod});
    };

    private applyCriteriaToAlternative(value: Alternative, criteria: Collection<Criterion>) {
        return {
            id: value.id,
            criteria: remapCollection(criteria, (cid, c) => {
                const current = value.criteria[cid];
                return {id: c.id, value: current ? current.value : 0}
            })
        }
    }

    private onAccept = () => {
        if (!this.state.selectedMethod) return;
        const params = this.state.selectedMethod.getParams(this.state.criteria);
        this.sendRequest(this.state.selectedMethod.methodName, params);
        console.log("PARAMS", params)
    };

    private sendRequest = async (preferenceFunction: string, methodParameters: any,) => {
        this.lastRequestId++;
        this.setState({decision: {}});
        const body = {
            preferenceFunction,
            methodParameters,
            criteria: Object.values(this.state.criteria),
            knownAlternatives: Object.values(this.state.alternatives).map(a => ({
                id: a.id,
                criteria: fromCollection(a.criteria, (_, v) => [v.id, v.value])
            })),
            choseToMake: Object.values(this.state.alternatives).map(a => a.id)
        };
        fetch(`${process.env.REACT_APP_API_ROOT}/decide`, {
            method: "POST",
            body: JSON.stringify(body),
            headers: {
                "Content-Type": "application/json"
            }
        }).then(r => r.json())
            .then(r => {
                if (r.error) {
                    this.setState({decision: {error: r.error}})
                } else {
                    this.setState({decision: {result: r.result}})
                }
            })
            .catch(e => this.setState({decision: {error: e.error || e.message}}))
    };

    renderMethod = () => this.state.selectedMethod && this.state.selectedMethod.getComponent(this.state.criteria);

    private clearMessageIfValid = () => {
        let id = this.lastRequestId;
        return () => {
            if (id === this.lastRequestId) this.setState({decision: {}})
        }
    };

    render() {
        return (
            <form noValidate autoComplete="off" className="query-form">
                <CriteriaContainerComponent payload={this.state.criteria} onUpdate={this.onCriteriaUpdated}/>
                <AlternativesContainerComponent
                    payload={this.state.alternatives}
                    onUpdate={this.onAlternativesUpdated}
                />
                <MethodsList
                    method={this.state.selectedMethod.methodName}
                    methodComponents={this.state.preferenceFunctions}
                    onMethodSelected={this.onMethodSelected}
                />
                {this.renderMethod()}
                <AcceptButton label="OK" onAccept={this.onAccept} enabled={!!this.state.selectedMethod}/>
                <ErrorMessage message={this.state.decision.error} closed={this.clearMessageIfValid()}/>
            </form>
        );
    }
}

export default QueryForm