var  wsLink
        
         
async  function  start() {
      let to = await document.querySelector('.friend-list').value;
      let divchat = document.createElement(`<div id="chat-box"></div>
        <div id="message-input">
                <input type="text" id="message" placeholder="Type a message...">
                <button id="send">Send</button>
        </div>`)
        divchat.style.position = "absolite"
        document.querySelector('body').appendChild(divchat)
    //creat div or block chat
         wsLink =  await `ws://localhost:8080/send?to=${to}`;
      
      if (user === ''|| to === '' || wsLink === '' ) {
          return;
      }
      console.log(wsLink);
      
      let ws = new WebSocket(wsLink)
   
      const chatBox = document.getElementById('message');
      const button = document.getElementById('send');
      button.addEventListener('click', () => {
          ws.send(chatBox.value);
          chatBox.value = '';
      });

      ws.onmessage = (message) => {
      const parsedMessage = JSON.parse(message.data);
      const chatBox = document.getElementById('chat-box');
     //check type dyal liwsal wax message wla xi 7ad tconecta wla deconecta
      chatBox.innerHTML += `<p>${parsedMessage.sender}: ${parsedMessage.message}</p>`;
        };
  };

  start()