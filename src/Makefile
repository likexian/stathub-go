TARGET=../bin/stathub

build:
	python template.py
	go build -v -o ${TARGET}

i686:
	python template.py
	GOOS=linux GOARCH=386 go build -v -ldflags '-w -s' -o ${TARGET}

x86_64:
	python template.py
	GOOS=linux GOARCH=amd64 go build -v -ldflags '-w -s' -o ${TARGET}

clean:
	if [ -f ${TARGET} ]; then rm -rf ${TARGET}; fi
