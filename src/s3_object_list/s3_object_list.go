package main

import (
	"time"
	"strconv"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"bufio"
	"os"
	"log"
	"path/filepath"
	"strings"
)

func CaseInsesitiveContains (s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	t := time.Now()
	bucket := flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	search := flag.String("search", "", "Search string to find in object paths")
	flag.Parse()
	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	name := dir + "/" + *bucket + "_" + *search + strconv.FormatInt(t.Unix(), 10) + ".log"
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var wg sync.WaitGroup
	keysCh := make(chan string, 10)

	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	params := &s3.ListObjectsInput{
		Bucket: bucket,
	}
	wg.Add(1)
	go func(param *s3.ListObjectsInput) {
		defer wg.Done()

		err := svc.ListObjectsPages(params,
			func(page *s3.ListObjectsOutput, last bool) bool {
				for _, object := range page.Contents {
					//keysCh <- fmt.Sprintf("%s:%s", *params.Bucket, *object.Key)
					keysCh <- fmt.Sprintf("%s", *object.Key)
				}
				return true
			},
		)
		if err != nil {
			fmt.Println("Error listing", *bucket, "Objects:", err)
		}
	}(params)
	go func() {
		wg.Wait()
		close(keysCh)
	}()
	w := bufio.NewWriter(f)
	for key := range keysCh {
		switch  {
		case *search == "" :
			fmt.Println(key)
			w.WriteString(key + "\n")
		case *search != "" :
			if CaseInsesitiveContains(key, *search) == true {
				fmt.Println(key)
				w.WriteString(key + "\n")
			} else {
				continue
			}
		}
	}
	w.Flush()
}
