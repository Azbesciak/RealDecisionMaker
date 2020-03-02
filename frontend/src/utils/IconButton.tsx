import React, {FunctionComponent} from "react";
import {Button} from "@material-ui/core";

interface OwnProps {
    label: string;
    icon: JSX.Element;
    disabled?: boolean;
    onClick: () => void;
}

const IconButton: FunctionComponent<OwnProps> = (props) => {
    return (
        <Button
            variant="outlined"
            startIcon={props.icon}
            disabled={props.disabled}
            onClick={props.onClick}
        >
            {props.label}
        </Button>
    );
};

export default IconButton;
