// package main
package authentication

import (
	// "github.com/aws/aws-sdk-go/aws/request"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "encoding/json"
	// "reflect"
	"strings"
	// "time"
	// "github.com/google/jsonapi"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main(){
	lambdaHandler()
}

func createSession()(*session.Session, error){
	creds := credentials.NewSharedCredentials("C:/credentials", "default")
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2"), Credentials: creds})

	return sess, err
}

func lambdaHandler(){
	sess, _ := createSession()
	userPoolID := "eu-west-2_MIPUHTLMd"
	// clientName := "emis-kiosk-3"
	kioskGUID := "emis:kiosk:3"
	allowedAuthType := "client_credentials"
	allowedContentType := "application/x-www-form-urlencoded"
	domainName := "https://patientflowkiosk.auth.eu-west-2.amazoncognito.com"

	cip := cognitoidentityprovider.New(sess)

	// createAppClient(cip, userPoolID, clientName, allowedAuthType)
	appPoolClients := getAppPoolClients(cip, userPoolID)
	kioskClientID := getKioskClientID(appPoolClients, kioskGUID)
	kioskAppClientDetails := getKioskAppClientDetails(cip, userPoolID, kioskClientID)
	kioskSecretKey := kioskAppClientDetails.ClientSecret
	accessToken := getAccessToken(kioskClientID, *kioskSecretKey, allowedAuthType, domainName, allowedContentType)

	fmt.Println(kioskClientID, kioskGUID, appPoolClients,domainName, allowedAuthType, kioskSecretKey, accessToken)
}

// GetAccessTokenForKiosk is used to get the access token for specific kiosk
func GetAccessTokenForKiosk(kioskGUID string)(string){

	sess, _ := createSession()
	userPoolID := "eu-west-2_MIPUHTLMd"
	allowedAuthType := "client_credentials"
	allowedContentType := "application/x-www-form-urlencoded"
	domainName := "https://patientflowkiosk.auth.eu-west-2.amazoncognito.com"

	cip := cognitoidentityprovider.New(sess)

	appPoolClients := getAppPoolClients(cip, userPoolID)
	kioskClientID := getKioskClientID(appPoolClients, kioskGUID)
	kioskAppClientDetails := getKioskAppClientDetails(cip, userPoolID, kioskClientID)
	kioskSecretKey := kioskAppClientDetails.ClientSecret
	accessToken := getAccessToken(kioskClientID, *kioskSecretKey, allowedAuthType, domainName, allowedContentType)
	return accessToken
}

// func createAppClient(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, kioskGUID string, allowedAuthType string){

// 	clientInput := cognitoidentityprovider.CreateUserPoolClientInput{
// 		UserPoolId: aws.String(userPoolID),
// 		ClientName: aws.String(kioskGUID),
// 		AllowedOAuthFlows: []*string{ aws.String(allowedAuthType)},
// 		AllowedOAuthScopes: []*string{ aws.String("https://gjpptao120.execute-api.eu-west-2.amazonaws.com/dev/organisation/kiosk_read")},
// 		GenerateSecret: aws.Bool(true),
// 		AllowedOAuthFlowsUserPoolClient: aws.Bool(true),
// 	}

// 	clientResponse, err := cip.CreateUserPoolClient(&clientInput)

// 	if err != nil{
// 		fmt.Println("Got error while creating cognito client: ", err)
// 	}

// 	fmt.Println(clientResponse)
// }

func getAppPoolClients(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string)([]*cognitoidentityprovider.UserPoolClientDescription){

	userPoolInput := cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId : aws.String(userPoolID),
	}

	listOfUserPoolClients, err := cip.ListUserPoolClients(&userPoolInput)

	if err != nil {
		fmt.Println("Got error while getting list of user pool clients: ", err)
	}

	return listOfUserPoolClients.UserPoolClients
}

func getKioskClientID(poolAppClients []*cognitoidentityprovider.UserPoolClientDescription, kioskGUID string)(string){

	replacer := strings.NewReplacer(":", "-")

	for i := range poolAppClients{
		if *poolAppClients[i].ClientName == replacer.Replace(kioskGUID){
			return *poolAppClients[i].ClientId
		}
	}

	return " "
}

func getKioskAppClientDetails(cip *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, kioskClientID string)(*cognitoidentityprovider.UserPoolClientType){

	userPoolInput := cognitoidentityprovider.DescribeUserPoolClientInput{
		UserPoolId : aws.String(userPoolID),
		ClientId: aws.String(kioskClientID),
	}

	responseUserPoolClientOutput, err := cip.DescribeUserPoolClient(&userPoolInput)

	if err != nil{
		fmt.Println("Got error while getting pool client data", err)
	}
	return responseUserPoolClientOutput.UserPoolClient
}

func getAccessToken(kioskClientID string, kioskSecretKey string, allowedAuthType string, domainName string, allowedContentType string) (string){
	accessToken := ""
	url := domainName + "/oauth2/token?grant_type=" + allowedAuthType
	req, err := http.NewRequest("POST", url, nil)
	if err != nil{
		fmt.Println("Got error in creating new http request", err)
	}

	req.SetBasicAuth(kioskClientID, kioskSecretKey)
	req.Header.Set("Content-Type", allowedContentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		fmt.Println("Got error while making call to cognito", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	accessToken = string(body)
	fmt.Println(accessToken)

	return accessToken
}