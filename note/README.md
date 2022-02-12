# 概要

- gngc を作成するにあたって調査したことなどのメモを残す。

## アプリケーションがエラーで落ちた際、標準のメッセージを出力しないようにするオプションに関して

- [cobra のソースコードを読んで調査してみた](https://hakiwata.jp/post/20220213/)

## rootCmd.Run と rootCmd.RunE の違いに関して

- [cobra のソースコードを読んで調査してみた](https://hakiwata.jp/post/20220213/)

## CLI アプリケーションが実行されるまでの流れに関して

- [cobra のソースコードを読んで調査してみた](https://hakiwata.jp/post/20220213/)

## `root.go` の `func Execute()` 内の実装に関して

- cobra の公式ドキュメント ([Create rootCmd](https://cobra.dev/#create-rootcmd)) には `func Execute()` 内の実装は以下のようなサンプルが紹介されている。

```golang
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    os.Exit(1)
  }
}
```

- [func CheckErr(msg interface{})](https://github.com/spf13/cobra/blob/master/cobra.go#L211) を用いて [func (c *Command) Execute() error](https://github.com/spf13/cobra/blob/master/command.go#L901) 内のエラーハンドリングをすることも可能である。[func CheckErr(msg interface{})](https://github.com/spf13/cobra/blob/master/cobra.go#L211) の内部の実装は以下である。以下の実装からもわかるようにどちらの実装でも大筋のハンドリングは同じである。

```golang
// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface{}) {
    if msg != nil {
        fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
```

## GitHub Actions の toml ファイルの内容に関して

- GitHub Actions の workflow で以下のような書き方をしていると、Job 内で以下のようなエラーが生じた。

```yaml
...
    - name: Setup .gngc.toml
      run: |
        cat << EOF > $HOME/.gngc.toml
        [GitHub]
        Token=${{ secrets.PERSONAL_ACCESS_TOKEN }}
        UserName=${{ secrets.USER_NAME }}

        [IFTTT]
        EventName=${{ secrets.IFTTT_EVENT_NAME }}
        Token=${{ secrets.IFTTT_TOKEN }}
        EOF
...
```

- Job が吐き出すエラーメッセージは以下である。

```bash
Error: While parsing config: (2, 9): no value can start with g
exit status 1
Error: Process completed with exit code 1.
```

- これは、ファイルを読み出すライブラリの [viper](https://github.com/spf13/viper/) の [func ReadInConfig()](https://github.com/spf13/viper/blob/v1.10.1/viper.go#L1464) 内で呼び出されている [func (v *Viper) ReadInConfig()](https://github.com/spf13/viper/blob/v1.10.1/viper.go#L1464) で生じていることは、[cmd/root.go](https://github.com/dilmnqvovpnmlib/gngc/blob/main/cmd/root.go) 内のデバッグのログから確認できた。[func ReadInConfig()](https://github.com/spf13/viper/blob/v1.10.1/viper.go#L1464) 自体は、[cmd/root.go](https://github.com/dilmnqvovpnmlib/gngc/blob/main/cmd/root.go) 内の [func ReadInConfig()](https://pkg.go.dev/github.com/spf13/viper#Viper.ReadInConfig) で呼び出されている。この関数が `err` を返すため、Job がエラーを吐くとこまで調査はできた。また、Job 内のエラーメッセージは [func (pe ConfigParseError) Error() string](https://github.com/spf13/viper/blob/a785a79f2240b55faa3c9fb488252ca9ea931339/util.go#L30) が吐き出している。

- 手元では上手く行くが、GitHub Actions 上では上手くいかない問題だったため、にっちもさっちもいかなった。ヤケクソで `$HOME/.gngc.toml` に認証情報を書き出す処理を適当な文字列 (hoge) に書き換えると、エラーメッセージが変わった。

```bash
 Error::: While parsing config: (2, 7): no value can start with h
exit status 1
Error: Process completed with exit code 1.
```

- このエラーメッセージの変化から、`toml ファイル` 内の `""` で囲われていない文字列は、変数として扱われるのではないかと言うことに気づいた。そのため、`toml ファイル` に書き出す認証情報を以下のように `""` で書き込むと Job が上手く走った。

```yaml
...
    - name: Setup .gngc.toml
      run: |
        cat << EOF > $HOME/.gngc.toml
        [GitHub]
        Token="${{ secrets.PERSONAL_ACCESS_TOKEN }}"
        UserName="${{ secrets.USER_NAME }}"

        [IFTTT]
        EventName="${{ secrets.IFTTT_EVENT_NAME }}"
        Token="${{ secrets.IFTTT_TOKEN }}"
        EOF
...
```

- エラーメッセージがわからなかったため、ソースコードを追ってみたが、無事に問題を解決することができて良かった。

## GitHub Actions 上で Golang のプログラムを実行する際にパッケージをインストール刷る処理を事前に明示的に書く必要があるか

- 結論から言うと、必要ない。プログラムの実行 (`go run main.go`) 前に依存関係をインストールしてくれる。

- 以下の GitHub Actions の Job のログからもその様子が確認できる。

```bash
Run go run main.go -n 
go: downloading github.com/mitchellh/go-homedir v1.1.0
go: downloading github.com/spf13/cobra v1.3.0
go: downloading github.com/spf13/viper v1.10.0
go: downloading github.com/shurcooL/githubv4 v0.0.0-20220115235240-a14260e6f8a2
go: downloading github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a
go: downloading golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
go: downloading github.com/fsnotify/fsnotify v1.5.1
go: downloading github.com/magiconair/properties v1.8.5
go: downloading github.com/mitchellh/mapstructure v1.4.3
go: downloading github.com/spf13/afero v1.6.0
go: downloading github.com/spf13/cast v1.4.1
go: downloading github.com/spf13/jwalterweatherman v1.1.0
go: downloading github.com/spf13/pflag v1.0.5
go: downloading github.com/subosito/gotenv v1.2.0
go: downloading gopkg.in/ini.v1 v1.66.2
go: downloading golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
go: downloading golang.org/x/sys v0.0.0-20211205182925-97ca703d548d
go: downloading golang.org/x/text v0.3.7
go: downloading github.com/hashicorp/hcl v1.0.0
go: downloading github.com/pelletier/go-toml v1.9.4
go: downloading gopkg.in/yaml.v2 v2.4.0
200 OK from IFTTT API
Congratulations! You've fired the *** event
```

## Golang のスタックトレーサに関して

- エラーが生じた箇所の行番号などを出力する関数に [func Caller](https://pkg.go.dev/runtime#Caller) がある。これを用いることでデバッグがしやすくなる。

### 参考

- [Go言語 runtime.Callerを使ってメッセージやerrorにソースファイル名、行番号を含める](https://qiita.com/h6591/items/468be2f4524ccc888795)

## ダウンロードしたバイナリに対してシンボリックリンクを張る

```bash
sudo ln -s ${HOME}/go/1.17.6/bin/gngc /usr/local/bin/gngc
```

## Go のファイルのフォーマットに関して

- 以下のコマンドを実行すると、配下の全てのファイルに対して `go fmt` を実行することができる。

```bash
go fmt ./...
```

- ちなみに `gofmt -l -s -w .` でも同様のことが実現できるらしい。

### 参考

- [go fmt をプロジェクト配下の全ファイルに対して実行したい](https://devlights.hatenablog.com/entry/2019/08/15/060851)
