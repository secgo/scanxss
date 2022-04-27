package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var LinksParam = []string{}
var wg sync.WaitGroup
var domain string

func main() {
	flag.StringVar(&domain, "u", "", "add url http://testphp.vulnweb.com/")
	flag.Parse()
	if domain == "" {
		fmt.Println(`
		-u enter domain Hint: [http://testphp.vulnweb.com]
		`)
		os.Exit(0)
	}
	ExtractLink("http://testphp.vulnweb.com/")
	fmt.Println("please wite...")
	time.Sleep(time.Second * 2)
	fmt.Printf("\n\n\n\n")
	for _, v := range LinksParam {
		// fmt.Println(v)
		wg.Add(1)
		go xssCheck(v)
	}
	wg.Wait()

}

func xssCheck(u string) {
	re := regexp.MustCompile(`Error|mysql|Unknown column`)
	defer wg.Done()
	f, err := os.Open("payload.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ur := changeParam(u, scanner.Text())
		req, err := http.NewRequest(http.MethodGet, ur, nil)
		if err != nil {
			log.Println(err)
		}
		Client := http.Client{}
		res, err := Client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		bd, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}
		if strings.Contains(string(bd), scanner.Text()) && !re.MatchString(string(bd)) {
			fmt.Printf("%sFound xss -> %s%s\n", Cyan, ur, Reset)
		} else if re.MatchString(string(bd)) {
			fmt.Printf("%sProbably sqli -> %s%s\n", Red, ur, Reset)

		} else {
			fmt.Println("Not Found")
		}
	} // end scanner
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func changeParam(u string, p string) string {
	ur, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	values := ur.Query()
	for v := range values {
		values.Del(v)
		values.Add(v, p)
	}
	ur.RawQuery = values.Encode()
	return ur.String()
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}
