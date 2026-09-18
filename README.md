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
 ```sh
git clone 
```


Need to separate db and add requests between services
