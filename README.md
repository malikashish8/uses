# 🔐 uses - USE Secure environment variables in dev

[![Release Action](https://github.com/malikashish8/uses/actions/workflows/release.yaml/badge.svg)](https://github.com/malikashish8/uses/actions/workflows/release.yaml)
[![CodeQL Scan](https://github.com/malikashish8/uses/actions/workflows/codeql.yaml/badge.svg?branch=master)](https://github.com/malikashish8/uses/actions/workflows/codeql.yaml)
[![Semgrep Scan](https://github.com/malikashish8/uses/actions/workflows/semgrep.yaml/badge.svg?branch=master)](https://github.com/malikashish8/uses/actions/workflows/semgrep.yaml)

Taking inspiration from [aws-vault](https://github.com/99designs/aws-vaults), `uses` makes `use` of OS provided `s`ecret management solutions to save secrets in the development environment. Grouping of secrets is made possible by a config file.

Having secrets lying around in environment variables in the development environment can be a nightmare as opensource packages are being [actively compromised](https://thehackernews.com/2022/05/pypi-package-ctx-and-php-library-phpass.html) to steal secrets. These packages can read all the environment variables. Good security hygiene dictates that no secrets are stored in environment variables (using configs such as ~/.bashrc and ~/.zshrc). `uses` helps to implement the least privilege principle by saving all the secrets in a password protected secret store.

## ⚡️ Installation

Install using [Homebrew](https://brew.sh/):

```bash
brew install malikashish8/tap/uses
```

Or download the binary from [releases](https://github.com/malikashish8/uses/releases) and add it to your path.

## 🧑‍💻 Usage

### Set Secret

Set a secret in secret store

```bash
❯ uses set GITHUB_TOKEN
Enter value: 
```

or if the secret is already in the environment

```bash
❯ uses set GITHUB_TOKEN=${GITHUB_TOKEN}
```

### Get Secret

Get a secret from secret store

```bash
❯ uses get GITHUB_TOKEN
sdknbowhlfownpns;s/dkfnbslsnwwn
```

Get can also run any command passed to it after setting the environment variable

```bash
❯ uses get GITHUB_TOKEN env
GITHUB_TOKEN=sdknbowhlfownpns;s/dkfnbslsnwwn
```

### List

Get a list of secrets managed by `uses`

```bash
❯ uses list
AWS_USER
GITHUB_TOKEN
```

### Projects

Group secrets and inject them as environment variables while running a command

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

Location of the config file is `/Users/<USER>/.uses/config.yaml`. `uses config` opens the config with default editor.

#### Same environment variable name for multiple projects

Sometimes multiple projects use same variable name but different values. Though secret key has to be unique, this can be achieved by using "key as variableName" syntax in the config file. For example:

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

### Enable Auto-completion

  1. zsh - `echo 'source <(uses completion zsh)' >>~/.zshrc`
  2. bash - `echo 'source <(uses completion bash)' >>~/.bashrc`

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
* [x] release on homebrew
* [ ] add more unit tests
