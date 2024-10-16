<div align="right">

![golangci-lint](https://github.com/yanosea/jrp/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/yanosea/jrp/actions/workflows/release.yml/badge.svg)

</div>

<div align="center">

# 🎲 jrp

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/yanosea/jrp?style=flat-square)](https://github.com/yanosea/jrp/releases/latest)
<br/>
[Coverage Report](https://yanosea.github.io/jrp/coverage.html)
<br/>
![demo](docs/demo.gif "demo")

</div>

## ℹ️ About

`jrp` is the CLI jokeey tool to generate Japanese random phrase(s).  
You can save the generated phrase(s) to the history and manage them.  
Also, you can save the generated phrase(s) to the favorite and manage them.

This tool uses [WordNet Japan sqlite database file](https://bond-lab.github.io/wnja/jpn/downloads.html).

## 💻 Usage

```
Usage:
  jrp [flags]
  jrp [command]

Available Subcommands:
  download,    dl,   d  📦 Download WordNet Japan sqlite database file from the official web site.
  generate,    gen,  g  ✨ Generate Japanese random phrase(s).
  history,     hist, h  📜 Manage the history of the "generate" command.
  favorite,    fav,  f  ⭐ Manage the favorited phrase(s) of the history of "generate" command.
  interactive, int,  i  💬 Generate Japanese random phrase(s) interactively.
  help                  🤝 Help for jrp.
  completion            🔧 Generate the autocompletion script for the specified shell.
  version               🔖 Show the version of jrp.

Flags:
  -n, --number       🔢 number of phrases to generate (default 1, e.g: 10)
  -p  --prefix       💬 prefix of phrase(s) to generate
  -s  --suffix       💬 suffix of phrase(s) to generate
  -d  --dry-run      🧪 generate phrase(s) without saving to the history
  -P, --plain        📝 plain text output instead of table output
  -i, --interactive  💬 generate Japanese random phrase(s) interactively
  -t, --timeout      ⏱️  timeout in seconds for the interactive mode (default 30, e.g: 10)
  -h, --help         🤝 help for jrp
  -v, --version      🔖 version for jrp

Arguments:
  number  🔢 number of phrases to generate (e.g: 10)
```

## 💬 Interactive mode

![demo_interactive](docs/demo_interactive.gif "demo_interactive")

`jrp` can generate Japanese random phrase(s) interactively.  
You can favorite, save, skip, and exit interactively while generating phrase(s).

To use this mode, run either command below.

```sh
# Those commands below are equivalent.
# And they have their aliases. Please check the help message.
jrp interactive
# or
jrp --interactive
# or
jrp generate interactive
# or
jrp generate --interactive
```

Press either key below for your action.

- `u`
  - Favorite, continue.
- `i`
  - Favorite, exit.
- `j`
  - Save, continue.
- `k`
  - Save, exit.
- `m`
  - Skip, continue.
- `other`
  - Skip, exit.

## 🌍 Environments

### 📁 Directory to store WordNet Japan sqlite database file

Default : `$XDG_DATA_HOME/jrp` or `$HOME/.local/share/jrp`

```sh
export JRP_WNJPN_DB_FILE_DIR=/path/to/your/directory
```

### 📁 Directory to store jrp sqlite database file

Default : `$XDG_DATA_HOME/jrp` or `$HOME/.local/share/jrp`

```sh
export JRP_DB_FILE_DIR=/path/to/your/directory
```

## 🔧 Installation

### 🐭 Using go

```sh
go install github.com/yanosea/jrp@latest
```

### 🍺 Using homebrew

```sh
brew tap yanosea/tap
brew install yanosea/tap/jrp
```

### 📦 Download from release

Go to the [Releases](https://github.com/yanosea/jrp/releases) and download the latest binary for your platform.

## ✨ Update

### 🐭 Using go

reinstall `jrp`!

```sh
go install github.com/yanosea/jrp@latest
```

### 🍺 Using homebrew

```sh
brew update
brew upgrade jrp
```

### 📦 Download from release

Download the latest binary from the [Releases](https://github.com/yanosea/jrp/releases) page and replace the old binary in your `$PATH`.

## 🧹 Uninstallation

### 🔧 Uninstall jrp

#### 🐭 Using go

```sh
rm $GOPATH/bin/jrp
sudo rm -fr $GOPATH/pkg/mod/github.com/yanosea/jrp@*
```

#### 🍺 Using homebrew

```sh
brew uninstall jrp
brew untap yanosea/tap/jrp
```

#### 📦 Download from release

Remove the binary you downloaded and placed in your `$PATH`.

### 🗑️ Remove data files

If you've set jrp envs, please replace `$HOME/.local/share/jrp` with envs you've set.  
These below commands are in the case of default. Ofcourse you can remove whole the directory.

#### 💾 Remove WordNet Japan sqlite database file

```sh
rm $HOME/.local/share/jrp/wnjpn.db
```

#### 💾 Remove jrp sqlite database file

```sh
rm $HOME/.local/share/jrp/jrp.db
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
