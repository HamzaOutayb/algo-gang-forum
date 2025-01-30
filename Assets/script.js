var page_posts = 1;
var page_comments = 1;
var page_messages = 1;
var page_users = 1;
var page_conversations = 1;
let nomoreposts = false;
let nomorecomment = false;
let nomoremessage = false;
let nomoreusers = false;
let nomoreconversations = false;
var NofetchComment = false


if (document.cookie) {
  GoToHomePage()
}

function Login_page(){
if ( document.querySelector('#register-button')){
    document.querySelector('#register-button').addEventListener('click', Register)
  }
if (document.querySelector('#login_button')) {
  document.querySelector('#login_button').addEventListener('click', Login)
}
if (document.querySelector('#signup_switch_button')) {
  document.querySelector('#signup_switch_button').addEventListener('click', function () {
      document.querySelector('.login-container').style.display = 'none'
      document.querySelector('.register-container').style.display = 'flex'
    })
  }
if (document.querySelector('#login_switch_button')){
  document.querySelector('#login_switch_button').addEventListener('click', function () {
      document.querySelector('.register-container').style.display = 'none'
      document.querySelector('.login-container').style.display = 'flex'
    })
  }
}


async function deleteCookie () {
  // await fetch('/logout', {
  //   method: 'POST',
  //   body: JSON.stringify({ session_token: document.cookie.split('=')[1] })
  // })
  document.cookie = 'session_token=;expires=Tue, 22 Aug 2001 12:00:00 UTC;'
  GoToLoginPage()
}

async function Login (Login_re,key_re) {
  let email =  document.querySelector('input#email')  
  let password = document.querySelector('input#password')
  const errorMessage = document.getElementById('errorMessage')
  let data = { email:  Login_re?.value || email.value, password:  key_re?.value || password.value }    
  
    let response = await fetch('/Signin', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
     
      errorMessage.classList.add("errorMessage")
      errorMessage.innerHTML = errorData
    } else {
      
      GoToHomePage()
    }
  
    // errorMessage.classList.add("errorMessage")
    // errorMessage.innerHTML = 'Network error occurred!'
    
}

async function Register () {
  let nickname = document.querySelector('input#nickname')
  let age = document.querySelector('input#age')
  let gender = document.querySelector('input#gender')
  let first_Name = document.querySelector('input#first_Name')
  let last_Name = document.querySelector('input#last_Name')
  let email = document.querySelector('input#email_re')
  let password = document.querySelector('input#password_re')
  let data = { nickname: nickname.value,age: age.value,email: email.value,gender: gender.value,first_Name: first_Name.value, last_Name: last_Name.value, password: password.value }
  
    let response = await fetch('/Signup', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage2')
      errorMessage.classList.add("errorMessage")
      errorMessage.innerHTML = errorData
    } else {
      Login(email,password)
    }
  
}

function GoToLoginPage() {
  if (document.querySelector("link[rel='stylesheet'][href='/Assets/post.css']")) {
    document.querySelector("link[rel='stylesheet'][href='/Assets/post.css']").href =  "/Assets/login.css"
  }
  document.body.innerHTML = ` <div class="content-spacer"></div>


    <div class="login-container">
        <p id="errorMessage"></p>
        <h2>Login</h2>
        <div class="input-group">
            <input type="email" id="email" name="email"  placeholder="Email:"required>
        </div>
        <div class="input-group">
            <input type="password" id="password" name="password"  placeholder="Password:" required>
        </div>
        <div class="button-group">
            <button type="submit" id="login_button">Login</button>
        </div>
        <div class="register-link">
           <button id="signup_switch_button" >SIGN UP</button>
        </div>
    </div>



    <div class="register-container">
        <p id="errorMessage2"></p>
        <h2>Sign up</h2>
        <div class="input-group">
            <input type="nickname" id="nickname" name="nickname"  placeholder="Nickname:" required>
        </div>

        <div class="input-group">
            <input type="age" id="age" name="age" placeholder="Age:" required>
        </div>
        <div class="input-group">
            <input type="gender" id="gender" name="gender" placeholder="Gender:" required>
        </div>
        <div class="input-group">
            <input type="first_Name" id="first_Name" name="first_Name" placeholder="First_Name:" required>
        </div>
        <div class="input-group">
            <input type="last_Name" id="last_Name" name="last_Name" placeholder="Last_Name:" required>
        </div>
        <div class="input-group">
            <input type="email" id="email_re" name="email" placeholder="Email:" required>
        </div>

        <div class="input-group">
            <input type="password" id="password_re" name="password" placeholder="Password:" required>
        </div>

        <button id="register-button">Create Account</button>

        <div class="register-link">
            <button id="login_switch_button">login</button>
        </div>
</div>
    <script src="/Assets/script.js" defer></script>`
    Login_page()
  
}

async function GoToHomePage() {
  document.body.innerHTML = ""


  let header = document.createElement('header');
  header.classList.add('header');
  header.innerHTML = `
      <div class="header-content">
          <h3 class="logo">ALGO GANG</h3>

          <div class="logout-container">
              <button class="logout-button" onclick="deleteCookie()">
                  <i class="fas fa-sign-out-alt"></i> Logout
              </button>
               <button type="submit" class="button-wrapper">
                      <i class="fas fa-plus-circle"></i> Create Post
                  </button>
          </div>
          </div>

  `

  document.body.appendChild(header)

  document.body.innerHTML +=  `
   <aside class="sidebar-left">
           <h2>Contact</h2>
        </aside>
      <main class="posts-container">
          <h1>Posts</h1>
          <ul>
          </ul>
      </main>
  `;

ShowCreatePost()
FetchChatWithConversations()
FetchConversations()
 if (document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']")) {
  document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']").href =  "/Assets/post.css"
}

GetAllPosts()
  
  
}
Login_page()


function ShowCreatePost() {
  document.querySelector('.button-wrapper').addEventListener('click', () => {
    document.body.innerHTML += ` <div class="content-spacer-create"></div>
   <p class="errorMessage"></p>

   <div class="create-post-container">
   <button class="X">X</button>
       <h2>Create Post</h2>
   
       <div>
           <input type="text" id="title" name="title" maxlength="50" minlength="10" placeholder="Title:"required>
       </div>
       <div>
           <textarea id="content" name="content" maxlength="500" minlength="10" placeholder="Content:" required></textarea>
       </div>
           <div>
               <label>Select Categories:</label>
               <div class="checkbox-container">
                   <input type="checkbox" class="checkbox" class="checkbox" id="javascript" name="category" value="javascript">
                   <label for="javascript" class="checkbox-label">javascript</label>

                   <input type="checkbox" class="checkbox" id="tech" name="category" value="tech">
                   <label for="tech" class="checkbox-label">Technology</label>

                   <input type="checkbox" class="checkbox" id="golang" name="category" value="golang">
                   <label for="golang" class="checkbox-label">golang</label>

                   <input type="checkbox" class="checkbox" id="rust" name="category" value="rust">
                   <label for="rust" class="checkbox-label">rust</label>

                   <input type="checkbox" class="checkbox" id="programming" name="category" value="programming">
                   <label for="programming" class="checkbox-label">programming</label>

                 
               </div>
           </div>
           <button type="submit" id="create-post-button">Create Post</button>
   </form>
   </div>`
   
   document.body.style.overflow = "hidden";
   document.querySelector(".X").addEventListener("click", () => {
     GoToHomePage()
     document.body.style.overflow = "auto";
   })
   document.querySelector("#create-post-button").addEventListener("click", () => {
     CreatePost()
   })
 })
}

async function FetchChatWithConversations() {
  await fetch("/ChatWithConversations/").then(response =>  response.json()).then(e => {
    let aside = document.querySelector('.sidebar-left')
     if (e){
       let listaside = document.createElement('div')
       listaside.classList.add('listaside')
     e.forEach((data)=> {
       listaside.innerHTML += `<button class="users" value="${data.id}">${data.nickname}</button>`
     })
     aside.appendChild(listaside)
   }
   })
  }

async function FetchConversations() {
  await fetch("/Conversations/").then(response =>  response.json()).then(e => {
    let aside = document.querySelector('.sidebar-left')
     if (e){
      aside.innerHTML += `<h2>Conversations</h2>`
      let listaside = document.createElement('div')
      listaside.classList.add('listaside')
    e.forEach((data)=> {
      listaside.innerHTML += `<button class="users" value="${data.id}">${data.nickname}</button>`
     })
     aside.appendChild(listaside)
   }
   })
}


async function GetAllPosts(page = 1) {
  if (nomoreposts) {
    NofetchComment = false
    return;
  }
  await fetch(`/api/post?page=${page}`) .then((response) => response.json()).then( async (e) => {
    if (e) {
    let ul = document.querySelector('ul')
    await e.forEach((data)=> {
        ul.innerHTML += `  <li class="post-item" data-post-id="${data.id}">
       
                  <div class="username">${data.author}</div>
                  <h3 class="post-title">${data.title}</h3>
                  <div class="category">Category: ${data.categories?.join(' - ') || "No Gategory" }</div>
                  <h4 class="content-preview">${data.content}</h4>
                  
                  <div class="post-date">${data.date }</div>

                  <!-- <div class="interaction-section"> -->
                  <div class="interaction-section">
                      <button class="like-post-btn ${data.isliked ? "like-reacted" : ""}" name="like_post" value="${data.id}" id="likes"
                          onclick="">
                          <i class="fas fa-thumbs-up"></i>
                          ${data.likes }
                      </button>
                      <button class="dislike-post-btn ${data.isdisliked ? "dislike-reacted" : ""}" name="deslike_post" value="${data.id}" id="likes"
                          onclick="">
                          <i class="fas fa-thumbs-down"></i>
                          ${data.dislikes}
                
                  </div>

                      <input type="text" name="comment" placeholder="Add a comment..." required>
                      <button type="submit" value="${data.id}" name="id-post">
                          <i class="fas fa-comment">add</i>
                      </button>
              </li>`
    })
  }
  }).catch(e => {
    nomoreposts = true; 
  })
  Likes_Posts()
  document.querySelector('h3.logo').addEventListener('click', GoToHomePage)
    GetSinglePost()
    InsertComment()
    ChatBox()
    let debounceTimer
    let done = false
    window.addEventListener("scroll", function() {
      if (window.scrollY + window.innerHeight >= document.body.scrollHeight - 100) {
        if (!done && !NofetchComment){
          GetAllPosts(++page_posts)
          clearTimeout(debounceTimer);
          done = true
        }
        debounceTimer = setTimeout(() => {
          done = false
        }, 1000);
          
      }
  });

}


function InsertComment() {
  let CommentBtn = document.querySelectorAll(`button[type="submit"][name="id-post"]`)
    CommentBtn.forEach(e => e.addEventListener("click", async (e) => {
      // let id = e.target.value;
      const id = e.target.closest('.post-item').getAttribute('data-post-id');
      const postsinput = document.querySelector(`.post-item[data-post-id="${id}"] input[name="comment"]`);
      const data = { postId: parseInt(id), content: String(postsinput.value) }
      console.log(data)
      await fetch('/comment', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },  
        body: JSON.stringify(data)
      }).catch(e => {
        console.log(e)
      })
      postsinput.value = '';
    }))
}

function GetSinglePost() {
  let showAllComments = document.querySelectorAll(`.post-item`)
    showAllComments.forEach(e => e.addEventListener("click", async (e) => {
      let id = await e.target.getAttribute("data-post-id");
     if (id) {
      let post = await fetch(`/api/post/${id}`).then(response => response.json())      
      if (post) {        
        document.querySelector("main > ul").innerHTML = `<li class="post-item" data-post-id="${post.id}">
                    <div class="username">${post.author}</div>
                    <h3>${post.title}</h3>
                    <div class="category">Category: ${post.categories?.join(' - ') || "No Gategory" }</div>
                    <p class="content-preview">${post.content }</p>
                    
                    <div class="post-date">${post.date }</div>

                    <!-- <div class="interaction-section"> -->
                    <div class="interaction-section">
                      <button class="like-post-btn ${post.isliked ? "like-reacted" : ""}" name="like_post" value="${post.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-up"></i>
                            ${post.likes }
                        </button>
                        <button class="dislike-post-btn ${post.isdisliked ? "dislike-reacted" : ""}" name="deslike_post" value="${post.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-down"></i>
                            ${post.dislikes}
                        </button>
                    </div>
                        <input type="text" name="comment" placeholder="Add a comment..." required>
                        <button type="submit" value="${post.id}" name="id-post">
                            <i class="fas fa-comment">add</i>
                        </button>
                </li>`
                NofetchComment = true
      }
      GetAllComment(id)
      
      
  let debounceTimer
    window.addEventListener("scroll", function() {
      if (window.scrollY + window.innerHeight >= document.body.scrollHeight - 100) {
       
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
          console.log('scrolling')
          GetAllComment(id,++page_comments);
        }, 1000);
          
      }
  });
    }
    }))
}

async function GetAllComment(id,page_comments = 1) {
  if (nomorecomment) {
    return
  }
  let Comment = await fetch(`/api/GetComments/${id}/?page=${page_comments}`).then(response => response.json())      
  if (Comment) { 
    const commentList = document.querySelector("main > ul")
    console.log(Comment);
    
   await Comment.forEach((comment) => { 
      commentList.innerHTML += `<li class="comment-item" data-comment-id="${comment.id}">
                <div class="username">${comment.author}</div>
                <p class="content-preview">${comment.content }</p>
                <div class="post-date">${comment.date }</div>
                <div class="interaction-section">
                    <button class="like-comment-btn ${comment.isliked ? "like-reacted" : ""}" name="like_post"  value="${comment.id}" id="likes">
                        <i class="fas fa-thumbs-up"></i>
                        ${comment.likes }
                    </button>
                    <button class="dislike-comment-btn ${comment.isdisliked ? "dislike-reacted" : ""}" name="deslike_post" value="${comment.id}" id="likes">
                        <i class="fas fa-thumbs-down"></i>
                        ${comment.dislikes}
                    </button>
            </li>`})
    Likes_Comments()
    InsertComment()
    Likes_Posts()    
  }
}
async function Likes_Comments() {
  document.querySelectorAll('.like-comment-btn').forEach(e => e.addEventListener('click', async (e) => {
    
    const id =  await e.target.closest('.comment-item').getAttribute('data-comment-id');
    
    const data = { thread_type: 'comment', thread_id: parseInt(id), react: 1 }
    let response = await fetch('/api/reaction', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
      })
      response = await response.json()
      const dislikeButton = await e.target.closest('.comment-item').querySelector('button.dislike-comment-btn')
       if (response.isliked){
          e.target.classList.add("like-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }else {
          e.target.classList.remove("like-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }
        if (response.isdisliked){
          dislikeButton.classList.add("dislike-reacted")
          dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }else {
          dislikeButton.classList.remove("dislike-reacted")
          dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }
    
    
    }))
    document.querySelectorAll('.dislike-comment-btn').forEach(e => e.addEventListener('click', async (e) => {
     
      const id = await e.target.closest('.comment-item').getAttribute('data-comment-id');
      
      const data = { thread_type: 'comment', thread_id: parseInt(id), react: 2 }
      let response = await fetch('/api/reaction', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
        })
        response = await response.json()
        
        const likeButton =  await e.target.closest('.comment-item').querySelector('button.like-comment-btn')
        if (response.isliked){
          likeButton.classList.add("like-reacted")
          likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }else {
          likeButton.classList.remove("like-reacted")
          likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }
        if (response.isdisliked){
          e.target.classList.add("dislike-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }else {
           e.target.classList.remove("dislike-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }

      }))
}


function ChatBox() {
  let users = document.querySelectorAll("button.users")
  
  users.forEach(e => e.addEventListener("click", async () => {
     var TO = e.innerHTML;
     
     document.querySelector("main").innerHTML = `
     <div class="chat-container">
     <button class="X">X</button>
     <h3 id="TO">${TO}</h3>
   <div class="chat-box" id="chatBox">
   </div>
   <div class="input-area">
     <input type="text " id="messageInput" class="message-input" placeholder="Type your message...">
     <button class="send-btn">Send</button>
   </div>
 </div>
     `
     await fetch("/api/chathistory", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
      },
      body: JSON.stringify({ message: TO }) 
     }).then(response => response.json()).then(data => {
      if (data) {
       const chatbox = document.querySelector("#chatBox")
       data.forEach(e => {
        if (e.Sender == TO) {
          chatbox.innerHTML += `
          <h4>${e.Sender} :</h4>
          <div class="message_to">
            <p>${e.Content}</p>
          </div>
          <h6>${e.Created_at}</h6></br>
          `
        }else{
          chatbox.innerHTML += `
         <h4>${e.Sender} :</h4>
          <div class="message_sender">
            <p>${e.Content}</p>
          </div>
          <h6>${e.Created_at}</h6></br>
          `
        }
       })
      }
      startchat()
      document.querySelector(".X").addEventListener("click", () => {
        GoToHomePage()
      })
     })
  }))
  
}

async  function  startchat() {
  let to = document.querySelector('#TO').innerHTML;
  const webs =  `ws://localhost:8080/chat?to=${to}`
  if (to === '') {
      return;
  }
  console.log("test",webs,to);
  
  let ws = new WebSocket(webs)

  const chatBox = document.getElementById('messageInput');
  const button = document.querySelector('.send-btn');
  button.addEventListener('click', () => {
      ws.send(chatBox.value);
      chatBox.value = '';
  });

  ws.onmessage = (message) => {
  const parsedMessage = JSON.parse(message.data);
  console.log(message)
    const chatBox = document.getElementById('chatBox');
    if (message.Sender == TO) {
    chatBox.innerHTML += ` <h4>${parsedMessage.Sender} :</h4>
            <div class="message_to">
              <p>${parsedMessage.Message}</p>
            </div>
            <h6>${parsedMessage.Date}</h6></br>`;
    }else {
      chatBox.innerHTML += ` <h4>${parsedMessage.Sender} :</h4>
      <div class="message_sender">
        <p>${parsedMessage.Message}</p>
      </div>
      <h6>${parsedMessage.Date}</h6></br>`;
    }
}
}



async function CreatePost() {
  const title = document.querySelector('#title').value
  const content = document.querySelector('#content').value
  const categories = Array.from(document.querySelectorAll('input[name="category"]:checked')).map(e => e.value)
  const data = { title: title, content: content, categories: categories }
  await fetch('/create_post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
    }).catch(e => {
      console.log(e)
      return
    })
    GoToHomePage()
    document.body.style.overflow = "auto";
}







 function Likes_Posts() {
  document.querySelectorAll('.like-post-btn').forEach(e => e.addEventListener('click', async (e) => {
    
    const id =  await e.target.closest('.post-item').getAttribute('data-post-id');
    
    const data = { thread_type: 'post', thread_id: parseInt(id), react: 1 }
    let response = await fetch('/api/reaction', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
      })
      response = await response.json()
      const dislikeButton = await e.target.closest('.post-item').querySelector('button.dislike-post-btn')
       if (response.isliked){
          e.target.classList.add("like-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }else {
          e.target.classList.remove("like-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }
        if (response.isdisliked){
          dislikeButton.classList.add("dislike-reacted")
          dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }else {
          dislikeButton.classList.remove("dislike-reacted")
          dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }
    
    
    }))
    document.querySelectorAll('.dislike-post-btn').forEach(e => e.addEventListener('click', async (e) => {
     
      const id = await e.target.closest('.post-item').getAttribute('data-post-id');
      
      const data = { thread_type: 'post', thread_id: parseInt(id), react: 2 }
      let response = await fetch('/api/reaction', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
        })
        response = await response.json()
        
        const likeButton =  await e.target.closest('.post-item').querySelector('button.like-post-btn')
        if (response.isliked){
          likeButton.classList.add("like-reacted")
          likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }else {
          likeButton.classList.remove("like-reacted")
          likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
        }
        if (response.isdisliked){
          e.target.classList.add("dislike-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }else {
           e.target.classList.remove("dislike-reacted")
          e.target.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
        }

      }))
    }
    // let is = false;
    // let done = false;
    // const width = window.innerWidth;
    // if (width < 768) {
    //   console.log('width', width);
      
    // Resize();
    // }
    
    // window.addEventListener('resize', Resize);

    // function Resize() {
    //    const width = window.innerWidth;
      
    //     if (width < 768) {
    //       if (!done) {
    //         done = true;
    //         const buttonaside = document.createElement('button');
    //         buttonaside.classList.add('buttonaside');
    //         document.body.appendChild(buttonaside);
      
         
    //         buttonaside.addEventListener('click', () => {
    //           if (is) {
    //             console.log('is0', is);
    //             document.querySelector('.sidebar-left').style.display = 'none';
    //           } else {
    //             console.log('is1', is);
    //             document.querySelector('.sidebar-left').style.display = 'block';
    //           }
    //           is = !is; 
    //         });
    //       }
    //     } else {
    //       done = false;
    //       is = false; 
    //       document.querySelector('.buttonaside')?.remove();
    //       document.querySelector('.sidebar-left')?.style.display = 'block'; 
    //     }
    //   }