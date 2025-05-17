# chikokulympic-api

## 概要
「遅刻リンピック (Chikoku Olympic)」のバックエンド API サービス
BaseURL：https://chikokulympic-api-143996483854.us-central1.run.app

## アーキテクチャ
このプロジェクトはクリーンアーキテクチャを踏襲している
```
.
├── cmd            // main.go 
├── config         // 環境変数、設定ファイル
├── domain         // ドメイン層（エンティティ，リポジトリインタフェース）
├── usecase        // ユースケース層
├── infrastructure // インフラ層（MongoDB 接続，外部サービス実装）
│   └── mongo
├── middleware     // ミドルウェア層
├── presentation   // プレゼンテーション層（HTTP ハンドラ／コントローラ）
│   └── v1         // バージョンごとにディレクトリ分割
└── docs           // ドキュメント
```

### 技術スタック
- **言語**：Go  
- **データベース**：MongoDB  
- **コンテナ**：Docker / Docker Compose  
- **ドキュメント**：OpenAPI (Swagger) (`docs/openapi.yaml`)

