# RBAC Distributed System

## Components(Container):
 - Frontend
 - Middleware
 - Loadbalancers
 - Auth service
 - Users service
 - Logs service

## Arquitecture (Docker compose)
  
 - Frontend
 - Middleware
 - userslb - usersservice(x3)
 - logslb - logsservice(x3)
 - authlb  - authservice(x3)
 - db

| Front layer  | Middleware    | LB's    |  Backend    | DB  |
| ------------ |:-------------:| -----:  | ----------- | --- |
|              |               | userslb | usersbe (x3)|     |
| Frontend     | middleware    | authlb  | authbe (x3) | db  |
|              |               | logslb  | logsbe (x3) |     |

## How to execute the project
### Prerequisites
 - Git
 - Apt package manager

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

Serves the server as on port 80, but on the docker compose 8080. This server uses a multiplexer: ServeMux. This permits recibing the requests and redirect them to the correct loadbalancer.  

### Load balancer



