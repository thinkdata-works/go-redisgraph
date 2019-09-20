# redisgraph-go-query - Redisgraph Go Query client

Support for querying redisgraph over the redigo client. This is derived from the [redisgraph-go](https://github.com/RedisGraph/redisgraph-go) client.

This library focuses on adding support for typing results and for now just supports simple scalar types.

## Building

```bash
$ docker-compose up -d
```

## Running tests

```bash
$ docker-compose run driver gotest -v ./rgraphquery
```

This library is set up with docker to ensure that the client will run against specific versions of redisgraph.

To check support against a new version, simply replace (or add a new) version to the docker compose file and run the tests against it.


# Reporting Issues

Please create a ticket in the Issues page for this repo, or reach out to us directly (see *Get in Touch*). 

# Contributing

Please see the Issues page for list of outstanding bugs. If you are looking to implement a feature outside of the list of bugs
please create an Issue for the feature that you would like to implement.

Assign it to yourself or comment that you will be taking it over. Then:

1. Fork this repo
2. Clone your repo and check out a branch based off of development. The name of the branch should be `feature/<ticket number>`
3. Implement your feature. Ensure that
    1. You have updated any version changes in the redisgraph image
    2. All existing tests are passing 
    3. You have added any necessary documentation
4. Open a pull request from your fork to our repo on the development branch. 
    1. Title: `feature/<ticket number> <General summary of work done>`
    2. Description: `Fixes #<ticket number>. <Additional details of work>`

# Get in touch

Please reach out to us at anytime at `support@namara.io`