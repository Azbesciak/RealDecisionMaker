import React, {FunctionComponent} from 'react';
import {Button} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";

interface OwnProps {
    label: string;
    onAdd: () => void;
}

type Props = OwnProps;

const AddButton: FunctionComponent<Props> = (props) => {
    return (
        <Button
            variant="outlined"
            startIcon={<AddIcon/>}
            onClick={props.onAdd}
        >
            {props.label}
        </Button>
    );
};

export default AddButton;
