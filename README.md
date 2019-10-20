# CSV2SQL
[![CircleCI](https://circleci.com/gh/sanshirookazaki/csv2sql.svg?style=svg)](https://circleci.com/gh/sanshirookazaki/csv2sql)
[![Coverage Status](https://coveralls.io/repos/github/sanshirookazaki/csv2sql/badge.svg?branch=master)](https://coveralls.io/github/sanshirookazaki/csv2sql?branch=master)
[![GoDoc](https://godoc.org/github.com/sanshirookazaki/csv2sql?status.svg)](https://godoc.org/github.com/sanshirookazaki/csv2sql)

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
    -S string
        Import specific tables (default: "")
    -s bool
        Separate CSV into 2 types. (default: false)
        if the first charactor in file name is not number, then add file name to table name
    -i bool
        Ignore 1st line when import in CSV (default: true)
```

## Example
Case1:
```
$ csv2sql -d user -p root ./csv
```

CSV files import to database. then table will be along directory
```
file                  　　-> import table
-----------------------------------------
csv
 └── user
      ├── 1.csv          -> user
      └── detail.csv     -> user
      └── task
            ├── 1.csv    -> user_task
            └── todo.csv -> user_task
```

<br>

Case2:
```
$ csv2sql -d user -p root -s ./csv
```

option "-s" works as follows
```
file                  　　-> import table
-----------------------------------------
csv
 └── user
      ├── 1.csv          -> user
      └── detail.csv     -> user_detail
      └── task
            ├── 1.csv    -> user_task
            └── todo.csv -> user_task_todo
```

<br>

Case3:
```
$ csv2sql -d user -p root -S task ./csv
```

option "-S" filtering words
```
file                  　　-> import table
-----------------------------------------
csv
 └── user
      ├── 1.csv
      └── detail.csv
      └── task
            ├── 1.csv    -> user_task
            └── todo.csv -> user_task
```
