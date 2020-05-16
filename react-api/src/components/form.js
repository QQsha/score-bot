import React from 'react'

class MyForm extends React.Component {
    constructor() {
      super();
      this.handleSubmit = this.handleSubmit.bind(this);
    }
  
    handleSubmit(event) {
      event.preventDefault();
      const data = new FormData(event.target);
      
      fetch('http://localhost:80/add_new_spam', {
        method: 'POST',
        body: data,
      });
      window.location.reload(false);
    }
  
    render() {
      return (
        <form onSubmit={this.handleSubmit}>
          <div class="row">
          <div class="col">
          <label htmlFor="spam">Enter spam word</label>
          <input id="spam" name="spam" type="text" class="form-control" />
          </div>
          <div class="col">
          <label htmlFor="ban">Enter ban duration</label>
          <input id="ban" name="ban" type="text" class="form-control" />
          </div>


  
          <button class="btn btn-primary">Save</button>
          </div>
        </form>
      );
    }
  }

//   <form>
//   <div class="row">
//     <div class="col">
//       <input type="text" class="form-control" placeholder="First name">
//     </div>
//     <div class="col">
//       <input type="text" class="form-control" placeholder="Last name">
//     </div>
//   </div>
// </form>

  export default MyForm