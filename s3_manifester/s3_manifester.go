package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	""
)

var (
	bucket  = flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	logFile = flag.String("file", "yes", "Save output to file instead of displaying on the screen")
	region  = flag.String("region", "us-east-1", "Region to connect to.")
	creds   = flag.String("creds", "default", "Credentials Profile to use")
	search  = flag.String("search", "", "Search string to find in object paths")
	akid    = flag.String("akid", "", "AWS Access Key")
	secKey  = flag.String("seckey", "", "AWS Secret Access Key")
	csv     = flag.String("csv", "no", "Create CSV log with full output")
	t       = time.Now()
	dir, _  = filepath.Abs(filepath.Dir(os.Args[0]))
	w       = bufio.NewWriter(os.Stdout)
)

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	flag.Parse()

	var svc = s3.New(session.New())

	if *akid != "" && *secKey != "" {
		svc = s3.New(session.New(&aws.Config{
			Region:      region,
			Credentials: credentials.NewStaticCredentials(*akid, *secKey, ""),
		}))
	} else {
		svc = s3.New(session.New(&aws.Config{
			Region:      region,
			Credentials: credentials.NewSharedCredentials("", *creds),
		}))
	}

	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}

	if *logFile == "yes" {
		var s string
		var ext string
		if *search != "" {
			s = *search
		}

		if *csv != "no" {
			ext = ".csv"
		} else {
			ext = ".log"
		}

		name := dir + "/" + *bucket + "_" + s + "_" + strconv.FormatInt(t.Unix(), 10) + ext

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
		case *search == "" && *csv == "no":
			fmt.Fprintln(w, *contentKeys.Key)
		case *search != "" && *csv == "no":
			if caseInsensitiveContains(*contentKeys.Key, *search) == true {
				fmt.Fprintln(w, *contentKeys.Key)
			} else {
				continue
			}
		case *search == "" && *csv == "yes":
			fmt.Fprintf(w, "%s,%d,%s\n", *contentKeys.Key, *contentKeys.Size, *contentKeys.ETag)
		case *search != "" && *csv == "yes":
			if caseInsensitiveContains(*contentKeys.Key, *search) == true {
				fmt.Fprintf(w, "%s,%d,%s\n", *contentKeys.Key, *contentKeys.Size, *contentKeys.ETag)
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

	for i := range prefixes {
		prefix := prefixes[i]
		wg.Add(1)
		go listObjectsWorker(objCh, prefix, bucket, svc)
	}

	go func() {
		wg.Wait()
		close(objCh)
	}()

	for obj := range objCh {
		switch {
		case *search == "" && *csv == "no":
			fmt.Fprintln(w, *obj.Key)
		case *search != "" && *csv == "no":
			if caseInsensitiveContains(*obj.Key, *search) == true {
				fmt.Fprintln(w, *obj.Key)
			} else {
				continue
			}
		case *search == "" && *csv == "yes":
			fmt.Fprintf(w, "%s,%d,%s\n", *obj.Key, *obj.Size, *obj.ETag)
		case *search != "" && *csv == "yes":
			if caseInsensitiveContains(*obj.Key, *search) == true {
				fmt.Fprintf(w, "%s,%d,%s\n", *obj.Key, *obj.Size, *obj.ETag)
			} else {
				continue
			}
		}
	}
	w.Flush()
}
