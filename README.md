# When - time zone converter for the terminal

Usage: `when ZONE [ZONE ...]`, where ZONE is a zone name, or part of it.

Example:
```
$ when brisbane utc
Saturday 2021-05-08 (EEST +03:00)
Zone      Δt   Time
Local          Sa  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
Brisbane +7:00  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 Su  1  2  3  4  5  6
UTC      -3:00 21 22 23 Sa  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20
```

## Installation

Download the binary from https://github.com/jupj/when/releases/

Or install from source (requires Go 1.16): `go get -u github.com/jupj/when`
