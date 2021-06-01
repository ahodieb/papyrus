# Development Instructions

I don't expect external contributions, but this should be useful when i try to update the code in 6 months 🤷‍♀️

## Requirements

* [Go](https://golang.org/)
* [Git](https://git-scm.com/)
* [Make](https://www.gnu.org/software/make/manual/html_node/Introduction.html)

```bash
brew install go git make
```

## Development

```bash
make build # builds a binary for current OS
make test # runs tests
make # default (runs tests and builds)
```

## Packages used

* [Cobra](https://github.com/spf13/cobra) for command line flags parsing
* [Viper](https://github.com/spf13/viper) for config parsing
* [Go-Cmp](https://github.com/google/go-cmp) for diffs
  