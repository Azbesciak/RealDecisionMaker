import {MethodFactory, SimpleWeights} from "./declarations";
import React from "react";
import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";

export abstract class SimpleWeightsMethodFactory implements MethodFactory<SimpleWeights> {
    protected constructor(private onUpdate: () => void) {
    }

    protected params: SimpleWeights = {weights: {}};
    readonly updateParams = (newWeights: Partial<SimpleWeights>) => {
        this.params = {...this.params, ...newWeights};
        this.onUpdate();
    };
    readonly abstract getComponent: (criteria: Collection<Criterion>) => JSX.Element;
    readonly abstract getParams: (criteria: Collection<Criterion>) => SimpleWeights;
    readonly abstract methodName: string;

}