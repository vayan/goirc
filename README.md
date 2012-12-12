Comment build le projet : 

* Set un env GOPATH qui link à la racine du projet  (ex: `export GOPATH=D:\My Documents\Git\goric` )

* dans votre dossier de projet 
** `mkdir bin`
** `mkdir src`
** `mkdir pkg`

* rajouter le dossier bin crée precedement a votre PATH (ex : `export PATH=$PATH:/path/to/goirc/bin` )

* `go get code.google.com/p/go.net/websocket`
* `go get github.com/gorilla/schema`
* `go get github.com/gorilla/mux`
* `go get github.com/thoj/go-ircevent`

* `cd src`
* `git clone git@git.vayan.fr:goirc.git`

* `cd src/goirc`
* `go install goirc`
* `goirc`
* et localhost:1111 dans votre navigateur