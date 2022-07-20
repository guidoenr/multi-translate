# Translate-Terminal
# author @guidoenr

install:
	go get main && go mod tidy >> /dev/null

run:
	go build main && go run main >> /dev/null