# GenRawID

`genrawid` is a niche command line tool and/or Go package to generate a unique consistent number from the imput.

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

## Why?

The main objective is to use [SQLite3](https://www.sqlite.org/) as a fast [KVS](https://en.wikipedia.org/wiki/Key%E2%80%93value_database) (Key-Value-Store) for [CAS](https://en.wikipedia.org/wiki/Content-addressable_storage) (Content-Addressable-Storage) usage.

`genrawid` generates a 64-bit/8-byte signed decimal number that can be used as [rawid in SQLite3](https://www.sqlite.org/lang_createtable.html#rowid).

> Searching for a record with a specific rowid, or for all records with rowids within a specified range is around **twice as fast** as a similar search made by specifying any other PRIMARY KEY or indexed value.
>
> <sub><sup>(From "[ROWIDs and the INTEGER PRIMARY KEY](https://www.sqlite.org/lang_createtable.html#rowid)" @ sqlite.org)</sup></sub>

So, in theory, if you know the rawid of the content, you can find it twice as fast.

On the other hand, **this command/package is of little use** when using SQLite3 as [RDB](https://en.wikipedia.org/wiki/Relational_database) (Relational Database) or when dealing with mutable content.

## CONTRIBUTING

- Bug/security report: [Issues](https://github.com/KEINOS/go-genrawid/issues)
- Anything fot the better: Feel free to PR!

### Statuses

[![go1.16+](https://github.com/KEINOS/go-genrawid/actions/workflows/go-versions.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/go-versions.yml)
[![golangci-lint](https://github.com/KEINOS/go-genrawid/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/golangci-lint.yml)
[![Codecov](https://github.com/KEINOS/go-genrawid/actions/workflows/code-coverage.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/code-coverage.yml)
[![CodeQL](https://github.com/KEINOS/go-genrawid/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/codeQL-analysis.yml)
[![Weekly Update](https://github.com/KEINOS/go-genrawid/actions/workflows/weekly-update.yml/badge.svg)](https://github.com/KEINOS/go-genrawid/actions/workflows/weekly-update.yml)
