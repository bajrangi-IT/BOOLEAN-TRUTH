# 🛡️ BOOLEAN TRUTH | Forensic Reality Auditor

TrustGuard is a high-performance backend API built with Go (Golang) that leverages Google Gemini AI to evaluate the safety, reliability, and truthfulness of text statements.

## ✨ Features
- **🌐 Real-time Fact-Checking**: Verify statements against global knowledge.
- **🚫 Safety Guard**: Detects Toxicity and Spam instantly.
- **📊 Detailed Audits**: Provides reasoning, detailed explanations, and source attribution for every result.
- **⚡ High Performance**: Fast and lightweight Gin-based REST API.
- **🎨 Premium UI**: Stunning, glassmorphism-inspired dark mode interface for auditing claims.

## 🏗️ Project Structure
```text
TrustGuard/
├── backend/
│   ├── main.go             # Server entry point
│   ├── controllers/        # API Request Handlers
│   ├── services/           # Gemini AI Integration
│   ├── .env                # API Keys & Config
│   └── go.mod              # Dependencies
├── frontend/
│   ├── index.html          # Main UI
│   ├── style.css           # Premium Aesthetics
│   └── app.js              # Frontend Logic
└── README.md
```

## 🚀 Getting Started

### 1. Prerequisites
- [Go](https://go.dev/dl/) (v1.20+)
- Gemini API Key from [Google AI Studio](https://aistudio.google.com/)

### 2. Setup Environment
1. Navigate to the `backend/` folder.
2. Open `.env` and paste your Gemini API Key:
   ```env
   GEMINI_API_KEY=your_api_key_here
   PORT=8080
   ```

### 3. Run the Backend
```bash
cd backend
go run main.go
```

### 4. Open the Frontend
Simply open `frontend/index.html` in any modern browser.

---
Built with ⚡ by Antigravity
