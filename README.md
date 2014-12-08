# SimpleReader
Simple reader for fb2 online

Based on martini server:
[Martini](https://github.com/go-martini)

Test site:
[SimpleReader](http://www.borscht-yourok.rhcloud.com)


## Simple Start
In linux terminal:

Clone source code
~~~
git clone https://github.com/YouROK/SimpleReader.git
~~~

Set environment
~~~
cd ./SimpleReader
export GOPATH=`pwd`
~~~

Get dependents
~~~
go get ./src/main/
~~~

Install and run
~~~
go install ./src/main/
./bin/main
~~~

Server run in port 9000, for test open in browser [http://localhost:9000/](http://localhost:9000/)
