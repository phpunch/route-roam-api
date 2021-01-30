## Route Roam API
API Backend for travel social media


### Features
- [x] Register
- [x] Login
- [x] Create post
- [ ] Edit post 
- [x] Delete post
- [x] Like post
- [x] Unlike post
- [x] Comment in the post
- [ ] Delete comment in the post

### Databases
- Redis for caching authentication token 
- PostgreSQL for storing users and posts
- Minio for keeping image files

Serving the service written in Go with http protocol and using JWT token for authentication

### How to run
1. Run docker services to serve databases 
```sh
docker-compose up -d
```

2. Execute the backend service
```
go run main.go
```

