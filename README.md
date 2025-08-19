# Write Base API Documentation

## Overview
Write Base is a Go REST API for a modern blogging platform. It lets you manage articles with rich content blocks, tags, claps, and views. It uses MongoDB for storage and integrates with Gemini for optional content generation.

Highlights:
- Clean architecture (controller → usecase → repository).
- Fast endpoints via ESR-aligned MongoDB indexes and $text search (status-prefixed, language set to none).
- Auth with JWT, role-based admin operations, and DTO validation.

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
- [Performance Benchmarks](#performance-benchmarks)
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

### Auth & Headers
- Most endpoints require a Bearer token: `Authorization: Bearer <JWT>`.
- Admin routes additionally require the user to have the admin role.

### Pagination
- Query parameters: `page` (default 1), `page_size` (default 20, capped by server policy).
- Responses include: `data`, `total`, `page`, `page_size`, and sometimes `total_pages`.

### Errors
- JSON error format: `{ "error": "message" }` with appropriate HTTP status code.

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
## Usage Notes

### Articles
- Create requires valid content blocks and tags. Slug is auto-generated from title if omitted.
- Publish checks that all tags are approved; otherwise returns 400.
- Get-by-ID returns drafts to authors; for others, only if the article is published. Views are recorded for published content.
- Search uses MongoDB `$text` and is limited to published content.

### Tags
- Users propose tags; admins approve/reject. Unapproved tags can be used in drafts but block publish.

### Claps & Views
- Claps are rate-limited per user per article; exceeding returns HTTP 429.
- Views increment on published content reads; both anonymous (via client IP) and authenticated views are supported.

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

## Testing & Coverage

Run tests with coverage:

```powershell
go test ./... -coverprofile=coverage.out
```

As of 2025-08-19:

- Overall coverage (all packages): 10.2%

Core module coverage (selected packages):
- config: 100.0%
- internal/domain: 100.0%
- internal/infrastructure: 87.3%
- internal/infrastructure/utils: 64.3%
- internal/delivery/http/router: 65.3%
- internal/delivery/http/controller: 49.9%
- internal/repository: 35.3%
- internal/usecase: 25.4%
- pkg/di: 18.3%

Notes:
- “Overall” includes non-core or generated code (e.g., mocks, ai stubs, cmd), which may skew totals.
- “Core modules” focus on the HTTP layer, business logic, repositories, infra, and domain.
- Coverage was measured on Windows PowerShell; per-file totals are available via `go tool cover -html=coverage.out`.

---

## Performance Benchmarks

This project includes reproducible Go benchmarks for the MongoDB article repository to measure the effect of database indexing.

Hardware/OS for the runs below: Windows 11, Intel i5-1155G7. Your results will vary; use the exact commands to reproduce on your machine.

### Index sets compared
- none: No indexes (baseline). Search is expected to fail without the text index.
- text: Only the compound text index on `{status, title, excerpt}`.
- full: ESR-oriented indexes plus the same compound text index (default when not set).

The text index keys: `{ status: 1, title: "text", excerpt: "text" }` with `default_language = none` and `language_override = none` to avoid unsupported language overrides on some servers.

### Reproducible commands (PowerShell)
Use a larger seed and longer benchtime to reduce noise.

```powershell
# Baseline: no indexes (Search will fail, skip it if needed)
$env:BENCH_SEED_COUNT = "20000"; $env:BENCH_INDEX_MODE = "none";
go test -run=^$ ./internal/repository -bench=BenchmarkRepo_ -benchmem -benchtime=5s -count=1

# Text-only: enable Search index only
$env:BENCH_SEED_COUNT = "20000"; $env:BENCH_INDEX_MODE = "text";
go test -run=^$ ./internal/repository -bench=BenchmarkRepo_ -benchmem -benchtime=5s -count=1

# Full ESR + Text: production-like indexing
Remove-Item Env:BENCH_INDEX_MODE -ErrorAction SilentlyContinue
Remove-Item Env:BENCH_NO_INDEX -ErrorAction SilentlyContinue
$env:BENCH_SEED_COUNT = "20000";
go test -run=^$ ./internal/repository -bench=BenchmarkRepo_ -benchmem -benchtime=5s -count=1
```

### Results summary (ns/op, docs/op)
From a representative run on the machine above:

- Full vs None
   - ListByAuthor: 28.41ms → 2.37ms (~+1098% faster)
   - Filter(Tags+Published): 37.20ms → 6.05ms (~+514% faster)
   - Trending: 33.36ms → 3.56ms (~+838% faster)
   - Popular: 28.76ms → 6.82ms (~+321% faster)
   - Create: 0.165ms → 0.694ms (slower due to index maintenance)
   - Search: baseline fails without text index; see Text-only vs Full below.

- Text-only vs Full
   - Search: 105.21ms (Full; $text + textScore; includes `status` equality)
   - Text-only mode enables Search but may be slower for other queries compared to Full due to missing ESR indexes.

Notes:
- docs/op is reported as 20 for list-style queries (page size = 20), confirming comparable work.
- Search without a text index is unsupported and will error; use Text or Full modes to measure Search.

### Why these indexes
- ESR rule: We build compound indexes to match Equality, Sort, then Range filters used by core queries.
   - Examples: `(status, stats.view_count, timestamps.published_at)` for Trending; `(status, tags, timestamps.created_at)` for tag filters; `(author_id, status, timestamps.created_at)` for author lists.
- Search uses a single text index with `status` as a prefix key to align with the equality filter (`published`).

These indexes drastically reduce scanned documents and improve latency for reads, at the cost of slightly slower writes (Create/Update) due to index updates.

### Interpreting and extending
- For tighter comparisons, run each benchmark multiple times (`-count=3`) and increase `-benchtime`.
- You can set `BENCH_SEED_COUNT` to scale dataset size.
- To analyze query plans, run MongoDB `explain()` from the shell on representative queries.

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
