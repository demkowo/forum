# Forum Service

## Overview
The **Forum Service** is a Golang-based API that manages comments, likes, dislikes, and complaints for an article-based discussion system. It is built with `Gin Gonic` for the HTTP router and `PostgreSQL` as the database. The service allows users to add and interact with comments, providing features like upvotes, downvotes, and complaint reporting.

## Features
- **Comment System**: Add, retrieve, delete, and count comments.
- **Like & Dislike System**: Users can upvote (like) and downvote (dislike) comments.
- **Complaint Handling**: Users can report inappropriate comments.
- **Soft Deletion**: Comments are soft-deleted to preserve discussion integrity.
- **Transaction Management**: Ensures atomicity in operations.

## Directory Structure
```
forum/
|-- app/          # Initiate app components and routes
│-- models/       # Contains data models
│-- repositories/ # Data access layer (PostgreSQL implementation)
│   ├── postgres/
│   │   ├── repository.go # Forum repository implementation
│-- handlers/     # HTTP handlers for API endpoints
|-- services/     # Business logic layer
│-- main.go       # Service entry point
```

## API Endpoints
| Method | Endpoint | Description |
|--------|-----------------------------------|------------------------------|
| `POST` | `/api/v1/comments/add` | Add a new comment |
| `DELETE` | `/api/v1/comments/delete/:comment_id` | Soft delete a comment |
| `GET`  | `/api/v1/comments/get/:comment_id` | Retrieve a specific comment |
| `GET`  | `/api/v1/comments/find` | Retrieve all comments |
| `GET`  | `/api/v1/comments/find/:article_id` | Retrieve comments for an article |
| `GET`  | `/api/v1/comments/count/:article_id` | Count comments for an article |
| `POST` | `/api/v1/likes/add` | Add a like to a comment |
| `DELETE` | `/api/v1/likes/delete` | Remove a like from a comment |
| `GET`  | `/api/v1/likes/count/:comment_id` | Count likes for a comment |
| `GET`  | `/api/v1/likes/find/:comment_id` | Retrieve likes for a comment |
| `POST` | `/api/v1/dislikes/add` | Add a dislike to a comment |
| `DELETE` | `/api/v1/dislikes/delete` | Remove a dislike from a comment |
| `GET`  | `/api/v1/dislikes/count/:comment_id` | Count dislikes for a comment |
| `GET`  | `/api/v1/dislikes/find/:comment_id` | Retrieve dislikes for a comment |
| `POST` | `/api/v1/complaints/add` | Report a comment |
| `DELETE` | `/api/v1/complaints/delete/:complaint_id` | Remove a complaint |
| `GET`  | `/api/v1/complaints/count/:comment_id` | Count complaints for a comment |
| `GET`  | `/api/v1/complaints/find/:comment_id` | Retrieve complaints for a comment |

## Database Schema
The service interacts with the following tables:

### `comments`
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    article_id UUID NOT NULL,
    thread_id UUID NOT NULL,
    parent_id UUID,
    author VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (parent_id) REFERENCES comments(id),
    FOREIGN KEY (article_id) REFERENCES articles(id),
    FOREIGN KEY (author) REFERENCES users(nickname)
);
```

### `likes`
```sql
CREATE TABLE likes (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
);
```

### `dislikes`
```sql
CREATE TABLE dislikes (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
);
```

### `complaints`
```sql
CREATE TABLE complaints (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    message TEXT,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
);
```

## Usage

### Add a Comment
```sh
curl -X POST http://localhost:8080/api/v1/comments/add -H "Content-Type: application/json" -d '{
    "article_id": "123e4567-e89b-12d3-a456-426614174000",
    "author": "john_doe",
    "content": "This is a comment."
}'
```

### Get a Comment by ID
```sh
curl -X GET http://localhost:8080/api/v1/comments/get/{comment_id}
```

### Like a Comment
```sh
curl -X POST http://localhost:8080/api/v1/likes/add -H "Content-Type: application/json" -d '{
    "comment_id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "456e7890-b12c-34d5-e678-910111213141"
}'
```

### Count Likes for a Comment
```sh
curl -X GET http://localhost:8080/api/v1/likes/count/{comment_id}
```

### Report a Comment
```sh
curl -X POST http://localhost:8080/api/v1/complaints/add -H "Content-Type: application/json" -d '{
    "comment_id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "456e7890-b12c-34d5-e678-910111213141",
    "message": "Inappropriate content."
}'
```

## Transactions & Error Handling
- All **write operations** (`AddComment`, `DeleteComment`, `AddLike`, etc.) use transactions to ensure atomicity.
- **Soft deletion** is implemented for comments to prevent accidental data loss.
- Errors are handled gracefully, returning appropriate HTTP status codes.

## Development Setup
### Prerequisites
- Golang (>=1.18)
- PostgreSQL
- `Gin Gonic`

### Install Dependencies
```sh
go mod tidy
```

### Run Service
```sh
go run main.go
```
