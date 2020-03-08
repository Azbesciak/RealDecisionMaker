import React, {Component} from "react";
import {Collection} from "./utils/ValuesContainerComponent";
import ErrorMessage from "./ErrorMessage";


interface DecisionError {
    error: string | null;
}

export interface NamedAlternative {
    id: string;
    criteria: Collection<number>;
}

interface AlternativeResult {
    alternative: NamedAlternative;
    value: number;
    betterThanOrSameAs: string[]
}

interface DecisionResult {
    result: AlternativeResult[];
}

export type Decision = Partial<DecisionError & DecisionResult>

export interface ResultProps {
    decision: Decision;
}

interface ResultViewState extends Decision {
    recentProps: Readonly<ResultProps> | null;
}

class ResultView extends Component<ResultProps, ResultViewState> {
    state: ResultViewState = {
        recentProps: null,
        error: null,
        result: []
    };
    private clearError = () => this.setState({error: null});

    static getDerivedStateFromProps(nextProps: Readonly<ResultProps>, nextContext: ResultViewState): ResultViewState {
        if (nextProps !== nextContext.recentProps) {
            return {
                recentProps: nextProps,
                ...nextProps.decision
            }
        }
        return nextContext;
    }

    render() {
        return (
            <React.Fragment>
                <ErrorMessage message={this.state.error} closed={this.clearError}/>
            </React.Fragment>
        )
    }
}

export default ResultView;