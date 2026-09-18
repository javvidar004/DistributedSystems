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



