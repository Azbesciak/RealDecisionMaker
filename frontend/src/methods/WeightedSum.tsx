import React from 'react';
import {MethodFactory, SimpleWeights} from "./declarations";
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import {SimpleWeightsComponent} from "./SimpleWeightsComponent";
import {fromCollection} from "../utils/utils";

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
    getParams = (criteria: Collection<Criterion>) => ({
        weights: fromCollection(criteria, (k, v) => [v.id, this.params.weights[k] || 0])
    })
}

export class WeightedSum extends SimpleWeightsComponent {
    keys = (criteria: Collection<Criterion>) => Object.entries(criteria).map(([id, v]) => ({id, name: v.id}));
}
