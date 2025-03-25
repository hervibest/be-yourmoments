# 🧩 Microservice Backend Architecture

This repository contains a **modular microservice backend architecture** designed to scale and support high concurrency. The system is built with a strong emphasis on **domain separation**, **IO vs CPU-intensive task isolation**, and **high availability**.

## 📦 Microservices Overview

### 1. Upload Service
- Handles image uploads and compression (**CPU-intensive**)
- Isolated for scalability and performance
- Prevents blocking of I/O-bound services like Photo Service

### 2. Photo Service
- Manages image metadata and access
- Designed as **I/O-focused** microservice
- Communicates with Upload Service via **gRPC**

### 3. Transaction Service
- Manages user transactions and related operations
- Independent domain to support future financial or logging features

### 4. User Service
- Handles authentication, user profiles, and user data
- Acts as a foundational service for all identity-based interactions

---

## ⚙️ Design Principles

- **Scalability**  
  CPU-heavy services (e.g., image compression in Upload Service) are separated to allow independent scaling.

- **Concurrency Optimization**  
  IO-bound (Photo Service) and CPU-bound (Upload Service) workloads are decoupled for efficient resource utilization.

- **Domain-Driven Design**  
  Clear separation of concerns across Upload, Photo, Transaction, and User domains for better maintainability and extensibility.

---

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: [Fiber](https://gofiber.io)
- **Communication**: gRPC
- **Database**: PostgreSQL
- **Cache/Session**: Redis
- **Service Discovery**: Consul

---

## 🤖 AI Integration (Coming Soon)

> Integration of AI-based **Face Recognition** and **Image Embedding** is currently in development.

This feature will enhance the system with intelligent photo processing and identity recognition, designed to integrate seamlessly with existing Upload and Photo services.

---

## 🚧 Project Status

This project is actively being developed and intended for future expansion with:
- AI capabilities
- Observability and monitoring tools
- CI/CD support
- Dockerized deployment

---

## 📂 Folder Structure (Optional Example)

