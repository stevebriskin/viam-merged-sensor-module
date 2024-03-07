bin: *.go */*.go go.*
	mkdir -p bin && CGO_ENABLED=0 go build -o bin/module cmd/main.go

module.tar.gz: bin
	tar -czf module.tar.gz run.sh bin

clean:
	rm -rf bin/* 
	rm -f module*.tar.gz  