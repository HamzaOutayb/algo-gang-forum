import { GoToHomePage } from "./home.js";


export function Login_page(){
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

export function GoToLoginPage() {
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
 
    
    
    
    
export    async function Login (Login_re,key_re) {
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
    
export    async function Register () {
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
    
    