document.querySelector('#register-button').addEventListener('click', () => {Register()})

document.querySelector('#login_button').addEventListener('click',  () => {Login()})

document
  .querySelector('#signup_switch_button')
  .addEventListener('click', function () {
    document.querySelector('.login-container').style.display = 'none'
    document.querySelector('.register-container').style.display = 'flex'
  })

document
  .querySelector('#login_switch_button')
  .addEventListener('click', function () {
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
    let response = await fetch('/Signin', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage')
      errorMessage.innerHTML = errorData
    } else {
      Home()
    }
  } catch (error) {
    const errorMessage = document.getElementById('errorMessage')
    errorMessage.innerHTML = 'Network error occurred!'
  }
}

async function Register () {
  console.log(document.querySelector('input#nickname'))
  let nickname = document.querySelector('input#nickname').value
  let age = document.querySelector('input#age').value
  let gender = document.querySelector('input#gender').value
  let first_Name = document.querySelector('input#first_Name').value
  let last_Name = document.querySelector('input#last_Name').value
  let email = document.querySelector('input#email').value
  let password = document.querySelector('input#password').value
  let data = { 
    nickname: nickname,
    age: age,
    email: gender,
    gender: first_Name,
    first_Name: last_Name, 
    last_Name: email,
    password: password
  }
  try {
    let response = await fetch('/Signup', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    if (!response.ok) {
      const errorData = await response.json()
      const errorMessage = document.getElementById('errorMessage')
      errorMessage.innerHTML = errorData
    } else {
      document.querySelector('.register-container').style.display = 'none'
      document.querySelector('.login-container').style.display = 'flex'
    }
  } catch (error) {
    const errorMessage = document.getElementById('errorMessage')
    errorMessage.innerHTML = 'Network error occurred!'
  }
}


function Home () {
  document.querySelector('.login-container').remove()
  document.querySelector('.register-container').remove()
  document.querySelector('main')
  const body = document.querySelector('body')
  let header = document.createElement('header')
  let aside = document.createElement('aside')
  let div = document.createElement('div')
  body.appendChild(header)
  body.appendChild(aside)
}