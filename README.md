# 🔐 uses - USE Secure environment variables in dev

[![Release Action](https://github.com/malikashish8/uses/actions/workflows/release.yaml/badge.svg)](https://github.com/malikashish8/uses/actions/workflows/release.yaml)
[![CodeQL Scan](https://github.com/malikashish8/uses/actions/workflows/codeql.yaml/badge.svg?branch=master)](https://github.com/malikashish8/uses/actions/workflows/codeql.yaml)
[![Semgrep Scan](https://github.com/malikashish8/uses/actions/workflows/semgrep.yaml/badge.svg?branch=master)](https://github.com/malikashish8/uses/actions/workflows/semgrep.yaml)

Taking inspiration from [aws-vault](https://github.com/99designs/aws-vaults), `uses` makes `use` of OS provided `s`ecret management solutions to save secrets in the development environment. Grouping of secrets is made possible by a config file.

Having secrets lying around in environment variables in the development environment can be a nightmare as opensource packages are being [actively compromised](https://thehackernews.com/2022/05/pypi-package-ctx-and-php-library-phpass.html) to steal secrets. These packages can read all the environment variables. Good security hygine dictates that no secrets are stored in environment variables (using configs such as ~/.bashrc and ~/.zshrc). `uses` helps to implement the least privilege principle by saving all the secrets in a secret store and explicitly allowing processes to access the required secrets.

## ⚡️ Installation

Install using [Homebrew](https://brew.sh/)

```bash
brew install malikashish8/tap/uses
```

Or download the binary from [releases](https://github.com/malikashish8/uses/releases) and add it to your path.

## 🧑‍💻 Usage

```text
❯ uses                   
NAME:
   uses - securely manage secrets in dev environment

USAGE:
   uses [global options] command [command options] [arguments...]

VERSION:
   v0.0.9

COMMANDS:
   project1    get secrets for project `project1` and run command
   webgoat     get secrets for project `webgoat` and run command
   config, c   open config file
   set, s      set a secret `name=value`
   get, g      get secret for a `name`
   list, l     list all secrets saved using `uses`
   remove, r   delete a `secret`
   completion  generate auto-complete commands for a shell
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Set

Set a secret in secret store

```bash
❯ uses set GITHUB_TOKEN
Overwrite value? (y/n) y
Enter value: 
INFO[0009] GITHUB_TOKEN saved
```

Set in one line (unsafe) 

```bash
❯ uses set GITHUB_TOKEN=sdknbowhlfownpns;s/dkfnbslsnwwn
```

Above will save the secret in prompt history. A safer way is to use `pbpaste` on Mac:

```bash
❯ uses set $(pbpaste)
```

or if the secret is already in the environment

```bash
❯ uses set GITHUB_TOKEN=${GITHUB_TOKEN}
```

### Get

Get a secret from secret store

```bash
❯ uses get GITHUB_TOKEN
sdknbowhlfownpns;s/dkfnbslsnwwn
```

Get can also run any command passed to it after setting the environment variable

```bash
❯ uses get GITHUB_TOKEN env | grep GITHUB_TOKEN
GITHUB_TOKEN=sdknbowhlfownpns;s/dkfnbslsnwwn
```

### Enable Auto-completion

#### ZSH

```bash
echo 'source <(uses completion zsh)' >>~/.zshrc
```

#### BASH

```bash
echo 'source <(uses completion bash)' >>~/.bashrc
```

### List

Get a list of secrets managed by `uses`

```bash
❯ uses list
[AWS_USER, GITHUB_TOKEN]
```

### Projects

Inject a number of environment variables while running a command

```bash
❯ uses webgoat code ~/projects/webgoat
INFO[0000] Starting child process: code /Users/u/projects/webgoat
```

This mapping of projects to environment variables is stored in a config file:

```yaml
project:
- name: webgoat
  secrets:
  - GITHUB_TOKEN
- name: project1
  secrets:
  - AWS_USER
  - GITHUB_TOKEN
```

Location of the config file can be found using

```bash
❯ uses config
config file location: /Users/u/.config/uses/config.yaml
```

#### Same environment variable name for multiple projects

Sometimes the same environment variable name is required to be set for multiple projects. This can be achieved by using "key as variableName" syntax in the config file:

```yaml
project:
- name: webgoat
   secrets:
   - GITHUB_TOKEN_WEBGOAT as GITHUB_TOKEN
- name: project1
   secrets:
   - GITHUB_TOKEN_PROJECT1 as GITHUB_TOKEN
```

Secrets stored by `uses` in the scenario are `GITHUB_TOKEN_WEBGOAT` and `GITHUB_TOKEN_PROJECT1`. But when using `uses webgoat` or `uses project1` the environment variable name is the same i.e. `GITHUB_TOKEN`.

## 🛠 Contributing

Contributions to the `uses` package are most welcome from engineers of all backgrounds and skill levels. In particular the addition of support for other popular operating systems would be appreciated.

This project will adhere to the [Go Community Code of Conduct](https://go.dev/conduct) in the Github.

To make a contribution:

* Fork the repository
* Make your changes on the fork
* Submit a pull request back to this repo with a clear description of the problem you're solving
* Ensure your PR passes all current (and new) tests

## 🌈 Bucket list

* [x] configure auto-complete
* [ ] make `uses` available for other OSes as well in addition to Mac Darwin
* [ ] release on homebrew
* [ ] add more unit tests
