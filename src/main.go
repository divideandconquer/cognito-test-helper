package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// AWS_REGION=us-east-1 go run main.go -login -clientID <client ID> -email foo@example.com -password Password1! -userPoolID <user pool ID>

func main() {
	register := flag.Bool("register", false, "set to register a new user")
	login := flag.Bool("login", false, "set to login an existing user")
	email := flag.String("email", "", "The email of the user to act on")
	password := flag.String("password", "", "The password of the user to act on")
	clientID := flag.String("clientID", "", "The client id to register the user with")
	userPoolID := flag.String("userPoolID", "", "The user pool id to register the user with")

	flag.Parse()

	if !(*register || *login) {
		fmt.Println("-register or -login must be set.")
	}

	if len(*email) == 0 {
		fmt.Println("-email must not be empty")
	}
	if len(*password) == 0 {
		fmt.Println("-password must not be empty")
	}

	if len(*clientID) == 0 {
		fmt.Println("-clientID must not be empty")
	}

	if len(*userPoolID) == 0 {
		fmt.Println("-userPoolID must not be empty")
	}

	sess := session.Must(session.NewSession())
	cog := cognitoidentityprovider.New(sess, nil)

	if *register {
		userID, err := registerUser(cog, userPoolID, clientID, email, password)
		if err != nil {
			fmt.Printf("Error registering user: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Registered user successfully with UUID: %s\n", userID)
	}

	if *login {
		token, err := loginUser(cog, userPoolID, clientID, email, password)
		if err != nil {
			fmt.Printf("Error authenticating user: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Printf("User logged in successfully with token: \n\n%s\n\n", token)
	}
	os.Exit(0)
}

func registerUser(cog *cognitoidentityprovider.CognitoIdentityProvider, userPoolID, clientID, email, password *string) (string, error) {

	input := cognitoidentityprovider.SignUpInput{
		ClientId: clientID,
		Username: email,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			&cognitoidentityprovider.AttributeType{
				Name:  aws.String("email"),
				Value: email,
			},
		},
		Password: password,
	}

	output, err := cog.SignUp(&input)
	if err != nil {
		return "", err
	}
	if output == nil || output.UserSub == nil {
		return "", fmt.Errorf("Response was empty")
	}

	confirmInput := cognitoidentityprovider.AdminConfirmSignUpInput{
		UserPoolId: userPoolID,
		Username:   email,
	}
	_, err = cog.AdminConfirmSignUp(&confirmInput)
	if err != nil {
		return "", fmt.Errorf("Couldnt auto confirm user: %s", err.Error())
	}
	return *output.UserSub, nil
}

func loginUser(cog *cognitoidentityprovider.CognitoIdentityProvider, userPoolID, clientID, email, password *string) (string, error) {
	input := cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeAdminNoSrpAuth),
		AuthParameters: map[string]*string{
			"USERNAME": email,
			"PASSWORD": password,
		},
		UserPoolId: userPoolID,
		ClientId:   clientID,
	}

	output, err := cog.AdminInitiateAuth(&input)
	if err != nil {
		return "", err
	}
	if output == nil || output.AuthenticationResult == nil || output.AuthenticationResult.AccessToken == nil {
		return "", fmt.Errorf("Response was empty")
	}

	return *output.AuthenticationResult.AccessToken, nil
}
