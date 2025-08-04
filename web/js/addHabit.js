let username = localStorage.getItem("username");
let id = localStorage.getItem("id");

document.getElementById("username-header").innerText = username;
console.log();

const form = document.getElementById("new-habit");
const titleForm = document.getElementById("title");
const descriptionForm = document.getElementById("description");
const freqForm = document.getElementById("freqs")
const dayForm = document.getElementById("days")

let title;
let description;
let freq;

form.addEventListener('submit', function(e){
    e.preventDefault();

    let errors = [];

    title = titleForm.value;
    description = descriptionForm.value;
    freq = freqForm.querySelector('input[name="freq"]:checked')?.value || "";
    day = dayForm.querySelector('input[name="day"]:checked')?.value || "";

    if (title.replaceAll(' ', '') === "") {
        errors.push("Title cannot be blank");
    }
    if (description.replaceAll(' ', '') === "") {
        errors.push("Description cannot be blank");
    }
    if (freq === "") {
        errors.push("Frequency cannot be blank");
    }

    if (day === "" && freq != "D") {
        errors.push("Day cannot be blank");
    }
    
    if (day != "" && freq == "D") {
        day = "";
    }

    if (errors.length > 0){
        alert(errors.join("\n"))
        return;
    }

    (async () =>{
        let response = await fetch("http://localhost:8080/habits", {
        method: "POST",
        credentials: "include", 
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            title,
            description,
            frequency: freq,
            day 
        })
    })
        let result = await response.json();

        if (!response.ok){
            alert(`Adding habit is failed : ${result.message}`)
            return;
        }

        console.log("Habit added successfully");
        window.location.href = "my-habits.html";
    })();


});