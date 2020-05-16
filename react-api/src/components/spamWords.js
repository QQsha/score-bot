
    import React from 'react';
    import DeleteButton from './deleteButton';

    const Contacts = ({ contacts }) => {
      return (
        <div>
          <center><h1>Spam List</h1></center>
          {contacts.map((contact) => (
            <div class="card ">
              <div class="card-body d-flex justify-content-around flex-row bd-highlight mb-3">
                <h5 class="card-title p-2 bd-highlight">{contact.word}</h5>
                <h5 class="card-title mb-2 text-muted p-2 bd-highlight">{contact.ban}</h5>
                <DeleteButton word={contact.word}/>
              </div>
            </div>
          ))}
        </div>
      )
    };

    export default Contacts