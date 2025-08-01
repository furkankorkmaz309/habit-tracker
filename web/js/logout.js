document.getElementById("logout").addEventListener("click", async function(e) {
    e.preventDefault();

    try {
        const response = await fetch("http://localhost:8080/logout", {
            method: "POST",
            credentials: "include"
        });

        const result = await response.json();

        if (!result.success) {
            alert("Logout failed: " + result.message);
            return;
        }

        localStorage.clear();
        window.location.href = "landing.html";

    } catch (err) {
        console.error("Logout error:", err);
    }
});
