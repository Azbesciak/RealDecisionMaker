import React from 'react';
import {CriterionTypeSelect} from "./CriterionTypeSelect";
import {ItemValueComponent} from "../utils/ItemValueComponent";
import {CriterionType} from "./CriterionType";
import {criterionNamePlaceholder} from "../utils/utils";

export interface Criterion {
    id: string;
    type: CriterionType;
}

export class CriterionComponent extends ItemValueComponent<Criterion> {
    onTypeChange = (type: CriterionType) => this.update({type});

    render() {
        return (
            <div className="criterion list-item">
                {this.getIdField()}
                <CriterionTypeSelect value={this.props.value.type} onChange={this.onTypeChange}/>
            </div>)
    }
}
