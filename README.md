# Habit Tracker - A Habit Tracking Application

This project is a full-stack web application that allows users to create and track their daily, weekly, or monthly habits. The backend is developed with Go (Golang), and the frontend is built with vanilla JavaScript, HTML, and CSS, without using any libraries or frameworks.

## ✨ Features

  * **User Management:** Secure signup, login, and logout functionalities.
  * **Habit Management (CRUD):**
      * Create new habits.
      * List all habits belonging to the user.
      * Edit existing habits (updating specific fields - PATCH).
      * Delete habits.
  * **Different Habit Frequencies:** Ability to define habits with different frequencies like daily, weekly, or monthly.

## 🛠️ Technical Architecture

The project consists of two modern, decoupled main layers:

### Backend (Go)

  * **Language:** Go (Golang)
  * **Web Server & Router:** Standard `net/http` library and `go-chi/chi` router.
  * **Database:** SQLite 3 with the `mattn/go-sqlite3` driver.
  * **Authentication:** JWT (JSON Web Tokens) and secure password hashing with `bcrypt`.

### Frontend (Vanilla JS)

  * **Structure:** Developed with pure HTML5, CSS3, and JavaScript (ES6+) without any frameworks.
  * **API Communication:** Asynchronous communication with the backend is established using the `Fetch API` and `async/await`.
  * **Styling:** Modern CSS layout techniques such as `Flexbox` and `CSS Grid` are used.

## 🛡️ Security Approach

Special attention has been given to security in this project.

1.  **Password Security:** User passwords are hashed with the **bcrypt** algorithm before being saved to the database. This makes it impossible to reverse the passwords.
2.  **Dual-Layer Session Security:**
      * When a user logs in, a standard **JWT** is created. This token contains the user's identity and permissions and is signed with a secret key (`SECRET_KEY`).
      * **As an additional security layer**, the signed JWT itself is fully **encrypted** on the server-side using the **AES-GCM** algorithm before being sent to the client.
      * This encrypted data is stored in an `HttpOnly` cookie. This ensures that even if the cookie is stolen, its content cannot be read without the encryption key (`ENCRYPTION_KEY`).
3.  **Authorization:** All API endpoints on the backend check at the database level whether the habit being acted upon belongs to the user making the request.

## 🚀 Setup and Running

You can follow the steps below to run the project on your local machine.

### Prerequisites

  * [Go](https://go.dev/doc/install) (version 1.20+ recommended)
  * A code editor (e.g., VS Code) and an extension like "Live Server" (for the frontend).

### Backend Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/furkankorkmaz309/habit-tracker.git
    cd habit-tracker
    ```

2.  **Set up Environment Variables:**
    Create a file named `.env` in the project's root directory and fill it with your own values based on the following content.

    ```env
    # .env.example

    # Path to the directory where the database file will be created
    DB_PATH=./internal/data/

    # Bcrypt hash cost factor
    COST_NUM=12

    # Secret key for signing the JWT (should be a complex string)
    SECRET_KEY="VERY_SECRET_KEY_12345"

    # Key for encrypting the JWT (MUST BE EXACTLY 32 BYTES!)
    ENCRYPTION_KEY="12345678901234567890123456789012" 
    ```

3.  **Run the Server:**

    ```bash
    go run ./cmd/habit-tracker/main.go
    ```

    The server will start running by default at `http://localhost:8080`.

### Frontend Setup

1.  If you are using VS Code, install the "Live Server" extension.
2.  Right-click on the `web/pages/landing.html` file in the project folder and select "Open with Live Server".
3.  The application will open in your browser at `http://localhost:5501` (or a similar port).

## 🗺️ API Endpoints

| Endpoint       | Method   | Protected | Description                                |
| -------------- | -------- | --------- | ------------------------------------------ |
| `/signup`      | `POST`   | No        | Creates a new user registration.           |
| `/login`       | `POST`   | No        | Logs in a user and returns a token.        |
| `/logout`      | `POST`   | No        | Ends the user's session.                   |
| `/habits`      | `POST`   | Yes       | Creates a new habit.                       |
| `/habits`      | `GET`    | Yes       | Lists all habits of the logged-in user.    |
| `/habits/{id}` | `GET`    | Yes       | Fetches the details of a specific habit.   |
| `/habits/{id}` | `PATCH`  | Yes       | Partially updates a specific habit.        |
| `/habits/{id}` | `DELETE` | Yes       | Deletes a specific habit.                  |

## 🎯 Future Enhancements

  * [ ] Replace `alert()` notifications in the frontend with user-friendly toast/notification components.
  * [ ] Add a "check-in" feature to track the daily completion status for each habit.
  * [ ] Create a statistical dashboard to display habit completion rates.

## Contact

Furkan Korkmaz - [GitHub](https://github.com/furkankorkmaz309) - [LinkedIn](https://www.linkedin.com/in/furkankorkmaz309)