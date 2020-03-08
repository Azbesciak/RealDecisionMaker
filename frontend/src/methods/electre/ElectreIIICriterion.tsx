import React, {FunctionComponent} from "react";
import {Card, CardContent, CardHeader, TextField} from "@material-ui/core";
import LinearFunction from "./LinearFunction";
import {ElectreCriterion, LinearFunctionParameters} from "./electre";
import {handleInputValueChange} from "../../utils/utils";

interface OwnProps {
    criterionName: string;
    params: ElectreCriterion;
    onChange: (criterion: ElectreCriterion) => void;
}

const ElectreIIICriterionComp: FunctionComponent<OwnProps> = (props) => {
    const handleChange = (name: string) => (params: LinearFunctionParameters) =>
        props.onChange({...props.params, [name]: params});

    const updateK = handleInputValueChange(value => props.onChange({...props.params, k: +(value || 0)}));
    return (
        <Card className="electre-criterion">
            <CardHeader title={props.criterionName}/>
            <CardContent className="electre-criterion-content">
                <div className="electre-k-field">
                    <label>k</label>
                    <TextField value={props.params.k} type={'number'} onChange={updateK}/>
                </div>
                <LinearFunction label={"q"} onChange={handleChange("q")} params={props.params.q}/>
                <LinearFunction label={"p"} onChange={handleChange("p")} params={props.params.p}/>
                <LinearFunction label={"v"} onChange={handleChange("v")} params={props.params.v}/>
            </CardContent>
        </Card>
    );
};

export default ElectreIIICriterionComp;
