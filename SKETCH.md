# ghg

githubから実行可能ファイルを取得するツール

    % ghg get Songmu/ghch
    % ghg get Songmu/ghch@v0.0.1
    % ghg get peco

## 配置ディレクトリ

- `~/.ghg/bin`

`~/.ghg` の場所は環境変数 `$GHG_HOME` で上書き可能

## 設定ファイル(必要?)

`~/.ghg/config.yml`

- bindir

## upgradeとか

- `-u` で無条件上書きになっている
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
