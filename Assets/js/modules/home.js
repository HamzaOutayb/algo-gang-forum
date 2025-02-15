import { Resize } from "../index.js";
import { FetchChatWithConversations } from "./sidebar.js";
import { StartWs } from "./chat.js";
import { deleteCookie } from "./auth.js";

export let ress = { is_resize: false, done_resize: false }
let page_posts = 1;
let page_comments = 1;
let nomoreposts = false;
let nomorecomment = false;
let NofetchComment = false
export let check_chat = { inchat: false, isstarted: false }


export async function CreatePost() {
  const title = document.querySelector('#title').value
  const content = document.querySelector('#content').value
  const error = document.querySelector('.errorMessage')
  const categories = Array.from(document.querySelectorAll('input[name="category"]:checked')).map(e => e.value)
  if (title.length < 2) {
    error.innerHTML = "Title must be at least 2 characters long"
    error.style.display = "block"
    return
  }
  if (content.length < 2) {
    error.innerHTML = "Content must be at least 2 characters long"
    error.style.display = "block"
    return
  }
  if (content.length > 500) {
    error.innerHTML = "Content must be less than 500 characters long"
    error.style.display = "block"
    return
  }
  if (title.length > 50) {
    error.innerHTML = "Title must be less than 100 characters long"
    error.style.display = "block"
    return
  }
  const data = { title: title, content: content, categories: categories }
  await fetch('/create_post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  }).then(response => {
    if (response.ok) {
      GoToHomePage()
      document.body.style.overflow = "auto";
    }
  }).catch(e => {
    error.innerHTML = e
    error.style.display = "block"
    return
  })

}


export function Likes_Posts() {
  document.querySelectorAll('.like-post-btn').forEach(e => e.addEventListener('click', async (e) => {
    const currentTarget = e.currentTarget

    const id = await e.target.closest('.post-item').getAttribute('data-post-id');
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

    if (response.isliked) {
      currentTarget.classList.add("like-reacted")
      currentTarget.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
    } else {
      currentTarget.classList.remove("like-reacted")
      currentTarget.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
    }
    if (response.isdisliked) {
      dislikeButton.classList.add("dislike-reacted")
      dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
    } else {
      dislikeButton.classList.remove("dislike-reacted")
      dislikeButton.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
    }


  }))
  document.querySelectorAll('.dislike-post-btn').forEach(e => e.addEventListener('click', async (e) => {
    const currentTarget = e.currentTarget

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

    const likeButton = await e.target.closest('.post-item').querySelector('button.like-post-btn')
    if (response.isliked) {
      likeButton.classList.add("like-reacted")
      likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
    } else {
      likeButton.classList.remove("like-reacted")
      likeButton.innerHTML = `<i class="fas fa-thumbs-up"></i> ${response.Like}`
    }
    if (response.isdisliked) {
      currentTarget.classList.add("dislike-reacted")
      currentTarget.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
    } else {
      currentTarget.classList.remove("dislike-reacted")
      currentTarget.innerHTML = `<i class="fas fa-thumbs-down"></i> ${response.Dislike}`
    }

  }))
}



export async function GetAllComment(id, page_comments = 1) {
  if (nomorecomment) {
    return
  }
  let Comment = await fetch(`/api/GetComments/${id}/?page=${page_comments}`).then(response => response.json())
  if (Comment) {
    const commentList = document.querySelector(".comment_container")
    await Comment.forEach((comment) => {
      commentList.innerHTML += `<li class="comment-item">
                  <div class="username">${comment.author}</div>
                  <p class="content-preview">${comment.content}</p>
                  <div class="post-date">${comment.date}</div>
              </li>`})
  }
  InsertComment()
  Likes_Posts()
}


export function GetSinglePost() {
  let showAllComments = document.querySelectorAll(`.post-item`)
  showAllComments.forEach(e => e.addEventListener("click", async (e) => {
    let id = await e.target.getAttribute("data-post-id");
    if (id) {
      let post = await fetch(`/api/post/${id}`).then(response => response.json())
      if (post) {
        document.querySelector("main > ul").innerHTML = `<li class="post-item" data-post-id="${post.id}">
                      <div class="username">${post.author}</div>
                      <h3>${post.title}</h3>
                      <div class="category">Category: ${post.categories?.join(' - ') || "No Gategory"}</div>
                      <p class="content-preview">${post.content}</p>
                      
                      <div class="post-date">${post.date}</div>
  
                      <!-- <div class="interaction-section"> -->
                      <div class="interaction-section">
                        <button class="like-post-btn ${post.isliked ? "like-reacted" : ""}" name="like_post" value="${post.id}" id="likes"
                              onclick="">
                              <i class="fas fa-thumbs-up"></i>
                              ${post.likes}
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
                          
                  </li>
                   <div class="comment_container"></div>`
        NofetchComment = true
      }
      GetAllComment(id)


      let debounceTimer
      window.addEventListener("scroll", function () {
        if (window.scrollY + window.innerHeight >= document.body.scrollHeight - 100) {

          clearTimeout(debounceTimer);
          debounceTimer = setTimeout(() => {
            console.log('scrolling')
            GetAllComment(id, ++page_comments);
          }, 1000);

        }
      });
    }
  }))
}



export function InsertComment() {
  let CommentBtn = document.querySelectorAll(`button[type="submit"][name="id-post"]`)
  CommentBtn.forEach(e => e.addEventListener("click", async (e) => {

    const id = e.target.closest('.post-item').getAttribute('data-post-id');
    const postsinput = document.querySelector(`.post-item[data-post-id="${id}"] input[name="comment"]`);
    const data = { postId: parseInt(id), content: String(postsinput.value) }
    if (data.content === "") {
      return;
    }

    await fetch('/comment', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(response => {
      if (response.ok) {
        console.log('Comment added successfully');

        const commentList = document.querySelector(".comment_container")
        commentList.innerHTML = `<li class="comment-item">
                  <div class="username">${localStorage.getItem("username")}</div>
                  <p class="content-preview">${data.content}</p>
                  <div class="post-date">Now</div>
              </li>` + commentList.innerHTML
      }
    }).catch(e => {
      console.log(e)
      return
    })
    postsinput.value = '';
  }))
}


export async function GetAllPosts() {
  if (nomoreposts) {
    NofetchComment = false
    return;
  }

  console.log("GetAllPosts");

  await fetch(`/api/post?page=${page_posts}`).then((response) => response.json()).then(async (e) => {
    if (e) {
      console.log(e);

      let ul = document.querySelector('ul')
      await e.forEach((data) => {
        ul.innerHTML += `  <li class="post-item" data-post-id="${data.id}">
         
                    <div class="username">${data.author}</div>
                    <h3 class="post-title">${data.title}</h3>
                    <div class="category">Category: ${data.categories?.join(' - ') || "No Gategory"}</div>
                    <h4 class="content-preview">${data.content}</h4>
                    
                    <div class="post-date">${data.date}</div>
  
                    <!-- <div class="interaction-section"> -->
                    <div class="interaction-section">
                        <button class="like-post-btn ${data.isliked ? "like-reacted" : ""}" name="like_post" value="${data.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-up"></i>
                            ${data.likes}
                        </button>
                        <button class="dislike-post-btn ${data.isdisliked ? "dislike-reacted" : ""}" name="deslike_post" value="${data.id}" id="likes"
                            onclick="">
                            <i class="fas fa-thumbs-down"></i>
                            ${data.dislikes}
                    </div>
  
                    <div class="comment_container"></div>
                </li>`
      })
    }
  }).catch(e => {
    console.log(e)
    nomoreposts = true
    return
  })
  Likes_Posts()
  document.querySelector('h3.logo').addEventListener('click', GoToHomePage)
  GetSinglePost()
  InsertComment()


  const width = window.innerWidth;
  if (width < 768) {
    console.log('width', width);

    Resize();
  }
  window.addEventListener('resize', Resize);
  let debounceTimer
  let done = false
  window.addEventListener("scroll", function () {
    if (window.scrollY + window.innerHeight >= document.body.scrollHeight - 100) {
      if (!done && !NofetchComment && !nomoreposts) {
        done = true
        page_posts++
        console.log("sad", page_posts)
        GetAllPosts()

        clearTimeout(debounceTimer);


      }
      debounceTimer = setTimeout(() => {
        done = false
      }, 1000);

    }
  });

}

export function ShowCreatePost() {
  document.querySelector('.button-wrapper').addEventListener('click', () => {
    document.body.innerHTML += ` <div class="content-spacer-create"></div>
     
  
     <div class="create-post-container">
     <button class="close_button_post">X</button>
         <h2>Create Post</h2>
     
         <div>
             <input type="text" id="title" name="title" maxlength="50"  placeholder="Title:"required>
         </div>
         <div>
             <textarea id="content" name="content" maxlength="500"  placeholder="Content:" required></textarea>
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
             <p class="errorMessage"></p>
     </div>`

    document.body.style.overflow = "hidden";
    document.querySelector(".close_button_post").addEventListener("click", () => {
      GoToHomePage()
      document.body.style.overflow = "auto";
    })
    document.querySelector("#create-post-button").addEventListener("click", () => {
      CreatePost()
    })
  })
}


export async function GoToHomePage() {

  ress.is_resize = false;
  ress.done_resize = false;
  document.body.innerHTML = ""



  let header = document.createElement('header');
  header.classList.add('header');
  header.innerHTML = `
        <div class="header-content">
            <h3 class="logo">ALGO GANG</h3>
            <div class="logout-container">
                <button class="logout-button" >
                    <i class="fas fa-sign-out-alt"></i> Logout
                </button>
                 <button type="submit" class="button-wrapper">
                        <i class="fas fa-plus-circle"></i> Create Post
                    </button>
            </div>
            </div>
    `

  document.body.appendChild(header)

  document.body.innerHTML += `
     <aside class="sidebar-left">
             <h2>Contact</h2>
             <div class="listaside"></div>
          </aside>
        <main class="posts-container">
        <div class="NOTIFICATION">
        </div>
            <h1>Posts</h1>
            <ul>
            </ul>
             <div class="chat-container">
    <button class="X">X</button>
    <h3 id="TO" value="" ></h3>
    
  <div class="chat-box" id="chatBox">
  
  </div>
  
  <div class="input-area">
    <input type="text " id="messageInput" class="message-input" maxlength="500" placeholder="Type your message...">
  </div>
  <span class="typing-indicator"></span>
</div>
        </main>
    `;
  document.querySelector("main").innerHTML += `
      
   
    `

  ShowCreatePost()
  await FetchChatWithConversations()
  if (document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']")) {
    document.querySelector("link[rel='stylesheet'][href='/Assets/login.css']").href = "/Assets/post.css"
  }

  await GetAllPosts()
  if (!check_chat.isstarted) {

    console.log("testtest")
    await StartWs()
    check_chat.isstarted = true
  }
  document.querySelector(".logout-button").addEventListener("click", deleteCookie)
  console.log("testtest")
}

