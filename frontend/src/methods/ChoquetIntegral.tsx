import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";
import {SimpleWeightsComponent} from "./SimpleWeightsComponent";
import React from "react";
import {SimpleWeightsMethodFactory} from "./SimpleWeightsMethodFactory";

export class ChoquetIntegralFactory extends SimpleWeightsMethodFactory {
    constructor(onUpdate: () => void) {
        super(onUpdate)
    }

    readonly methodName = "choquetIntegral";
    getComponent = (criteria: Collection<Criterion>) => (
        <ChoquetIntegral criteria={criteria} methodParameters={this.params} onChange={this.updateParams}/>
    );
    getParams = (criteria: Collection<Criterion>) => ({
        weights: Object.fromEntries(powerSet(Object.keys(criteria))
            .map(composedKeys => [composedKeys.map(k => criteria[k].id).join(","), this.params.weights[composedId(composedKeys)] || 0]))
    })
}

function powerSet<T>(ids: T[]): T[][] {
    const obj: any = {};
    //This loop is to take out all duplicate number/letter
    for (let i = 0; i < ids.length; i++) {
        obj[ids[i]] = true;
    }
    //variable array will have no duplicates
    const array: string[] = Object.keys(obj);
    let result = [[]];
    for (let i = 0; i < array.length; i++) {
        //this line is crucial! It prevents us from infinite loop
        const len = result.length;
        for (let x = 0; x < len; x++) {
            result.push(result[x].concat(array[i] as any))
        }
    }
    // removes first EMPTY element
    result.shift();
    return result;
}

export class ChoquetIntegral extends SimpleWeightsComponent {
    keys = (criteria: Collection<Criterion>) => {
        const keys = powerSet(Object.keys(criteria));
        return keys.map(composedKey => ({
            id: composedId(composedKey),
            name: composedKey.map(k => criteria[k].id).join(", ")
        }));
    };
}

function composedId(composedKey: string[]) {
    return composedKey.join(",")
}
