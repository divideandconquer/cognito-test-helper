# Cognito test helper

This script can register and login cognito users for testing purposes.  This can be helpful if you are building services that require a Cognito JWT to authenticate users.

In order for it to work hassle free the cognito user group and client app configuration has a few caveats:

* The user group should have MFA and verification disabled.
* The user group should only require emails/passwords.
* The client app should not have a client secret.
* The client app should have `ADMIN_NO_SRP_AUTH` enabled.

Production services probably shouldn't authenticate this test user group's JWTs.

**Note** that this script relies on standard AWS authentication so it will require the proper AWS ID/secret in your .aws folder or set as environment variables.

## Usage

```sh
# Register a user (returns the user's UUID) 
AWS_REGION=us-east-1 go run main.go -register -clientID <client ID> -email foo@example.com -password Password1! -userPoolID <user pool ID>

# Login a user (returns a valid access token for the user)
AWS_REGION=us-east-1 go run main.go -login -clientID <client ID> -email foo@example.com -password Password1! -userPoolID <user pool ID>

# Register and login in 1 call
AWS_REGION=us-east-1 go run main.go -login -register -clientID <client ID> -email foo@example.com -password Password1! -userPoolID <user pool ID>
```


