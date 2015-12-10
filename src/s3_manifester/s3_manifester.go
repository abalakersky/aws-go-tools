package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	bucket = flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	logFile = flag.String("file", "yes", "Save output to file instead of displaying on the screen")
	region = flag.String("region", "us-east-1", "Region to connect to.")
	creds  = flag.String("creds", "default", "Credentials Profile to use")
	search = flag.String("search", "", "Search string to find in object paths")
	t      = time.Now()
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	w = bufio.NewWriter(os.Stdout)
)

func caseInsesitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	flag.Parse()
	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}


	if *logFile == "yes" {
		var s string
		if *search != "" {
			s = *search
		}

		name := dir + "/" + *bucket + "_" + s + "_" + strconv.FormatInt(t.Unix(), 10) + ".log"

		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}



	topLevel, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket, Delimiter: aws.String("/")})
	if err != nil {
		log.Println("Failed to list Top Level objects", err)
		return
	}
	for _, contentKeys := range topLevel.Contents {
		switch {
		case *search == "":
			fmt.Fprintln(w, *contentKeys.Key)
		case *search != "":
			if caseInsesitiveContains(*contentKeys.Key, *search) == true {
				fmt.Fprintln(w, *contentKeys.Key)
			} else {
				continue
			}
		}
	}

	var prefixes []string
	for _, commonPrefix := range topLevel.CommonPrefixes {
		prefixes = append(prefixes, *commonPrefix.Prefix)
	}

	objCh := make(chan *s3.Object, 100)
	var wg sync.WaitGroup

	listObjectsWorker := func(objCh chan<- *s3.Object, prefix string, bucket *string, svc s3iface.S3API) {
		params := &s3.ListObjectsInput{
			Bucket: bucket,
			Prefix: &prefix,
		}
		err := svc.ListObjectsPages(params,
			func(page *s3.ListObjectsOutput, last bool) bool {
				for _, object := range page.Contents {
					objCh <- object
					//				objCh <- fmt.Sprintf("%s", *object.Key)
				}
				return true
			},
		)

		if err != nil {
			fmt.Println("failed to list objects by prefix", prefix, err)
		}
		wg.Done()
	}

	wg.Add(len(prefixes))

	for i := range prefixes {
		prefix := prefixes[i]
		go listObjectsWorker(objCh, prefix, bucket, svc)
	}

	go func() {
		wg.Wait()
		close(objCh)
	}()

	for obj := range objCh {
		switch {
		case *search == "":
			fmt.Fprintln(w, *obj.Key)
			//				fmt.Println(*obj.Key)
		case *search != "":
			if caseInsesitiveContains(*obj.Key, *search) == true {
				fmt.Fprintln(w, *obj.Key)
				//				fmt.Println(*obj.Key)
			} else {
				continue
			}
		}
	}
	w.Flush()
}
