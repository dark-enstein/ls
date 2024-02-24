all: clean build

build:
	@go build -o cargs main.go
	@cp cargs ~/.adhoc/bin
	@rm cargs

clean:
	@[ -f ~/.adhoc/bin/cargs ] && ( rm ~/.adhoc/bin/cargs ) || printf ""

bench:
	@go test -bench=. -count 20 -run=^#