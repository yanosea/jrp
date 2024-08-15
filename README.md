<div align="right">

![golangci-lint](https://github.com/yanosea/jrp/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/yanosea/jrp/actions/workflows/release.yml/badge.svg)

</div>

<div align="center">

# 🎲 jrp

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/yanosea/jrp?style=flat-square)](https://github.com/yanosea/jrp/releases/latest)

</div>

## ℹ️ About

`jrp` is the CLI tool to generate random Japanese phrase(s). (It's jokeey tool!)  
This tool uses [WordNet Japan](https://bond-lab.github.io/wnja/jpn/downloads.html) sqlite database file.

## 💻 Usage

```
Usage:
  jrp [flags]
  jrp [command]

Available Subcommands:
  download    📥 Download Japanese Wordnet sqlite3 database file from the official site.
  generate    ✨ Generate Japanese random phrase(s).
  help        🤝 Help of jrp.
  completion  🔧 Generate the autocompletion script for the specified shell.
  version     🔖 Show the version of jrp.

Flags:
  -n, --number    🔢 number of phrases to generate (default 1). You can abbreviate "generate" sub command.
  -h, --help      🤝 help for jrp
  -v, --version   🔖 version for jrp

Arguments:
  number  🔢 number of phrases to generate (e.g: 10).
```

## 🔧 Installation

### 🍺 Using homebrew

```
brew tap yanosea/tap
brew install yanosea/tap/jrp
```

## ✨ Update

```
brew update
brew upgrade jrp
```

## 📃 License

[🔓MIT](./LICENSE)

## 🖊️ Author

[🏹 yanosea](https://github.com/yanosea)

## 🔥 Motivation

I love the smart phone application [PhrasePlus!](https://www.phraseplus.org)  
I wanted to run an application with equivalent functionality to this in the terminal, so I created it!

## 🤝 Contributing

Feel free to point me in the right direction🙏
