# 🛍️ Coolshop API

[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Framework-green)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue)](https://www.postgresql.org/)
[![Swagger](https://img.shields.io/badge/API-Docs-yellow)](#📁-api-documentation)

Coolshop API is a **RESTful API** built with **Golang, Gin, PostgreSQL, Redis, and JWT authentication**.  
It provides **user authentication, session management, and user data operations** for an e-commerce system.

---

## 🚀 Features
✅ **User Authentication** - Sign up, login, logout, password update  
✅ **JWT-based Authorization** - Secure API using **access & refresh tokens**  
✅ **Role-based Access Control** - Protected API endpoints  
✅ **Swagger API Documentation** - Interactive API docs  
✅ **Structured Logging** - Powered by **Uber's `zap`**  
✅ **Error Handling Middleware** - Centralized error handling  
✅ **Dockerized Setup** - Easy deployment

---


## 🛠️ Installation Guide

### **1️⃣ Install Go**

If you don't have Go installed, download and install it:

- Go: [Install Go](https://go.dev/doc/install/)

Verify installation:

```sh
  go version
```

### **2️⃣ Install Swag for Swagger Documentation**

Install `swag` globally:

```sh
  go install github.com/swaggo/swag/cmd/swag@latest
```

Verify installation:

```sh
  swag --version
```
### **3️⃣ Installation of Docker**
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)


## 📁 API Documentation
Coolshop API provides **interactive API documentation** via **Swagger**.

📌 **Access Swagger Docs** after running the server:  
📞 **http://localhost:8080/swagger/index.html**

To regenerate Swagger docs:
```sh
  swag init -g cmd/main.go
```

---

## 🛠️ Installation Guide

### **1️⃣ Clone the Repository**
```sh
    git clone https://github.com/kemalkochekov/CoolShop.git
    cd CoolShop
```

### **2️⃣ Install Dependencies**
```sh
  go mod tidy
```

### 3️⃣ Run Docker
```sh
  docker-compose up
```

---

## 📂 Project Structure

```
📦 Coolshop
 ┣ 📂 cmd                     # Main entry point
 ┃ ┗ 📜 main.go               # Starts the API server
 ┣ 📂 internal                
 ┃ ┣ 📂 app                   # App initialization & routes
 ┃ ┣ 📂 auth                  # JWT authentication & middleware
 ┃ ┣ 📂 cache                 # Redis caching
 ┃ ┣ 📂 config                # App configuration
 ┃ ┣ 📂 connection            # Database & Redis connections
 ┃ ┣ 📂 migrations            # Database migrations
 ┃ ┣ 📂 model                 # Data models
 ┃ ┣ 📂 user                  # User business logic (handlers, repository, usecase)
 ┣ 📂 logger                  # Logging with Uber Zap
 ┣ 📂 middleware              # Custom middlewares (logging, error handling)
 ┣ 📂 pkg                     # Utility packages (errors, constants, request validation)
 ┣ 📂 docs                    # Swagger API documentation
 ┣ 📂 tmp                     # Temporary files
 ┣ 📜 docker-compose.yml      # Docker configuration
 ┣ 📜 Makefile                # Task automation
 ┣ 📜 go.mod                  # Go module dependencies
 ┗ 📜 README.md               # Project documentation


---

## 🏠 Tech Stack
- **Language:** [Golang](https://go.dev/)
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Caching:** [Redis](https://redis.io/)
- **Authentication:** JWT (JSON Web Token)
- **Logging:** [Uber Zap](https://github.com/uber-go/zap)
- **API Docs:** [Swagger](https://swagger.io/)

---

## 🔒 Authentication & Authorization
### **🔑 Access Token (JWT)**
- Used for **authenticated requests**
- Sent in the `Authorization` header:
  ```http
  Authorization: YOUR_ACCESS_TOKEN
  ```

### **🔄 Refresh Token**
- Stored in **cookies**
- Used to **refresh the access token** when it expires.

---

## 🏰 Docker Setup
Run the project using **Docker**:
```sh
  docker-compose up --build
```
Stops and removes all containers:
```sh
  docker-compose down
```

---

## 🌐 API Endpoints

### **🛠️ Authentication**
| Method | Endpoint         | Description              | Auth Required |
|--------|-----------------|--------------------------|--------------|
| `POST` | `/auth/sign_up`  | User registration       | ❌ No |
| `POST` | `/auth/login`    | User login (returns JWT) | ❌ No |
| `POST` | `/auth/refresh`  | Refresh access token    | ✅ Yes |

### **👤 User Management**
| Method   | Endpoint             | Description           | Auth Required |
|----------|----------------------|----------------------|--------------|
| `GET`    | `/user/{id}`          | Get user by ID      | ✅ Yes |
| `PUT`    | `/user/update_password` | Update password | ✅ Yes |
| `DELETE` | `/user/`              | Delete user account | ✅ Yes |
| `POST`   | `/user/logout`        | Logout user         | ✅ Yes |

### **🚀 System Health**
| Method | Endpoint    | Description          |
|--------|------------|----------------------|
| `GET`  | `/health`  | API health check    |

---

## 👨‍💻 Contributing
🚀 We welcome contributions!  
To contribute:
1. **Fork the repository**.
2. **Create a new branch**:
   ```sh
   git checkout -b feature-new
   ```
3. **Make your changes & commit**:
   ```sh
   git commit -m "Add new feature"
   ```
4. **Push the branch**:
   ```sh
   git push origin feature-new
   ```
5. **Submit a Pull Request** 🚀

---


## ⭐ Support
If you find this project useful, **give it a star ⭐ on GitHub!**
