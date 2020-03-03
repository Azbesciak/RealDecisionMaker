import React from 'react';
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import {SimpleWeightsComponent} from "./SimpleWeightsComponent";
import {fromCollection} from "../utils/utils";
import {SimpleWeightsMethodFactory} from "./SimpleWeightsMethodFactory";

export class WeightedSumFactory extends SimpleWeightsMethodFactory {
    constructor(onUpdate: () => void) {
        super(onUpdate)
    }

    readonly methodName = "weightedSum";
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
