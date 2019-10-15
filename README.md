# CSV2SQL
[![CircleCI](https://circleci.com/gh/sanshirookazaki/csv2sql.svg?style=svg)](https://circleci.com/gh/sanshirookazaki/csv2sql)

A CLI tool to csv import to database along directory.

## Install

```
go get -u github.com/sanshirookazaki/csv2sql
```

## Usage

```
$ csv2sql [OPTIONS] [CSV_DIR]

OPTIONS:
    -d string
        Name of Database (default: "")
    -h string
        Host of Database (default: "127.0.0.1")
    -P int
        Database port number (default: 3306)
    -u string
        Database user (default: "root")
    -p string
        Database password (default: "")
```

## Example
Running:
```
$ csv2sql -d user -p admin ./user
```

CSV files import to database. then table will be along directory
```
file             -> import table

user
 ├── 1.csv       -> user
 ├── 2.csv       -> user
 └── item
     └── 1.csv   -> user_item
```
