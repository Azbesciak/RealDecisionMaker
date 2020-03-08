import React, {FunctionComponent} from "react";
import {TextField} from "@material-ui/core";
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
        <div className="linear-function">
            <div className="function-label">{props.label}</div>
            <div className="function-coefficient">
                <TextField
                    value={params.a}
                    label={"a"}
                    type={'number'} onChange={updateWeight("a")}/>
            </div>
            <div className="function-coefficient">
                <TextField
                    value={params.b}
                    label={"b"}
                    type={'number'} onChange={updateWeight("b")}/>
            </div>
        </div>
    );
};

export default LinearFunction;
