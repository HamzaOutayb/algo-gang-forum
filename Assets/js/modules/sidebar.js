import { ChatBox } from "./chat.js";
import { list } from "./chat.js";


export function Status(data) {
  console.log(data);
  
    document.querySelectorAll('.users').forEach(e => {
  
      for (let key in data) {
        if (e.value == key) {
          if (data[key]) {
            const statusElement = e.querySelector('.status');
            if (statusElement) {
              statusElement.innerHTML = 'online';
              statusElement.classList.remove("status_offline")
              statusElement.classList.add("status_online")
            }
          } else {
            const statusElement = e.querySelector('.status');
            if (statusElement) {
              statusElement.innerHTML = 'offline';
              statusElement.classList.remove("status_online")
              statusElement.classList.add("status_offline")
  
            }
          }
        }
      }
    });
  }
  

export async function FetchChatWithConversations() {
  console.log("FetchChatWithConversations");
  await fetch("/ChatWithConversations/").then(response => response.json()).then(data  => {
    console.log(data);
    if (data) {
      let listaside = document.querySelector('.listaside')
      listaside.innerHTML = ""
      if (listaside && listaside.innerHTML === "") {
        data.slice().reverse().forEach((data) => {
          listaside.innerHTML += `<button class="users" value="${data.friendid}">${data.nickname}
          <p class="status status_offline">offline</p>
          </button>`
        })
      }
    }
FetchConversations()
   }).catch (e => {
     console.log(e);
   })
}

  
 export async function FetchConversations() {
   await fetch("/Conversations/").then(response => response.json()).then(data => {
     let listaside = document.querySelector('.listaside')
     if (data && listaside) {
       data.slice().reverse().forEach((data) => {
         listaside.innerHTML += `<button class="users" value="${data.friendid}">${data.nickname}
     <p class="status status_offline">offline</p>
     </button>`
       })
     }
   })
   Status(list.liststatus)
   ChatBox()
  }