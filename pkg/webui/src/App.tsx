import React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    useParams,
} from "react-router-dom";
import { JobList } from './JobList';
import { JobView, JobViewProps } from './JobView';
import { PiroServiceClient } from './api/piro_pb_service';
import { WithStyles, ThemeProvider, withStyles } from '@material-ui/styles';
import { CssBaseline, createMuiTheme, createStyles } from '@material-ui/core';
import { GithubPage } from './GithubPage';
import { StartJob } from './StartJob';
import { PiroUIClient } from './api/piro-ui_pb_service';
import { BranchList } from './BranchList';
import { IsReadonly } from './components/IsReadonly';

export interface AppProps extends WithStyles<typeof styles> { }

let url = `${window.location.protocol}//${window.location.host}`;
console.log("server url", url);
const client = new PiroServiceClient(url);
const uiClient = new PiroUIClient(url);

const JobViewWithName: React.SFC<Partial<JobViewProps>> = props => {
    const {name} = useParams();
    return <JobView client={client} jobName={name!} defaultView={props.defaultView} {...props} />
}

const AppImpl: React.SFC<AppProps> = (props) => {
    const { classes } = props;

    return <ThemeProvider theme={theme}>
        <div className={classes.root}>
            <CssBaseline />
            <div className={classes.app}>
                <Router>
                    <Switch>
                        <Route path="/job/:name/raw">
                            <IsReadonly uiClient={uiClient}>
                                <JobViewWithName client={client} defaultView="raw-logs" />
                            </IsReadonly>
                        </Route>
                        <Route path="/job/:name/results">
                            <IsReadonly uiClient={uiClient}>
                                <JobViewWithName client={client} defaultView="results" />
                            </IsReadonly>
                        </Route>
                        <Route path="/job/:name/logs">
                            <IsReadonly uiClient={uiClient}>
                                <JobViewWithName client={client} defaultView="logs" />
                            </IsReadonly>
                        </Route>
                        <Route path="/job/:name">
                            <IsReadonly uiClient={uiClient}>
                                <JobViewWithName client={client} />
                            </IsReadonly>
                        </Route>
                        <Route path="/github">
                            <GithubPage />
                        </Route>
                        <Route path="/start">
                            <StartJob client={client} uiClient={uiClient} />
                        </Route>
                        <Route path="/branches">
                            <IsReadonly uiClient={uiClient}>
                                <BranchList client={client} />
                            </IsReadonly>
                        </Route>
                        <Route path="/">
                            <IsReadonly uiClient={uiClient}>
                                <JobList client={client} />
                            </IsReadonly>
                        </Route>
                    </Switch>
                </Router >
                <footer className={classes.footer}>
                    <img src="/piro-small.png" alt="Bhojpur Piro logo" />
                </footer>
            </div>
        </div>
    </ThemeProvider>
}


const theme = function () {
    let theme = createMuiTheme({
        palette: {
            primary: {
                light: '#63ccff',
                main: '#39355B',// '#009be5',
                dark: '#006db3',
            },
        },
        typography: {
            fontFamily: [
                'Dosis',
                'sans-serif'
            ].join(', '),
            fontSize: 16,
            h5: {
                fontWeight: 500,
                fontSize: 26,
                letterSpacing: 0.5,
            },
        },
        shape: {
            borderRadius: 8,
        },
        props: {
            MuiTab: {
                disableRipple: true,
            },
        },
        mixins: {
            toolbar: {
                minHeight: 48,
            },
        },
    });

    theme = {
        ...theme,
        overrides: {
            MuiDrawer: {
                paper: {
                    backgroundColor: '#18202c',
                },
            },
            MuiButton: {
                label: {
                    textTransform: 'none',
                },
                contained: {
                    boxShadow: 'none',
                    '&:active': {
                        boxShadow: 'none',
                    },
                },
            },
            MuiTabs: {
                root: {
                    marginLeft: theme.spacing(1),
                },
                indicator: {
                    height: 3,
                    borderTopLeftRadius: 3,
                    borderTopRightRadius: 3,
                    backgroundColor: theme.palette.common.white,
                },
            },
            MuiTab: {
                root: {
                    textTransform: 'none',
                    margin: '0 16px',
                    minWidth: 0,
                    padding: 0,
                    [theme.breakpoints.up('md')]: {
                        padding: 0,
                        minWidth: 0,
                    },
                },
            },
            MuiIconButton: {
                root: {
                    padding: theme.spacing(1),
                },
            },
            MuiTooltip: {
                tooltip: {
                    borderRadius: 4,
                },
            },
            MuiDivider: {
                root: {
                    backgroundColor: '#404854',
                },
            },
            MuiListItemText: {
                primary: {
                    fontWeight: theme.typography.fontWeightMedium,
                },
            },
            MuiListItemIcon: {
                root: {
                    color: 'inherit',
                    marginRight: 0,
                    '& svg': {
                        fontSize: 20,
                    },
                },
            },
            MuiAvatar: {
                root: {
                    width: 32,
                    height: 32,
                },
            },
        },
    };
    return theme;
}();

const drawerWidth = 256;

const styles = createStyles({
    root: {
        display: 'flex',
        minHeight: '100vh',
    },
    drawer: {
        [theme.breakpoints.up('sm')]: {
            width: drawerWidth,
            flexShrink: 0,
        },
    },
    app: {
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
    },
    footer: {
        padding: theme.spacing(2),
        background: '#eaeff1',
        textAlign: 'center'
    },
});

export default withStyles(styles)(AppImpl);
