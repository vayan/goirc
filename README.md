
#last commit

* better readme
* TODO formating
* Move import JS


Comment build le projet : 

* Set un env GOPATH qui link à la racine du projet  
 (ex: `export GOPATH=D:\My Documents\Git\go` )

* dans votre dossier de projet 
* `mkdir bin`
* `mkdir src`
* `mkdir pkg`

* rajouter le dossier bin crée precedement a votre PATH 
 (ex : `export PATH=$PATH:/path/to/go/bin` )

* `go get code.google.com/p/go.net/websocket`
* `go get github.com/gorilla/schema`
* `go get github.com/gorilla/mux`
* `go get github.com/thoj/go-ircevent`

* `cd src`
* `git clone git@bitbucket.org:vayan/goirc.git`

* `cd goirc`
* `go install goirc`
* `goirc`
* et localhost:port dans votre navigateur