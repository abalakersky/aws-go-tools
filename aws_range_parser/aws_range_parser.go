package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

var (
	region   = flag.String("region", "", "Region to check on.")
	service  = flag.String("service", "", "Service to look up information for")
	url      = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	services = []string{}
	regions  = []string{}
)

type awsEntry struct {
	Prefixes []awsPrefixes `json:"prefixes"`
}

type awsPrefixes struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}

func u() {
	fmt.Printf("This script is used to display AWS specific IP ranges that could be used for Firewall or Security Group configurations. These ranges specify public IPs that AWS uses for a each public facing service.\n")
	fmt.Printf("\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("aws_range_parser [-h] [-region REGION_NAME] -service SERVICE_NAME\n")
	fmt.Printf("\n")
	fmt.Printf("Service:\n")
	fmt.Printf("    Valid values: %v\n", services)
	fmt.Printf("\n")
	fmt.Printf("Region:\n")
	fmt.Printf("    Valid values: %v\n", regions)
	fmt.Printf("\n")
	fmt.Printf("Notes:\n")
	fmt.Printf("Please remember that some services, such as CloudFront and Route53 are Global and as such use only GLOBAL as their region. Their information can be gathered with or without specifying region name.\n")
}

func main() {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		var a awsEntry
		err = json.Unmarshal(contents, &a)

		if err != nil {
			log.Fatal(err)
		}
		var s []string
		var r []string
		for i := range a.Prefixes {
			s = append(s, a.Prefixes[i].Service)
			r = append(r, a.Prefixes[i].Region)
		}
		services = removeDuplicatesUnordered(s)
		sort.Strings(services)
		regions = removeDuplicatesUnordered(r)
		sort.Strings(regions)

		flag.Usage = u
		flag.Parse()
		if flag.NFlag() == 0 {
			u()
		}

		for i := range a.Prefixes {
			switch {
			case *service != "" && *region != "":
				if a.Prefixes[i].Service == *service && a.Prefixes[i].Region == *region {
					fmt.Printf("%s\n", a.Prefixes[i].IPPrefix)
				}
			case *service != "" && *region == "":
				if a.Prefixes[i].Service == *service {
					fmt.Printf("%s\n", a.Prefixes[i].IPPrefix)
				}
			}
		}
		if *service == "" && *region != "" {
			fmt.Printf("\nService is a required value\n\n")
		}

	}
}
