   Returns the raw Markdown content for the GoTel project as a string.
    This function avoids rendering the Markdown in the chat interface.
    """
    # We use a raw string (r"""...""") to handle backslashes and special characters correctly.
    return r"""# 🏨 GoTel API (Hotel Reservation System)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v3-000000?style=flat&logo=gofiber)](https://gofiber.io/)
[![Database](https://img.shields.io/badge/MongoDB-4.4+-47A248?style=flat&logo=mongodb)](https://www.mongodb.com/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> **Note:** This is an **Educational Project** developed to demonstrate deep understanding of Go, RESTful API architecture, and NoSQL database integration.

**GoTel** is a high-performance, full-featured backend system for hotel reservation management. It provides comprehensive functionality for managing users, hotels, rooms, and bookings, featuring a secure JWT-based authentication system and an admin panel.

---

## 🚀 Features

* **Modular Architecture:** Clean separation of concerns (Handlers, Stores, Models) for maintainability.
* **Secure Authentication:** Robust **JWT** implementation for secure session management.
* **Fiber v3 Framework:** Built with the fastest web framework in the Go ecosystem.
* **MongoDB Integration:** Flexible NoSQL data storage using the official MongoDB driver.
* **Dockerized:** Fully containerized setup for easy deployment via Docker & Docker Compose.
* **Swagger UI:** Interactive and complete API documentation.
* **Security Best Practices:** Includes CORS management, rigorous input validation, and Role-Based Access Control (Admin/User).
* **Task Automation:** Integrated `Taskfile` for streamlined build and run commands.

---

## 🛠 Tech Stack

* **Language:** [Go (Golang)](https://go.dev/)
* **Web Framework:** [Fiber v3](https://gofiber.io/)
* **Database:** [MongoDB](https://www.mongodb.com/)
* **Cache:** [Redis](https://redis.io/) (Planned/Implemented)
* **Documentation:** [Swagger (Swaggo)](https://github.com/swaggo/swag)
* **Containerization:** Docker & Docker Compose
* **Task Runner:** [Task](https://taskfile.dev/)

---

## 🏁 Getting Started

Follow these steps to set up the project locally.

### Prerequisites

* [Docker](https://www.docker.com/) & Docker Compose
* [Go](https://go.dev/) (version 1.22 or higher)
* [Task](https://taskfile.dev/) (Optional, recommended for running commands)

### 1. Clone the Repository

```bash
git clone https://github.com/raminfathi/GoTel.git
cd GoTel
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory (or copy `.env.example`) and set the following values:

```env
HTTP_LISTEN_ADDRESS=:5000
MONGO_DB_NAME=gotel
MONGO_DB_URL=mongodb://localhost:27017
JWT_SECRET=your_secret_string_here
```

### 3. Run with Docker (Recommended)

Use the configured Taskfile to spin up the application and database containers:

```bash
task docker
# Or without Task: docker-compose up --build
```

### 4. Database Seeding

To populate the database with the initial admin user and sample hotel data, run:

```bash
task seed
```

> **🔑 Default Admin Credentials:**
>
> * **Email:** `admin@admin.com`
> * **Password:** `admin_admin`

---

## 📖 API Documentation (Swagger)

Once the application is running, you can explore and test the API endpoints via Swagger UI:

👉 **[http://localhost:5000/swagger/index.html](http://localhost:5000/swagger/index.html)**

---

## 📂 Project Structure

```text
GoTel/
├── api/             # HTTP Handlers (Controllers) and Middleware
├── cmd/             # Application Entry Points (Main, Seed)
├── db/              # Database Access Layer (Repositories)
├── types/           # Data Models and Structs
├── docs/            # Swagger Generated Files
├── Dockerfile       # Docker Image Configuration
├── docker-compose.yml
└── Taskfile.yml     # Automation Scripts
```

---

## 🧪 Useful Commands

| Command | Description |
| :--- | :--- |
| `task run` | Run the project locally (without Docker) |
| `task build` | Build the binary executable |
| `task test` | Run unit tests |
| `task docker` | Build and run Docker containers |
| `task seed` | Populate database with seed data |

---



