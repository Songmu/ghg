# ghg

githubから実行可能ファイルを取得するツール

    % ghg get Songmu/ghch
    % ghg get Songmu/ghch@v0.0.1
    % ghg get peco

## 配置ディレクトリ

- `/usr/local/bin`
- `~/.ghg/bin`

前者で行きたい気持ちもあるけど、コンフリクト時にめんどうなので後者が無難かも。(現状カレントディレクトリ)
設定ファイルで書き換え可能にする(?)

## 設定ファイル

`~/.ghg/config.yml`

- bindir

## upgradeとか

- 現状無条件で上書き
  - 上書きしないほうが良い？
- `-u` で無条件上書き
- バージョン比較？
  - ゆくゆく

## ダウンロード履歴

- 必要？
- path:
  - repo_url:
  - artifact_url:
  - version:

## サブコマンド

    % ghg bin
