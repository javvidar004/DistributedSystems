# RBAC Distributed System

## Index

- [Overview](#overview)
- [Architecture](#architecture)
- [Services](#services)
  - [Auth service](#auth-service)
  - [Users service and list endpoint](#users-service-and-list-endpoint)
  - [Logs service](#logs-service)
- [Registration flow: auth to users over gRPC](#registration-flow-auth-to-users-over-grpc)
- [Logging flow](#logging-flow)
- [Key implementation details](#key-implementation-details)
- [How to run](#how-to-run)
- [Docker composition](#docker-composition)

## Overview

This project implements a distributed RBAC system designed to act as an identity and authorization layer for other applications. The core idea is that external apps can send authentication requests to the auth service, receive a JWT, and validate it without needing to implement user management logic themselves.

The system is split into independent services:

- Frontend: administrative interface
- NGINX gateway: routes requests and balances traffic
- Auth service: login, registration, JWT issuance, password update
- Users service: user profile storage and listing
- Logs service: audit and action history
- Databases: separate persistence for auth, users, and logs

## Architecture

The platform runs as a microservice architecture behind Docker Compose and Consul-based service discovery.

- Auth instances are load-balanced behind the auth layer.
- Users instances are load-balanced behind the users layer.
- Logs instances are load-balanced behind the logs layer.
- NGINX acts as the API gateway and forwards HTTP and gRPC traffic to the correct backend.
- Kafka is used to decouple event logging from the app services that generate those events.

| Layer     | Components    |
| ---       | ---           |
| Frontend  | Users/admin UI |
| Middleware | NGINX + Consul service discovery |
| Backend   | authbe x3, usersbe x3, logsbe x3 |
| Database  | authdb, usersdb, logsdb |
| Messaging | Kafka         |

## Services

A general pattern used throughout the project is the separation between the transport layer and the business logic layer. Each service is organized so that:

- the HTTP handlers in the `controllers` package receive REST requests
- the gRPC servers in the `grpcserver` package receive RPC requests
- the real business operations live in the `service` package and are shared by both protocols

This design is intentional. It allows the same logic to be reused regardless of protocol. For example, the same `AuthService.RegisterAndReplicate` method can be called from an HTTP endpoint or from a gRPC RPC, while the handlers and gRPC server methods only adapt the input/output format to the protocol being used. In other words, protocol-specific code is kept thin, and the business rules remain centralized in one place.

This pattern makes the system easier to maintain because it avoids duplicating logic across REST and gRPC implementations. The handlers are responsible for binding JSON payloads, validating request format, and returning HTTP responses; the gRPC server does the same for protobuf messages; and both delegate to the `service` layer for the actual work.

### Auth service

The auth service is the entry point for authentication and registration. It manages:

- user creation in the Auth database
- password hashing with bcrypt
- JWT creation and validation
- registration flows that replicate the user to the users service
- event logging for login, register, update password, and delete user actions

The service exposes the following HTTP endpoints:

- POST /login
- POST /validate
- POST /register
- PUT /update
- DELETE /delete/:userID
- GET /heartbeat

The main service object is `AuthService`, initialized with the PostgreSQL connection and the gRPC Users client. This service handles the identity side of the system while the users service stores the wider profile data.

### Users service and list endpoint

The users service stores user profile records such as:

- name
- last name
- work position
- salary
- email

It exposes the following HTTP endpoints:

- POST /user
- PUT /user/:userID
- DELETE /user/:userID
- GET /user/:userID
- GET /users
- GET /heartbeat

The list endpoint is:

- GET /users

This is implemented by `GetUsersHandler`, which calls `usersSvc.GetUsers(ctx)`. The response is a list of user records in JSON, for example:

```json
{
  "users": [
    {
      "id": "...",
      "email": "user@example.com",
      "name": "John",
      "last_name": "Doe",
      "work_position": "Engineer",
      "salary": 50000
    }
  ]
}
```

The users service also supports:

- `POST /user` to create a user record
- `PUT /user/:userID` to update a profile
- `DELETE /user/:userID` to remove a record
- `GET /user/:userID` to fetch a single user by ID

At the gRPC layer, the service implements `CreateUser`, `UpdateUser`, `DeleteUser`, `GetUser`, `GetUsers`, and `Heartbeat` through the generated `UsersService` RPC interface defined in the proto file.

### Logs service

The logs service centralizes audit and access tracking. It stores records with:

- timestamp
- username
- action
- status

The main endpoints are:

- GET /log
- POST /log
- GET /heartbeat

It also exposes gRPC methods such as:

- `GetLogs`
- `NewLog`
- `Heartbeat`

## gRPC logic and proto files by service

### Auth service gRPC flow

The auth service defines its RPC contract in [auth/proto/auth.proto](auth/proto/auth.proto). It exposes:

- `Login`
- `Register`
- `UpdatePassword`
- `DeleteUser`
- `Heartbeat`
- `LogsService.NewLog`

The generated Go bindings are created under [auth/pb/auth.pb.go](auth/pb/auth.pb.go) and [auth/pb/auth_grpc.pb.go](auth/pb/auth_grpc.pb.go). Those files are not handwritten; they are generated from the proto definitions and provide the client/server interfaces used by Go code.

The logic behind the gRPC implementation is:

1. The service is created in `main.go` with `pb.RegisterAuthServiceServer(...)`.
2. The server implementation is in `auth/grpcserver/auth_server.go`.
3. Every RPC method receives a proto request, maps it to the internal service input, calls the business logic, and converts the result back into a proto response.

For example, `Register` does this:

```go
res, err := s.authSvc.RegisterAndReplicate(ctx, service.RegisterInput{
    Name:         req.GetName(),
    LastName:     req.GetLastName(),
    Email:        req.GetEmail(),
    WorkPosition: req.GetWorkPosition(),
    Salary:       req.GetSalary(),
    Password:     req.GetPassword(),
})
```

This allows other internal services or clients to register users through a strongly typed RPC rather than only HTTP endpoints.

### Users service gRPC flow

The users service proto file is [users/proto/users.proto](users/proto/users.proto). It defines the complete user-management contract, including:

- `CreateUser`
- `UpdateUser`
- `DeleteUser`
- `GetUser`
- `GetUsers`
- `Heartbeat`
- `LogsService.NewLog`

The generated code is under [users/pb/users.pb.go](users/pb/users.pb.go) and [users/pb/users_grpc.pb.go](users/pb/users_grpc.pb.go).

The actual implementation pattern is:

1. `users/main.go` registers the `UsersService` server with gRPC.
2. `users/grpcserver/users_server.go` implements the RPC methods.
3. Each method translates the proto request to `service.CreateUserInput`, `UpdateUserInput`, etc.
4. The business logic runs in [users/service/users_service.go](users/service/users_service.go).
5. The result is transformed back into a proto response object.

For example, `GetUsers` is implemented as:

```go
res, err := s.usersSvc.GetUsers(ctx)
users := make([]*pb.User, 0, len(res.Users))
for _, u := range res.Users {
    users = append(users, userToProto(u))
}
return &pb.GetUsersResponse{Users: users}, nil
```

This is why the list endpoint and the gRPC list method both share the same core business logic: the HTTP layer calls the service directly, while the gRPC layer translates the proto request and then calls the same service.

### Logs service gRPC flow

The logs service proto file is [logs/proto/logs.proto](logs/proto/logs.proto). It defines:

- `GetLogs`
- `NewLog`
- `Heartbeat`

The generated bindings are under [logs/pb/logs.pb.go](logs/pb/logs.pb.go) and [logs/pb/logs_grpc.pb.go](logs/pb/logs_grpc.pb.go).

The gRPC logic for logs is:

1. The service creates a gRPC server in `logs/main.go`.
2. `logs/grpcserver/logs_server.go` listens for `NewLog` and `GetLogs` RPCs.
3. Those methods call `LogsService.NewLog` and `LogsService.GetLogs` from [logs/service/logs_service.go](logs/service/logs_service.go).
4. The service persists log rows to the database and returns the stored object back as a proto `Log` response.

This is important because the logs service is also consuming Kafka messages asynchronously. The Kafka consumer writes to the same `LogsService.NewLog` method, so the process is consistent whether logs arrive from HTTP, gRPC, or Kafka.

## How the proto files are used in practice

The generated proto files are the contract between services. They define the exact message types and RPC names used for inter-service communication.

In this project the pattern is:

- the `.proto` files define the service schema
- `protoc` generates Go code in the `pb` folder
- application code calls those generated client methods
- server implementations receive the request and route it to the internal service logic

Examples in this codebase:

- Auth calls users through `usersClient.CreateUser(...)` using the generated type `pb.CreateUserRequest`.
- Logs are created by the `NewLog` RPC from the logs service contract.
- Users expose `GetUsers` and `GetUser` over both HTTP and gRPC, with the same business behavior in the service layer.

This separation keeps the services decoupled while still allowing strong typing and consistent request/response contracts across the distributed system.

## Registration flow: auth to users over gRPC

When a user registers through the auth service, the registration process is not confined to the auth database alone. The auth service intentionally replicates the user to the users service over gRPC.

The flow is:

1. The client sends a POST request to /register.
2. `RegisterHandler` binds the payload and calls `authSvc.RegisterAndReplicate(...)`.
3. `AuthService.RegisterAndReplicate` validates the payload, checks duplicate email, hashes the password, and creates the user in the auth database.
4. If the auth DB write succeeds, it invokes `replicateUserToUsersService(...)`.
5. That method sends a gRPC request to the users service via `usersClient.CreateUser(...)`.
6. The request contains the profile fields:
   - `Id`
   - `Email`
   - `Name`
   - `LastName`
   - `WorkPosition`
   - `Salary`
7. The users service receives the request in `UsersServer.CreateUser` and stores the user profile in the users database.
8. If the users service call fails, the auth service deletes the newly created auth user to avoid inconsistent state.

This is the key distributed pattern in the project: the auth service owns login credentials, while the users service owns the user profile data.

The actual gRPC call looks like this:

```go
_, err := s.usersClient.CreateUser(grpcCtx, &pb.CreateUserRequest{
    Id:           userID,
    Email:        input.Email,
    Name:         input.Name,
    LastName:     input.LastName,
    WorkPosition: input.WorkPosition,
    Salary:       input.Salary,
})
```

## Logging flow

The system logs authentication and user-management actions. In practice, the logging flow is a combination of Kafka event publication and a logs service that persists those events.

### Producer side

Each auth or users action calls a `newLogRequest(...)` helper. It serializes a payload such as:

```json
{
  "timestamp": "2026-10-29T12:00:00Z",
  "username": "user@example.com",
  "action": "register",
  "status": "success"
}
```

That message is then sent to Kafka on the `logs_topic` topic.

### Consumer side

The logs service starts a Kafka consumer (`consumeLogsFromKafka`) and listens for messages on `logs_topic`. When a message is received, it unmarshals the payload and calls:

```go
logsSvc.NewLog(context.Background(), service.NewLogInput{
    Timestamp: payload.Timestamp,
    Username:  payload.Username,
    Action:    payload.Action,
    Status:    payload.Status,
})
```

This writes the audit record to the logs database.

### gRPC exposure for logs

The logs service also exposes gRPC methods, especially `GetLogs` and `NewLog`, so other services can query or record logs directly. In other words, the project uses Kafka for asynchronous log publishing and gRPC for direct service interaction when needed.

## Key implementation details

- Passwords are hashed with bcrypt before being stored.
- JWTs are used for authentication and validation.
- Service discovery is handled by Consul.
- NGINX is the public gateway and load balancer.
- Each backend service keeps a separate database.
- gRPC is used for service-to-service communication between auth and users; Kafka decouples event logging.
- Failures are handled explicitly: the auth service rolls back the auth user if replication to the users service fails.

## How to run

### Prerequisites

- Git
- Docker and Docker Compose
- Access to a Unix-like environment or WSL for the startup script

```sh
git clone https://github.com/javvidar004/DistributedSystems.git
cd DistributedSystems/
sudo su
chmod 700 serverConfig.sh
./serverConfig.sh
```

## Structure Breakdown

### Middleware
Middleware uses a defined list of endpoints depending on the prefix of the path of the request. Is is defined:
- /auth/
- /users/
- /logs/

NOTE: It is important to have in mind that defining endpoints in Nginx the definition of the location ends in "/" and after the instruction of "proxy_pass http://backend_pool/" also has to end in an "/" this way we ensure that the request sent to the backend does not sends the prefix, this one just helps know to which type of backend it is send the request but the location of the backend is after each of the 3 prefix listed before.  

The stack includes:

- frontend
- gateway
- auth backend instances
- users backend instances
- logs backend instances
- auth database
- users database
- logs database
- Kafka
- Consul

The project is designed so each service can be scaled individually while remaining behind the same gateway and discovery layer.

## Summary

This system is a distributed identity and audit platform built around authentication, service isolation, gRPC communication, and Kafka-based log streaming. The auth service issues JWTs and validates credentials, the users service stores rich profile data, and the logs service stores operational audit information for monitoring and security review.


