import React from 'react';
import {Criterion, CriterionComponent} from "./CriterionComponent";
import {ValuesContainerComponent} from "../utils/ValuesContainerComponent";
import {defaultCriterion} from "./CriterionType";


class CriteriaContainerComponent extends ValuesContainerComponent<Criterion> {
    label = "criterion";
    newItemFactory = () => ({id: "", type: defaultCriterion});

    createNewComponent(key: string, value: Criterion): JSX.Element {
        return (<CriterionComponent key={key} onChange={v => this.update(key, v)} value={value}/>);
    }

}

export default CriteriaContainerComponent;
