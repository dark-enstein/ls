all: clean build

build:
	@go build -o cargs main.go
	@mv cargs ~/.adhoc/bin

clean:
	@[ -f ~/.adhoc/bin/cargs ] && ( rm ~/.adhoc/bin/cargs ) || printf ""

bench:
	@go test -bench=. -count 20 -run=^#