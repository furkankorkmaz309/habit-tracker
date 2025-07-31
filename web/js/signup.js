const form = document.getElementById("form");
const usernameForm = document.getElementById("username");
const emailForm = document.getElementById("email");
const passwordForm = document.getElementById("password");
const rePasswordForm = document.getElementById("re-password");

let username;
let email;
let password;
let rePassword; 

form.addEventListener('submit', function(e){
    e.preventDefault();

    username = usernameForm.value;
    email = emailForm.value;
    password = passwordForm.value;
    rePassword = rePasswordForm.value;
    
    let errors = [];

    if (username === "") {
        errors.push("Username cannot be blank");
    }
    if (email === "") {
        errors.push("Email cannot be blank");
    }
    if (password === "") {
        errors.push("Password cannot be blank");
    }
    if (rePassword === "") {
        errors.push("Re-password cannot be blank");
    }
    if (password !== rePassword) {
        errors.push("Passwords are not same");
    }

    if (errors.length > 0){
        alert(errors.join("\n"))
        return;
    }

    (async() => {
        let response = await fetch("http://localhost:8080/signup",{
            method: "POST",
            headers: {"Content-type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                username: username,
                email: email,
                password: password,
            }),
        })

        let result = await response.json();

        if (!response.ok){
            alert(`Signup failed : ${result.message}`)
            return;
        }

        console.log("Signup succeed");
        window.location.href = "login.html"
        }
    )();
});