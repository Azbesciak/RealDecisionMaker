import {Collection} from "../../utils/ValuesContainerComponent";

export interface LinearFunctionParameters {
    a?: number;
    b?: number;
}

export interface ElectreCriterion {
    k: number;
    q: LinearFunctionParameters;
    p: LinearFunctionParameters;
    v: LinearFunctionParameters;
}

export interface ElectreIIIParams {
    electreCriteria: Collection<ElectreCriterion>;
    electreDistillation?: LinearFunctionParameters;
}

export function defaultDistillationFun(): LinearFunctionParameters {
    return {a: -.15, b: .3}
}

export function blankDistillationFun(): LinearFunctionParameters {
    return {a: undefined, b: undefined}
}

export function blankElectreCriterion(): ElectreCriterion {
    return {
        k: 0,
        q: blankDistillationFun(),
        p: blankDistillationFun(),
        v: blankDistillationFun()
    }
}
