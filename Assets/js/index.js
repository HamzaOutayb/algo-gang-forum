import { ShowCreatePost, Likes_Posts, GetSinglePost, InsertComment, ress, GoToHomePage } from "./modules/home.js";
import { Login_page } from "./modules/auth.js";
import { deleteCookie } from "./modules/auth.js";
import authUtils from "./utils/auth.js";

export let idtime




window.history.pushState({}, '', "/");


export function Listeners() {
  ShowCreatePost()
  Likes_Posts()
  GetSinglePost()
  InsertComment()
}


export function Resize() {
  const width = window.innerWidth;

  if (width < 768) {
    if (!ress.done_resize) {
      ress.done_resize = true;
      const buttonaside = document.createElement('button');
      buttonaside.classList.add('buttonaside');
      buttonaside.innerHTML = '<i class="fas fa-bars"></i>';
      document.body.appendChild(buttonaside);


      buttonaside.addEventListener('click', () => {
        if (ress.is_resize) {
          document.querySelector('.sidebar-left').style.display = 'none';
        } else {
          document.querySelector('.sidebar-left').style.display = 'block';
        }
        ress.is_resize = !ress.is_resize;
      });
    }
  } else {
    ress.done_resize = false;
    ress.is_resize = false;
    document.querySelector('.buttonaside')?.remove();
    const sidebar = document.querySelector('.sidebar-left');
    if (sidebar) {
      sidebar.style.display = 'block';
    }
  }
}

document.addEventListener('DOMContentLoaded', async () => {
 
  Login_page()
})