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
NOTE: It is important to have in mind that defining endpoints in Nginx the definition of the location ends in "/" and after the instruction of "proxy_pass http://backend_pool/" also has to end in an "/" this way we ensure that the request sent to the backend does not sends the prefix, this one just helps know to which type of backend it is send the request but the location of the backend is after each of the 3 prefix listed before.  

For the function of LB, it is defined and upstream pool for each of the types of backends. This way we ensure all 3 copies of the backend are used by the Middleware.  
  
IMPORTANT: The container of the middleware should be the last one to be up, this because it checks that all backends defined on the config file are up, if any of the backends is down the server will exit with code 1.






