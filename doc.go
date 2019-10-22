/*
A CLI tool to csv import to database along directory.

Usage

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
            if the first character in file name is not number, then add file name to table name
        -i bool
            Ignore 1st line when import in CSV (default: false)
        -a bool
            Auto completion with file name when lack of csv columns (default: false)
        -sn bool
            If csv columns is camelcase, convert to snakecase (default: false)

Example

Case1:

    $ csv2sql -d user -p root ./csv


CSV files import to database. then table will be along directory

    file                     -> import table
    -----------------------------------------
    csv
     └── user
          ├── 1.csv          -> user
          └── detail.csv     -> user
          └── task
                ├── 1.csv    -> user_task
                └── todo.csv -> user_task


Case2:

    $ csv2sql -d user -p root -s ./csv

option "-s" works as follows

    file                     -> import table
    -----------------------------------------
    csv
     └── user
          ├── 1.csv          -> user
          └── detail.csv     -> user_detail
          └── task
                ├── 1.csv    -> user_task
                └── todo.csv -> user_task_todo


Case3:

    $ csv2sql -d user -p root -S task ./csv

option "-S" filtering words

    file                     -> import table
    -----------------------------------------
    csv
     └── user
          ├── 1.csv
          └── detail.csv
          └── task
                ├── 1.csv    -> user_task
                └── todo.csv -> user_task
*/
package main
