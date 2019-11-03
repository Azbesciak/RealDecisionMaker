import React, {FunctionComponent, useState} from 'react';
import {createStyles, FormControl, InputLabel, makeStyles, MenuItem, Select, Theme} from "@material-ui/core";
import {handleInputValueChange} from "../utils/utils";
import {criteriaTypes, CriterionType} from "./CriterionType";
import {ItemValue} from "../utils/item-value";


const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        formControl: {
            margin: theme.spacing(1),
            minWidth: 120,
            textTransform: "capitalize"
        },
        selectEmpty: {
            marginTop: theme.spacing(2),
        },
    }),
);
export const CriterionTypeSelect: FunctionComponent<ItemValue<CriterionType>> = (props) => {
    const classes = useStyles();
    const handleChange = handleInputValueChange((value: CriterionType) => {
        props.onChange(value)
    });
    return (
        <FormControl className={classes.formControl}>
            <InputLabel>Type</InputLabel>
            <Select
                value={props.value}
                onChange={handleChange}
            >
                {criteriaTypes.map(t => (<MenuItem key={t} value={t}>{t}</MenuItem>))}
            </Select>
        </FormControl>
    );
};
