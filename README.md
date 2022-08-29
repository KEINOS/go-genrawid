> __Warning__
> This package has been deprecated. We have found a faster, simpler and more reliable way. See [issue #12](https://github.com/KEINOS/go-genrawid/issues/12).

# GenRawID

`genrawid` is a niche command line tool and/or Go package to generate a unique consistent number from the input.

```shellsession
$ genrawid --string "abcdefgh"
-2474118025671277174
$ # 'abcdefgh'           --> The data to store in SQLite3
$ # -2474118025671277174 --> The primary key/rawid of stored data
```

```shellsession
$ echo -n 'abcdefgh' > ./sample.txt
$ genrawid ./sample.txt
-2474118025671277174
$ # 'abcdefgh'           --> The data to store in SQLite3
$ # -2474118025671277174 --> The primary key/rawid of stored data
```

`genrawid` is similar to a hash function, but **its value is a combination of a hash and a checksum**.

The hash value is the digest of the input and the checksum is the CRC of the complete hash value.

By default, [BLAKE3](https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE3)-512 is used for the hash algorithm and [CRC-32C](https://en.wikipedia.org/wiki/Cyclic_redundancy_check#Polynomial_representations_of_cyclic_redundancy_checks) (CRC-32 with Castagnoli polynomial) for the checksum.

```text
8 Bytes = The first 4 Bytes of the hash + 4 Bytes of the checksum of the hash
```

- See benchmark of BLAKE3 and CRC-32 comparing to other hash algorithms:
    - https://github.com/KEINOS/go-blake3-example/blob/main/bench_results/bench_results_stats.txt

## Why?

The main objective is to use [SQLite3](https://www.sqlite.org/) as a fast [KVS](https://en.wikipedia.org/wiki/Key%E2%80%93value_database) (Key-Value-Store) for [CAS](https://en.wikipedia.org/wiki/Content-addressable_storage) (Content-Addressable-Storage) usage.

`genrawid` generates a 64-bit/8-byte signed decimal number that can be used as [rawid in SQLite3](https://www.sqlite.org/lang_createtable.html#rowid).

> Searching for a record with a specific rowid, or for all records with rowids within a specified range is around **twice as fast** as a similar search made by specifying any other PRIMARY KEY or indexed value.
>
> <sub><sup>(From "[ROWIDs and the INTEGER PRIMARY KEY](https://www.sqlite.org/lang_createtable.html#rowid)" @ sqlite.org)</sup></sub>

So, in theory, if you know the rawid of the content, you can find it twice as fast.

On the other hand, **this command/package is of little use** when using SQLite3 as [RDB](https://en.wikipedia.org/wiki/Relational_database) (Relational Database) or when dealing with mutable content.

## Install

Download the binary for your operating system and architecture and place it as an executable in your PATH.

- [Latest Releases](https://github.com/KEINOS/go-genrawid/releases/latest) (Windows, macOS, Linux, RaspberryPi)

For [Homebrew/Linuxbrew](https://brew.sh/) users:

```bash
brew install KEINOS/apps/genrawid
```

## Statuses

- Unit Tests/Code Coverage

[![go1.16+](https://github.com/KEINOS/go-genrawid/actions/workflows/go-versions.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/go-versions.yml)
[![PlatformTests](https://github.com/KEINOS/go-genrawid/actions/workflows/platform-test.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/platform-test.yml)
[![codecov](https://codecov.io/gh/KEINOS/go-genrawid/branch/main/graph/badge.svg?token=cFoXdcwtaj)](https://codecov.io/gh/KEINOS/go-genrawid)
[![Go Report Card](https://goreportcard.com/badge/github.com/KEINOS/go-genrawid)](https://goreportcard.com/report/github.com/KEINOS/go-genrawid)

- Secrurity/Vulnerability Check

[![golangci-lint](https://github.com/KEINOS/go-genrawid/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/golangci-lint.yml)
[![CodeQL](https://github.com/KEINOS/go-genrawid/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/codeQL-analysis.yml)
[![Weekly Update](https://github.com/KEINOS/go-genrawid/actions/workflows/weekly-update.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/weekly-update.yml)

## CONTRIBUTING

[![Go Reference](https://pkg.go.dev/badge/github.com/KEINOS/go-genrawid.svg)](https://pkg.go.dev/github.com/KEINOS/go-genrawid/ "View document")
[![Opened Issues](https://img.shields.io/github/issues/KEINOS/go-genrawid?color=lightblue&logo=github)](https://github.com/KEINOS/go-genrawid/issues "opened issues")
[![PR](https://img.shields.io/github/issues-pr/KEINOS/go-genrawid?color=lightblue&logo=github)](https://github.com/KEINOS/go-genrawid/pulls "Pull Requests")

- [GolangCI Lint](https://golangci-lint.run/) rules: [.golangci-lint.yml](https://github.com/KEINOS/go-genrawid/blob/main/.golangci.yml)
- To run tests in a container:
  - `docker-compose --file ./.github/docker-compose.yml run v1_17`
  - This will run:
    - `go test -cover -race ./...`
    - `golangci-lint run`
    - `golint ./...`
- Branch to PR: `main`
  - It is recommended that [DraftPR](https://github.blog/2019-02-14-introducing-draft-pull-requests/) be done first to avoid duplication of work.

## License

- [MIT](https://github.com/KEINOS/go-genrawid/blob/main/LICENSE), Copyright (c) [KEINOS and the GenRawID contributors](https://github.com/KEINOS/go-genrawid/graphs/contributors).
