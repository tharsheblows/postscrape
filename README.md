# postscrape - scrape http POST responses

PostScrape lets you scrape targeted html nodes in enumerateable http POST requests. It uses [goquery](https://github.com/PuerkitoBio/goquery) for the targeting. Currently it's command line only with the following inputs:

* **Enter output filename without extension -- eg out rather than out.txt:** The output files will be built on what you enter here and have the form output-0.txt, output-1.txt etc where "output" is the bit you've entered
* **Enter the key to iterate (eg "page" or "paged"):** What is the iteratable key to use? eg page=1, page=2, etc.
* **Enter number of pages in total to scan (just press enter if this is irrelevant):** If it's only one you can skip this. Otherwise, put in the total number of pages to scrape.
* **Limit how hard you hit the other server and how much is in each output file by indicating how many pages to get at once. How many concurrent pages?** How many pages to get at once? I usually do a max of 20.
* **Now add in the query args one by one in the form "key=value" or you can add them all at once like "key=value&key2=value2". When you're done, type "done" If there's a page argument, it'll be added to the query as page=__ whatever the number is, you don't need to add it here:** You can either add each key/value pair on a separate line or add them all at once in the form of a query string. Leave out the iterable key. Type "done" to finish this bit.
* **Enter the full url (without any arguments):** Make sure it starts with http or https.
* **Enter what you'd like to get from the page (jQuery selector):** This makes further manipulation a bit easier. You could use "html" if you wanted to get all of the page I guess. This is my compromise between just getting the data I think I need and getting everything.

## License

The [MIT license](LICENSE). [goquery's license is BSD 3-Clause license](https://github.com/PuerkitoBio/goquery).
