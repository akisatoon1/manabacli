# manaba CLI

[`github.com/akisatoon1/manaba`](https://github.com/akisatoon1/manaba) ライブラリを使って、
コマンドラインから manabaのレポート提出ページにファイルをアップロードする CLI です。

## 使い方

```
manaba <ファイル名> <URL>
```

- `<ファイル名>`: manaba にアップロードするローカルファイルのパス
- `<URL>`: 提出先の manaba レポート提出ページ URL

例:

```
manaba report.pdf https://manaba.example.ac.jp/ct/course_xxxx_report_yyyy
```

成功すると `アップロードに成功しました: report.pdf` と表示されます。
（このツールはアップロードのみを行い、提出は実行しません。）

## 認証情報の設定

`~/.manabacli/config`というファイルを作成し、manaba のユーザー名・パスワードを保存します。

ファイルの作成
```
mkdir -p ~/.manabacli && touch ~/.manabacli/config && chmod 600 ~/.manabacli/config
```

ファイルの中身
```
username=あなたのユーザー名
password=あなたのパスワード
```

- 1行1項目の `key=value` 形式です。
- `#` で始まる行と空行は無視されます。
- パスワードは平文で保存されるため、`chmod 600` でファイル権限を制限することを推奨します。

## 学校の PC にインストールする方法

Go をインストールできない・管理者権限（sudo）が使えない学校の PC でも、
[Releases](https://github.com/akisatoon1/manabacli/releases) に置いてあるビルド済みバイナリを
ダウンロードして PATH を通すだけで利用できます。

1. **バイナリをダウンロードする。** ホームディレクトリ配下の `~/.local/bin`（無ければ作成）に
   最新リリースの `manaba` を保存します。

   ```
   mkdir -p ~/.local/bin
   curl -L -o ~/.local/bin/manaba https://github.com/akisatoon1/manabacli/releases/latest/download/manaba
   chmod +x ~/.local/bin/manaba
   ```

   `curl` が使えない場合は、ブラウザで [Releases](https://github.com/akisatoon1/manabacli/releases)
   を開き、最新リリースの `manaba` を `~/.local/bin/` に保存して `chmod +x ~/.local/bin/manaba` を実行してください。

2. **`~/.local/bin` に PATH を通す。** シェルの設定ファイル（bash なら `~/.bashrc`、zsh なら `~/.zshrc`）に
   次の1行を追記します。管理者権限は不要です。

   ```
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   ```

   追記したら設定を再読み込みします（または端末を開き直します）。

   ```
   source ~/.bashrc
   ```

3. **確認する。** どこからでも `manaba` コマンドが実行できれば成功です。

   ```
   manaba
   ```

   Usage が表示されればインストール完了です。あとは「[認証情報の設定](#認証情報の設定)」を済ませれば使えます。

> ダウンロードしたバイナリは Linux 向けです。後でアップデートするときは、同じ手順で `~/.local/bin/manaba` を上書きしてください。

## ビルド

Go 1.22 以上が必要です。

```
go build -o manaba .
```

生成された `manaba` バイナリを PATH の通った場所に置くと、どこからでも実行できます。

```
go install .          # $GOBIN もしくは ~/go/bin にインストール
# または
sudo cp manaba /usr/local/bin/
```

## エラーと終了コード

| 状況 | 終了コード | 出力先 |
| --- | --- | --- |
| 引数の数が不正 | 2 | 標準エラー（Usage を表示） |
| 設定ファイル欠落・不備 | 1 | 標準エラー |
| ファイルが見つからない／ディレクトリ指定 | 1 | 標準エラー |
| ログイン失敗（認証情報誤り等） | 1 | 標準エラー |
| アップロード失敗 | 1 | 標準エラー |
| 成功 | 0 | 標準出力 |
