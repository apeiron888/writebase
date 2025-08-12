# Write Base API Documentation

## Overview
The **Write Base API** is a **Go-based RESTful API** for managing articles, tags, claps, and views.  

It provides endpoints for creating, updating, retrieving, and managing articles with rich content blocks, as well as handling user interactions like **claps** and **views**.  

The API uses **MongoDB** as the database and integrates with the **Gemini API** for content generation.

---

## Table of Contents
- [Architecture](#architecture)
- [Installation](#installation)
- [API Endpoints](#api-endpoints)
  - [Article Endpoints](#article-endpoints)
  - [Tag Endpoints](#tag-endpoints)
  - [Clap Endpoints](#clap-endpoints)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Scenarios](#scenarios)
- [Environment Variables](#environment-variables)
- [Contributing](#contributing)
- [License](#license)

---

## Architecture
The API follows a **clean architecture** pattern with the following layers:

- **Controller**: Handles HTTP requests and responses using the **Gin** framework.
- **Usecase**: Contains business logic and orchestrates interactions between repositories and policies.
- **Repository**: Manages database operations (**MongoDB**).
- **Domain**: Defines core data models, interfaces, and constants.
- **Policy**: Enforces business rules and validations.
- **Utils**: Provides helper functions like UUID and slug generation.

---

## Installation

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd write_base
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables**  
   Create a `.env` file in the root directory with the following:
   ```env
   MONGODB_URI=mongodb://localhost:27017
   MONGODB_NAME=write_base
   JWT_SECRET=your_jwt_secret
   SERVER_PORT=8080
   GEMINI_API_KEY=your_gemini_api_key
   ```

4. **Run the Application**
   ```bash
   go run main.go
   ```

5. **Build for Production**
   ```bash
   go build -o write_base
   ./write_base
   ```

---

## API Endpoints

### Article Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| **POST** | `/articles/new` | Create a new article | User |
| **PUT** | `/articles/:id` | Update an existing article | User |
| **DELETE** | `/articles/:id` | Soft delete an article | User |
| **PATCH** | `/articles/:id/restore` | Restore a soft-deleted article | User |
| **GET** | `/articles/:id` | Retrieve an article by ID | User |
| **POST** | `/articles/:id/publish` | Publish an article | User |
| **POST** | `/articles/:id/unpublish` | Unpublish an article | User |
| **POST** | `/articles/:id/archive` | Archive an article | User |
| **POST** | `/articles/:id/unarchive` | Unarchive an article | User |
| **GET** | `/articles/:id/stats` | Get article statistics | User |
| **GET** | `/articles/stats/all` | Get all article stats for a user | User |
| **GET** | `/:slug` | Retrieve an article by slug | Optional |
| **GET** | `/authors/:author_id/articles` | List articles by author | User |
| **GET** | `/articles/trending` | List trending articles (last 7 days) | User |
| **GET** | `/articles/new` | List newest articles | User |
| **GET** | `/articles/popular` | List popular articles | User |
| **POST** | `/authors/:author_id/articles/filter` | Filter articles by author | User |
| **POST** | `/articles/filter` | Filter articles for all users | User |
| **GET** | `/search?q=<query>` | Search articles by query | User |
| **GET** | `/article/tags?tags=<tag1>,<tag2>` | List articles by tags | User |
| **DELETE** | `/me/trash` | Empty user's trash | User |
| **DELETE** | `/articles/trash/:id` | Permanently delete article from trash | User |
| **POST** | `/generateslug` | Generate slug from title | User |
| **POST** | `/articles/generatecontent` | Generate article content using Gemini API | User |

**Admin Endpoints**
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| **GET** | `/admin/articles` | List all articles | Admin |
| **DELETE** | `/admin/articles/:id/delete` | Hard delete an article | Admin |
| **POST** | `/admin/articles/:id/unpublish` | Unpublish an article | Admin |

---

### Tag Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| **POST** | `/tags/new` | Create a new tag | User |
| **GET** | `/tags` | List tags by status | User |
| **PATCH** | `/tags/:id/approve` | Approve a tag | Admin |
| **PATCH** | `/tags/:id/reject` | Reject a tag | Admin |
| **DELETE** | `/tags/:id` | Delete a tag | Admin |

---

### Clap Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| **POST** | `/articles/:id/clap` | Add a clap to an article | User |

---

## Data Models

### Article
```go
type Article struct {
    ID            string
    Title         string
    Slug          string
    AuthorID      string
    ContentBlocks []ContentBlock
    Excerpt       string
    Language      string
    Tags          []string
    Status        ArticleStatus
    Stats         ArticleStats
    Timestamps    ArticleTimes
}
```
- **ContentBlock**: Types include `heading`, `paragraph`, `image`, `code`, `video_embed`, `list`, `divider`.
- **ArticleStatus**: `draft`, `scheduled`, `published`, `archived`, `deleted`.
- **ArticleStats**: Tracks `ViewCount` and `ClapCount`.
- **ArticleTimes**: Tracks `CreatedAt`, `UpdatedAt`, `PublishedAt`, `ArchivedAt`.

---

### Tag
```go
type Tag struct {
    ID        string
    Name      string
    Status    TagStatus
    CreatedBy string
    CreatedAt time.Time
}
```
- **TagStatus**: `pending`, `approved`, `rejected`.

---

### Clap
```go
type Clap struct {
    ID        string
    UserID    string
    ArticleID string
    Count     int
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

### View
```go
type View struct {
    ID        string
    UserID    string
    ArticleID string
    ClientIP  string
    CreatedAt time.Time
}
```

---

## Error Handling
Standard error format:
```go
type Error struct {
    Code    string
    Message string
}
```
**Common error codes:**
- `GEN001`: Internal server error
- `USER001`: Unauthorized
- `ARTICLE001`: Invalid article payload
- `ARTICLE002`: Article not found
- `CLAP001`: Clap limit exceeded
- `TAG001`: Tag not found

---

## Scenarios

### 1. Creating and Publishing an Article
```bash
# Create article
curl -X POST http://localhost:8080/articles/new \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{
    "title": "My First Article",
    "slug": "my-first-article",
    "content_blocks": [{"type": "paragraph", "order": 1, "content": {"paragraph": {"text": "Hello, world!"}}}],
    "excerpt": "A brief introduction",
    "language": "en",
    "tags": ["tech", "intro"]
}'
```
Response:
```json
{ "data": { "id": "<article_id>" } }
```

```bash
# Publish article
curl -X POST http://localhost:8080/articles/<article_id>/publish \
-H "Authorization: Bearer <token>"
```
Response:
```json
{ "data": { "id": "<article_id>", "status": "published" } }
```

---

### 2. Filtering Articles by Author
```bash
curl -X POST http://localhost:8080/authors/<author_id>/articles/filter \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{
    "filter": {"statuses": ["published"], "tags": ["tech"]},
    "pagination": {"page": 1, "page_size": 10}
}'
```
Response:
```json
{ "data": [...], "total": 5, "page": 1, "page_size": 10 }
```

---

### 3. Adding a Clap
```bash
curl -X POST http://localhost:8080/articles/<article_id>/clap \
-H "Authorization: Bearer <token>"
```
Response:
```json
{ "view_count": 10, "clap_count": 5 }
```

---

### 4. Admin Operations
```bash
# List all articles
curl -X GET http://localhost:8080/admin/articles?page=1&page_size=20 \
-H "Authorization: Bearer <admin_token>"
```
Response:
```json
{ "data": [...], "total": 100, "page": 1, "page_size": 20 }
```

```bash
# Hard delete article
curl -X DELETE http://localhost:8080/admin/articles/<article_id>/delete \
-H "Authorization: Bearer <admin_token>"
```
Response:
```
204 No Content
```

---

## Environment Variables
| Variable         | Description                           | Required |
|------------------|---------------------------------------|----------|
| `MONGODB_URI`    | MongoDB connection URI                | Yes |
| `MONGODB_NAME`   | MongoDB database name                 | Yes |
| `JWT_SECRET`     | Secret key for JWT authentication     | Yes |
| `SERVER_PORT`    | Port for the HTTP server              | Yes |
| `GEMINI_API_KEY` | API key for Gemini content generation | Yes |

---

## Contributing
1. Fork the repository.
2. Create a feature branch:  
   ```bash
   git checkout -b feature/<feature_name>
   ```
3. Commit your changes:  
   ```bash
   git commit -m 'Add feature'
   ```
4. Push to your branch:  
   ```bash
   git push origin feature/<feature_name>
   ```
5. Create a pull request.

---

## License
MIT License
