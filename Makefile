install:
	dep ensure
	packr build -o antoragen

build:
	packr build -o antoragen