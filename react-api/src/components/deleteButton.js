import React from 'react'

        // function fetch() {
        //     return new Promise(resolve => setTimeout(() => resolve(42), 1000));
        //   }
        
        //   function fetchAPI(param) {
        //     // param is a highlighted word from the user before it clicked the button
        //     return fetch("http://localhost:80/delete_spam?param=" + param);
        //   }
  
//   class DeleteButton2 extends React.Component {
//     state = { result: null };
  
//     toggleButtonState = () => {
//       let selectedWord = window.getSelection().toString();
//       fetchAPI(selectedWord).then(result => {
//         this.setState({ result });
//       });
//     };
  
//     render() {
//       return (
//         <div>
//           <button onClick={this.toggleButtonState}> Click me </button>
//           <div>{this.state.result}</div>
//         </div>
//       );
//     }
//   }

  class DeleteButton extends React.Component {
    constructor(props) {
      super(props);
      this.eventClick = this.eventClick.bind(this);
    }

    eventClick() {
      fetch('http://localhost:80/delete_spam?spam='+this.props.word, {
        method: 'GET',
      });
      window.location.reload(false);
    }
    render() {
      return <button onClick={this.eventClick}>delete</button>;
    }
  }

  export default DeleteButton