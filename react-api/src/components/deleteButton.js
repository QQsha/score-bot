import React from 'react'

  class DeleteButton extends React.Component {
    constructor(props) {
      super(props);
      this.eventClick = this.eventClick.bind(this);
    }

    eventClick() {
      fetch('https://chelsea-score-bot.herokuapp.com/delete_spam?spam='+this.props.word, {
        method: 'GET',
      });
      window.location.reload(false);
    }
    render() {
      return <button onClick={this.eventClick}>delete</button>;
    }
  }

  export default DeleteButton