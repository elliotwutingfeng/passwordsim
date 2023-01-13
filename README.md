# passwordsim

[![Go Reference](https://img.shields.io/badge/go-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/elliotwutingfeng/passwordsim)
[![Go Report Card](https://goreportcard.com/badge/github.com/elliotwutingfeng/passwordsim?style=for-the-badge)](https://goreportcard.com/report/github.com/elliotwutingfeng/passwordsim)
[![Codecov Coverage](https://img.shields.io/codecov/c/github/elliotwutingfeng/passwordsim?color=bright-green&logo=codecov&style=for-the-badge&token=1FMR3I0ZXO)](https://codecov.io/gh/elliotwutingfeng/passwordsim)

[![GitHub license](https://img.shields.io/badge/LICENSE-BSD--3--CLAUSE-GREEN?style=for-the-badge)](LICENSE)

## Summary

**passwordsim** lets you search for passwords similar to your specified password in any passwords dataset. The similarity metric used is the [Damerau-Levenshtein](https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance) distance.

## Use cases

- Choosing strong passwords
- Guessing password variants

## Requirements

Tested on Linux x64

- Fast multicore CPU
- At least 32 GB RAM recommended for larger datasets (multi-GB)
- Go 1.19

## Installation

```sh
go get github.com/elliotwutingfeng/passwordsim
```

## Setup

passwordsim executable will be created in the folder **'dist/'**

```bash
make build_cli
```

## Usage

Search for passwords similar to '[correct horse battery staple](https://xkcd.com/936)' in **'test/passwords.txt'**, with a normalised Damerau-Levenshtein distance score no larger than **0.3**.

```bash
dist/passwordsim -f 'test/passwords.txt' -o 'output.txt' -p 'correct horse battery staple' -t 0.3
```

For Windows systems, use back slashes `\`.

### Terminal output

```text
Searching for similar passwords...
3 similar passwords found in 'test/passwords.txt'. Threshold: 0.30
Results saved to file 'output.txt'.
```

### File output

**Filename:** output.txt

```text
correct horse battery staple 0
incorrect horse battery staple 0.06666666666666667
incorrect horse battery st@ple 0.1
```

The number on the right of each password is its [normalised Damerau-Levenshtein distance score](https://github.com/lmas/Damerau-Levenshtein) relative to your specified password. Scores range from 0 to 1. The smaller the number, the more similar the password. 0 implies exact match.

## Password datasets

- [SecLists](https://github.com/danielmiessler/SecLists)
- [Xato.net 10 million passwords](https://xato.net/today-i-am-releasing-ten-million-passwords-b6278bbe7495)
- [Weekpass](https://weakpass.com)
