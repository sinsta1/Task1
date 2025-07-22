# 🛒 BlinkIt Web Scraper (GoLang)

This is a web scraping tool written in Go that fetches product data from BlinkIt's subcategory APIs (example: "Snacks & Munchies > Nachos").

## ✅ Features

- Scrapes product data from BlinkIt's public API
- Works with category and subcategory IDs
- Supports setting latitude and longitude
- Saves output as a CSV: `blinkit_products.csv`

## 📂 Files

- `main.go` — The Go scraper script
- `blinkit_products.csv` — Clean CSV output of product data

## 🧮 Data Fields Extracted

- Product ID
- Product Name
- Image URL

## 📍 How It Works

- Makes a `POST` request to:
