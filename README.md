Comment build le projet : 

* `git clone git@git.vayan.fr:goirc.git`
* `cd goirc`
* `mkdir bin`
* rajouter le dossier bin crée precedement a votre PATH (ex : `export PATH=$PATH:/path/to/goirc/bin` )
* Set un env GOPATH qui link à la racine du projet  (ex: `expot GOPATH=D:\My Documents\Git\goric` )
* `cd src/goirc/`
* `go install`
* `goirc`
* et localhost:1111 dans votre navigateur