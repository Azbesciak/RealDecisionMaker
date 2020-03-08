import React from 'react';
import {ItemValueComponent} from "../utils/ItemValueComponent";
import {Collection} from "../utils/ValuesContainerComponent";
import {TextField} from "@material-ui/core";
import {criterionNamePlaceholder, handleInputValueChange} from "../utils/utils";


export interface Alternative {
    id: string;
    criteria: Collection<AlternativeCriterion>;
}

export interface AlternativeCriterion {
    id: string;
    value: number;
}

export class AlternativeComponent extends ItemValueComponent<Alternative> {
    updateWeight = (criterionId: string) => handleInputValueChange(valueStr => {
        const {id, criteria} = this.props.value;
        const value = (+valueStr);
        this.props.onChange({
            id,
            criteria: {
                ...criteria, [criterionId]: {id: criteria[criterionId].id, value}
            }
        })
    });

    render() {
        return (
            <div>
                {this.getIdField()}
                {Object.entries(this.props.value.criteria || {}).map(([k, v], i) => (
                    <TextField key={k} value={v.value} label={v.id || criterionNamePlaceholder(i) } required
                               type={'number'}
                               onChange={this.updateWeight(k)}/>
                ))}
            </div>
        )
    }
}
