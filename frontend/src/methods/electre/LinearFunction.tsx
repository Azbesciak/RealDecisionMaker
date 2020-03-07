import React, {FunctionComponent} from "react";
import {Grid, TextField} from "@material-ui/core";
import {blankDistillationFun, LinearFunctionParameters} from "./electre";
import {handleInputValueChange} from "../../utils/utils";

interface LinearFunctionComponentParams {
    label: string;
    params?: LinearFunctionParameters;
    onChange: (params: LinearFunctionParameters) => void;
}

const LinearFunction: FunctionComponent<LinearFunctionComponentParams> = (props) => {
    const params = props.params || blankDistillationFun();
    const updateWeight = (field: string) => handleInputValueChange(valueStr => {
        const value = +(valueStr || 0);
        const numbers = {...params, [field]: value};
        props.onChange(numbers)
    });
    return (
        <Grid spacing={5} container direction={"row"} style={{flexGrow: 1}}>
            <Grid item xs={1}>{props.label}</Grid>
            <Grid item xs={5}>
                <TextField
                    value={params.a}
                    label={"a"}
                    type={'number'} onChange={updateWeight("a")}/>
            </Grid>
            <Grid item xs={5}>
                <TextField
                    value={params.b}
                    label={"b"}
                    type={'number'} onChange={updateWeight("b")}/>
            </Grid>
        </Grid>
    );
};

export default LinearFunction;
