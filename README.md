# When - time zone converter

Usage: `when ZONE [ZONE ...]`, where ZONE is a zone name, or part of it.

Example:
```
$ when brisbane utc 
Saturday 2021-04-24 (EEST +03:00)
Zone      Δt   Time
Local          Sat 1  2  3   4  5  6  7  8  9  10 11 12 *13 14 15 16 17  18 19 20 21 22 23 
Brisbane +7:00 7   8  9  10  11 12 13 14 15 16 17 18 19 *20 21 22 23 Sun 1  2  3  4  5  6  
UTC      -3:00 21  22 23 Sat 1  2  3  4  5  6  7  8  9  *10 11 12 13 14  15 16 17 18 19 20 
```
