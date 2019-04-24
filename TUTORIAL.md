# Oya tutorial

## What is Oya

Oya lets you bootstrap CI/CD for your project using ready-made packs supporting various workflows and CI/CD tools.

## Install Oya and its dependencies

    curl https://raw.githubusercontent/bilus/oya/master/scripts/setup.sh | bash

## Our application

The application we will work with is a Hello world of HTTP servers written in Golang.

To get you quickly started, simply clone the following repository: TODO

It contains `hello/main.go` with the source code of the HTTP server.

## Deploying it to GKE

I'll show you how to deploy the Hello world app as a Docker container running on the Google Kubernetes Engine (TODO: One-sentence explanation).

The first step is initializing the repository:

    oya init

All it does is generate an empty Oyafile in the root directory.

Let's dockerize our app using an Oya pack that makes it easy. First, install the pack:

    oya get github.com/tooploox/oya/packs/docker [--alias docker]

Then, generate the assets necessary to put the application into a Docker image:

    oya generate docker hello

The command invokes an Oya generator named "docker" passing it the path to the Hello world application as an argument. The command will generate a few files:

    hello/
       Oyafile
       Dockerfile

Oyafile contains Oya configuration, Dockerfile contains the instructions of how to build the Docker image.

> In general syntax of the `generate` Oya command is: `oya generate <generator> <args>`. You can list the available generators like so: `oya generate list`.

Let's build the image:

    oya run docker.build

    oya docker build

> If you run this command in the root directory, it'll build docker images for all applications you used the `oya generate docker` command on. In our case, it's just the Hello world app. You can also run it in a specific app's directory.

HERE NEXT

We need to tell Oya this is what we want to do:

    oya get github.com/tooploox/oya/packs/gke

The command installs the gke pack. We'll also need the docker pack to simplify putting our app in a docker container:

    oya get github.com/tooploox/oya/packs/docker

Let's generate code necessary to deploy our app in a docker container:

    cd hello
    oya generate g

## Adding CI/CD
