/*
A CLI tool to csv import to database along directory.


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
        -n string
            Newline code (default: "\n")
        -N bool
            If csv value is empty, set null (default: false)
        -s bool
            Separate CSV into 2 types. (default: false)
            if the first character in file name is not number, then add file name to table name
        -i bool
            Ignore 1st line when import in CSV (default: false)
        -a bool
            Auto completion with file name when lack of csv columns (default: false)
        -sn int
            convert columns into snakecase (default: 0)
                0: Do nothing
                1: Snakecase      (e.g. "testUser123Id" -> "test_user_123_id")
                2: Ignore number  (e.g. "testUser123Id" -> "test_user123_id")
        -dry-run bool
            dry run (default: false)
        -f bool
            Force run: ignore error (default: false)

## Examples

Case1:
    $ csv2sql -d todo -i ./examples

CSV files import to database, then table will be along directory.

    file                  -> table
    -----------------------------------------
    examples
    ├── user
    │   ├── 1.csv         -> user
    │   ├── detail.csv    -> detail
    │   └── task
    │       ├── 1.csv     -> user_task
    │       └── 2.csv     -> user_task
    └── user.csv          -> examples


Case2:

    $ csv2sql -d todo -i -s ./examples

option "-s", works as follows

    file                  -> table
    -----------------------------------------
    examples
    ├── user
    │   ├── 1.csv         -> user
    │   ├── detail.csv    -> user_detail
    │   └── task
    │       ├── 1.csv     -> user_task
    │       └── 2.csv     -> user_task
    └── user.csv          -> user


Case3:

    $ csv2sql -d todo -i -S user/task ./examples

option "-S", csv path filtering

    file                  -> table
    -----------------------------------------
    examples
    ├── user
    │   ├── 1.csv
    │   ├── detail.csv
    │   └── task
    │       ├── 1.csv     -> user_task
    │       └── 2.csv     -> user_task
    └── user.csv



Case4:

option "-a", auto compretion with file name.

e.g. examples/user/task/1.csv

    id,task
    1,homework


result is follow. user_id is complemented by file name "1" (1.csv).

    +----+---------+----------+
    | id | user_id | task     |
    +----+---------+----------+
    |  1 |       1 | homework |
    +----+---------+----------+

*/
package main
