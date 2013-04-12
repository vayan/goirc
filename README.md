# GO#IRC

- - -

## Build Requirements

* Go                     <[http://golang.org](http://golang.org)>
* websocket              <[code.google.com/p/go.net/websocket](code.google.com/p/go.net/websocket)>
* Gorilla                <[http://www.gorillatoolkit.org/](http://www.gorillatoolkit.org/)>
* go-ircevent            <[github.com/thoj/go-ircevent](github.com/thoj/go-ircevent)>
* Go-SQL-Driver/MySQL    <[github.com/Go-SQL-Driver/MySQL](github.com/Go-SQL-Driver/MySQL)>

* MySQL                  <[https://www.mysql.com/](https://www.mysql.com/)>


## Comment setup le projet :

* Set un env GOPATH qui link à la racine du projet
 (ex: `export GOPATH=/path/to/go/` )

* dans votre dossier de projet :

```
#!shell

mkdir bin
mkdir src
mkdir pkg
```

* rajouter le dossier bin crée precedement a votre PATH
 (ex : `export PATH=$PATH:/path/to/go/bin` )

* Recuperer les dependances :

```
#!shell

go get code.google.com/p/go.net/websocket
go get bitbucket.org/vayan/gomin
go get github.com/gorilla/schema
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/thoj/go-ircevent
go get github.com/Go-SQL-Driver/MySQL
cd src
git clone git@bitbucket.org:vayan/goirc.git
cd goirc
```

* Creer le fichier de config `conf.json` (exemple dans `conf.json.exemple`)

* Done !

## Comment lancer GO#IRC :

```
#!shell

go install goirc/goirc
goirc
```

Puis ouvrez un navigateur localhost:port (localhost:1112 par default)