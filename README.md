## What it is?

It is a space repo which is basically a middleware written on top of external oms's


## Prerequisite

One need to install the latest version of Golang i.e 1.18  from https://golang.org/dl/ and setup GOROOT and GOPATH.

### 1. ApiHelpers
Basically contains the helper functions used in returning api responses, HTTP status codes, default messages etc.

### 2. Controllers
Contains handler functions for particular route to be called when an api is called.

### 3. Helpers
Contains helper functions used in all apis

### 4. Middlewares
Middleware to be used for the project

### 5. Models
Database tables to be used as models struct

### 6. Resources
Resources contains all structures other than models which can be used as responses

### 7. Routers
Resources define the routes for your project

### 8. Seeder
It is optional, but if you want to insert lots of dummy records in your database, then you can use seeder.

### 9. Services
All the core apis for your projects should be within services.

### 10. Storage
It is generally for storage purpose.

### 11. Templates
Contains the HTML templates used in your project

### 12. .env
Contains environment variables.


## Steps to run Follow

- swag init

-- for linux
- make build

-- for other system
- go build 

in order to run binary generated 
- ./space  

        
## API with versioning

# For using version 1 api
```127.0.0.1:8099/api/v1/user-list```

# For using version 2 api
```127.0.0.1:8099/api/v2/user-list``` 


