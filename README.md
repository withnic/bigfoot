# Bigfoot

* It navigates chrome browser.

# Need

Please install libraries
refs : https://agouti.org/

```
$ brew install phantomjs
$ brew install chromedriver
$ brew install selenium-server-standalone
```

# Usage

```
// single browswer 5sec display
$./bigfoot data/test.csv
```

or 

```
// display 10sec
$./bigfoot -t=10 data/test.csv
```

or 

```
// multi browswer , browser limit default 10
$./bigfoot -m data/test.csv
```

or 

```
// multi browswer, browswer limit 2
$./bigfoot -par=2 -m data/test.csv
```

test.csv

```
https://www.google.co.jp/
https://yahoo.co.jp/
```

# Help

```
$./bigfoot -h
```