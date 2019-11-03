import * as React from 'react';
import {Component} from 'react';
import AddButton from "./AddButton";
import {List, ListItem} from "@material-ui/core";
import {RemoveButtonComponent} from "./RemoveButtonComponent";

export interface ValuesContainerState<T> {
    payload: Collection<T>;
}

export interface Collection<T> {
    [key: string]: T
}

export interface ValuesContainerProps<T> extends ValuesContainerState<T> {
    onUpdate: (items: Collection<T>) => void;
}

export abstract class ValuesContainerComponent<T> extends Component<ValuesContainerProps<T>, ValuesContainerState<T>> {
    private counter = 0;
    abstract newItemFactory: () => T;
    addItem = () => this.update("" + ++this.counter, this.newItemFactory());

    update(id: string, value: T) {
        this.props.onUpdate({...this.props.payload, [id]: value})
    }

    removeItem = (k: keyof T) => () => {
        const {[k]: key, ...payload} = this.props.payload;
        this.props.onUpdate(payload);
    };

    abstract label: string;

    abstract createNewComponent(key: string, value: T): JSX.Element

    render() {
        return (
            <React.Fragment>
                <AddButton label={`Add ${this.label}`} onAdd={this.addItem}/>
                <List>
                    {Object.entries(this.props.payload).map(([k, v]) => (
                        <ListItem key={k}>
                            {this.createNewComponent(k, v)}
                            <RemoveButtonComponent onRemove={this.removeItem(k as any)}/>
                        </ListItem>
                    ))
                    }
                </List>
            </React.Fragment>
        );
    }
}