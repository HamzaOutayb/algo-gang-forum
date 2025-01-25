if (document.cookie) {
  GoToHomePage()
}

function Start(){
  if ( document.querySelector('#register-button')){
    document.querySelector('#register-button').addEventListener('click', () => {Register()})
  }
if (document.querySelector('#login_button')) {
  document.querySelector('#login_button').addEventListener('click', () => { Login()})
}
  if (document.querySelector('#signup_switch_button')) {
  document.querySelector('#signup_switch_button').addEventListener('click', function () {
      document.querySelector('.login-container').style.display = 'none'
      document.querySelector('.register-container').style.display = 'flex'
    })
  }
  if (document
    .querySelector('#login_switch_button')){
  document
    .querySelector('#login_switch_button')
    .addEventListener('click', function () {
      document.querySelector('.register-container').style.display = 'none'
      document.querySelector('.login-container').style.display = 'flex'
    })
  }
}


async function deleteCookie () {
  await fetch('/logout', {
    method: 'POST',
    body: JSON.stringify({ session_token: document.cookie.split('=')[1] })
  })
  document.cookie = 'session_token=;expires=Tue, 22 Aug 2001 12:00:00 UTC;'
  GoToLoginPage()
}

async function Login (Login_re,key_re) {
  let email =  document.querySelector('input#email')  
  let password = document.querySelector('input#password')
  let data = { email:  Login_re || email.value, password:  key_re || password.value }  
  try {
    let response = await fetch('/signin', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage')
      errorMessage.innerHTML = errorData
    } else {
      GoToHomePage()
    }
  } catch (error) {
    const errorMessage = document.getElementById('errorMessage')
    errorMessage.innerHTML = 'Network error occurred!'
    }
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
  try {
    let response = await fetch('/signup', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage')
      errorMessage.innerHTML = errorData
    } else {
      Login(email,password)
    }
  } catch (error) {
    const errorMessage = document.getElementById('errorMessage')
    errorMessage.innerHTML = 'Network error occurred!'
  }
}

function GoToLoginPage() {
  if (document.querySelector("link[rel='stylesheet'][href='/Assets/post.css']")) {
    document.querySelector("link[rel='stylesheet'][href='/Assets/post.css']").href =  "/Assets/login.css"
  }
  document.body.innerHTML = ` 
  
  
  <div class="content-spacer"></div>


    <div class="login-container">
        <p id="errorMessage"></p>
        <h2>Login</h2>
        <div class="input-group">
            <label for="email">Email:</label>
            <input type="email" id="email" name="email" required>
        </div>
        <div class="input-group">
            <label for="password">Password:</label>
            <input type="password" id="password" name="password" required>
        </div>
        <div class="button-group">
            <button type="submit" id="login_button">Login</button>
        </div>
        <div class="register-link">
           <button id="signup_switch_button" >SIGN UP</button>
        </div>
    </div>



    <div class="register-container">
    
        <h2>Sign up</h2>
        <p id="errorMessage"></p>
        <div class="input-group">
            <label for="nickname">Nickname:</label>
            <input type="nickname" id="nickname" name="nickname" required>
        </div>

        <div class="input-group">
            <label for="age">Age:</label>
            <input type="age" id="age" name="age" required>
        </div>
        <div class="input-group">
            <label for="gender">Gender:</label>
            <input type="gender" id="gender" name="gender" required>
        </div>
        <div class="input-group">
            <label for="first_Name">First_Name:</label>
            <input type="first_Name" id="first_Name" name="first_Name" required>
        </div>
        <div class="input-group">
            <label for="last_Name">Last_Name:</label>
            <input type="last_Name" id="last_Name" name="last_Name" required>
        </div>
        <div class="input-group">
            <label for="email">Email:</label>
            <input type="email" id="email_re" name="email" required>
        </div>

        <div class="input-group">
            <label for="password">Password:</label>
            <input type="password" id="password_re" name="password" required>
        </div>

        <button id="register-button">Create Account</button>

        <div class="register-link">
            <button id="login_switch_button">login</button>
        </div>
</div>
    <script src="/Assets/script.js" defer></script>`
    Start()
  
}

async function GoToHomePage() {
  document.body.innerHTML = ""


  let header = document.createElement('header');
  header.classList.add('header');
  header.innerHTML = `
      <div class="header-content">
          <h3>ALGO GANG<h3>
          <div class="nav-links">
              <a href="/">Home</a>
              <a href="/about">About</a>
          </div>
          <div class="logout-container">
              <button class="logout-button" onclick="deleteCookie()">
                  <i class="fas fa-sign-out-alt"></i> Logout
              </button>
          </div>
          <div class="button-wrapper">
                  <button type="submit">
                      <i class="fas fa-plus-circle"></i> Create Post
                  </button>
              </div>
          </div>
      </div>
  `

  document.body.appendChild(header)
  document.body.innerHTML +=  `
   <aside class="sidebar-left">
           <h2>Contact<h2><br>

        </aside>
      <main class="posts-container">
          <h1>Posts</h1>
          <ul>
          </ul>
      </main>
  `;

  document.querySelector('.button-wrapper > button').addEventListener('click', (e) => {
     document.body.innerHTML += ` <div class="content-spacer-create"></div>
    <p class="errorMessage">{{.}}</p>

    <div class="create-post-container">
    <button class="X">X</button>
        <h2>Create Post</h2>
    
        <div>
            <label for="title">Title:</label>
            <input type="text" id="title" name="title" maxlength="50" minlength="10" required>
        </div>
        <div>
            <label for="content">Content:</label>
            <textarea id="content" name="content" maxlength="500" minlength="10" required></textarea>
        </div>
            <div>
                <label>Select Categories:</label>
                <div class="checkbox-container">
                    <input type="checkbox" id="news" name="category" value="news">
                    <label for="news" class="checkbox-label">News</label>

                    <input type="checkbox" id="tech" name="category" value="tech">
                    <label for="tech" class="checkbox-label">Technology</label>

                    <input type="checkbox" id="lifestyle" name="category" value="lifestyle">
                    <label for="lifestyle" class="checkbox-label">Lifestyle</label>

                    <input type="checkbox" id="education" name="category" value="education">
                    <label for="education" class="checkbox-label">Education</label>

                    <input type="checkbox" id="health" name="category" value="health">
                    <label for="health" class="checkbox-label">Health</label>

                    <input type="checkbox" id="entertainment" name="category" value="entertainment">
                    <label for="entertainment" class="checkbox-label">Entertainment</label>
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
  })
 
await fetch("/contact").then(response =>  response.json()).then(e => {
 let aside = document.querySelector('.sidebar-left')
  if (e){
  e.Contact_list.forEach((data)=> {
    aside.innerHTML += `<button class="users">${data}</button>`
  })
}


 })

 if (document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']")) {
  document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']").href =  "/Assets/post.css"
}
   await fetch("/post") .then((response) => response.json()).then( (e) => {
      if (e) {
      let ul = document.querySelector('ul')
      e.forEach((data)=> {
          ul.innerHTML += `  <li class="post-item" data-post-id="${data.id}">
                    <div class="username">${data.user_id}</div>
                    <h3>${data.title}</h3>
                    <div class="category">Category: ${data.categories }</div>
                    <p class="content-preview">${data.content }</p>
                    
                    <div class="post-date">${data.date }</div>

                    <!-- <div class="interaction-section"> -->
                    <div class="interaction-section">
                        <button class="like-comment-btn" name="like_post" value="${data.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-up"></i>
                            ${data.likes }
                        </button>
                        <button class="dislike-comment-btn" name="deslike_post" value="${data.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-down"></i>
                            ${data.dislikes}
                        </button>
                            <button class="show-all-comments" value="${data.id}" name="id-post">Show all
                                Comment</button>
                    </div>

                        <input type="text" name="comment" placeholder="Add a comment..." required>
                        <button type="submit" value="${data.id}" name="id-post">
                            <i class="fas fa-comment">add</i>
                        </button>
                </li>`
      })
    }
    })
    let showAllComments = document.querySelectorAll(`.show-all-comments`)
    showAllComments.forEach(e => e.addEventListener("click", async (e) => {
      let id = e.target.value;
      console.log(id);
      let post = await fetch(`/api/post/${id}`).then(response => response.json())
      console.log(post);
      
      if (post) {
        document.querySelector("main > ul").innerHTML = `<li class="post-item" data-post-id="${post.id}">
                    <div class="username">${post.user_id}</div>
                    <h3>${post.title}</h3>
                    <div class="category">Category: ${post.categories }</div>
                    <p class="content-preview">${post.content }</p>
                    
                    <div class="post-date">${post.date }</div>

                    <!-- <div class="interaction-section"> -->
                    <div class="interaction-section">
                        <button class="like-comment-btn" name="like_post" value="${post.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-up"></i>
                            ${post.likes }
                        </button>
                        <button class="dislike-comment-btn" name="deslike_post" value="${post.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-down"></i>
                            ${post.dislikes}
                        </button>
                        <input type="text" name="comment" placeholder="Add a comment..." required>
                        <button type="submit" value="${post.id}" name="id-post">
                            <i class="fas fa-comment">add</i>
                        </button>
                </li>`
      }
    }))
    

 function loop() {
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
         <input type="text" id="messageInput" class="message-input" placeholder="Type your message...">
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
              <div class="message_to">
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
  loop()
}
Start()

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
    const chatBox = document.getElementById('chatBox');
    chatBox.innerHTML += ` <h4>${parsedMessage.Sender} :</h4>
            <div class="message_to">
              <p>${parsedMessage.Message}</p>
            </div>
            <h6>${parsedMessage.Date}</h6></br>`;
      };
}