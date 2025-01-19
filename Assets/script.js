document.querySelector('#register-button').addEventListener('click', Register())

document.querySelector('#login_button').addEventListener('click', Login())

document
  .querySelector('#signup_switch_button')
  .addEventListener('click', function () {
    GoToHomePage()

   // document.querySelector('.login-container').style.display = 'none'
   // document.querySelector('.register-container').style.display = 'flex'
  })

document.querySelector('#login_switch_button').addEventListener('click', function () {
    document.querySelector('.register-container').style.display = 'none'
    document.querySelector('.login-container').style.display = 'flex'
  })

async function deleteCookie () {
  await fetch('/logout', {
    method: 'POST',
    body: JSON.stringify({ Session: document.cookie.split('=')[1] })
  })
  document.cookie = 'session=;expires=Tue, 22 Aug 2001 12:00:00 UTC;'
  window.location.href = '/'
}

async function Login () {
  let email = document.querySelector('input#email')  
  let password = document.querySelector('input#password')
  let data = { email: email.value, password: password.value }
  
  try {
    let response = await fetch('/signin', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage')
      errorMessage.innerHTML = errorData
    
  } catch (error) {
    const errorMessage = document.getElementById('errorMessage')
    errorMessage.innerHTML = 'Network error occurred!'
  }
}

async function Register () {
  let nickname = document.querySelector('input#Nickname')
  let age = document.querySelector('input#Age')
  let gender = document.querySelector('input#Gender')
  let first_Name = document.querySelector('input#First_Name')
  let last_Name = document.querySelector('input#Last_Name')
  let email = document.querySelector('input#email')
  let password = document.querySelector('input#password')
  let data = { nickname: nickname.value,age: age.value,email: gender.value,gender: first_Name.value,first_Name: last_Name.value, last_Name: email.value, password: password.value }
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

 function GoToHomePage() {
  document.body.innerHTML = ""


  let header = document.createElement('header');
  header.classList.add('header');
  header.innerHTML = `
      <div class="header-content">
          <a href="/">
              <img src="/styles/logo.png" alt="logo">
          </a>
          <nav class="nav-links">
              <a href="/">Home</a>
              <a href="/about">About</a>
          </nav>
          <div class="logout-container">
              <button class="logout-button" onclick="deleteCookie()">
                  <i class="fas fa-sign-out-alt"></i> Logout
              </button>
          </div>
      </div>
  `

  let forumContainer = document.createElement('div');
  forumContainer.classList.add('forum-container');
  forumContainer.innerHTML = `
      <main class="posts-container">
          <form action="/create" method="post">
              <div class="button-wrapper">
                  <button type="submit">
                      <i class="fas fa-plus-circle"></i> Create Post
                  </button>
              </div>
          </form>

          <h1>Posts</h1>
          <ul>
          </ul>
      </main>
  `;

  document.body.appendChild(header)
    document.body.appendChild(forumContainer)
    document.querySelector("link[rel='stylesheet']").href = "post.css"
  //  document.head.appendChild(document.createElement('link').href = "header.css")

    fetch("/post").then((response) => response.json()).then((e) => {
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
                            <button class="show-all-comments" name="id-post" onclick="window.location.href='/comment?id_comment={{ .ID }}&page=1'">Show all
                                Comment</button>
                    </div>

                    <form class="comment-form" action="/newcomment" method="POST">
                        <input type="text" name="comment" placeholder="Add a comment..." required>
                        <button type="submit" value="${data.id}" name="id-post">
                            <i class="fas fa-comment"></i>
                        </button>
                    </form>
                </li>`
      })
    })
}