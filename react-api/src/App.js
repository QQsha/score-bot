    // src/App.js

    import React, { Component } from 'react';
    import SpamList from './components/spamWords';
    import MyForm from './components/form';
    

    class App extends Component {
      state = {
        spamWords: []
      }
      componentDidMount() {
        fetch('https://chelsea-score-bot.herokuapp.com/get_spam')
        .then(res => res.json())
        .then((data) => {
          this.setState({ spamWords: data })
        })
        .catch(console.log)
      }
      render() {
        return (
          <div>
          <SpamList contacts={this.state.spamWords} />
          <br></br>
          <MyForm />
          
          </div>
     
          
        );
      }
      sad() {
        return (
          <MyForm  />
        );
      }
    }

    export default App;