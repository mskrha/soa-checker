[![Go Report Card](https://goreportcard.com/badge/github.com/mskrha/soa-checker)](https://goreportcard.com/report/github.com/mskrha/soa-checker)

## soa-checker

### Description
Simple tool to compare zone serial number on all NS and hidden master, if used.

### Build
```shell
git clone https://github.com/mskrha/soa-checker.git
cd soa-checker/source
make
```

### Usage
```shell
$ soa-checker -zone example.com -master 127.0.0.1
```
```shell
$ soa-checker -zone example.com -serial 1234567890
```
