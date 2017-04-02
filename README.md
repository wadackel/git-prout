# git-prout

[![Travis](https://img.shields.io/travis/tsuyoshiwada/git-prout.svg?style=flat-square)](https://travis-ci.org/tsuyoshiwada/git-prout)
[![GitHub release](http://img.shields.io/github/release/tsuyoshiwada/git-prout.svg?style=flat-square)](https://github.com/tsuyoshiwada/git-prout/releases)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/tsuyoshiwada/git-prout/blob/master/LICENSE)

> `git-prout` is a CLI tool using Golang. You can easily checkout GitHub Pull Request locally.



## Table of Contents

- [Demo](#demo)
- [Install](#install)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)



## Demo

![Demo Animation](./docs/demo.gif)



## Install

### Homebrew

```bash
$ brew tap tsuyoshiwada/git-prout
$ brew install git-prout
```

If you are in another platform, you can download binary from [release page](todo) and place it in `$PATH` directory.

### Golang

Or you can use `go get`.

```bash
$ go get -u github.com/tsuyoshiwada/git-prout
```



## Usage

```bash
$ git-prout [<flags>] <number>

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
      --debug            Enable debug mode.
  -r, --remote="origin"  Reference of remote.
  -f, --force            Force execute pull or checkout.
  -q, --quiet            Silencey any progress and errors.
      --version          Show application version.

Args:
  <number>  ID number of pull request

```


### Tips

You can execute `git-prout` as a git subcommand.

```bash
$ git prout 123
```



## Contribute

1. Fork (https://github.com/tsuyoshiwada/git-prout)
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test` command and confirm that it passes
1. Create new Pull Request :)

Bugs, feature requests and comments are more than welcome in the [issues](https://github.com/tsuyoshiwada/git-prout/issues).



## License

[MIT © tsuyoshiwada](./LICENSE)
