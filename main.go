package main

import (
    "fmt"
    "html/template"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
)

type FileEntry struct {
    Name  string
    Href  string
    Size  string
    Mtime string
    IsDir bool
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Index of {{.Title}} - Shaan OS Repositories</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --bg-primary: #0a0a0a;
            --bg-secondary: #1a1a1a;
            --bg-hover: #252525;
            --fg-primary: #bfff00;
            --fg-secondary: #80aa00;
            --fg-muted: #6b8e23;
            --border: #2a4a0a;
            --link: #bfff00;
            --link-hover: #d4ff33;
            --text: #e0e0e0;
            --text-dim: #9ca3af;
        }

        body {
            font-family: 'Source Code Pro', 'Consolas', 'Monaco', monospace;
            background-color: var(--bg-primary);
            color: var(--text);
            min-height: 100vh;
            padding: 2rem;
            background-image: 
                linear-gradient(to right, rgba(191, 255, 0, 0.03) 1px, transparent 1px),
                linear-gradient(to bottom, rgba(191, 255, 0, 0.03) 1px, transparent 1px);
            background-size: 30px 30px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            margin-bottom: 3rem;
            padding-bottom: 2rem;
            border-bottom: 2px solid var(--border);
        }

        h1 {
            font-size: 2.5rem;
            color: var(--fg-primary);
            font-weight: 700;
            margin-bottom: 0.5rem;
            text-shadow: 0 0 20px rgba(191, 255, 0, 0.3);
            letter-spacing: -0.02em;
        }

        .subtitle {
            color: var(--text-dim);
            font-size: 0.95rem;
            font-weight: 400;
        }

        .breadcrumb {
            color: var(--fg-muted);
            font-size: 0.9rem;
            margin-top: 0.5rem;
        }

        .file-list {
            background-color: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 8px;
            overflow: hidden;
            margin-bottom: 2rem;
        }

        .file-header {
            display: grid;
            grid-template-columns: 1fr auto auto;
            gap: 2rem;
            padding: 1rem 1.5rem;
            background-color: rgba(191, 255, 0, 0.05);
            border-bottom: 1px solid var(--border);
            font-weight: 600;
            color: var(--fg-secondary);
            font-size: 0.85rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        ul {
            list-style: none;
        }

        li {
            display: grid;
            grid-template-columns: 1fr auto auto;
            gap: 2rem;
            padding: 1rem 1.5rem;
            border-bottom: 1px solid var(--border);
            transition: all 0.2s ease;
            align-items: center;
        }

        li:last-child {
            border-bottom: none;
        }

        li:hover {
            background-color: var(--bg-hover);
            border-left: 3px solid var(--fg-primary);
            padding-left: calc(1.5rem - 3px);
        }

        .file-name {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            min-width: 0;
        }

        .file-icon {
            flex-shrink: 0;
            font-size: 1.2rem;
        }

        a {
            color: var(--link);
            text-decoration: none;
            font-weight: 500;
            transition: color 0.2s ease;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        a:hover {
            color: var(--link-hover);
            text-shadow: 0 0 10px rgba(191, 255, 0, 0.5);
        }

        .meta {
            font-size: 0.85rem;
            color: var(--text-dim);
            font-family: 'Courier New', monospace;
            white-space: nowrap;
        }

        .size {
            text-align: right;
            min-width: 80px;
        }

        .mtime {
            min-width: 140px;
        }

        footer {
            margin-top: 3rem;
            padding-top: 2rem;
            border-top: 1px solid var(--border);
        }

        .footer-title {
            color: var(--fg-secondary);
            font-size: 0.9rem;
            font-weight: 600;
            margin-bottom: 1rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .links {
            display: flex;
            flex-wrap: wrap;
            gap: 1.5rem;
        }

        .footer-link {
            color: var(--text-dim);
            text-decoration: none;
            transition: all 0.2s ease;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            border: 1px solid transparent;
            font-size: 0.9rem;
        }

        .footer-link:hover {
            color: var(--fg-primary);
            border-color: var(--border);
            background-color: rgba(191, 255, 0, 0.05);
        }

        .footer-link::before {
            content: '→ ';
            opacity: 0;
            transition: opacity 0.2s ease;
        }

        .footer-link:hover::before {
            opacity: 1;
        }

        @media (max-width: 768px) {
            body {
                padding: 1rem;
            }

            h1 {
                font-size: 1.8rem;
            }

            .file-header {
                display: none;
            }

            li {
                grid-template-columns: 1fr;
                gap: 0.5rem;
            }

            .meta {
                font-size: 0.75rem;
            }

            .links {
                flex-direction: column;
                gap: 0.5rem;
            }
        }

        @keyframes glow {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.6; }
        }

        .status-indicator {
            display: inline-block;
            width: 8px;
            height: 8px;
            background-color: var(--fg-primary);
            border-radius: 50%;
            margin-right: 0.5rem;
            animation: glow 2s ease-in-out infinite;
            box-shadow: 0 0 10px var(--fg-primary);
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1><span class="status-indicator"></span>Shaan OS Repositories</h1>
            <div class="subtitle">Shaan Vision Operating System</div>
            <div class="breadcrumb">📂 {{.Title}}</div>
        </header>

        <div class="file-list">
            <div class="file-header">
                <div>Name</div>
                <div>Size</div>
                <div>Modified</div>
            </div>
            <ul>
            {{range .Items}}
                <li>
                    <div class="file-name">
                        <span class="file-icon">{{if .IsDir}}📁{{else}}📄{{end}}</span>
                        <a href="{{.Href}}">{{.Name}}</a>
                    </div>
                    <span class="meta size">{{.Size}}</span>
                    <span class="meta mtime">{{.Mtime}}</span>
                </li>
            {{end}}
            </ul>
        </div>

        <footer>
            <div class="footer-title">Quick Links</div>
            <div class="links">
                <a href="https://os.shaanvision.com.tr/" class="footer-link" target="_blank" rel="noopener noreferrer">ShaanOS Home</a>
                <a href="https://pkgs.os.shaanvision.com.tr/" class="footer-link" target="_blank" rel="noopener noreferrer">Packages Browser</a>
                <a href="https://www.shaanvision.com.tr/linkler" class="footer-link" target="_blank" rel="noopener noreferrer">Shaan Vision Links</a>
            </div>
        </footer>
    </div>
</body>
</html>
`))

var notFoundTmpl = template.Must(template.New("404").Parse(`
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - Sayfa Bulunamadı | Shaan OS</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Source Code Pro', 'Consolas', monospace;
            background-color: #0a0a0a;
            color: #e0e0e0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            background-image: 
                linear-gradient(to right, rgba(191, 255, 0, 0.03) 1px, transparent 1px),
                linear-gradient(to bottom, rgba(191, 255, 0, 0.03) 1px, transparent 1px);
            background-size: 30px 30px;
        }

        .error-container {
            text-align: center;
            max-width: 600px;
            padding: 2rem;
        }

        .error-code {
            font-size: 8rem;
            font-weight: 700;
            color: #bfff00;
            text-shadow: 0 0 40px rgba(191, 255, 0, 0.5);
            line-height: 1;
            margin-bottom: 1rem;
        }

        h1 {
            font-size: 2rem;
            color: #bfff00;
            margin-bottom: 1rem;
        }

        p {
            color: #9ca3af;
            font-size: 1.1rem;
            margin-bottom: 2rem;
            line-height: 1.6;
        }

        .btn {
            display: inline-block;
            padding: 1rem 2rem;
            background-color: transparent;
            color: #bfff00;
            text-decoration: none;
            border: 2px solid #bfff00;
            border-radius: 4px;
            font-weight: 600;
            transition: all 0.3s ease;
            font-size: 1rem;
        }

        .btn:hover {
            background-color: #bfff00;
            color: #0a0a0a;
            box-shadow: 0 0 20px rgba(191, 255, 0, 0.5);
            transform: translateY(-2px);
        }

        .decoration {
            font-size: 3rem;
            margin-bottom: 1rem;
            opacity: 0.5;
        }
    </style>
</head>
<body>
    <div class="error-container">
        <div class="decoration">⚠️</div>
        <div class="error-code">404</div>
        <h1>Sayfa Bulunamadı</h1>
        <p>Aradığınız kaynak mevcut değil veya taşınmış olabilir.</p>
        <a href="/" class="btn">← Ana Sayfaya Dön</a>
    </div>
</body>
</html>
`))

func formatSize(size int64) string {
    if size < 1024 {
        return fmt.Sprintf("%d B", size)
    } else if size < 1024*1024 {
        return fmt.Sprintf("%.1f KB", float64(size)/1024)
    } else if size < 1024*1024*1024 {
        return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
    }
    return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
}

func generateIndex(path string) error {
    entries, err := os.ReadDir(path)
    if err != nil {
        return err
    }

    var files []FileEntry
    for _, entry := range entries {
        if entry.Name() == "index.html" {
            continue
        }

        info, err := entry.Info()
        if err != nil {
            continue
        }

        href := entry.Name()
        if entry.IsDir() {
            href += "/"
        }

        files = append(files, FileEntry{
            Name:  entry.Name(),
            Href:  href,
            Size:  "-",
            Mtime: info.ModTime().Format("2006-01-02 15:04"),
            IsDir: entry.IsDir(),
        })

        if !entry.IsDir() {
            files[len(files)-1].Size = formatSize(info.Size())
        }
    }

    // Dizinleri önce, sonra dosyaları alfabetik sırala
    sort.Slice(files, func(i, j int) bool {
        if files[i].IsDir != files[j].IsDir {
            return files[i].IsDir
        }
        return files[i].Name < files[j].Name
    })

    f, err := os.Create(filepath.Join(path, "index.html"))
    if err != nil {
        return err
    }
    defer f.Close()

    return tmpl.Execute(f, struct {
        Title string
        Items []FileEntry
    }{
        Title: filepath.Base(path),
        Items: files,
    })
}

func walk(base string) {
    filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
        if err != nil || !d.IsDir() {
            return nil
        }
        fmt.Println("Generating index for:", path)
        return generateIndex(path)
    })
}

func generate404(base string) error {
    f, err := os.Create(filepath.Join(base, "404.html"))
    if err != nil {
        return err
    }
    defer f.Close()
    return notFoundTmpl.Execute(f, nil)
}

func main() {
    walk("public")
    if err := generate404("public"); err != nil {
        fmt.Fprintf(os.Stderr, "Error generating 404 page: %v\n", err)
    }
    fmt.Println("✓ All indices generated successfully!")
}
