all: build

clean:
	rm -rf pkg bin

test: clean
	gb test

build: test
	gb build all

update:
	gb vendor update --all
