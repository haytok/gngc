# gngc (Get and Notify GitHub Contributions)

A simple command line application to get GitHub contributions and notify them to [IFTTT](https://ifttt.com/).

# How to setup

- Create `$HOME/.gngc.toml`. Or, you can also optionally specify a configuration file. (such as gngc --config config/config.toml)

```bash
cat << EOF > $HOME/.gngc.toml
[GitHub]
UserName = "UserName"
Token = "Token"

[IFTTT]
EventName = "EventName"
Token = "Token"
EOF
```

- Install

```bash
go install https://github.com/dilmnqvovpnmlib/gngc@latest
```

# How to use

```bash
> gngc -h
A simple command line application to get GitHub contributions and notify them to IFTTT.

Usage:
  gngc [flags]

Flags:
      --config string   config file (default is $HOME/.gngc.toml)
  -h, --help            help for gngc
  -n, --notify          Get GitHub contributions and notify them to IFTTT.
```

```bash
> gngc
2022 年 02 月 07 日のコミット数は 5 です！
```

```bash
> gngc -n
200 OK
Congratulations! You've fired the tools event
```
