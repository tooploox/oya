# Oya HOWTO
Examples : Web server
We have a simple http server written in Go.


    $ tree
    .
    ├── public
    │   └── index.html
    └── app.go
    
    1 directory, 2 files
    
    $ cat app.go
    package main
    
    import (
            "fmt"
            "net/http"
            "os"
    )
    
    func main() {
            port := os.Getenv("PORT")
            if port == "" {
                    port = "80"
            }
            fileServer := http.FileServer(http.Dir("public/"))
            http.Handle("/", fileServer)
    
            fmt.Printf("Starting server :%v\n", port)
            http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
    }
    

Oya can help with maintaining this project.


# 1. Oya helps with project build and run

Final code can be found [here]( https://github.com/tooploox/oya-example/tree/master/sample)

1.1 Install Oya
1.2 Project
1.3 Tasks
1.4 Values
1.5 Render
1.6 .oya files
1.7 Secrets


## 1.1 Install Oya


    $ curl https://oya.sh/get | bash

This will install latest version of oya. You can say which version you want with `$ curl https://oya.sh/get | bash -s v0.0.7` for example. For release versions check [releases](https://github.com/tooploox/oya/releases) .

Now we can setup Oya Project.

## 1.2 Project

Typically, a project will contain an `Oyafile` in its root directory. An `Oyafile` is a YAML file containing Oya configuration and the  `Project:` directive to indicate that this is the project root. You can create Oyafile by hand or with init command, let’s do this:


    $ oya init
    $ cat Oyafile
    Project: project

Let’s rename our project to `OyaExample`


## 1.3 Tasks

Oya task is a bash script defined as a Oyafile key.
Tasks in Oyafile are defined as a yaml keys, with pipe at the beginning line and bash code in following.
Our server is golang app so before starting it needs to be compiled. 
So lets add `build` and `start` tasks.


    $ cat Oyafile
    Project: OyaExample
    
    build: |
      go build app.go
    
    start: |
      ./server

To list available tasks use tasks command:


    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run start

Execute task


    oya run build
    oya run start


## 1.4 Values

Values are defined in `Oyafile` as a list under uppercase key `Values`. Values can be accessed from tasks with `${Oya[value_name]}`.
As a default our server starts at 80 port, but it might change. Let’s keep port as a Value in `Oyafile` to easier management.


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      port: 3000
    
    build: |
      go build server.go
    
    start: |
      export PORT=${Oya[port]}
      ./server


## 1.5 Render

Oya can also render a template files or even whole directory. Oya uses Plush templating system You can find out more [here](https://github.com/gobuffalo/plush).

We can generate `index.html` for our http server. We need a template file.


    $ cat template/index.html
    <h1><%= title %></h1>

Let’s add our title under  `Value:`


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      port: 3000
      title: Hello from Oya
    ...

Render output into `public/` dir so our server can see it.


    $ oya render template/index.html -o public
    $ cat public/index.html
    <h1>Hello from Oya</h1>

Let’s add a task `generate` to do it for us.


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      port: 3000
      title: Hello from Oya
    
    build: |
      go build app.go
    
    start: |
      export PORT=${Oya[port]}
      ./app
    
    generate: |
      oya render template/index.html -o public

`$ oya run generate` will generate `index.html` for us.


## 1.6 .oya files

Inside a project you can have many files with named `*.oya` they will be read as a Value files, expected syntax is a pair of  `key: value`. 

One exceptions is `secrets.oya` which needs to be encrypted before use. 

## 1.7 Secrets

Ok, but what if i have some confidential values ?
Oya uses [SOPS](https://github.com/mozilla/sops) to help with secrets management. 
First you need to configure SOPS for [encryption methods](https://github.com/mozilla/sops/blob/master/README.rst#usage).
For our example we can use PGP key.

    $ export SOPS_PGP_FP="317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202"

Oya secrets commands:

    $ oya secrets --help
    ...
      edit        Edit secrets file
      encrypt     Encrypt secrets file
      view        View secrets
    ...

**First run:**

First you need to create `secrets.oya` file. with `key: value` in each line. And encrypt it. (You can also go straight to edit with `$ oya secrets edit secrets.oya`).


    $ cat secrets.oya
    magical_spell: hokus pokus czary mary
    $ oya secrets encrypt secrets.oya

Done your secrets are safe. Check how secrets.oya looks like.


    $ cat secrets.oya
    {
            "data": "ENC[AES256_GCM,data:XXXX=,tag:XXXX==,type:str]",
            "sops": {
                    ...
                    "pgp": [...],
                    ...
            }
    }%

We can see only sops metadata, our data are safe and encrypted.

To view or edit use:


    $ oya secrets view secrets.oya
    magical_spell: hokus pokus czary mary
    $ oya secrets edit secrets.oya

You can access secret value with  `${Oya[magical_spell]}` from task.


    $ cat Oyafile
    Project: OyaExample
    ...
    spell: |
      echo ${Oya[magical_spell]}
    
    $ oya run spell
    hokus pokus czary mary


# 2. Project structure - Recursive Oya.

[Example code](https://github.com/tooploox/oya-example/tree/master/recurse)

Ok, but what if i want to create frontend APP in other language, how can we separate it.
We have a golang server and http files in one place. In Oyafile we’re building golang and generating html. It is not a best practice. 
We can separate it into `backend/` and `frontend/`. 


    $ tree
    .
    ├── Oyafile
    ├── backend
    │   ├── Oyafile
    │   └── app.go
    └── frontend
        ├── Oyafile
        └── template
            └── index.html
    
    5 directories, 7 files

Each of them will have own `Oyafile`, and thanks to recursive tasks each file can have task with the same name in our case `build`.

## Backend

Here is how our backend Oyafile looks like (note that there is no `Project:` for subdirectories):


    $ cat ./backend/Oyafile
    build: |
      echo "Compiling server"
      go build -o ../build/server app.go


## Frontend


    $ cat ./frontend/Oyafile
    Values:
      title: Hello from Oya
    
    build: |
      echo "Rendering index.html"
      oya render template/index.html -o ../build/public

Our frontend holds only template file and Values necessary to render it.


## Project Oyafile


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      port: 3000
    
    build: |
      echo "Preparing build/ folder"
      rm -rf build/ && mkdir -p build/public
    
    start: |
      export PORT=${Oya[port]}
      cd ./build/ && ./server


## Recursive run

Now let’s see what tasks we have. To do it for whole project including subdirectories we need to use `-r or --recurse`  flag.


    $ oya tasks -r
    # in ./Oyafile
    oya run build
    oya run start
    
    # in ./backend/Oyafile
    oya run build
    
    # in ./frontend/Oyafile
    oya run build

As you can see we have three `build` tasks one per Oyafile. We can now run them all.


    $ oya run -r build
    Preparing build/ folder
    Compiling server
    Rendering index.html

And now we can start the app with `$ oya run start`.

# 3. Packs how to use them
## Import pack

Pack is a Oya project with general purpose tasks which can be easily shared and used inside other projects. Oya installs pack in your home `~/.oya` directory. Each time you rune oya command dependencies will be resolved and installed.

For better portability our server can be run in Docker container. We can use Docker pack for it. To do it first we need to import it to our project.


    $ oya import github.com/tooploox/oya-packs/docker
    $ cat Oyafile
    Project: OyaExample
    Require:
      github.com/tooploox/oya-packs/docker: v0.0.6
    Import:
      docker: github.com/tooploox/oya-packs/docker
    ...

Import will add importing pack under `Import:`, key of imported pack is his alias and can be accessed by this name, (you can change it if needed). 

**Packs versioning**

Import will automatically resolve dependencies with newest versions and add them under `Require:`.

Imported pack added bunch of new commands into our project


    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run docker.build
    oya run docker.generate
    oya run docker.run
    oya run docker.stop
    oya run docker.version
    oya run start

We can easily generate Dockerfile, update it, and build our project.


# 4. Pack development - sharing oya’s

Each oya project is a pack itself, all you need to do is push it to git and tag it version `name/v0.0.1`. Import it as before with `$ oya import github.com/tooploox/oya-packs/name` oya should automatically resolve newest version and add `Require` to your oyafile.

…. 😕 need decent example…. from life … something what everyone does all over again …. for each setup … hmmm.. help

