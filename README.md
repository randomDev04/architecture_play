# Go Backend Fundamentals: From Zero to Understanding

## Part 1: How HTTP Servers Actually Work

### The Absolute Basics

When you type `google.com` in your browser:
1. Browser sends an **HTTP Request** to Google's server
2. Server processes the request
3. Server sends back an **HTTP Response**

That's it. Everything else is details.

### HTTP Request Structure
```
GET /api/tasks HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGc...
Content-Type: application/json

{"title": "Learn Go"}  ← Body (optional, not for GET)
```

### HTTP Response Structure
```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 45

{"id": 1, "title": "Learn Go", "done": false}
```

---

## Part 2: Your First Go HTTP Server (No Framework)

### Step 1: The Simplest Possible Server

Create `main.go`:

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // This function runs for EVERY request to the server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

**Run it:**
```bash
go run main.go
# Visit: http://localhost:8080
```

**What's happening?**
1. `http.ListenAndServe(":8080", nil)` - Opens port 8080, waits for requests
2. When request comes in, Go creates a **goroutine** (lightweight thread) to handle it
3. Finds matching handler function
4. Executes the function
5. Returns response

**Key Insight:** Each request is handled in its own goroutine. This is why Go is fast - it can handle thousands of requests simultaneously without creating OS threads.

---

### Step 2: Understanding Request and Response

```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // r *http.Request contains EVERYTHING about the incoming request
    fmt.Printf("Method: %s\n", r.Method)        // GET, POST, PUT, DELETE
    fmt.Printf("Path: %s\n", r.URL.Path)        // /api/tasks
    fmt.Printf("Query: %s\n", r.URL.RawQuery)   // ?page=1&limit=10
    fmt.Printf("Headers: %v\n", r.Header)       // All HTTP headers
    
    // w http.ResponseWriter is how you WRITE the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // 200
    fmt.Fprintf(w, `{"message": "Request received"}`)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

**Test with curl:**
```bash
curl -X POST "http://localhost:8080/test?page=1" \
  -H "Authorization: Bearer token123" \
  -d '{"name": "Shubham"}'
```

---

### Step 3: Manual Routing (Why Frameworks Exist)

**Problem:** `http.HandleFunc("/", handler)` matches EVERYTHING.

How do you handle different routes?

**Method 1: Manual routing (the hard way)**

```go
package main

import (
    "fmt"
    "net/http"
)

func router(w http.ResponseWriter, r *http.Request) {
    // Manual routing - check path and method
    path := r.URL.Path
    method := r.Method

    if path == "/api/tasks" && method == "GET" {
        getTasks(w, r)
    } else if path == "/api/tasks" && method == "POST" {
        createTask(w, r)
    } else if path == "/api/users" && method == "GET" {
        getUsers(w, r)
    } else {
        // 404 Not Found
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, `{"error": "Not Found"}`)
    }
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `{"tasks": []}`)
}

func createTask(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `{"message": "Task created"}`)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `{"users": []}`)
}

func main() {
    http.HandleFunc("/", router)
    http.ListenAndServe(":8080", nil)
}
```

**Problems with this approach:**
1. ❌ No URL parameters (`/api/tasks/123` - how to extract `123`?)
2. ❌ Messy if-else chains
3. ❌ No middleware support (auth, logging)
4. ❌ Every project reinvents the wheel

**This is why frameworks exist.**

---

### Step 4: Enter ServeMux (Go's Basic Router)

Go has a built-in router: `http.ServeMux`

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    // Register routes
    mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            fmt.Fprintf(w, `{"tasks": []}`)
        } else if r.Method == "POST" {
            fmt.Fprintf(w, `{"message": "Task created"}`)
        } else {
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `{"users": []}`)
    })

    fmt.Println("Server on :8080")
    http.ListenAndServe(":8080", mux)
}
```

**Better, but still limited:**
- No URL parameters
- Still manually checking HTTP methods
- No middleware

---

### Step 5: Why Use Gin/Echo/Chi?

**Frameworks add:**
1. **Better routing** (with URL parameters)
2. **Middleware support** (auth, logging, CORS)
3. **JSON helpers** (auto marshal/unmarshal)
4. **Validation**
5. **Error handling**

**Example with Gin:**

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default() // Includes logger and recovery middleware

    // Clean routing
    r.GET("/api/tasks", getTasks)
    r.POST("/api/tasks", createTask)
    r.GET("/api/tasks/:id", getTask)  // URL parameter!
    r.PUT("/api/tasks/:id", updateTask)
    r.DELETE("/api/tasks/:id", deleteTask)

    r.Run(":8080")
}

func getTasks(c *gin.Context) {
    c.JSON(200, gin.H{"tasks": []string{}})  // Auto JSON!
}

func getTask(c *gin.Context) {
    id := c.Param("id")  // Extract URL parameter
    c.JSON(200, gin.H{"id": id})
}

func createTask(c *gin.Context) {
    var task struct {
        Title string `json:"title" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&task); err != nil {  // Auto validate!
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, gin.H{"message": "Created", "title": task.Title})
}

func updateTask(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Updated"})
}

func deleteTask(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Deleted"})
}
```

**Much cleaner!**

---

## Part 3: Database Choices

### The Options

| Database | Type | When to Use | When NOT to Use |
|----------|------|-------------|-----------------|
| **PostgreSQL** | Relational (SQL) | - Complex relationships<br>- ACID transactions needed<br>- Structured data<br>- Reporting/analytics | - Extreme write speed needed<br>- Unstructured data<br>- Horizontal scaling (without complexity) |
| **MySQL** | Relational (SQL) | - Similar to PostgreSQL<br>- Simpler setup<br>- Read-heavy workloads | - Complex JSON queries<br>- Advanced features (PostgreSQL better) |
| **MongoDB** | Document (NoSQL) | - Flexible schema<br>- Rapid prototyping<br>- Unstructured data<br>- Easy horizontal scaling | - Complex joins needed<br>- Strong consistency required<br>- Financial data |
| **Redis** | Key-Value (In-memory) | - Caching<br>- Sessions<br>- Real-time data<br>- Pub/sub | - Primary database<br>- Persistent data (can, but risky)<br>- Complex queries |
| **SQLite** | Relational (File-based) | - Embedded apps<br>- Local storage<br>- Prototyping<br>- Mobile apps | - Multi-user web apps<br>- High concurrency<br>- Large scale |

### For Your Task Manager: PostgreSQL

**Why?**
1. **Relationships:** Users → Tasks (one-to-many)
2. **ACID transactions:** Task creation + notification should be atomic
3. **Data integrity:** Foreign keys prevent orphaned data
4. **Querying:** Complex filters (completed tasks, overdue, by user)
5. **Industry standard:** 90% of backend jobs use relational DBs

### SQL vs NoSQL Decision Tree

```
Do you have clear relationships between data? (users, posts, comments)
├─ YES → SQL (PostgreSQL/MySQL)
└─ NO → Do you need flexible schema?
    ├─ YES → NoSQL (MongoDB)
    └─ NO → Do you need extreme speed for simple lookups?
        ├─ YES → Redis/DynamoDB
        └─ NO → Still use PostgreSQL (safest choice)
```

---

## Part 4: Core Backend Concepts

### 1. Request Lifecycle

```
Client Request
    ↓
Load Balancer (optional)
    ↓
Web Server (your Go app)
    ↓
Middleware (auth, logging, etc.)
    ↓
Route Handler
    ↓
Business Logic
    ↓
Database Query
    ↓
Response Formatting
    ↓
Client Response
```

### 2. Middleware Pattern

**What:** Code that runs BEFORE your handler

**Why:** Don't repeat auth/logging in every handler

```go
// Middleware function
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            w.WriteHeader(401)
            fmt.Fprintf(w, `{"error": "Unauthorized"}`)
            return
        }
        
        // Verify token...
        
        next(w, r)  // Call the actual handler
    }
}

// Usage
http.HandleFunc("/api/tasks", authMiddleware(getTasks))
```

**Middleware Chain:**
```
Request → Logger → Auth → Rate Limiter → Your Handler → Response
```

### 3. Database Connection Pooling

**Problem:** Opening a DB connection is SLOW (50-100ms)

**Solution:** Connection pool - reuse connections

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("postgres", "connection_string")
    if err != nil {
        panic(err)
    }
    
    // THIS IS CRITICAL FOR PERFORMANCE
    db.SetMaxOpenConns(25)      // Max connections to DB
    db.SetMaxIdleConns(5)       // Keep 5 idle connections ready
    db.SetConnMaxLifetime(5*60) // Recycle connections every 5 min
}
```

**Why it matters:**
- Without pool: Each request opens new connection (SLOW)
- With pool: Request grabs existing connection (FAST)

**Tradeoff:**
- Too few connections → requests wait
- Too many connections → database overwhelmed

**Rule of thumb:** `MaxOpenConns = (CPU cores * 2) + disk drives`

### 4. Error Handling

**Go forces you to handle errors explicitly**

```go
// Bad (panics if error)
rows := db.Query("SELECT * FROM tasks")

// Good
rows, err := db.Query("SELECT * FROM tasks")
if err != nil {
    log.Printf("Query failed: %v", err)
    http.Error(w, "Internal Server Error", 500)
    return
}
defer rows.Close()  // Always cleanup!
```

**HTTP Status Codes (Know These):**
- `200 OK` - Success
- `201 Created` - Resource created (POST)
- `204 No Content` - Success, no response body (DELETE)
- `400 Bad Request` - Client sent invalid data
- `401 Unauthorized` - No/invalid auth token
- `403 Forbidden` - Authenticated but not allowed
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server crashed

### 5. JSON Handling

```go
// Struct tags define JSON mapping
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    IsCompleted bool      `json:"is_completed"`
    CreatedAt   time.Time `json:"created_at"`
}

// Marshal (Go → JSON)
task := Task{ID: 1, Title: "Learn Go"}
jsonBytes, err := json.Marshal(task)
// Result: {"id":1,"title":"Learn Go","is_completed":false}

// Unmarshal (JSON → Go)
var task Task
err := json.Unmarshal(jsonBytes, &task)
```

### 6. Concurrency (Goroutines)

**Why Go is fast:**

```go
// Sequential (SLOW)
func handleRequests() {
    for _, req := range requests {
        processRequest(req)  // Blocks until done
    }
}

// Concurrent (FAST)
func handleRequests() {
    for _, req := range requests {
        go processRequest(req)  // Runs in parallel!
    }
}
```

**Each request gets its own goroutine automatically:**
```go
http.HandleFunc("/api/tasks", getTasks)
// Go creates goroutine for each request to getTasks
```

**Key:** Goroutines are cheap (~2KB memory). You can have 100,000+ goroutines.

### 7. Context (Request Cancellation)

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Long database query
    rows, err := db.QueryContext(ctx, "SELECT * FROM tasks")
    
    // If client disconnects, query is cancelled!
}
```

**Why:** If user closes browser, don't waste server resources.

---

## Part 5: Putting It All Together

### Architecture of a Real Backend

```
┌─────────────────────────────────────────────┐
│  Client (React Native App)                  │
└────────────┬────────────────────────────────┘
             │ HTTP Request
             ↓
┌─────────────────────────────────────────────┐
│  Load Balancer (nginx/HAProxy)              │
│  - Distributes load across servers          │
└────────────┬────────────────────────────────┘
             │
        ┌────┴────┐
        ↓         ↓
    ┌───────┐ ┌───────┐
    │Server1│ │Server2│  (Your Go Apps)
    └───┬───┘ └───┬───┘
        │         │
        └────┬────┘
             ↓
┌─────────────────────────────────────────────┐
│  Database Connection Pool                   │
└────────────┬────────────────────────────────┘
             ↓
┌─────────────────────────────────────────────┐
│  PostgreSQL Database                        │
│  - Primary (writes)                         │
│  - Replica (reads)                          │
└─────────────────────────────────────────────┘
             ↑
             │ Cache queries
┌─────────────────────────────────────────────┐
│  Redis (Cache Layer)                        │
│  - Session storage                          │
│  - API response cache                       │
└─────────────────────────────────────────────┘
```

### Request Flow Example

```
1. User clicks "Get Tasks" in React Native app
2. App sends: GET /api/tasks
3. Load balancer picks Server1
4. Server1 middleware chain:
   ├─ Logger middleware (logs request)
   ├─ Auth middleware (verifies JWT)
   └─ Rate limiter (checks request count)
5. Handler function executes:
   ├─ Check Redis cache
   ├─ If miss: Query PostgreSQL
   ├─ Store in Redis for 5 min
   └─ Return JSON response
6. Response travels back to client
```

---

## Next Steps

Now you understand:
- ✅ How HTTP servers work (net/http)
- ✅ Why frameworks exist (Gin/Echo simplify routing)
- ✅ Database choices (PostgreSQL for your case)
- ✅ Core concepts (middleware, pooling, errors, JSON, concurrency)

**Ready to build?**

Your homework before we code:
1. Draw the architecture of your Task Manager (boxes and arrows)
2. Define your database schema (users table, tasks table)
3. List your API endpoints (GET /api/tasks, POST /api/tasks, etc.)

Once you've done this, we'll write actual code.