import React from 'react';
import {ValuesContainerComponent} from "../utils/ValuesContainerComponent";
import {Alternative, AlternativeComponent} from "./AlternativeComponent";


export class AlternativesContainerComponent extends ValuesContainerComponent<Alternative> {
    readonly classNames = "alternatives-container";
    readonly label = "alternative";

    newItemFactory = () => ({id: "", criteria: {}});
    createNewComponent(key: string, value: Alternative): JSX.Element {
        return (<AlternativeComponent key={key} onChange={v => this.update(key, v)} value={value}/>);
    }
}
