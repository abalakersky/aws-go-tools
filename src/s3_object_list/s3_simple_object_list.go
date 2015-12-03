package main

import (
//	"encoding/json"
	"flag"
	"fmt"
	"log"
//	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
//	"io/ioutil"
//	"os"
)


type directories struct {
	CommonPrefixes []commonPrefixes `json:"CommonPrefixes"`
}

type commonPrefixes struct {
	Prefix string `json:"Prefix"`
}

type files struct {
	Keys []contents `json:"Contents"`
}

type contents struct {
	Key string `json:"Key"`
}

func main() {
//	t := time.Now()
	bucket := flag.String("bucket", "", "Bucket Name to list objects from")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	flag.Parse()
	//	fmt.Printf("Your bucket name is %q", *bucket)
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

//	fmt.Printf("Here are the objects in %q bucket on %s\n\n", *bucket, t.Format(time.RFC1123))
	fmt.Println(result)
	for _, w := range result.CommonPrefixes {
		fmt.Printf("%v\n", w)
	}
//	fmt.Println(result.CommonPrefixes)
//	fmt.Println(result.Prefix)
//	fmt.Println(result.Contents)
//	fmt.Printf("%+v\n", result.CommonPrefixes)

//	for range result.CommonPrefixes {
//		fmt.Printf("Directory is %s\n", result.CommonPrefixes)
//	}
//	j, err := json.Marshal(result.CommonPrefixes)
//	fmt.Println(string(j))
//
//	var a commonPrefixes
//	var b files
//	err = json.Unmarshal(j, &a)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(a)
//	err = json.Unmarshal(contents, &b)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var p []string
//	var r []string
//	for i := range a.CommonPrefixes {
//		p = append(p, a.CommonPrefixes[i].Prefix)
//	}
//	for i := range b.Keys {
//		r = append(r, b.Keys[i].Key)
//	}
//
//	fmt.Printf("%v\n", p)
//	fmt.Printf("%v\n", r)

	//	for _, object := range result.Contents {
	//		fmt.Printf("%s\n", *object.Key)
	//		fmt.Printf(*result)
	//	}
}
