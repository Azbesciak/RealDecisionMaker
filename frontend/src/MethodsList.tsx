import React, {FunctionComponent} from 'react';
import {Paper, Tab, Tabs} from '@material-ui/core';
import {camelCaseToNormal, isUndefined} from "./utils/utils";
import {makeStyles} from "@material-ui/core/styles";


interface OwnProps {
    methodComponents: { [key: string]: any };
    method?: string;
    onMethodSelected: (methodName: string) => void
}

const useStyles = makeStyles({
    root: {
        flexGrow: 1,
    },
});

type Props = OwnProps;
const MethodsList: FunctionComponent<Props> = (props) => {
    const classes = useStyles();
    const methods = [...Object.keys(props.methodComponents)];
    const index = props.method ? methods.indexOf(props.method) : 0;
    const handleChange = (event: React.ChangeEvent<{}>, newIndex: number) => {
        props.onMethodSelected(methods[newIndex]);
    };
    return (
        <Paper className={classes.root}>
            <Tabs value={index}
                  onChange={handleChange}
                  indicatorColor="primary"
                  textColor="primary"
                  centered
            >
                {methods.map(k => (
                    <Tab key={k} label={camelCaseToNormal(k)}/>
                ))}
            </Tabs>
        </Paper>
    );
};

export default MethodsList;
