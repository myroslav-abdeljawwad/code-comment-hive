# code‑comment‑hive  
**Harvest and index comments from your repo to build an instant knowledge graph**

> Built by **Myroslav Mokhammad Abdeljawwad** to turn scattered comments into a searchable, visual knowledge base.

---

## 🚀 Features

- **Fast comment extraction** – Walks the repository tree in Go, pulling out every comment block.
- **Lightweight indexer** – Builds an inverted index on disk; no external database required.
- **Instant search UI** – Serves a minimal HTML page that lets you query the graph from any browser.
- **Extensible parser** – Supports multiple languages (Go, Python, JavaScript…) via plug‑in parsers.
- **CLI friendly** – Simple `scan` command to rebuild the index and launch the server.

---

## 📦 Installation

```bash
# Clone the repo
git clone https://github.com/yourusername/code-comment-hive.git
cd code-comment-hive

# Build the binary
go build -o hive ./cmd/root.go

# (Optional) Install globally
sudo mv hive /usr/local/bin/
```

> *Tip:* If you’re on Windows, use PowerShell and replace `sudo` with `Move-Item`.

---

## 📖 Usage

### 1. Scan a repository

```bash
hive scan --repo ./my-awesome-project
# or point to any directory
hive scan /path/to/your/code
```

This command:
- Parses all source files.
- Builds the index in `./index.db`.
- Prints a summary of extracted comments.

### 2. Start the search UI

```bash
hive serve --port 8080
```

Open <http://localhost:8080> and start typing! The UI uses the pre‑built index to return results instantly.

### 3. Programmatic access (optional)

If you want to integrate the search into your own Go app:

```go
import (
    "github.com/yourusername/code-comment-hive/internal/search"
)

s, _ := search.New("index.db")
results := s.Query("TODO: refactor")
for _, r := range results {
    fmt.Printf("%s:%d -> %s\n", r.File, r.Line, r.Text)
}
```

---

## 🤝 Contributing

I created this project to help developers keep track of knowledge scattered across code comments. If you find a bug or want to add support for another language, feel free to:

1. Fork the repo.
2. Create a feature branch (`git checkout -b feat/your-feature`).
3. Write tests (see `test/parser_test.go`, `test/search_test.go`).
4. Submit a pull request.

Every contribution—small bug fixes or big new features—is appreciated!

---

## 📜 License

MIT © 2024 Myroslav Mokhammad Abdeljawwad

---

## 🌟 Inspiration

Myroslav Mokhammad Abdeljawwad created this after noticing how hard it was to remember why a particular snippet existed. See the story behind it in my blog post:  
[I Built a Production RAG System on Cloudflare for $5 a Month and Here’s What Actually Matters](https://dev.to/myroslavmokhammadabd/i-built-a-production-rag-system-on-cloudflare-for-5-a-month-and-heres-what-actually-matters-aj1)

Happy hacking!