# Usage
# Installation & Getting Started with ProDocs

Welcome to **ProDocs**! This guide will help you install ProDocs and get your documentation site up and running in minutes.

---

## 🚀 Prerequisites

- **Go** (version 1.18 or higher recommended)
- A terminal or command prompt
- (Optional) `git` if you want to clone the repo

---

## 🛠️ Installation

### 1. Download the ProDocs Source

You can either **clone the repository** or **download the source** directly.

#### Using Git

```sh
git clone https://github.com/sanjeev-prodevans/prodocs.git
cd prodocs
```

#### Or Download

- Download and extract the latest release from [GitHub Releases](https://github.com/sanjeev-prodevans/prodocs/releases).

---

### 2. Build the ProDocs Server

Make sure you are in the `prodocs` directory. Then run:

```sh
go build -o prodocs
```
This will create an executable called `prodocs` in your directory.

---

### 3. Initialize Your Docs

- Create a directory named `content/` in the project root.
- Add your Markdown files (e.g., `home.md`, `getting-started.md`) to the `content/` directory.
- Edit or create the `config.yaml` file to define your site structure, title, theme, and navigation.

**Example `config.yaml`:**

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
  - title: Installation
    path: installation.md
  - title: API Reference
    children:
      - title: Authentication
        path: api/auth.md
      - title: Data
        path: api/data.md
```

---

### 4. Start the Server

You can now start your documentation server:

```sh
./prodocs
```

By default, ProDocs will look for `config.yaml` and the `content/` folder in the same directory.  
You can specify custom paths using command-line flags:

```sh
./prodocs -config=myconfig.yaml -content=mydocs/
```

---

### 5. Open Your Documentation

Open your web browser and go to:

```
http://127.0.0.1:8080
```

You'll see your beautiful documentation site live!

---

## ⚡ Quick Start Example

1. **Write docs:** Edit Markdown files in `content/`.
2. **Edit navigation:** Update `config.yaml` to organize your sidebar and structure.
3. **Customize:** Change the theme and logo to match your brand.
4. **Search:** Use the search bar in the navbar to quickly find pages.
5. **Enjoy:** Your docs are live and easy to share!

---

## 💡 Tips

- Use headings (`#`, `##`, `###`) in your Markdown for automatic Table of Contents and deep-linking.
- Add code blocks and tables—ProDocs will style and highlight them for you.
- Organize your docs with nested navigation using the `children` key in `config.yaml`.

---

## 🆘 Need Help?

- Read more in the other docs pages.
- Check out the [ProDocs GitHub repository](https://github.com/sanjeev-prodevans/prodocs) for updates and support.
- Open an issue or join the community for questions and feedback.

---

Happy documenting with **ProDocs** from prodevans!  Documentation Made Easy

Special Thanks to **Gagan Pattnayak** and **Deepak Mishra** for his guidance and support in creating this documentation tool. Your expertise has been invaluable in shaping ProDocs into a powerful solution for developers everywhere!, Sanjeev Senapati, Shibram Mishra and the ProDocs team.

