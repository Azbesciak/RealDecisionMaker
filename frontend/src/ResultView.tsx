import React, {Component} from "react";
import ErrorMessage from "./ErrorMessage";
import {Decision} from "./Result";
import DecisionView, {Dimensions} from "./DecisionView";


export interface ResultProps {
    decision: Decision;
}

interface ResultViewState extends Decision {
    recentProps: Readonly<ResultProps> | null;
    viewProps: Dimensions;
}

class ResultView extends Component<ResultProps, ResultViewState> {
    state: ResultViewState = {
        recentProps: null,
        error: null,
        result: [],
        viewProps: {width: 400, height: 400}
    };
    private clearError = () => this.setState({error: null});

    componentDidMount(): void {
        const updateViewProps = () => this.setState({viewProps: {width: window.innerWidth * .9, height: 500}});
        window.addEventListener("resize", updateViewProps);
        updateViewProps()
    }

    static getDerivedStateFromProps(nextProps: Readonly<ResultProps>, nextContext: ResultViewState): ResultViewState {
        if (nextProps !== nextContext.recentProps) {
            return {
                recentProps: nextProps,
                viewProps: nextContext.viewProps,
                ...nextProps.decision
            }
        }
        return nextContext;
    }

    render() {
        return (
            <React.Fragment>
                <ErrorMessage message={this.state.error} closed={this.clearError}/>
                {this.state.result && this.state.result.length > 0 &&
                <DecisionView result={this.state.result} dimensions={this.state.viewProps}/>}
            </React.Fragment>
        )
    }
}

export default ResultView;