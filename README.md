# get

`get` was created because I couldn't install `go` packages through the command `go get github.com/<AUTHOR>/<PACKAGE>`. The reason was because `git clone` doesn't work for me during to network/security constraints.

## Usage

Once installed you can install a `go` package by running the following:

```
$GOPATH/bin/get github.com/<AUTHOR>/<PACKAGE>
```
At the end of the execution, you will be prompted to run `go install` manually e.g.

```
Now you need to run the following command to install the package:

cd $GOPATH/src/github.com/lwheng/get && go install && cd -
```

##

Please star this repo if you see the irony
