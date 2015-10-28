# safron

Quick-and-dirty static file web server, which you can obtain from
[the releases page](https://github.com/ripta/safron/releases), or
build it yourself:

```
  $ brew install go
  $ go get github.com/ripta/safron
```

By default, it make the current working directory available on port 8080. More
options are available through `-h`.

Example:

```
  $ pwd
  /Users/rpasay/Public
  $ safron
  2015/10/27 17:24:13 Safron version 1
  2015/10/27 17:24:13 Listening to http://0.0.0.0:8080
  2015/10/27 17:24:13 Serving /Users/rpasay/Public
  2015/10/27 17:24:13 ^C to exit
```

