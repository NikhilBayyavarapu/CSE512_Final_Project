Please follow the below instructions to smoothly run the code in the repository.

## Environment Setup

### 1) Install Go: 
Visit URL: https://go.dev/doc/install

This Go Language tool installation is optional: https://pkg.go.dev/golang.org/x/tools/gopls#section-readme. But this will be a useful tool have when debugging and writing code.

### 2) Install Node and NPM:
Visit URL: https://nodejs.org/en/download/package-manager

### 3) Install Docker:
Install Docker by visting the URL: https://www.docker.com/get-started/

### 4) Install http-server:
Run ```npm install -g http-server``` to server http files locally using a server.

### 5) Optional Installations
- Go Language Sever and Analysis Tool: https://pkg.go.dev/golang.org/x/tools/gopls#section-readme. 
- MongoDB Compass: https://www.mongodb.com/try/download/compass


### 6) Dependency Installations:
Run ```npm install``` to install dependencies required by Javascript Code in the repository.

In case there is a problem running the Go code and it looks like a dependency issue, please run the following commands:

``` go mod vendor```

``` go mod tidy```