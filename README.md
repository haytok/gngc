# gngc (Get and Notify GitHub Contributions)

A simple command line application to get GitHub contributions and notify them to IFTTT.

# How to setup

- Create .env file

```bash
cat << EOF > hoge
USER_NAME=USER_NAME
API_TOKEN=API_TOKEN
IFTTT_TOKEN=IFTTT_TOKEN
IFTTT_EVENT_NAME=IFTTT_EVENT_NAME
EOF
```

- Build main.go

```bash
go build main.go
```

- Set binary

```bash
sudo ln -s ${PWD}/gngc/main /usr/local/bin/gngc
```

# How to use

```bash
> gngc -h
Get GitHub contributions.

Usage:
  gngc [flags]

Flags:
  -h, --help     help for gngc
  -n, --notify   Get GitHub contributions and notify them to IFTTT.
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
