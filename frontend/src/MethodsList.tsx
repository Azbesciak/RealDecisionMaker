import React, {FunctionComponent} from 'react';
import {createStyles, List, ListItem, makeStyles, Theme} from '@material-ui/core';


interface OwnProps {
    methodComponents: { [key: string]: any }
}

type Props = OwnProps;
const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            width: '100%',
            maxWidth: 360,
            backgroundColor: theme.palette.background.paper,
        },
    }),
);

const MethodsList: FunctionComponent<Props> = (props) => {
    const classes = useStyles();
    return (
        <div className={classes.root}>
            <List>
                {Object.entries(props.methodComponents).map(([k, v]) => (
                    <ListItem button key={k}>{k}</ListItem>
                ))}
            </List>
        </div>
    );
};

export default MethodsList;
