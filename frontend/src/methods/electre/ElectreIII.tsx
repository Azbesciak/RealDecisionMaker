import {Method, MethodFactory} from "../declarations";
import {Collection} from "../../utils/ValuesContainerComponent";
import {Criterion} from "../../criteria/CriterionComponent";
import React, {Component} from "react";
import {criterionNamePlaceholder, fromCollection} from "../../utils/utils";
import {
    blankElectreCriterion,
    defaultDistillationFun,
    ElectreCriterion,
    ElectreIIIParams,
    LinearFunctionParameters
} from "./electre";
import ElectreIIICriterionComp from "./ElectreIIICriterion";
import LinearFunction from "./LinearFunction";

export class ElectreIIIFactory implements MethodFactory<ElectreIIIParams> {
    constructor(private onUpdate: () => void) {
    }

    readonly methodName = "electreIII";

    protected params: ElectreIIIParams = {electreCriteria: {}, electreDistillation: defaultDistillationFun()};
    readonly updateParams = (newWeights: Partial<ElectreIIIParams>) => {
        this.params = {...this.params, ...newWeights};
        this.onUpdate();
    };
    readonly getComponent = (criteria: Collection<Criterion>) => (
        <ElectreIII criteria={criteria} methodParameters={this.params} onChange={this.updateParams}/>
    );
    readonly getParams = (criteria: Collection<Criterion>) => ({
        electreCriteria: fromCollection(criteria, (k, c) => [c.id, this.params.electreCriteria[k] || blankElectreCriterion()]),
        distillationFun: this.params.electreDistillation || defaultDistillationFun()
    });

}


export class ElectreIII extends Component<Method<ElectreIIIParams>> {
    private handleCriterionChange = (criterionId: string) => (criterion: ElectreCriterion) =>
        this.props.onChange({
            electreCriteria: {...this.props.methodParameters.electreCriteria, [criterionId]: criterion}
        });

    private handleDistillationChange = (electreDistillation: LinearFunctionParameters) => this.props.onChange({
        electreDistillation
    });

    render() {
        return (
            <React.Fragment>
                <div className="electre-criteria">
                    {Object.entries(this.props.criteria || {}).map(([id, c], i) => (
                        <ElectreIIICriterionComp
                            key={id}
                            criterionName={c.id || criterionNamePlaceholder(i)}
                            params={this.props.methodParameters.electreCriteria[id] || blankElectreCriterion()}
                            onChange={this.handleCriterionChange(id)}/>
                    ))}
                </div>
                <LinearFunction
                    params={this.props.methodParameters.electreDistillation}
                    onChange={this.handleDistillationChange}
                    label={"distillation function"}/>
            </React.Fragment>
        )
    }
}
