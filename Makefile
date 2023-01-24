all: amd64 arm

amd64:
	env GOOS=linux GOARCH=amd64 go build -o ./compiled/amd_pitester

arm:
	env GOOS=linux GOARCH=arm go build -o ./compiled/arm_pitester

clean:
	rm -rf ./compiled/*
