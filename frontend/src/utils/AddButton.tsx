import React, {FunctionComponent} from 'react';
import AddIcon from "@material-ui/icons/Add";
import IconButton from "./IconButton";

interface OwnProps {
    label: string;
    onAdd: () => void;
}

const AddButton: FunctionComponent<OwnProps> = (props) => {
    return (<IconButton label={props.label} icon={<AddIcon/>} onClick={props.onAdd}/>);
};

export default AddButton;
