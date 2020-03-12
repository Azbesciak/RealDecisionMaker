import React, {useState} from 'react';
import Header from "./Header"
import {ThemeProvider} from '@material-ui/styles';
import {createMuiTheme, CssBaseline, useMediaQuery} from "@material-ui/core";
import QueryForm from "./QueryForm";
import ResultView, {Decision} from "./ResultView";
import {lightGreen, orange} from "@material-ui/core/colors";

const App: React.FC = () => {
    const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
    const [state, setState] = useState({decision: {}});
    const theme = React.useMemo(
        () => createMuiTheme({
            palette: {
                type: prefersDarkMode ? 'dark' : 'light',
                primary: lightGreen,
                secondary: orange
            },
        }),
        [prefersDarkMode],
    );
    const onResult = (decision: Decision) => {
        setState({decision})
    };
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            <Header title={"Real Decision Maker"}/>
            <QueryForm onResult={onResult}/>
            <ResultView decision={state.decision}/>
        </ThemeProvider>
    );
};

export default App;
