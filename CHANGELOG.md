# Changelog

All notable changes to this project will be documented in this file.  


---

## [1.0.0] - 2025-08-15
### Added
- **Movie Service**  
  - Added RabbitMQ to consume seat lock
  - Dockerfile for containerized deployment.

- **Booking Service**  
  - Lock and unlock seat endpoints.
  - Seat Lock Handler
  - Basic seat availability check
  - Added RabbitMQ to stream seat lock details to movie service

- **Shared Infrastructure**  
  - `docker-compose.yml` for local multi-service setup.
  - `.env.example` with configuration placeholders.
  - `.gitignore` and `.dockerignore` for clean repository management.
  - `README.md` with project overview and setup instructions.

### Changed
- Improved folder structure to follow Go best practices:

### Fixed
- RabbitMQ code for movie-service
- Fixed seat lock check

---

[1.0.0]: https://github.com/levii0203/movie-reservation
