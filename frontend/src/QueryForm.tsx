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
import {ChoquetIntegralFactory} from "./methods/ChoquetIntegral";
import {OWAFactory} from "./methods/OWA";
import {ElectreIIIFactory} from "./methods/electre/ElectreIII";
import {Decision} from "./Result";


interface DecisionMakerQuery {
    preferenceFunctions: { [key: string]: any };
    criteria: Collection<Criterion>;
    alternatives: Collection<Alternative>;
    selectedMethod: MethodFactory<any>;
}

export interface QueryFromProps {
    onResult: (decision: Decision) => void;
}

class QueryForm extends React.Component<QueryFromProps, DecisionMakerQuery> {
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
    };

    private sendRequest = async (preferenceFunction: string, methodParameters: any,) => {
        this.lastRequestId++;
        this.props.onResult({});
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
                this.props.onResult(r);
            })
            .catch(e => this.props.onResult({error: e.error || e.message}))
    };

    renderMethod = () => this.state.selectedMethod && this.state.selectedMethod.getComponent(this.state.criteria);

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
            </form>
        );
    }
}

export default QueryForm