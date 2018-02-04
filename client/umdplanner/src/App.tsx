import * as React from 'react';
import './App.css';
import {
  ApolloClient,
  gql,
  graphql,
  ApolloProvider,
  createNetworkInterface,
  ChildProps,
} from 'react-apollo';

const networkInterface = createNetworkInterface({ 
  uri: 'http://localhost:3001/query',
});

const client = new ApolloClient({
  networkInterface,
});

const cmsc250Query = gql`
  query ChannelsListQuery {
    class(code: "CMSC250") {
      code
      title
    }
  }
`;

type InputProps = {
  episode: string
};

export class ClassFormat extends React.Component<ChildProps<InputProps, Response>, {}> {
  render() {
    const { loading, out, error } = this.props.data;
    if (loading) {
      return <p>Loading ...</p>;
    }
    if (error) {
    return <p>{error.message}</p>;
    }
    return (
        <ul>
          {out.code}
        </ul>
      );
  }
}

const ClassWithData = graphql(cmsc250Query)(ClassFormat);

const logo = require('./logo.svg');

class App extends React.Component {
  render() {
    return (
      <ApolloProvider client={client}>
        <div className="App">
          <div className="App-header">
            <img src={logo} className="App-logo" alt="logo" />
            <h2>Welcome to React</h2>
          </div>
          <p className="App-intro">
            To get started, 2edit <code>src/App.tsx</code> and save to reload.
          </p>
          <ClassWithData />
        </div>
      </ApolloProvider>
    );
  }
}

export default App;
