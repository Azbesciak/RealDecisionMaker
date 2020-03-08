import React, {Component} from "react";
import {Identifiable, ItemValue} from "./item-value";
import {handleInputValueChange} from "./utils";
import {TextField} from "@material-ui/core";

export abstract class ItemValueComponent<T extends Identifiable> extends Component<ItemValue<T>> {
    getIdField(placeholder?: string) {
        return (
            <TextField
                label="Name"
                required
                placeholder={placeholder}
                value={this.props.value.id} onChange={this.handleIdChange}
            />
        )
    }

    handleIdChange = handleInputValueChange(id => this.update({id} as Partial<T>));

    update(update: Partial<T>) {
        const newState = {...this.props.value, ...update};
        this.props.onChange(newState);
    }
}