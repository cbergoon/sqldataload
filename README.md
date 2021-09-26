<h1 align="center">sqldataload</h1>
<p align="center">
<a href="https://godoc.org/github.com/cbergoon/sqldataload"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

`sqldataload` reads and executes a directory of sql files on a provided SQL Server database. 

#### Documentation 

See the docs [here](https://godoc.org/github.com/cbergoon/sqldataload).

#### Install
```
go install github.com/cbergoon/sqldataload
```

#### Example Usage

```
USAGE: Usage: sqldataload [--directory=<directory>] <username>:<password>@<address>:<port>/<database>
```

To generate execute sql files provided `table-data` directory in specified directory use: 

```
$ sqldataload --directory=./table-data '<username>:<password>@127.0.0.1:1433/AdventureWorks'
```

To use windows authentication:

```
$ sqldataload --directory=./table-data '<domain>\<username>:<password>@127.0.0.1:1433/AdventureWorks'
```
#### License
This project is licensed under the MIT License.







