// This gets pages from a site which can be iterated over either using a POST request
// I use it to get most of the page, in particular the parent div of the items I really want (eg I'll get the ul when I want the li)
// The impetus was having a local copy makes realising I've forgotten something I've needed when I've done a targeted scrape much less painful

package main

import (
	"bufio" // for the Reader / Writer
	"fmt"
	"math"
	"net/http" // to do the GET / POST request
	"net/url"  // to get the POST form in the proper format
	"os"       // Create output file
	"strconv"  // convert int to string
	"strings"  // string manipulation

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color" // color the command line
)

// color functions
var yellowf = color.New(color.FgYellow).PrintfFunc()
var yellowln = color.New(color.FgYellow).PrintlnFunc()
var greenf = color.New(color.FgGreen).PrintfFunc()
var greenln = color.New(color.FgGreen).PrintlnFunc()

// Used to check the error. I forget when I got this, maybe something other than panic is better?
func check(e error, m string) {
	color.Set(color.FgHiRed)
	if e != nil {
		fmt.Println(m) // print a message
		panic(e)       // panic
	}
	color.Unset()
}

// This scrapes the url, gets the relevant selector's content and passes it to the ch channel
func scrape(u string, p string, selector string, ch chan string, chFinished chan bool) {

	fmt.Printf("Scraping url with POST args %s \n", p)

	// I think there's probably a better way to get the query in proper format -- I am using this for POST requests
	uv, err := url.ParseQuery(p)
	check(err, "ERROR: Could not do url.ParseQuery for \""+p+"\"")

	// POST to the desired url
	resp, err := http.PostForm(
		u,
		uv)
	em := "ERROR: Failed to scrape \"" + u + "\""
	check(err, em)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	// get the response ready for goquery
	doc, err := goquery.NewDocumentFromResponse(resp)
	check(err, "ERROR: There was an issue making the doc. \n")

	// get the html in the selector from the doc
	text, err := doc.Find(selector).Html()
	check(err, "ERROR: There was an issue selection " + selector + " from the doc. \n")

	// send it to the channel
	ch <- text
}

// This is the post data input
func addPostData(p []string) []string {

	// bufio manages getting std input
	getInput := bufio.NewReader(os.Stdin)

	// read the line
	o, _ := getInput.ReadString('\n')

	// take off the newline
	o = strings.Replace(o, "\n", "", -1)

	// if you're done, quit
	if o == "done" {
		return p
	}
	// add in the query to the array
	p = append(p, o)

	fmt.Println(" ")
	greenln("Please put in another query arg, either one at a time or as a query string")
	greenf("(eg key=value or key=value&key2=value2) -- or if you're done, enter \"done\": ")

	// add in another key/value pair or be done
	p = addPostData(p)

	return p // I mean, it should never get here but it won't run without this, I must be doing something wrong
}

func outputFiles(f string) (fp *os.File) {
	f = f + ".txt"
	// create the file
	fp, err := os.Create(f)
	check(err, "ERROR: The file wasn't created properly.")

	// say that the file has been created
	yellowln("Created file " + f)

	return fp
}

func main() {

	// new Reader for std input
	getInput := bufio.NewReader(os.Stdin)

	// this will be the path to the output file
	greenf("Enter output filename without extension -- eg out rather than out.txt: ")
	o, _ := getInput.ReadString('\n')

	// What's the iterable key?
	greenf("Enter the key to iterate (eg \"page\" or \"paged\"): ")
	itKey, _ := getInput.ReadString('\n')
	itKey = strings.Replace(itKey, "\n", "", -1)

	yellowln("Will iterate over " + itKey)

	// how many pages should it scan?
	greenf("Enter number of pages in total to scan (just press enter if this is irrelevant): ")
	numPages, _ := getInput.ReadString('\n')
	numPages = strings.Replace(numPages, "\n", "", -1)

	yellowln("Will grab " + numPages + " pages")

	n, e := strconv.Atoi(numPages)
	check(e, "Problem with int to string converstion. \n")

	// how many pages should each output file contain?
	greenln("Limit how hard you hit the other server and how much is in each output file by indicating how many pages to get at once.")
	greenf("How many concurrent pages? ")
	numConcPages, _ := getInput.ReadString('\n')
	numConcPages = strings.Replace(numConcPages, "\n", "", -1)

	yellowln("Will grab " + numConcPages + " pages at once")

	nc, e := strconv.Atoi(numConcPages)
	check(e, "Problem with int to string converstion. \n")

	// the slice with the post data
	var pd []string

	fmt.Println(" ")
	greenln("Now add in the query args one by one in the form \"key=value\" or you can add them all at once like \"key=value&key2=value2\". When you're done, type \"done\"")
	greenln("If there's a page argument, it'll be added to the query as page=__ whatever the number is, you don't need to add it here:")

	// add the post data
	pd = addPostData(pd)
	yellowln("Here's the post data: ")
	fmt.Println(pd)

	// which url should it use?
	greenf("Enter the full url (without any arguments): ")
	u, _ := getInput.ReadString('\n')
	u = strings.Replace(u, "\n", "", -1)

	// what should we grab from the page?
	greenf("Enter what you'd like to get from the page (jQuery selector): ")
	s, _ := getInput.ReadString('\n')
	s = strings.Replace(s, "\n", "", -1)

	// Channels
	chContent := make(chan string)
	chFinished := make(chan bool)

	// total number of output files
	tn := int(math.Ceil(float64(n) / float64(nc)))
	var p string
	var content string
	var foundContent []string
	var max int

	// so
	for count := 0; count < tn; count++ {
		chCount := nc // the channel is going to hold nc entries until the last iteration
		if (count+1)*nc > n {
			max = n                    // don't go over the number of pages you want to get
			chCount = n - nc*(count+1) // on the last iteration the chCount will be the number of pages you want minus the number of pages already gotten (n*(count+1))
		} else {
			max = (count + 1) * nc // the max, except for the last iteration maybe, is this
		}
		// Kick off the scrape process (concurrently)
		for i := count * nc; i < max; i++ {
			page := strconv.Itoa(i + 1)
			p = strings.Join(pd, "&")
			p = p + "&" + itKey + "=" + page
			p = string(p)
			go scrape(u, p, s, chContent, chFinished)
		}

		foundContent = nil
		content = ""

		// Subscribe to both channels
		for c := 0; c < chCount; {
			select {
			case content = <-chContent:
				foundContent = append(foundContent, content)
			case <-chFinished:
				c++
			}
		}

		// We're done! Print the results...

		of := strings.TrimSpace(o) + "-" + strconv.Itoa(count) // I'm going to name the files eg output-0, output-1, output-2 ....
		f := outputFiles(of)
		defer f.Close()
		w := bufio.NewWriter(f)

		for i := range foundContent {
			n4, err := w.WriteString(foundContent[i] + "\n")
			if err != nil {
				color.Red("failed writing: " + foundContent[i] + "\n")
			}
			yellowf("wrote %d bytes\n", n4)
			// Use Flush to ensure all buffered operations have been applied to the underlying writer.
			w.Flush()
		}
	}

	close(chContent)
}
