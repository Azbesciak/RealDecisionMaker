import {Collection} from "./utils/ValuesContainerComponent";

export interface DecisionError {
    error: string | null;
}

export interface NamedAlternative {
    id: string;
    criteria: Collection<number>;
}

export interface AlternativeResult {
    alternative: NamedAlternative;
    value: number;
    betterThanOrSameAs: string[]
}

export interface DecisionResult {
    result: AlternativeResult[];
}

export type Decision = Partial<DecisionError & DecisionResult>