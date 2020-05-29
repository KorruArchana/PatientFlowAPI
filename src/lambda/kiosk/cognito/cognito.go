package main

import (
	"bytes"
	"encoding/json"
	"reflect"

	// "strings"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/jsonapi"
)

func main() {
	lambdaHandler()
}

func createSession() (*session.Session, error) {
	creds := credentials.NewSharedCredentials("C:/credentials", "default")
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2"), Credentials: creds})

	return sess, err
}

func lambdaHandler() {
	sess, _ := createSession()
	userPoolID := "eu-west-2_MIPUHTLMd"
	clientName := "emis-kiosk-4"
	kioskGUID := "emis:kiosk:4"
	allowedAuthType := "client_credentials"
	domainName := "https://patientflowkiosk.auth.eu-west-2.amazoncognito.com"

	cip := cognitoidentityprovider.New(sess)

	// createAppClient(cip, userPoolID, clientName, allowedAuthType)
	appPoolClients := getAppPoolClients(cip, userPoolID)
	// kioskClientID := getKioskClientID(appPoolClients, kioskGUID)
	// kioskAppClientDetails := getKioskAppClientDetails(cip, userPoolID, kioskClientID)
	// kioskSecretKey := kioskAppClientDetails.ClientSecretKey
	// accesstoken := getAccessToken(kioskClientID, kioskSecretKey, allowedAuthType, domainName)

	fmt.Println(kioskGUID, appPoolClients, domainName, clientName, allowedAuthType)
	// fmt.Println(kioskSecretKey)
}

// type poolClients struct{
// 	userPoolClientDescription []*userPoolClientDescription `jsonapi:"relation,UserPoolClientDescription"`
// }

type appPoolClients struct {
	allPoolClients poolClients `jsonapi:"relation,UserPoolClientDescription"`
}

type poolClients []struct {
	ClientID   string `jsonapi:"primary,ClientId"`
	UserPoolID string `jsonapi:"attr,UserPoolId"`
	ClientName string `jsonapi:"attr,ClientName"`
}

type appClient struct {
	AllowedOAuthFlows               []*string  `json:AllowedOAuthFlows`
	AllowedOAuthFlowsUserPoolClient bool       `json:AllowedOAuthFlowsUserPoolClient`
	AllowedOAuthScopes              []*string  `json:AllowedOAuthScopes`
	ClientID                        string     `json:ClientId`
	ClientName                      string     `json:ClientName`
	ClientSecretKey                 string     `json:ClientSecret`
	CreationDate                    *time.Time `json:CreationDate`
	LastModifiedDate                *time.Time `json:LastModifiedDate`
	RefreshTokenValidity            int        `json:RefreshTokenValidity`
	UserPoolID                      string     `json:UserPoolId`
}

func createAppClient(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, kioskGUID string, allowedAuthType string) {

	clientInput := cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId:                      aws.String(userPoolID),
		ClientName:                      aws.String(kioskGUID),
		AllowedOAuthFlows:               []*string{aws.String(allowedAuthType)},
		AllowedOAuthScopes:              []*string{aws.String("https://gjpptao120.execute-api.eu-west-2.amazonaws.com/dev/organisation/kiosk_read")},
		GenerateSecret:                  aws.Bool(true),
		AllowedOAuthFlowsUserPoolClient: aws.Bool(true),
	}

	clientResponse, err := cip.CreateUserPoolClient(&clientInput)

	if err != nil {
		fmt.Println("Got error while creating cognito client: ", err)
	}

	fmt.Println(clientResponse)
}

func getAppPoolClients(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string) []poolClients {

	userPoolInput := cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: aws.String(userPoolID),
	}

	listOfUserPoolClients, err := cip.ListUserPoolClients(&userPoolInput)

	if err != nil {
		fmt.Println("Got error while getting list of user pool clients: ", err)
	}

	// firstClient := appPoolClients(listOfUserPoolClients.UserPoolClients)

	// fmt.Println(firstClient)

	poolAppClients := []poolClients{}
	appClient := new(poolClients)

	// appPoolClients := []poolClients{}
	// newappPoolClients := appPoolClients(firstClient)
	// fmt.Println(newappPoolClients)
	// bytesData := jsonapi.MarshalPayload(listOfUserPoolClients.UserPoolClients[0], )
	bytesData, _ := json.Marshal(listOfUserPoolClients.UserPoolClients)
	// err = jsonapi.MarshalPayload(listOfUserPoolClients.UserPoolClients[0], appClient)

	// jsonapi.UnmarshalPayload(bytes.NewReader(bytesData), appClient)

	// err = jsonapi.UnmarshalPayload(bytes.NewReader(bytesData),appClient)
	responsePoolAppClients, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(bytesData), reflect.TypeOf(appPoolClients{}))
	if err != nil {
		fmt.Println("Got error while unmarshalling poolClients: ", err)
	}

	fmt.Println(responsePoolAppClients)
	fmt.Println("FirstAppClient: ", appClient)
	return poolAppClients
}

// func getKioskClientID(poolAppClients []poolClients, kioskGUID string)(string){

// 	replacer := strings.NewReplacer(":", "-")

// 	for i,_ := range poolAppClients{
// 		if poolAppClients[i].ClientName == replacer.Replace(kioskGUID){
// 			return poolAppClients[i].ClientID
// 		}
// 	}

// 	return " "
// }

func getKioskAppClientDetails(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, kioskClientID string) *appClient {

	userPoolInput := cognitoidentityprovider.DescribeUserPoolClientInput{
		UserPoolId: aws.String(userPoolID),
		ClientId:   aws.String(kioskClientID),
	}

	_, responseUserPoolClientOutput := cip.DescribeUserPoolClientRequest(&userPoolInput)

	// if err != nil {
	// 	fmt.Println("Got error while getting list of user pool clients: ", err)
	// }

	appClientDetails := new(appClient)

	requestByte, _ := json.Marshal(responseUserPoolClientOutput)

	err := jsonapi.UnmarshalPayload(bytes.NewReader(requestByte), &appClientDetails)
	if err != nil {
		fmt.Println("Got error while unmarshalling poolClients: ", err)
	}

	fmt.Println("USerPoolCLient: ", responseUserPoolClientOutput)

	return appClientDetails
}

// func getAccessToken(kioskClientID, kioskSecretKey, allowedAuthType, domainName)string{

// 	return ""
// }
