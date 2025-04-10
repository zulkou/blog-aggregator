# Blog Aggregator
## Description
A go based app to scrape posts of given blog feeds. This project is one of the guided projects in [boot.dev](https://www.boot.dev/courses/build-static-site-generator-python).

## How to Use This Project
To run this project, you will need to install [go](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/).
After both app installed you can clone this project locally and run it with provided shell below.
```bash
$ git clone https://github.com/zulkou/blog-aggregator.git
$ cd blog-aggregator
```
If you want to install this project locally, you can use `go install <pkgname>`.
```bash
$ go install github.com/zulkou/blog-aggregator
```
### Available Commands
```bash
$ blog-aggregator register <usrname>        # registered user is auto logged in
$ blog-aggregator login <usrname>
$ blog-aggregator reset                     # reset database
$ blog-aggregator users                     # list available users
$ blog-aggregator addfeed <url> <feedname>  # need to be logged in to add new feed
$ blog-aggregator feeds                     # list all available feeds
$ blog-aggregator agg <interval>            # will start scraping at given interval
$ blog-aggregator follow <url>              # current user will follow feed with given url
$ blog-aggregator following                 # list all feeds current user following
$ blog-aggregator unfollow <url>            # current user will unfollow feed with given url
$ blog-aggregator browse <limit>            # will list posts from followed feeds with given limit
```
