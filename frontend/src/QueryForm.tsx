import React from 'react';
import {createStyles, makeStyles, Theme} from "@material-ui/core";
import MethodsList from "./MethodsList";
import CriteriaContainerComponent from "./criteria/CriteriaContainerComponent";
import {AlternativesContainerComponent} from "./alternatives/AlternativesContainerComponent";
import {Criterion} from "./criteria/CriterionComponent";
import {Collection} from "./utils/ValuesContainerComponent";
import {Alternative} from "./alternatives/AlternativeComponent";
import {remapCollection} from "./utils/utils";
import {WeightedSumFactory} from "./methods/WeightedSum";
import {MethodFactory} from "./methods/declarations";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            display: 'flex',
            flexWrap: 'wrap',
        },
        textField: {
            marginLeft: theme.spacing(1),
            marginRight: theme.spacing(1),
            width: 200,
        },
    }),
);

interface DecisionMakerQuery {
    preferenceFunctions: { [key: string]: any };
    criteria: Collection<Criterion>;
    alternatives: Collection<Alternative>;
    selectedMethod?: MethodFactory<any>;
}

class QueryForm extends React.Component<any, DecisionMakerQuery> {
    private functions = [
        new WeightedSumFactory(() => this.setState({}))
    ];
    state: DecisionMakerQuery = {
        preferenceFunctions: {},
        criteria: {},
        alternatives: {}
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
        const selectedMethod = this.functions.find(m => m.methodName === method);
        this.setState({selectedMethod});
        console.log("METHOD SELECTED", method);
    }

    private applyCriteriaToAlternative(value: Alternative, criteria: Collection<Criterion>) {
        return {
            id: value.id,
            criteria: remapCollection(criteria, (cid, c) => {
                const current = value.criteria[cid];
                return {id: c.id, value: current ? current.value : 0}
            })
        }
    }

    renderMethod = () => this.state.selectedMethod && this.state.selectedMethod.getComponent(this.state.criteria);

    render() {
        return (
            <form noValidate autoComplete="off">
                <CriteriaContainerComponent payload={this.state.criteria} onUpdate={this.onCriteriaUpdated}/>
                <AlternativesContainerComponent
                    payload={this.state.alternatives}
                    onUpdate={this.onAlternativesUpdated}
                />
                <MethodsList methodComponents={this.state.preferenceFunctions}
                             onMethodSelected={this.onMethodSelected}/>
                {this.renderMethod()}
            </form>
        );
    }
}

export default QueryForm