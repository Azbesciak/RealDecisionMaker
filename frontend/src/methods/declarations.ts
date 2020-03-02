import {Collection} from "../utils/ValuesContainerComponent";
import {Criterion} from "../criteria/CriterionComponent";

export interface Method<T> {
    criteria: Collection<Criterion>;
    methodParameters: T;
    onChange: (update: Partial<T>) => void;
}

export interface MethodFactory<T> {
    readonly methodName: string;
    getComponent: (criteria: Collection<Criterion>) => JSX.Element;
    getParams: (criteria: Collection<Criterion>) => T;
}

export interface SimpleWeights {
    weights: Collection<number>;
}