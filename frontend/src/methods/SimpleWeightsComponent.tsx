import React, {Component} from "react";
import {Method, SimpleWeights} from "./declarations";
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import {handleInputValueChange} from "../utils/utils";
import Grid from "@material-ui/core/Grid";
import {TextField} from "@material-ui/core";

export interface CriterionWeightInfo {
    id: string;
    name: string;
}

export abstract class SimpleWeightsComponent extends Component<Method<SimpleWeights>> {

    abstract keys: (criteria: Collection<Criterion>) => CriterionWeightInfo[];

    private updateWeight = (criterionId: string) => handleInputValueChange(valueStr => {
        const weights = this.props.methodParameters.weights;
        const value = (+valueStr);
        if (value >= 0)
            this.props.onChange({
                weights: {...weights, [criterionId]: value}
            })
    });

    render() {
        return (
            <div className="weights-grid">
                {this.keys(this.props.criteria || {}).map(info => (
                    <TextField
                        className="weight-input"
                        key={info.id}
                        value={this.props.methodParameters.weights[info.id] || 0}
                        label={info.name} required
                        inputProps={{min: "0"}}
                        type={'number'} onChange={this.updateWeight(info.id)}/>
                ))}
            </div>
        )
    }
}
