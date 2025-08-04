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

    let id = localStorage.getItem("habit-id")
    title = titleForm.value;
    description = descriptionForm.value;
    freq = freqForm.querySelector('input[name="freq"]:checked')?.value || "";
    day = dayForm.querySelector('input[name="day"]:checked')?.value || "";



    if (title.replaceAll(' ', '') === "") {
        title = ""
    }
    if (description.replaceAll(' ', '') === "") {
        description = ""
    }

    if (freq == "D") {
        day = "";
    }
    if (day === "" && freq != "D") {
        alert("Day cannot be blank");
        return
    }


    (async () =>{
        let response = await fetch(`http://localhost:8080/habits/${id}`, {
        method: "PATCH",
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
            alert(`Editing habit is failed : ${result.message}`)
            return;
        }

        localStorage.removeItem("habit-id");
        console.log("Habit edited successfully");
        window.location.href = "my-habits.html";
    })();


});