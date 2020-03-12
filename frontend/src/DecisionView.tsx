import React, {FunctionComponent} from "react"
import {AlternativeResult} from "./Result";
import {Graph} from "react-d3-graph";
import {Theme, useTheme} from "@material-ui/core";

export interface Dimensions {
    width: number;
    height: number;
}

interface DecisionViewProps {
    result: AlternativeResult[];
    dimensions: Dimensions;
}

const myConfig = {
    "automaticRearrangeAfterDropNode": true,
    "collapsible": true,
    "directed": true,
    "focusAnimationDuration": 0.75,
    "focusZoom": 1,
    "height": 400,
    "width": 800,
    "highlightDegree": 1,
    "highlightOpacity": 1,
    "linkHighlightBehavior": false,
    "maxZoom": 8,
    "minZoom": 0.1,
    "nodeHighlightBehavior": true,
    "panAndZoom": false,
    "staticGraph": false,
    "staticGraphWithDragAndDrop": false,
    "d3": {
        "alphaTarget": 0.05,
        "gravity": -100,
        "linkLength": 100,
        "linkStrength": 1,
        "disableLinkForce": false
    },
    "node": {
        "color": "#d3d3d3",
        "fontColor": "black",
        "fontSize": 8,
        "fontWeight": "normal",
        "highlightColor": "SAME",
        "highlightFontSize": 8,
        "highlightFontWeight": "bold",
        "highlightStrokeColor": "SAME",
        "highlightStrokeWidth": "SAME",
        "labelProperty": "id",
        "mouseCursor": "pointer",
        "opacity": 1,
        "renderLabel": true,
        "size": 200,
        "strokeColor": "none",
        "strokeWidth": 1.5,
        "svg": "",
        "symbolType": "circle"
    },
    "link": {
        "color": "#d3d3d3",
        "fontColor": "black",
        "fontSize": 8,
        "fontWeight": "normal",
        "highlightColor": "SAME",
        "highlightFontSize": 8,
        "highlightFontWeight": "normal",
        "labelProperty": "label",
        "mouseCursor": "default",
        "opacity": .5,
        "renderLabel": false,
        "semanticStrokeWidth": false,
        "strokeWidth": 1.5,
        "markerHeight": 6,
        "markerWidth": 6
    }
};

const getConfig = (theme: Theme, {width, height}: Dimensions) => {
    const color = theme.palette.primary[theme.palette.type];
    const fontColor = theme.palette.text.primary;
    const highlightStrokeColor = theme.palette.secondary[theme.palette.type];
    return {
        ...myConfig,
        width, height,
        node: {...myConfig.node, color, fontColor, highlightStrokeColor},
        link: {...myConfig.link, color, fontColor}
    }
};

// https://williaster.github.io/data-ui/?selectedKind=network&selectedStory=Default%20network&full=0&addons=0&stories=1&panelRight=0
const DecisionView: FunctionComponent<DecisionViewProps> = ({result, dimensions}) => {
    const config = getConfig(useTheme(), dimensions);
    const nodes = result.map((r, i) => ({id: r.alternative.id}));
    const data = {
        nodes: nodes,
        links: result.flatMap(r1 => r1.betterThanOrSameAs.map(r2 => ({
            source: r1.alternative.id,
            target: r2
        })))
    };
    return (<Graph
        id="graph-id" // id is mandatory, if no id is defined rd3g will throw an error
        data={data}
        config={config}
    />);
};
export default DecisionView;