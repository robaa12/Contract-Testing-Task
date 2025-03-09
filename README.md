# Microservice Contract Testing with Keploy

This project demonstrates contract testing between microservices using Keploy. It consists of three microservices:

1. User Service - Manages user data
2. Order Service - Handles order processing
3. Payment Service - Processes payments

## Architecture

- **User Service**: Provides APIs for user management
- **Order Service**: Handles order management and communicates with User Service
- **Payment Service**: Processes payments using Stripe

## Prerequisites

- Go 1.20+
- Docker and Docker Compose
- PostgreSQL
- Keploy CLI

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/robaa12/keploy-contract-testing.git
   cd keploy-contract-testing
   ```

2. Start the services with Docker Compose:
   ```
   docker-compose up -d
   ```

3. The services will be available at:
   - User Service: http://localhost:8080
   - Order Service: http://localhost:8081
   - Payment Service: http://localhost:8082

## Contract Testing

This project uses Keploy for contract testing between the microservices. The contract tests ensure that any changes to one service don't break the communication with dependent services.

To run contract tests:

```bash
keploy contract test 
```

## Services and Dependencies

- **User Service**: Independent service with PostgreSQL database
- **Order Service**: Depends on User Service
- **Payment Service**: Integration with Stripe API
