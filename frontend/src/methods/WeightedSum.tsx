import React, {Component} from 'react';
import {handleInputValueChange} from "../utils/utils";
import {TextField} from "@material-ui/core";
import {Method, MethodFactory, SimpleWeights} from "./declarations";
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import makeStyles from "@material-ui/core/styles/makeStyles";

export class WeightedSumFactory implements MethodFactory<SimpleWeights> {
    constructor(private onUpdate: () => void) {
    }

    readonly methodName = "weightedSum";
    private params: SimpleWeights = {weights: {}};
    private updateParams = (newWeights: Partial<SimpleWeights>) => {
        this.params = {...this.params, ...newWeights};
        this.onUpdate();
    };
    getComponent = (criteria: Collection<Criterion>) => (
        <WeightedSum criteria={criteria} methodParameters={this.params} onChange={this.updateParams}/>
    );
    getParams = () => this.params;
}

const useStyles = makeStyles(theme => ({
    root: {
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
}));

export class WeightedSum extends Component<Method<SimpleWeights>> {

    updateWeight = (criterionId: string) => handleInputValueChange(valueStr => {
        const weights = this.props.methodParameters.weights;
        const value = (+valueStr);
        if (value >= 0)
            this.props.onChange({
                weights: {...weights, [criterionId]: value}
            })
    });

    render() {
        return (
            <Grid container spacing={3}>
                {Object.entries(this.props.criteria || {}).map(([k, v]) => (
                    <Grid item key={k}>
                        <Paper>
                            <TextField value={this.props.methodParameters.weights[k] || 0} label={v.id} required
                                       inputProps={{min: "0"}}
                                       type={'number'} onChange={this.updateWeight(k)}/>
                        </Paper>
                    </Grid>
                ))}
            </Grid>
        )
    }
}
