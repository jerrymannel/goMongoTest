docker build -t jerrymannel/gomongotest:1.13 .

# building for go 1.10
docker build -t jerrymannel/gomongotest:1.10 -f Dockerfile_go1.10 .

# building for go 1.11
docker build -t jerrymannel/gomongotest:1.11 -f Dockerfile_go1.11 .

# building for go 1.13
docker build -t jerrymannel/gomongotest:1.13 -f Dockerfile_go1.13 .