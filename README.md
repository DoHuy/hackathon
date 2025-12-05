# Hackathon API Service

Dự án Backend API được xây dựng bằng **Go (Golang)** và **Echo Framework**.
Hệ thống cung cấp tính năng xác thực người dùng (JWT) và tải lên tệp tin.

## 🚀 Tính năng

- **Authentication**: Đăng ký, Đăng nhập sử dụng JWT (HS256).
- **File Upload**:
  - Upload ảnh (JPG, PNG, GIF).
  - Validate Magic Bytes (chống fake đuôi file).
  - Giới hạn dung lượng (Configurable).
- **Database**: PostgreSQL lưu trữ User và File Metadata.
- **Architecture**: Modular (Handler -> Service -> Repository).
- **Observability**: Structured Logging với Zerolog.
- **Documentation**: Swagger UI tự động.

## 🛠 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Labstack Echo v4
- **Database**: PostgreSQL
- **ORM**: GORM
- **Logging**: Zerolog
- **Docs**: Swaggo

## 📂 Cấu trúc dự án

```text
hackathon/
├── config/         # Đọc cấu hình từ env.ini
├── database/       # Khởi tạo kết nối GORM
├── docs/           # Swagger generated files
├── dto/            # Data Transfer Objects (Request/Response struct)
├── handlers/       # HTTP Handlers (Controller layer)
├── middleware/     # Custom middlewares (JWT, Logger, SizeLimit)
├── models/         # Database Models
├── pkg/            # Packages dùng chung (Logger util)
├── repositories/   # Data Access Layer (Tương tác DB)
├── services/       # Business Logic Layer
├── main.go         # Entry point
├── env.ini         # Configuration file
├── Dockerfile      # Docker build instructions
└── docker-compose.yml
```

## Hướng dẫn cài đặt Local

1.  **Clone mã nguồn**

    ```bash
    git clone https://github.com/dohuy/hackathon.git
    cd hackathon
    ```

2.  **Cài đặt Go**

    Đảm bảo bạn đã cài đặt Go phiên bản 1.21 trở lên.

3.  **Cài đặt các thư viện**

    ```bash
    go mod tidy
    ```

4.  **Cấu hình file `env.ini`**

    Sao chép file `env.example.ini` thành `env.ini` và cập nhật các thông tin cho môi trường local của bạn.

    ```bash
    cp env.example.ini env.ini
    ```

    Đối với môi trường local, bạn có thể cần cập nhật phần `[database]` để sử dụng PostgreSQL trên máy của bạn.

5.  **Chạy ứng dụng**

    ```bash
    go run main.go
    ```

    API sẽ chạy ở địa chỉ `http://localhost:8080`. Tài liệu Swagger có ở `http://localhost:8080/swagger/index.html`.

## Hướng dẫn cài đặt với Docker

1.  **Cài đặt Docker**

    Đảm bảo Docker và Docker Compose đã được cài đặt trên máy của bạn.

2.  **Chạy bằng Docker Compose**

    Lệnh này sẽ build và khởi chạy các service `app` và `db` ở chế độ nền.

    ```bash
    docker-compose up -d --build
    ```

3.  **Kiểm tra các service**

    Bạn có thể kiểm tra các container đang chạy bằng lệnh sau:

    ```bash
    docker-compose ps
    ```

4.  **Dừng các service**

    Để dừng và xóa các container, network và volume, chạy lệnh:

    ```bash
    docker-compose down
    ```
