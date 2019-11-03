import React, {FunctionComponent} from 'react';
import DeleteIcon from "@material-ui/icons/Delete";
import {IconButton} from "@material-ui/core";

interface RemoveButtonProps {
    onRemove: () => void;
}

export const RemoveButtonComponent: FunctionComponent<RemoveButtonProps> = props => (
    <IconButton aria-label="delete" size="small" onClick={props.onRemove}>
        <DeleteIcon fontSize="small"/>
    </IconButton>
);
