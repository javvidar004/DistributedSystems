# RBAC Distributed System

## Components(Container):
 - Frontend
 - Middleware
 - Loadbalancers
 - Auth service
 - Users service
 - Logs service

## Arquitecture (Docker compose)

                         - authlb  - authservice(x3)
 - Frontend - Middleware - userslb - usersservice(x3) - db
                         -  logslb - logsservice(x3)


Need to separate db and add requests between services
