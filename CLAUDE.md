# CLAUDE.md

このファイルは、このリポジトリで作業する Claude Code 向けのガイドです。

## 概要
ToDo 管理のサンプルアプリ。Claude Code を使った開発フローの練習用。

- `backend/`: Go（標準 `net/http`）+ SQLite（`modernc.org/sqlite`、cgo 不要）による REST API
- `frontend/`: React + Vite + TypeScript の SPA。開発時は Vite の proxy で `/api` をバックエンドに転送する

## コマンド

### Backend（`backend/` で実行）
- 起動: `go run ./cmd/server`（環境変数 `PORT` 既定 8080、`DB_PATH` 既定 `todo.db`）
- テスト: `go test ./...`
- 静的解析: `go vet ./...`
- フォーマット: `gofmt -w .`

### Frontend（`frontend/` で実行）
- 依存インストール: `npm ci`
- 開発サーバー: `npm run dev`（http://localhost:5173）
- テスト: `npm test -- --run`
- Lint: `npm run lint`
- ビルド: `npm run build`

## アーキテクチャ
- `backend/internal/todo`: `Todo` 型、バリデーション、`Store` インターフェースと SQLite 実装
- `backend/internal/httpapi`: HTTP ハンドラ。依存するのは `todo.Store` インターフェースだけ
- `backend/cmd/server`: 依存の組み立てとサーバー起動
- `frontend/src/api.ts`: 型付きの API クライアント。コンポーネントからの API 呼び出しはすべてここを通す

## API
| Method | Path | 説明 |
|---|---|---|
| GET | /api/todos | 一覧（作成日時の降順） |
| POST | /api/todos | 作成 `{title}` |
| GET | /api/todos/{id} | 1件取得 |
| PATCH | /api/todos/{id} | 部分更新 `{title?, completed?}` |
| DELETE | /api/todos/{id} | 削除 |
| GET | /healthz | ヘルスチェック |

エラーレスポンスは `{"error": "message"}`。タイトルは前後の空白を除いて 1〜200 文字。

## 規約
- 変更したら、コミット前にバックエンドとフロントエンド両方のテストと lint を通す
- バックエンドの機能を追加したら `store_test.go` と `handler_test.go` にテストを足す
