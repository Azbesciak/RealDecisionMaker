import React, {FunctionComponent} from 'react';
import DoneIcon from "@material-ui/icons/Done";
import IconButton from "./IconButton";

interface OwnProps {
    label: string;
    onAccept: () => void;
    enabled: boolean;
}

const AcceptButton: FunctionComponent<OwnProps> = (props) => {
    return (<IconButton label={props.label} icon={<DoneIcon/>} onClick={props.onAccept} disabled={!props.enabled}/>);
};

export default AcceptButton;
