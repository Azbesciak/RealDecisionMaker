import {SimpleWeightsMethodFactory} from "./SimpleWeightsMethodFactory";
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import {SimpleWeightsComponent} from "./SimpleWeightsComponent";
import React from "react";

export class OWAFactory extends SimpleWeightsMethodFactory {
    constructor(onUpdate: () => void) {
        super(onUpdate)
    }

    readonly methodName = "owa";
    getComponent = (criteria: Collection<Criterion>) => (
        <OWA criteria={criteria} methodParameters={this.params} onChange={this.updateParams}/>
    );
    getParams = (criteria: Collection<Criterion>) => ({
        weights: Object.fromEntries(Object.keys(criteria).map((k, i) => [_id(i), this.params.weights[k] || 0]))
    })
}

export class OWA extends SimpleWeightsComponent {
    keys = (criteria: Collection<Criterion>) => Object.entries(criteria).map(([id], i) => ({id, name: _id(i)}));
}

function _id(index: number) {
    return "" + (index + 1);
}
