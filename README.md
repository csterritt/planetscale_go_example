# PlanetScale Go web server example

This is an example program, written in go, which serves as a very simple web
server which accesses the [PlanetScale](https://www.planetscale.com) database
service. It is run under a docker container so I can deploy it easily.

To run it, install the `pscale` PlanetScale cli tool, the current go compiler
(I'm using go1.16.4) and do the following:

First, create a database on PlanetScale and give it a name. Mine is named
`firstexample`.

First, you can try running this directly.

Run:

~~~
go build
~~~

to create the `ps_ws_ex` executable, then run:

~~~
PORT=4000 pscale connect firstexample main --execute ./ps_ws_ex
~~~

You can then use (e.g.) the `curl` command line tool to hit the service from another
terminal window:

~~~
curl -X POST http://localhost:4000/get-reminder
~~~

----

Next, we'd like to run this in production, for which we'll need a service token.

Stop the program started above.

To generate the service token from the command line, run:

~~~
pscale service-token create
~~~

This generates a TOKEN_NAME and a TOKEN_VALUE.

Next, run:

~~~
pscale service-token add-access TOKEN_NAME connect_production_branch --database firstexample
~~~

Next, update the `Dockerfile` in this directory, and put the TOKEN_NAME and TOKEN_VALUE
into lines 30 and 31 respectively.

Finally, run:

~~~
docker build -t pscale_test .
docker run --rm -w /app -p 4000:8080 pscale_test
~~~
