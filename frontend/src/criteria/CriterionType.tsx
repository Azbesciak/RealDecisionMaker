export enum CriterionType {
    gain = "gain",
    cost = "cost"
}

export const criteriaTypes = [CriterionType.gain, CriterionType.cost];
export const defaultCriterion = CriterionType.gain;