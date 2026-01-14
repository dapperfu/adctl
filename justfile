export bin := "adctl"

set dotenv-load := false

# test for quoted args to work, didn't do anything. doesn't matter much.
#set positional-arguments

default:
    just --list

coverage:
    go test ./cmd -coverprofile=coverage.out
    go tool cover -html=coverage.out

# TODO Need to clean all this up.

# run *ARGS: mac-notest
run *ARGS: build
    ./$bin {{ ARGS }}

qbuild:
    go build -o $bin .

qrun *ARGS: qbuild
    ./$bin {{ ARGS }}

qinstall: qbuild
    cp ./$bin ~/bin/

# TODO I hate that I need to install this in my path first
#  but I can't figure out how to get tescript to use ./adctl and stop searching my path
#  also tried 'env PATH=$PATH:$PWD' and that didn't work

# urgh.
test:
    go test ./cmd

testv:
    go test ./cmd -test.v

testall: test testcli

# TODO: removed a lot of tests because they're in testscripts now.
testcli: mac
    ./$bin status
    ./$bin log get 42 | jq '.data | length'

#    ./$bin status enable
#    ./$bin status
#    ./$bin status disable
#    ./$bin status
#    ./$bin status disable 15s
#    ./$bin status
#    ./$bin status enable
#    ./$bin status
#    ./$bin status toggle
#    ./$bin status
#    ./$bin status toggle
#    ./$bin status
#    ./$bin log get | jq '.oldest'

fmt:
    just --unstable --fmt
    goimports -l -w .
    go fmt

mac: test
    go build -o $bin .

mac-notest:
    go build -o $bin .

build: test
    go build -o $bin .

clean:
    go clean -testcache
    go mod tidy
    rm -f $bin 
    rm -rf dist

install: mac
    cp ./$bin ~/bin/

# TODO: prompt for a tag here?
# not for now
# git tag -a v0.1.0 -m "first release"
# git push origin v0.1.0
# takes two arguments. first is tag (v0.1.0), second is tag description.

# TODO: what do I do if I have uncommitted changes?
release arg1: testall
    rm -rf dist/
    git tag {{ arg1 }}
    git push origin {{ arg1 }}
    go build -o $bin .
