# 📦 microservices‑learn

> Một minh họa đơn giản về microservices architecture viết bằng **Go**, gồm hai dịch vụ độc lập:
>
> - **user-service**: quản lý thông tin người dùng  
> - **order-service**: quản lý đơn hàng và gọi REST API đến `user-service`

<!-- Badges: CI / Go‑version / License -->
[![CI](https://github.com/RibunLoc/microservices‑learn/actions/workflows/ci.yml/badge.svg)](https://github.com/RibunLoc/microservices‑learn/actions/workflows/ci.yml)
[![Go version](https://img.shields.io/github/go‑mod/go‑version/RibunLoc/microservices‑learn?style=flat&logo=go)](https://pkg.go.dev/mod/github.com/RibunLoc/microservices-learn)
[![License: MIT](https://img.shields.io/github/license/RibunLoc/microservices‑learn?style=flat)](/LICENSE)

## ⚙️ Mục lục  
‑ [Tổng quan](#tổng‑quan)  
‑ [Kiến trúc](#kiến‑trúc)  
‑ [Chạy local](#chạy‑local)  
‑ [Docker Compose](#docker‑compose)  
‑ [API Endpoints](#api‑endpoints)  
‑ [Ví dụ (curl)](#ví‑dụ-curl)  
‑ [Contributing & Code style](#contributing‑code‑style)  
‑ [License](#license)

---

## Tổng quan  
Repos này nhằm mục đích **mô hình hoá kiến trúc microservices cơ bản**, với các đặc điểm:

- Tách rõ domain (user vs order), mỗi service phát triển độc lập bằng Go  
- Giao tiếp **REST nội bộ**: `order-service → user-service` để đồng bộ thông tin  
- Hoàn toàn dễ mở rộng thêm: API Gateway, service registry (Consul, etcd), messaging (Kafka/RabbitMQ)  

