# ついったーのクローン

ある試験課題として作成したアプリです。役目を終えたためここに供養します。

## アプリの概要

### 機能一覧

1. ログイン

1. タイムライン参照

1. ユーザ検索とフォロー

### 構成

- Next.jsでSSGした静的ファイルをgolangにホスティングさせています

- APIサーバもgolangで実装しており、SPAがAPIを叩くことで通信しています

- データベースはPostgreSQLです

### 利用ライブラリ等

- サーバサイド：gorilla/mux, form3tech-oss/jwt-go
    - DBアクセスには標準ライブラリを利用

- クライアントサイド：Next.js, TypeScript

## 構築したアプリのディレクトリ構成

```
/workspace/
  |-- .devcontainer/ # Reopen in Contaier するための設定
  |-- .vscode/       # F5でGoをデバッグする設定と、ESLintの設定
  |
  |-- initdb/        # DB初期化SQL
  |
  |-- client/        # Next.jsによるクライアントのソース
  |    |-- out/          # SSGされた静的ファイルが出力されるディレクトリ
  |    |-- pages/
  |        |-- index.tsx  # ログイン画面
  |        |-- home.tsx   # 主要機能すべて（タイムライン、ツイート、ユーザ検索、フォロー）
  |        |-- signup.tsx # ユーザ新規登録画面（検証の便利のため）
  |
  |-- server/        # Goによるサーバのソース
  |    |-- controller/    # APIのパスの定義、パラメータ取得とチェック、
  |    |                  # modelおよびrepositoryの呼び出し、レスポンス返却
  |    |
  |    |-- model/         # モデルの定義と業務ロジック
  |    |
  |    |-- repository/    # DBアクセス
  |    |
  |    |-- util/          # 多数のソースから呼ばれるロジック
  |    |
  |    |-- server.go      # サーバのエントリーポイント
  |
  |-- docker-compose.yml  # デプロイ用composeファイル
  |-- app.Dockerfile      # デプロイ用アプリ用Dockerfile
  |-- db.Dockerfile       # デプロイ用DB用Dockerfile
```

## 開発環境立ち上げ方法

※当開発環境はVSCodeのdevcontainer機能の利用を前提としています。

1. 事前準備

    1. VSCodeの`Remote Development`拡張をインストールする

    1. ssh-agentにGitHubの鍵を登録しておく（HTTPSでも構わない）

1. ローカルで実施する手順

    1. ローカルに当リポジトリをクローンする

    1. VSCodeでクローンしたディレクトリを開く

    1. `.env`を用意する

        1 `cp .devcontainer/.env.sample .devcontainer/.env`

        1 `.env`を各自編集する

    1. 左下の`><`をクリックする

    1. 選択肢から`Reopen in Container`を開く

1. `Reopen in Container`後の手順

    1. コンテナの`/workspace`ディレクトリに当リポジトリをクローンする
        - `git clone git@github.com:shu000/too-simple-twitter.git /workspace`
        - 以下の「ビルド・実行方法」に従うとアプリが起動する

## ビルド・実行方法

以下により、`localhost:5000`でHTMLホスティングおよびAPIサーバが起動します

```
# クライアントをビルド
cd /workspace/client
yarn
yarn export

# サーバーを起動
cd /workspace
go run ./server/server.go
# (上記に代わりF5でも動作・デバッグ可能です)
```

## データベースについて

#### データの参照方法

開発環境ではAdminerのコンテナが立ち上がります。`localhost:8080`にアクセスすると`.env`で設定した認証情報でログインできます。

#### テーブル初期化方法

Adminerで`initdb/init.sql`を実行してください

#### その他

検証用のユーザ作成ページを用意しています。`localhost:5000/app/signup`にアクセスすると任意のユーザが作成できます。

## デプロイ方法

Dockerおよびdocker-composeがインストールされたサーバにて、以下を実行してください。

```
# 当該ソースを取得
git clone https://github.com/shu000/too-simple-twitter
cd too-simple-twitter

# 環境変数を設定
# 開発環境と同様、任意の値で書き換えてください
cp .env.sample .env

# コンテナ立ち上げ
docker-compose -f docker-compose.yml up --build
```

## TODO

1. クライアントのホットリロードサーバ（:3000）からAPIサーバ（:5000）へのプロキシ
