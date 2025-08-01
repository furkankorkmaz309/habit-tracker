const form = document.getElementById("form");
const usernameForm = document.getElementById("username");
const passwordForm = document.getElementById("password");

let username;
let password;

form.addEventListener('submit', function(e){
    e.preventDefault();

    username = usernameForm.value;
    password = passwordForm.value;

    let errors = [];

    if (username.replaceAll(' ', '') === "") {
        errors.push("Username cannot be blank");
    }
    if (password.replaceAll(' ', '') === "") {
        errors.push("Password cannot be blank");
    }

    if (errors.length > 0){
        alert(errors.join("\n"))
        return;
    }

    (async () =>{
        let response = await fetch("http://localhost:8080/login",{
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                username,
                password
            })
        })

        let result = await response.json();

        if (!result.success){
            alert(`Login failed : ${result.message}`)
            return;
        }

        console.log("login succeed");
        
        localStorage.setItem("username", username);
        window.location.href = "dashboard.html";
    })();

});