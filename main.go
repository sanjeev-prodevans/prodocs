package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/yuin/goldmark"
	// headingid "github.com/yuin/goldmark-extensions/headingid"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed assets/*
var assetsFS embed.FS

var (
	configPath string
	contentDir string
)

type Page struct {
	Title    string  `yaml:"title"`
	Path     string  `yaml:"path,omitempty"`
	Children []*Page `yaml:"children,omitempty"`
}

type Config struct {
	SiteName        string `yaml:"site_name"`
	SiteDescription string `yaml:"site_description"`
	Theme           struct {
		PrimaryColor string `yaml:"primary_color"`
		AccentColor  string `yaml:"accent_color"`
		Logo         string `yaml:"logo"`
	} `yaml:"theme"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Pages []*Page `yaml:"pages"`
}

type DocEntry struct {
	Path    string // URL path (e.g., "/intro")
	Title   string // Title from the page config
	Content string // Raw Markdown content
}

var (
	cfg       Config
	mdParser  goldmark.Markdown
	baseTmpl  *template.Template
	accessLog *os.File
	errorLog  *os.File
	docIndex  []DocEntry
)

func logAccess(format string, v ...interface{}) {
	if accessLog != nil {
		fmt.Fprintf(accessLog, "%s ACCESS: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, v...))
	}
}
func logError(format string, v ...interface{}) {
	if errorLog != nil {
		fmt.Fprintf(errorLog, "%s ERROR: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, v...))
	}
}

func loadConfig() {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
}

func setupMarkdown() {
	mdParser = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			highlighting.NewHighlighting(),
			// extension.NewHeadingID(), // For heading anchors: id="..."
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

func renderMarkdown(mdPath string) (string, error) {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return "", err
	}
	var sb bytes.Buffer
	if err := mdParser.Convert(data, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// Goldmark-compatible anchor generator (matches heading IDs)
func anchorFromTitle(title string) string {
	var buf strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			buf.WriteRune(r)
			lastDash = false
		case unicode.IsSpace(r) || r == '-':
			if !lastDash {
				buf.WriteRune('-')
				lastDash = true
			}
		}
		// ignore everything else
	}
	anchor := buf.String()
	anchor = strings.Trim(anchor, "-")
	return anchor
}

func extractTOC(mdPath string) template.HTML {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	headingRe := regexp.MustCompile(`^(#{1,3})\s+(.*)`)
	var toc []string
	for _, line := range lines {
		m := headingRe.FindStringSubmatch(line)
		if m != nil {
			level := len(m[1])
			text := m[2]
			anchor := anchorFromTitle(text)
			toc = append(toc, fmt.Sprintf(`<li class="toc-level-%d"><a href="#%s">%s</a></li>`, level, anchor, text))
		}
	}
	if len(toc) == 0 {
		return ""
	}
	return template.HTML(`<div class="toc"><b>On this page</b><ul>` + strings.Join(toc, "\n") + `</ul></div>`)
}

func buildPageRoutes(pages []*Page, prefix string, urlToFile map[string]string, urlToTitle map[string]string) {
	for _, page := range pages {
		pageUrl := prefix
		if page.Path != "" {
			base := strings.TrimSuffix(page.Path, ".md")
			if base == "index" && prefix == "" {
				pageUrl = "/"
			} else {
				pageUrl = "/" + strings.ReplaceAll(base, "\\", "/")
			}
			urlToFile[pageUrl] = filepath.Join(contentDir, page.Path)
			urlToTitle[pageUrl] = page.Title
		}
		if len(page.Children) > 0 {
			buildPageRoutes(page.Children, pageUrl, urlToFile, urlToTitle)
		}
	}
}

func renderNavSidebar(pages []*Page, currentPath string) template.HTML {
	var buf strings.Builder
	buf.WriteString(`<ul class="sidebar">`)
	for _, page := range pages {
		if page.Path != "" {
			base := strings.TrimSuffix(page.Path, ".md")
			var url string
			if base == "index" {
				url = "/"
			} else {
				url = "/" + strings.ReplaceAll(base, "\\", "/")
			}
			currentClass := ""
			if url == currentPath {
				currentClass = ` class="active"`
			}
			buf.WriteString(fmt.Sprintf(`<li><a href="%s"%s>%s</a>`, url, currentClass, page.Title))
		} else {
			buf.WriteString(fmt.Sprintf(`<li><span>%s</span>`, page.Title))
		}
		if len(page.Children) > 0 {
			buf.WriteString(string(renderNavSidebar(page.Children, currentPath)))
		}
		buf.WriteString(`</li>`)
	}
	buf.WriteString(`</ul>`)
	return template.HTML(buf.String())
}

func loadTemplates() {
	tmplData, err := templatesFS.ReadFile("templates/base.html")
	if err != nil {
		log.Fatalf("template parse error: %v", err)
	}
	baseTmpl, err = template.New("base.html").Parse(string(tmplData))
	if err != nil {
		log.Fatalf("template parse error: %v", err)
	}
}

// --- SEARCH CAPABILITY ---

func buildDocIndex(urlToFile, urlToTitle map[string]string) {
	docIndex = nil
	for url, file := range urlToFile {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		docIndex = append(docIndex, DocEntry{
			Path:    url,
			Title:   urlToTitle[url],
			Content: string(content),
		})
	}
}

func searchDocs(query string) []DocEntry {
	if query == "" {
		return nil
	}
	query = strings.ToLower(query)
	var results []DocEntry
	for _, entry := range docIndex {
		if strings.Contains(strings.ToLower(entry.Title), query) ||
			strings.Contains(strings.ToLower(entry.Content), query) {
			results = append(results, entry)
		}
	}
	return results
}

// --- MAIN ---

func main() {
	flag.StringVar(&configPath, "config", "config.yaml", "Path to config.yaml")
	flag.StringVar(&contentDir, "content", "content", "Path to content directory")
	flag.Parse()

	var err error
	accessLog, err = os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to open access.log: %v", err)
	}
	defer accessLog.Close()
	errorLog, err = os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to open error.log: %v", err)
	}
	defer errorLog.Close()

	loadConfig()
	setupMarkdown()
	loadTemplates()

	urlToFile := make(map[string]string)
	urlToTitle := make(map[string]string)
	buildPageRoutes(cfg.Pages, "", urlToFile, urlToTitle)
	buildDocIndex(urlToFile, urlToTitle)

	assetsSubFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Fatalf("could not sub assets: %v", err)
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsSubFS))))

	// SEARCH HANDLER
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		results := searchDocs(q)
		data := map[string]interface{}{
			"Results":         results,
			"Query":           q,
			"SiteName":        cfg.SiteName,
			"SiteDescription": cfg.SiteDescription,
			"Theme":           cfg.Theme,
			"Sidebar":         renderNavSidebar(cfg.Pages, ""),
			"PageTitle":       "Search Results",
			"TOC":             "",
			"Content":         "",
		}
		err := baseTmpl.Execute(w, data)
		if err != nil {
			logError("Template error in search: %v", err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if reqPath != "/" && strings.HasSuffix(reqPath, "/") {
			reqPath = strings.TrimSuffix(reqPath, "/")
		}
		clientIP := r.RemoteAddr

		mdPath, ok := urlToFile[reqPath]
		pageTitle := urlToTitle[reqPath]
		if !ok {
			logAccess("MISS  %s %s (not found)", clientIP, reqPath)
			logError("404 NOT FOUND %s for %s", reqPath, clientIP)
			http.NotFound(w, r)
			return
		}
		mdHTML, err := renderMarkdown(mdPath)
		if err != nil {
			logAccess("FAIL  %s %s (render error)", clientIP, reqPath)
			logError("Render error for %s from %s: %v", reqPath, clientIP, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		toc := extractTOC(mdPath)
		data := map[string]interface{}{
			"SiteName":        cfg.SiteName,
			"SiteDescription": cfg.SiteDescription,
			"Theme":           cfg.Theme,
			"Content":         template.HTML(mdHTML),
			"PageTitle":       pageTitle,
			"Sidebar":         renderNavSidebar(cfg.Pages, reqPath),
			"TOC":             toc,
			"Query":           "",
			"Results":         nil,
		}
		err = baseTmpl.Execute(w, data)
		if err != nil {
			logAccess("FAIL  %s %s (template error)", clientIP, reqPath)
			logError("Template render error for %s from %s: %v", reqPath, clientIP, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		logAccess("HIT   %s %s (success)", clientIP, reqPath)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("ProDocs server running at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
