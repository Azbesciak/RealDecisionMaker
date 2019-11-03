import React from 'react';
import {createStyles, makeStyles, Theme} from "@material-ui/core";
import MethodsList from "./MethodsList";
import CriteriaContainerComponent from "./criteria/CriteriaContainerComponent";
import {AlternativesContainerComponent} from "./alternatives/AlternativesContainerComponent";
import {Criterion} from "./criteria/CriterionComponent";
import {Collection} from "./utils/ValuesContainerComponent";
import {Alternative} from "./alternatives/AlternativeComponent";
import {remapCollection} from "./utils/utils";

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
}

class QueryForm extends React.Component<any, DecisionMakerQuery> {
    state = {
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

    private applyCriteriaToAlternative(value: Alternative, criteria: Collection<Criterion>) {
        return {
            id: value.id,
            criteria: remapCollection(criteria, (cid, c) => {
                const current = value.criteria[cid];
                return {id: c.id, value: current ? current.value : 0}
            })
        }
    }

    render() {
        return (
            <form noValidate autoComplete="off">
                <CriteriaContainerComponent payload={this.state.criteria} onUpdate={this.onCriteriaUpdated}/>
                <AlternativesContainerComponent
                    payload={this.state.alternatives}
                    onUpdate={this.onAlternativesUpdated}
                />
                <MethodsList methodComponents={this.state.preferenceFunctions}/>
            </form>
        );
    }
}

export default QueryForm