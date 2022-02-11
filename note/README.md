# 概要

- gngc を作成するにあたって調査したことなどのメモを残す。

## アプリケーションがエラーで落ちた際、標準のメッセージを出力しないようにするオプションに関して

```golang
var rootCmd = &cobra.Command{
...
    SilenceErrors: true,
	SilenceUsage:  true,
...
}
```

## `root.go` の `Execute 関数` 内の実装に関して

- cobra の公式ドキュメント ([Create rootCmd](https://cobra.dev/#create-rootcmd)) には `Execute 関数` 内の実装は以下のようなサンプルが紹介されている。

```golang
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    os.Exit(1)
  }
}
```

- [CheckErr 関数](https://github.com/spf13/cobra/blob/master/cobra.go#L211) を用いて [*Command の Execute 関数](https://github.com/spf13/cobra/blob/master/command.go#L901) 内のエラーハンドリングをすることも可能である。[CheckErr 関数](https://github.com/spf13/cobra/blob/master/cobra.go#L211) の内部の実装は以下である。以下の実装からもわかるようにどちらの実装でも大筋のハンドリングは同じである。

```golang
// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface{}) {
    if msg != nil {
        fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
```

## rootCmd.Run と rootCmd.RunE の違いに関して

## CLI アプリケーションが実行されるまでの流れに関して

## GitHub Actions の toml ファイルの内容に関して

## Golang のスタックトレーサに関して

## ダウンロードしたバイナリに対してシンボリックリンクを張る

```bash
sudo ln -s ${HOME}/go/1.17.6/bin/gngc /usr/local/bin/gngc
```
