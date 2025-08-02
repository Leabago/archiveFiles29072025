# Archive Service API 📦
App for creating and managing file archives from URLs.

## Run
```
go run ./cmd/app/
```

## Swagger
```
http://localhost:8080/doc/
```
## Config
```
config/config.yaml
```

## 🌟 API Examples

### 1. Send a set of links. Zip created in runtime.

**Request**:
```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/upload' \
  -H 'accept: application/zip' \
  -H 'Content-Type: application/json' \
  -d '{
  "links": [
    "https://www.sample-videos.com/img/Sample-jpg-image-50kb.jpg",
    "https://s24.q4cdn.com/216390268/files/doc_downloads/test.pdf"
  ]
}'
```

### 2. Create a new task
**Request**:
```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/task' \
  -H 'accept: application/json' \
  -d ''
 ```
 
 
### 3.  Adds a new link to the task's links array
**Request**:
```bash
curl -X 'PATCH' \
  'http://localhost:8080/api/v1/task/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "link": "https://sample-videos.com/img/Sample-jpg-image-50kb.jpg"
}'
 ```
 
### 4. Download archive by link.
**Request**:
```bash
curl -X 'GET' \
  'http://localhost:8080/api/v1/arch/2' \
  -H 'accept: application/zip'
}'
 ```
 
