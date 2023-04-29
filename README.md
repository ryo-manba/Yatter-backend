# Yatter-backend

Twitter・Mastodonに似た仮想SNSサービスYatterのバックエンドAPIをGo言語で実装しています。

## 主な機能
- アカウント作成
- ステータス（投稿）
- タイムライン

## アーキテクチャ
主に以下の要素で構成されています。

- app: 依存性の注入 (DI) コンテナが扱われているパッケージです。
- config: サーバーの設定がまとめられているパッケージです。
- domain: ドメイン層で、コアビジネスロジックが含まれています。
- handler: インターフェース層およびアプリケーション層で、HTTPリクエストハンドラが含まれています。
- dao: インフラストラクチャ層で、ドメイン/リポジトリの実装が含まれています。
- ddl: データベース定義マスタが含まれています。

## 使用ライブラリ
- HTTP: chi
- DB: sqlx

## 開発環境
- Go
- Docker / docker-compose

## 開発環境のセットアップ

**docker-composeを使って開発環境を立ち上げる**
```bash
docker-compose up -d
```

**開発環境をシャットダウンする**
```bash
docker-compose down
```

**Swagger UIでAPI仕様を確認**
開発環境を立ち上げ、Webブラウザで localhost:8081 にアクセス

