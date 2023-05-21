package lambda

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

//func main() {
//	lambda.Start(Handler)
//}

//set GOOS=linux
//set GOARCH=amd64
//set CGO_ENABLED=0
//
// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
//
//C:\Users\hteren\go\bin\build-lambda-zip.exe -o lambda-handler.zip main

// build-lambda-zip --output .\target\main.zip .\target\main

//set GOOS=linux
//go build -o main
//mkdir target
//move .\main .\target\main
//build-lambda-zip --output .\target\main.zip .\target\main

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	//request.PathParameters
	//request.Headers

	//rCtx := request.RequestContext
	//if rCtx.Authorizer == nil || rCtx.Authorizer.JWT == nil {
	//	return events.APIGatewayV2HTTPResponse{}, errors.New("nil auth")
	//}

	r := map[string]any{
		"headers":      request.Headers,
		"path-params":  request.PathParameters,
		"query-params": request.QueryStringParameters,
		"raw-path":     request.RawPath,
		"raw-query":    request.RawQueryString,
		//"claims": rCtx.Authorizer.JWT.Claims,
		//"scopes": rCtx.Authorizer.JWT.Scopes,
	}

	res, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(res),
	}, nil
}
