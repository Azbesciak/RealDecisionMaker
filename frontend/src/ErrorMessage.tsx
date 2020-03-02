import React, {FunctionComponent} from 'react';
import MuiAlert from '@material-ui/lab/Alert';
import {Snackbar} from "@material-ui/core";

interface OwnProps {
    message?: string | null;
    closed: () => void;
}

const ErrorMessage: FunctionComponent<OwnProps> = (props) => {
    return (
        <Snackbar open={!!props.message} autoHideDuration={6000} onClose={props.closed}>
            <MuiAlert severity="error" elevation={6} variant={"filled"}>{props.message}</MuiAlert>
        </Snackbar>
    );
};

export default ErrorMessage