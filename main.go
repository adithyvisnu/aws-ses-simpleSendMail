package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	KEY_ID          = "YOUR AWS SES KEY_ID HERE"
	SECRET_KEY      = "YOUR AWS SES SECRET_KEY HERE"
	TOKEN           = "YOUR AWS SES TOKEN HERE"
	RETURN_PATH_ARN = "YOUR AWS SES VERIFIED RETURN PATH ARN EMAIL ACCOUNT"
)

func main() {
	awsConf := &aws.Config{}
	awsConf.Endpoint = aws.String("email.us-east-1.amazonaws.com")
	awsConf.Region = aws.String("us-east-1")
	awsConf.MaxRetries = aws.Int(3)
	awsConf.Credentials = credentials.NewStaticCredentials(KEY_ID, SECRET_KEY, TOKEN)

	awsSess := session.Must(session.NewSession(awsConf))

	sesSess := ses.New(awsSess, awsConf)

	var from, subject, textbody, htmlbody, charset string
	from = "adithyavisnu@helio.id"
	subject = "Hello this is from AWS!"
	textbody = "it's good to be true"
	htmlbody = "<b>it's good to be true<b>"
	charset = "UTF-8"

	var to []string = []string{"fahrudin@helio.id", "ganiamri@helio.id"}
	var cc []string = []string{""}

	sendEmailOutput, err := sesSess.SendEmail(
		&ses.SendEmailInput{
			Source: aws.String(from),
			Destination: &ses.Destination{
				ToAddresses: aws.StringSlice(to),
				CcAddresses: aws.StringSlice(cc),
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(charset),
						Data:    aws.String(htmlbody),
					},
					Text: &ses.Content{
						Charset: aws.String(charset),
						Data:    aws.String(textbody),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(subject),
				},
			},
			ReturnPath:    aws.String(from),
			ReturnPathArn: aws.String(RETURN_PATH_ARN),
		},
	)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("Email Sent to address: ")
	fmt.Println(sendEmailOutput)
}
