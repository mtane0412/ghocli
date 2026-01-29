# ghoプロジェクト固有設定

## プロジェクト概要

Ghost Admin API CLIツール（Go言語実装）

## 開発言語

- Go

## 品質チェックコマンド

### Lint
```bash
make lint
```

### 型チェック
```bash
make type-check
```

### テスト
```bash
make test
```

### カバレッジ
```bash
make test-coverage
```

### ビルド
```bash
make build
```

## 開発時の注意事項

- `make lint` を実行するには golangci-lint のインストールが必要
- ビルドされたバイナリは `./gho` として出力される
- カバレッジレポートは `coverage.html` に出力される
