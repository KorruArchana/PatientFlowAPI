package utils

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/emisgroup/PatientFlowAPI/src/internal/constant"
	"github.com/google/uuid"
)

// GenerateID is used to generate new unique Id
func GenerateID(module string) string {
	newuuid, _ := uuid.NewRandom()
	return constant.EmisERN + ":" + module + ":" + newuuid.String()
}

// GeneratePK is used to create PK
func GeneratePK(module string) string {
	return constant.EmisERN + ":" + module
}

// CreateSession is used to create new session
func CreateSession() (*session.Session, error) {
	creds := credentials.NewEnvCredentials()
	// creds := credentials.NewSharedCredentials("C:/credentials", "default")
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2"), Credentials: creds})

	return sess, err
}

// ServerError returns internal server error proxy response
func ServerError(err error, errorSource string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       errorSource + ":" + err.Error(),
	}
}

// ClientError returns client specific error proxy response
func ClientError(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       "Request body has no data: " + http.StatusText(status),
	}
}

// GetAPIGateWayResponse returns proxy response
func GetAPIGateWayResponse(httpStatus int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: httpStatus,
		Headers: map[string]string{
			"Content-Type":                "application/vnd.api+json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: body,
	}
}
