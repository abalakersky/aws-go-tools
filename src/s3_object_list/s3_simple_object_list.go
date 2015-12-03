package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	bucket := flag.String("bucket", "", "Bucket Name to list objects from")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	flag.Parse()
	if *bucket == "" {
		fmt.Println("\n" + "You must specify bucket name" + "\n")
		return
	}

	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	result, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket, Delimiter: aws.String("/")})
	if err != nil {
		log.Println("Failed to list objects", err)
		return
	}

//	fmt.Println(result)
	for _, commonPrefix := range result.CommonPrefixes {
		fmt.Println("Prefix:", *commonPrefix.Prefix)
	}
	for _, keys := range result.Contents {
		fmt.Println("Key:", *keys.Key)
	}

}
