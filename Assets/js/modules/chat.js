
import { FetchChatWithConversations, Status } from "./sidebar.js";
import { check_chat } from "./home.js";
import { GoToLoginPage } from "./auth.js";
import authUtils from "../utils/auth.js";

export let list = {
  liststatus: {},
}
let idtime
async function reload_sidebar() {
  let listaside = document.querySelector('.listaside')
  listaside.innerHTML = ""
  if (listaside && listaside.innerHTML === "") {
    await FetchChatWithConversations()
  }
}



export let ws
export async function StartWs() {
  ws = new WebSocket('ws://localhost:8080/chat');

  ws.onopen = () => {
    console.log('Connected');
  };

  ws.onmessage = async (message) => {
    authUtils.isLoggedIn().then(async (loggedIn) => {
      if (!loggedIn) {
        ws?.close()
        GoToLoginPage()
        document.cookie = 'session_token=;expires=Tue, 22 Aug 2001 12:00:00 UTC;'
        return
      }
    })



    let id_to = document.querySelector("#TO")?.getAttribute("value");
    let username_to = document.querySelector("#TO")?.innerHTML;
    try {


      const parsedData = JSON.parse(message.data);

      if (parsedData.Status) {
        console.log(parsedData.Status);
        list.liststatus = parsedData.Status
        await reload_sidebar()
      } else if (parsedData.message) {
        await reload_sidebar()
        if (!check_chat.inchat) {
          const NOTIFICATION = document.querySelector(".NOTIFICATION")
          NOTIFICATION.innerHTML = `
          <h4 >NEW MESSAGE FROM ${parsedData.sender}</h4>`
          NOTIFICATION.style.display = "block"
          setTimeout(() => {
            NOTIFICATION.style.display = "none"
            NOTIFICATION.innerHTML = ""
          }, 1500)
        } else {
          if (parsedData.sender == username_to) {
            const chatBox = document.getElementById('chatBox');
            console.log("received a message")
            if (chatBox) {
              chatBox.innerHTML += ` <div class=${parsedData.to == id_to ? "Message_TO" : "Message_From"}>
                  <h4 >${parsedData.sender}</h4>
                    <span>${parsedData.message}</span>
                  <h6>${parsedData.Date.split(".")[0]}</h6>
                  </div></br>`;
            }
          }
        }
      } else if (parsedData.istyping) {
        if (username_to == parsedData.sender && check_chat.inchat) {
          let typing = document.querySelector(".typing-indicator")
          if (typing) {
            if (idtime) {
              clearTimeout(idtime)
            }
            typing.innerHTML = `${parsedData.sender}  is typing<img src="Assets/JVX7.gif" alt="loding"> `
            idtime = setTimeout(() => {
              typing.innerHTML = ``
            }, 1500)
          }
        }
      }

    } catch (error) {
      console.error('Error parsing message data:', error);
    }
  };

  ws.onclose = () => {
    console.log('Connection closed');
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };
}









let page = 1

export function ChatBox() {
  let users = document.querySelectorAll("button.users")
  console.log(users);

  users.forEach(e => e.addEventListener("click", async () => {
    let chatContainer = document.querySelector(".chat-container")
    let chatbox = chatContainer.querySelector("#chatBox")
    chatbox.innerHTML = ""
    chatbox.removeEventListener("scroll", fetchhistory)
    console.log(chatContainer.querySelector("#chatBox"));
    let TO = e.innerHTML.split("<")[0].trim();
    let TO_id = e.value
   page = 1
    chatContainer.querySelector("#TO").setAttribute("value", TO_id)
    chatContainer.querySelector("#TO").innerHTML = TO
    chatContainer.style.display = "flex"
    const data = { message: TO, to: TO_id }


    await fetch(`/api/chathistory/${page}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    }).then(response => response.json()).then(data => {
      if (data) {
        const chatbox = document.querySelector("#chatBox")
        data.slice().reverse().forEach(e => {
          chatbox.innerHTML += `
          <div class=${e.Sender == TO ? "Message_From" : "Message_TO"}>
          <h4 >${e.Sender}</h4>
          <span>${e.Content}</span>
          <h6>${e.Created_at}</h6>
          </div></br>
          `
        })
        chatbox.scrollTop = chatbox.scrollHeight
        console.log(page)
        chatbox.addEventListener("scroll", fetchhistory)
      }
      startchat()
    })
  }))

}
async function fetchhistory() {
  let chatbox = document.querySelector("#chatBox")
  let TO = document.querySelector(".chat-container #TO").innerHTML.split("<")[0].trim();
  let TO_id = document.querySelector(".chat-container #TO").getAttribute("value");
      let scrollHeight = chatbox.scrollHeight
      if (chatbox.scrollTop == 0) {
        const data = { message: TO, to: TO_id }
        fetch(`/api/chathistory/${++page}`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify(data)
        }).then(response => response.json()).then(data => {
          if (data) {
            console.log(TO,TO_id,data);
            const chatbox = document.querySelector("#chatBox")
            data.forEach(e => {
              chatbox.innerHTML = `
              <div class=${e.Sender == TO ? "Message_From" : "Message_TO"}>
                  <h4>${e.Sender}</h4>
                  <span>${e.Content}</span>
                  <h6>${e.Created_at}</h6>
              </div></br>
              ` + chatbox.innerHTML;
            });
            chatbox.scrollTop = chatbox.scrollHeight - scrollHeight
          }
        })
      }
    
}



async function startchat() {
  document.body.style.overflow = "hidden";
  document.querySelector(".X").addEventListener("click", async () => {
    check_chat.inchat = false
    let chatContainer = document.querySelector(".chat-container")
    chatContainer.querySelector("#TO").setAttribute("value", "")
    chatContainer.querySelector("#TO").innerHTML = ""
    chatContainer.querySelector("#chatBox").innerHTML = ""
    chatContainer.querySelector("#chatBox").removeEventListener("scroll", fetchhistory)
    chatContainer.style.display = "none"
    document.body.style.overflow = "auto";
  })
  check_chat.inchat = true

  const chatBox = document.getElementById('messageInput');
  chatBox.addEventListener('keydown', async (event) => {
    if (event.key === 'Enter') {
      if (chatBox.value.trim() === '') {
        return;
      }
      if (chatBox.value.length > 500) {
        return
      }

      function escapeHTML(str) {
        return str.replace(/[&<>"']/g, (char) => {
          const escapeChars = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#39;',
          };
          return escapeChars[char];
        });
      }

      let id_to = document.querySelector("#TO").getAttribute("value");
      const DM = escapeHTML(chatBox.value)

      ws.send(JSON.stringify({ message: DM, to: parseInt(id_to) }));

      const ChatBox = document.getElementById('chatBox');
      if (ChatBox) {
        ChatBox.innerHTML += ` <div class="Message_TO">
            <h4 >${localStorage.getItem("username")}</h4>
              <span>${DM}</span>
            <h6>${new Date().toTimeString()}</h6>
            </div></br>`;
      }
      chatBox.value = '';
      await reload_sidebar()
    } else {
      let id_to = document.querySelector("#TO").getAttribute("value");
      ws.send(JSON.stringify({ istyping: true, to: parseInt(id_to) }));
    }
  })
}