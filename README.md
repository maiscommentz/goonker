# Goonker - Tic Tac Toe Reinvented
A modern implementation of the classic Tic Tac Toe game with multiplayer support, hard AI opponent and a unique nested quiz mechanic, built using Ebiten and Go, and compiled to WebAssembly for web deployment. This project was developed as part of the "Elegant programming in Go" course miniproject at HEIA-FR, Switzerland.

🎮 Access the game online: [https://maiscommentz.github.io/goonker/](https://maiscommentz.github.io/goonker/)


## 🛠️ Technologies
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![WebAssembly](https://img.shields.io/badge/-WebAssembly-654FF0?logo=webassembly&logoColor=white&style=for-the-badge)
![Docker](https://img.shields.io/badge/-Docker-2496ED?logo=docker&logoColor=white&style=for-the-badge)


## ✨ Key Features
1. **Classic Gameplay** - Full implementation of Tic Tac Toe with standard rules built using Ebiten library.

2. **Game Modes** - Play against an AI bot or another player in multiplayer mode.

3. **Nested Quiz** - A player can conquer a cell by answering a question correctly.

4. **Web Interface** - Interactive web application compiled to WebAssembly for optimal performance.

5. **Client-Server Architecture** - Strict separation of client and server logic with a robust Go backend featuring a hub system for managing games and players, all communicated via structured packets.

6. **Containerized Deployment** - Docker for simple and reproducible deployment.

## 📁 Project Structure
```
.
├── client/              # Client application
│   ├── assets/
│   ├── ui/
│   ├── game.go
│   ├── main.go
│   └── network.go
├── server/              # Backend server
│   ├── hub/             # Game management
│   ├── logic/           # Game logic
│   └── main.go
├── common/              # Shared client-server code
│   └── packets.go
├── web/                 # Web resources
│   ├── demo.wasm        # (generated at build time)
│   ├── index.html
│   └── wasm_exec.js
└── Dockerfile           # Docker configuration
```


## 🚀 Installation and Local Setup

### Prerequisites
- Go 1.24 or higher
- Docker (optional)

### Build and Run
Launch the server and client:
```bash
go run ./server
go run ./client
```

Compile the client to WebAssembly:
```bash
GOOS=js GOARCH=wasm go build -o ./web/demo.wasm ./client
```

### Using Docker
```bash
docker build -t goonker .
docker run -p 8080:8080 goonker
```


## 🔄 Continuous Integration and Deployment (CI/CD)
The project uses automated workflows to ensure code quality and facilitate deployment. Here are the key elements:

- **Linting and Testing**: On every push or pull request, the codebase is automatically linted and tested to maintain high standards.

- **Client Deployment**: The client is compiled to WebAssembly on each push to the main branch, ensuring the latest version is always ready for deployment. It is then deployed to GitHub Pages, making the game accessible online.

- **Server Deployment**: A docker image of the backend server is built and pushed to the repository registry. It is then deployed to a Self-Hosted Server using SSH. This ensures the server is always running the latest version.

This process ensures fast, reliable updates without manual intervention.

## ❤️ Contributors
- Filipe Casimiro Ferreira - [@maiscommentz](https://github.com/maiscommentz)
- Louis Pasquier - [@louis-pasquier](https://github.com/louis-pasquier)