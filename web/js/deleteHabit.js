form.addEventListener('submit', function(e){
    e.preventDefault();
    
    (async () =>{
        let response = await fetch(`http://localhost:8080/habits/${id}`, {
        method: "DELETE",
        credentials: "include", 
        headers: {"Content-Type": "application/json"},
    })
        let result = await response.json();

        if (!response.ok){
            alert(`Deleting habit is failed : ${result.message}`)
            return;
        }

        localStorage.removeItem("habit-id");
        console.log("Habit deleted successfully");
        window.location.href = "my-habits.html";
    })();
});