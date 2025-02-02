document.querySelector('#register-button').addEventListener('click', Register());

document.querySelector('#login_button').addEventListener('click', Login());
document.querySelector('#signup_switch_button').addEventListener('click', function() {
    document.querySelector('.login-container').style.display = 'none';
    document.querySelector('.register-container').style.display = 'flex';
});
document.querySelector('#login_switch_button').addEventListener('click', function() {
    document.querySelector('.register-container').style.display = 'none';
    document.querySelector('.login-container').style.display = 'flex';
});


async function deleteCookie() {
    await fetch("/logout", {
        method: "POST",
        body: JSON.stringify({ Session: document.cookie.split('=')[1] })
    });
    document.cookie = "session=;expires=Tue, 22 Aug 2001 12:00:00 UTC;";
    window.location.href='/';
}

async function Login() {
    let email = document.querySelector("input#email")
    let password = document.querySelector("input#password")
    let data = { email: email.value, password: password.value };
    try {
        let response = await fetch("/signin", {
            method: "POST",
            body: JSON.stringify(data)
        })
        if (!response.ok) {
            const errorData = await response.json();
            const errorMessage =  document.getElementById('errorMessage');
            errorMessage.innerHTML =  errorData;    
                            
        } else {
            window.location.href = "/"
        }
    } catch (error) {
        const errorMessage = document.getElementById('errorMessage');
        errorMessage.innerHTML = 'Network error occurred!';
    }
}

async function Register() {
    let user = document.querySelector("input#username")
    let email = document.querySelector("input#email")
    let password = document.querySelector("input#password")
    let data = { email: email.value, password: password.value, user: user.value };
    try {
        let response = await fetch("/signup", {
            method: "POST",
            body: JSON.stringify(data)
        })
        if (!response.ok) {
            const errorData = await response.json();
            const errorMessage = document.getElementById('errorMessage');
            errorMessage.innerHTML =  errorData;                  
        } else {
            window.location.href = "/login" 
        }
    } catch (error) {
        const errorMessage = document.getElementById('errorMessage');
        errorMessage.innerHTML = 'Network error occurred!';
    }
}