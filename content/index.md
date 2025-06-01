# Welcome to ProDocs

ProDocs is a minimalist, blazing-fast, and beautiful documentation server written in Go. It is designed to help you organize, search, and present your documentation with ease. Below you'll find a showcase of ProDocs' key features and usage examples.

---

## 🚀 Features

- **Markdown-powered Content:**  
  Write your documentation in Markdown and ProDocs will render it with beautiful styling, code highlighting, and automatic anchor links for headings.

- **Automatic Table of Contents:**  
  Each page can display a table of contents (TOC) for easy navigation, generated from document headings.

- **Sidebar Navigation:**  
  Configure your documentation structure in a YAML config file for an instant sidebar menu with support for nested sections.

- **Responsive & Modern UI:**  
  Enjoy a clean, Material-inspired layout that adapts smoothly to all device sizes.

- **Instant Search:**  
  Use the search bar in the navbar to quickly find relevant pages or keywords in your documentation.

- **Code Highlighting:**  
  Syntax-highlighted code blocks with a light lime background for easy reading.

- **Heading Anchors:**  
  All headings have unique anchor IDs for easy deep linking and TOC navigation.

- **Customizable Theme:**  
  Easily change colors and logo via the configuration file.

- **Fast & Secure:**  
  ProDocs is built with Go and serves static content quickly and safely.

---

## 📝 Example Markdown Features

### Headings with Anchors

Click any Table of Contents entry or link like [Code Example](#code-example) to jump to that section.

---

### Code Example

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, ProDocs!")
}
```

---

### Table Example

| Feature            | Description                                |
|--------------------|--------------------------------------------|
| Markdown Support   | Write docs in plain Markdown               |
| Responsive Design  | Works on desktop and mobile                |
| Fast Search        | Instant search results                     |

---

### Blockquotes

> ProDocs makes documentation simple and beautiful.

---

### Lists

#### Unordered List
- Easy setup
- Instant navigation
- Custom themes

#### Ordered List
1. Download ProDocs
2. Add your Markdown files
3. Start the server

---

## 🔍 Searching

Use the search bar in the top-right corner to find any keyword, topic, or function name in your docs.

---

## 🎨 Theming

You can customize the primary color, accent color, and logo from the `config.yaml` file.

---

## ⚙️ Configuration

Define your documentation structure and navigation in a simple YAML file, for example:

```yaml
site_name: My Project Docs
site_description: Beautiful docs made easy!
theme:
  primary_color: "#2196f3"
  accent_color: "#ffeb3b"
  logo: "assets/logo.svg"
server:
  host: "127.0.0.1"
  port: 8080
pages:
  - title: Home
    path: home.md
  - title: Getting Started
    path: getting-started.md
  - title: API Reference
    children:
      - title: Authentication
        path: api/auth.md
      - title: Data
        path: api/data.md
```

---

## 📚 How to Use ProDocs

1. **Write Markdown files** in the `content/` directory.
2. **Configure your docs** in `config.yaml`.
3. **Run ProDocs:**  
   ```sh
   go run main.go
   ```
4. **Open your browser** and enjoy your docs!

---

## ❤️ Enjoy ProDocs!

Need help or want to contribute? Visit the [GitHub repository](https://github.com/yourorg/prodocs) or check the sidebar for more topics!