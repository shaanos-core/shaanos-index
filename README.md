# 🧩 ShaanOS Index Generator

> 🇹🇷 [Türkçe sürüm için tıklayın](README-TR.md)

**ShaanOS Index Generator** is a lightweight static directory index generator for [Shaan OS](https://dl-os.shvn.tr).  
It scans the `public/core/` directory created by GitHub Actions and generates modern, neon-style `index.html` pages for each folder.

---

## 🚀 Features

- 📂 Generates index pages for **core/x86_64** and **core/x86**
- 🧩 Written in **Go**, no dependencies
- 💚 ShaanOS-style neon green theme
- ⚙️ Shows file size, modification time, and folder icons
- 🪄 Auto-generated `404.html`
- 🌐 Fully static, compatible with Cloudflare Pages / GitHub Pages

---

## 🏗️ Usage Scenario

This tool is used inside the ShaanOS GitHub Actions build pipeline:

> `.github/workflows/build.yml` runs on each `packages/**` change.  
> Built `.apk` packages are copied into `public/core/{arch}/`.  
> Then `shaanos-index` generates index pages for these folders.

---

## 📁 Directory Structure

```bash
public/
└── core/
    ├── x86_64/
    │   ├── APKINDEX.tar.gz
    │   ├── shaan-base-1.0-r1.apk
    │   └── ...
    └── x86/
        ├── APKINDEX.tar.gz
        └── ...
```

The generated output:

```bash
public/core/index.html
public/core/x86_64/index.html
public/core/x86/index.html
public/404.html
```

---

## ⚙️ Manual Run

```bash
go run main.go
# or
go build -o shaanos-index
./shaanos-index
```

By default, it scans `public/` and creates an `index.html` in every directory.

---

## 🧱 GitHub Actions Example

```yaml
- name: Generate directory indexes
  run: |
    go run main.go
```

This step is executed right after package builds are done.

---

## 🌈 Theme

- Neon green (`#bfff00`) highlights  
- Black background (`#0a0a0a` / `#1a1a1a`)  
- "Source Code Pro" font  
- Fully responsive design  

Live example: [https://dl-os.shvn.tr/core/x86_64](https://dl-os.shvn.tr/core/x86_64)

---

## 💡 Why only “core”?

Because the build workflow currently copies `.apk` files only to `public/core/$ARCH/`:

```bash
cp "$apk" /workspace/public/core/$TARGET_ARCH/
```

---

## 🧠 Technical Info

| Field | Value |
|--------|--------|
| Language | Go |
| Template Engine | `html/template` |
| Main Functions | `walk()`, `generateIndex()`, `generate404()` |
| Dependencies | None |
| Output | Static HTML |

---

## 📜 License

MIT License © 2025 [Shaan Vision](https://www.shaanvision.com.tr)
