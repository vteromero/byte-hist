# byte-hist
byte-hist is a simple tool that prints a histogram of the bytes appearing in
a file. It also gives you further information about the byte distribution.
## Installation
To install byte-hist, simply run:
```
$ go get bitbucket.org/vteromero/byte-hist
```
Once the `get` completes, you should find the `byte-hist` executable inside
`$GOPATH/bin`.
## Usage
These are some examples of how you can run `byte-hist`:
```
$ ./byte-hist file
$ ./byte-hist -format=x file
$ echo "hello world!" | ./byte-hist
```
For the complete available options, just type:
```
$ ./byte-hist -help
```

