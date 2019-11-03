import React from 'react';
import Header from "./Header"
import {ThemeProvider} from '@material-ui/styles';
import {createMuiTheme, CssBaseline, useMediaQuery} from "@material-ui/core";
import QueryForm from "./QueryForm";


const App: React.FC = () => {
    const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

    const theme = React.useMemo(
        () => createMuiTheme({
            palette: {
                type: prefersDarkMode ? 'dark' : 'light',
            },
        }),
        [prefersDarkMode],
    );
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            <Header title={"Real Decision Maker"}/>
            <QueryForm/>
        </ThemeProvider>
    );
};

export default App;
