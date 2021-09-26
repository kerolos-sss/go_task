# Page Details

## How to run

after installing the go environment 
make sure you are using go v 1.17
```
go run main.go
```
- a webpage will open 
- supply the needed url to get info about
- press submit button 
- wait for the response


## includes

a web application which takes a website URL as an input and provides general information
about the contents of the page:
- HTML Version
- Page Title
- Headings count by level
- Amount of internal and external links
- Amount of inaccessible links
- If a page contains a login form

## assumptions

- considered the DOCTYPE tag text as the html version
- elements will be detected only on the page initial html response.
- links only means 'a' tags.
- generated on user interaction eg( paged 'a' tags, etc. ... ) will not be detected/collected.
- links will be gathered exactly so we may find the same link but ending by a '/' as different links
- inaccessibility is determined by an html error codes only
- I will be using a simple way of detection for the login form that can make false positives and negatives
    - https://stackoverflow.com/questions/14975000/how-can-i-detect-a-login-form-in-a-webpage-using-javascript


## testing extra commands

```
go test ./scraper -v
```
starting again Sum 11:43 am
