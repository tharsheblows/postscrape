# postscrape - scrape http POST responses

PostScrape lets you scrape targeted html nodes in enumerateable http POST requests. It uses [goquery](https://github.com/PuerkitoBio/goquery) for the targeting. Currently it's command line only with the following inputs:

* _Enter output filename without extension -- eg out rather than out.txt:_ The output files will be built on what you enter here and have the form output-0.txt, output-1.txt etc where "output" is the bit you've entered
* _Enter the key to iterate (eg "page" or "paged"):_ What is the iteratable key to use? eg page=1, page=2, etc.
* _Enter number of pages in total to scan (just press enter if this is irrelevant):_ If it's only one you can skip this. Otherwise, put in the total number of pages to scrape.
* _Limit how hard you hit the other server and how much is in each output file by indicating how many pages to get at once. How many concurrent pages?_ How many pages to get at once? I usually do a max of 20.
* _Now add in the query args one by one in the form "key=value" or you can add them all at once like "key=value&key2=value2". When you're done, type "done" If there's a page argument, it'll be added to the query as page=__ whatever the number is, you don't need to add it here:_ You can either add each key/value pair on a separate line or add them all at once in the form of a query string. Leave out the iterable key. Type "done" to finish this bit.
* _Enter the full url (without any arguments):_ Make sure it starts with http or https.
* _Enter what you'd like to get from the page (jQuery selector):_ This makes further manipulation a bit easier. You could use "html" if you wanted to get all of the page I guess. This is my compromise between just getting the data I think I need and getting everything.

## License

The [MIT license](LICENSE). [goquery's license is BSD 3-Clause license](https://github.com/PuerkitoBio/goquery).
