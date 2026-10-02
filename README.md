# Example01

Claude Code を使った開発フローを練習するための ToDo サンプルアプリ。

- `backend/`: Go（標準 `net/http`）+ SQLite の REST API
- `frontend/`: React + Vite + TypeScript

## 起動方法

```sh
# ターミナル1: API サーバー（http://localhost:8080）
cd backend
go run ./cmd/server

# ターミナル2: フロントエンド（http://localhost:5173）
cd frontend
npm ci
npm run dev
```

開発サーバーは `/api` へのリクエストを `localhost:8080` に転送します。

## テスト

```sh
cd backend && go vet ./... && go test ./...
cd frontend && npm run lint && npm test -- --run && npm run build
```

API の仕様と開発の決まりごとは [CLAUDE.md](CLAUDE.md) を参照してください。
