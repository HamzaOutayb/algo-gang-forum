async  function  start() {
      let to = await document.querySelector('#TO').value;
      if (to === '') {
          return;
      }
      
      let ws = new WebSocket(`ws://localhost:8080/send?to=${to}`)
   
      const chatBox = document.getElementById('messageInput');
      const button = document.querySelector('send-btn');
      button.addEventListener('click', () => {
        if (chatBox.value) {
          ws.send(chatBox.value);
          chatBox.value = '';
        }  
        });

      ws.onmessage = (message) => {
        const parsedMessage = JSON.parse(message.data);
      const chatBox = document.getElementById('chatBox');
      chatBox.innerHTML += ` <h4>${parsedMessage.Sender} :</h4>
              <div class="message_to">
                <p>${parsedMessage.Content}</p>
              </div>
              <h6>${parsedMessage.Date}</h6></br>`;
        };
    
  };

