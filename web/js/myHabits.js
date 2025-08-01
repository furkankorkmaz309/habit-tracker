let username = localStorage.getItem("username");
let id = localStorage.getItem("id");

document.getElementById("username-header").innerText = username;
console.log();

(async () => {
    let response = await fetch("http://localhost:8080/habits", {
        method:"GET",
        credentials: "include"
    });

    let result = await response.json();

    if (!result.success){
        alert(`Request failed : ${result.message}`);
        return;
    }

    let habitsData = result.data


    for (let i = 0 ; i < habitsData.length ; i++) {
        let habits = document.querySelector(".habits");

        let habitDiv = document.createElement("div")
        habitDiv.className = "habit"

        let lSide = document.createElement("div")
        lSide.className = "l-side"
        
        let titleDiv = document.createElement("h1")
        titleDiv.innerText = habitsData[i].title;

        
        let descDiv = document.createElement("h4")
        descDiv.innerText = habitsData[i].description;
        

        let rSide = document.createElement("div")
        rSide.className = "r-side"
        
        let dateDiv = document.createElement("span")

        const formattedDate = FormatDate(habitsData[i].created_at);
        dateDiv.innerHTML = `${formattedDate} <span style="margin-left: 37px"> (${habitsData[i].day}) <span style="margin-left: 37px"> (${habitsData[i].frequency})`;
        
        let historyDiv = document.createElement("div")
        historyDiv.className = "history"
        historyDiv.innerText = "History:"

        let pDiv = document.createElement("p")
        pDiv.className = "done-bar"
        pDiv.innerText = "❓";


        let buttonsDiv = document.createElement("div")
        buttonsDiv.className = "buttons"
        
        let aDivEdit = document.createElement("a")
        aDivEdit.className = "edit-button"
        aDivEdit.innerHTML= `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#1f1f1f"><path d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"/></svg>`
        
        habitDiv.dataset.id = habitsData[i].id;
        aDivEdit.addEventListener('click', function(e){
            localStorage.setItem("habit-id", habitsData[i].id);
            window.location.href = "edit-habit.html";
        })

        let aDivDelete = document.createElement("a")
        aDivDelete.className = "delete-button"
        aDivDelete.innerHTML= `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#1f1f1f"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`
        aDivDelete.addEventListener('click', function(e){
            e.preventDefault();

            (async () =>{
                let response = await fetch(`http://localhost:8080/habits/${habitsData[i].id}`, {
                method: "DELETE",
                credentials: "include", 
                headers: {"Content-Type": "application/json"},
            })
                let result = await response.json();

                if (!response.ok){
                    alert(`Deleting habit is failed : ${result.message}`)
                    return;
                }

                console.log("Habit deleted successfully");
            })();
        });

        rSide.appendChild(dateDiv)
        historyDiv.appendChild(pDiv)
        rSide.appendChild(historyDiv)

        buttonsDiv.appendChild(aDivEdit)
        buttonsDiv.appendChild(aDivDelete)
        rSide.appendChild(buttonsDiv)

        habitDiv.appendChild(rSide)


        lSide.appendChild(titleDiv)
        lSide.appendChild(descDiv)
        habitDiv.appendChild(lSide)

        habits.appendChild(habitDiv)
    }


})();


function FormatDate(date) {
    const rawDate = new Date(date);
    const day = String(rawDate.getDate()).padStart(2, '0');
    const month = String(rawDate.getMonth() + 1).padStart(2, '0');
    const year = rawDate.getFullYear();
    const hours = String(rawDate.getHours()).padStart(2, '0');
    const minutes = String(rawDate.getMinutes()).padStart(2, '0');

    const formattedDate = `${day}-${month}-${year} ${hours}:${minutes}`;

    return formattedDate
}