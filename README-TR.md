# 🧩 ShaanOS Dizin Oluşturucu

> 🇬🇧 [For English version, click here](README.md)

**ShaanOS Index Generator**, [Shaan OS](https://dl-os.shvn.tr) için geliştirilmiş hafif bir statik dizin oluşturucudur.  
GitHub Actions tarafından oluşturulan `public/core/` dizinini tarar ve her klasör için modern, neon tarzı `index.html` dosyaları üretir.

---

## 🚀 Özellikler

- 📂 **core/x86_64** ve **core/x86** dizinleri için index oluşturur  
- 🧩 **Go** diliyle yazılmış, bağımlılıksız  
- 💚 ShaanOS tarzı neon yeşil tema  
- ⚙️ Dosya boyutu, değişim zamanı, klasör simgeleri  
- 🪄 Otomatik `404.html` oluşturma  
- 🌐 Cloudflare Pages / GitHub Pages ile tamamen uyumlu

---

## 🏗️ Kullanım Senaryosu

Bu araç, ShaanOS’un GitHub Actions derleme sürecinde kullanılır:

> `.github/workflows/build.yml` dosyası `packages/**` değişikliklerinde çalışır.  
> Derlenen `.apk` dosyaları `public/core/{arch}/` dizinine kopyalanır.  
> Ardından `shaanos-index` bu klasörler için index sayfaları oluşturur.

---

## 📁 Klasör Yapısı

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

Çalıştırıldığında:

```bash
public/core/index.html
public/core/x86_64/index.html
public/core/x86/index.html
public/404.html
```

oluşturulur ✅

---

## ⚙️ Manuel Çalıştırma

```bash
go run main.go
# veya
go build -o shaanos-index
./shaanos-index
```

Varsayılan olarak `public/` dizinini tarar ve her klasör için bir `index.html` oluşturur.

---

## 🧱 GitHub Actions Örneği

```yaml
- name: Dizin indexlerini oluştur
  run: |
    go run main.go
```

Bu adım, paketler derlendikten sonra çalışır.

---

## 🌈 Tema

- Neon yeşili (`#bfff00`) vurgu rengi  
- Siyah arka plan (`#0a0a0a` / `#1a1a1a`)  
- “Source Code Pro” yazı tipi  
- Mobil uyumlu, responsive tasarım  

Canlı örnek: [https://dl-os.shvn.tr/core/x86_64](https://dl-os.shvn.tr/core/x86_64)

---

## 💡 Neden sadece “core”?

Çünkü mevcut derleme süreci `.apk` dosyalarını yalnızca `public/core/$ARCH/` dizinine kopyalıyor:

```bash
cp "$apk" /workspace/public/core/$TARGET_ARCH/
```

---

## 🧠 Teknik Bilgiler

| Alan | Değer |
|------|--------|
| Dil | Go |
| Şablon Sistemi | `html/template` |
| Ana Fonksiyonlar | `walk()`, `generateIndex()`, `generate404()` |
| Bağımlılıklar | Yok |
| Çıktı | Statik HTML |

---

## 📜 Lisans

MIT Lisansı © 2025 [Shaan Vision](https://www.shaanvision.com.tr)
