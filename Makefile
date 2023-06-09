WORKDIR := $(PWD)

BUILD_LAMBDA ?=CGO_ENABLED=0 GOOS=linux go build -o main; mkdir build; move .\main .\build\main;	build-lambda-zip --output .\build\main.zip .\target\main
#SET GOOS=linux; SET CGO_ENABLED=0; go build -o main; mkdir build; move .\main .\build\main; build-lambda-zip --output .\build\main.zip .\build\main

#build: DIR
#build:
# 	GOOS=linux CGO_ENABLED=0 go build -o main
#	mkdir build
#	move .\main .\target\main
#	build-lambda-zip --output .\target\main.zip .\target\main



#//mockgen -package mock -destination currency/exchanger/mock/mock.go  github.com/glbter/currency-ex/currency/exchanger CurrencyRater,CurrencySeriesRater,AllCurrencyRater

#//mockgen -package mock -destination stocks/mock/mock.go  github.com/glbter/currency-ex/stocks PortfolioRepository,TickerRepository

#mockgen -package mock -destination pkg/sql/mock/mock.go  github.com/glbter/currency-ex/pkg/sql DB

build:
	mkdir target; \
	$(BUILD_LAMBDA)


#build-zip-file: compile-go
#	cd tmp && zip my-lambda-function.zip my-lambda-function && cd ..
#
#compile-go:
#	mkdir tmp; \
#		cd my-lambda-function; \
#		GOOS=linux go build -o ../tmp/
#
#clean-up:
#	rm tmp/my-lambda-function

#goplantuml -recursive ./stocks ./currency ./pkg > ./diagram.puml
