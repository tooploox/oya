# Oya

## Usage

1. Install oya and its dependencies:

        curl https://raw.githubusercontent/bilus/oya/master/scripts/setup.sh | sh

1. Initialize project to use a certain CI/CD tool and workflow. Example:

        oya init jenkins-monorepo

   It boostraps configuration for Jenkins pipelines supporting the 1.a workflow
   (see Workflows below), an Oyafile and supporting scripts and compatible
   generators.

1. Define a task you can run:

        mkdir app1
        cat > Oyafile
        build: echo "Hello, world"

1. Run the task:

        oya run build
        Hello, world

The task in the above example is called "build" but there's nothing special about the name. In fact, a task name is any camel case identifier as long as it starts with a lower-case letter. You can have as many different tasks as you like.

## Plugins

oya vendor p github.com/bilus/oya/packs/circleci-helm-platform

installs into vendor/
symlinks vendor/p to it
symlinks all in vendor/p/bin/ to oya/bin/p

cd delivery/broadcasts/Oyafile
oya p/generate/docker

delivery/broadcasts/Oyafile

Import:
  - github.com/bilus/oya/packs/jenkins-monorepo

--

Path:
  jm: github.com/bilus/oya/packs/jenkins-monorepo/bin

buildDocker:
  jm/buildDocker

buildChart:
  jm/buildChart


## How it works

A directory is included in the build process if contains an Oyafile regardless of how deeply nested it is. You can use Oyafiles in these directories to define their own tasks.

For example, to set up a CI/CD pipeline in a mono-repository containing several
microservices, you'd put each microservice in its directory, each with its own
Oyafile containing the tasks necessary to support the CI/CD workflow.


Imagine you have the following file structure:

```yaml
# ./Oyafile

build: |
  echo "Top-level directory"
```

```yaml
# ./subdir/Oyafile

build: |
  echo "Sub-directory"
```

When you run `oya run build`, Oya first walks the directory tree, starting from
the current directory, to build the **changeset**: the list of directories that
are marked as changed. In the above example it would be, as you probably
guessed, `.` (the top-level directory) and `subdir` (the sub-directory).

Finally, Oya executes the task you specified for every directory marked as
changed, starting from the top directory. Going back to our example, it would
generate the following output:

```
Top-level directory
Sub-directory
```

As you say, tasks and their corresponding scripts are defined in `Oyafile`s.
Their names must be camel-case yaml identifiers, starting with a lower-case
letter. Built-in tasks start with capital letters.

More realistic example of an `Oyafile`:

```
build: docker build .
test: pytest
```

## Changesets

TODO

   * `Changeset` -- (optional) modifies the current changeset (see Changesets).

Oya first walks all directories to build the changeset: a list of directories
containing an Oyafile that are marked as "changed".

It then walks the list,
running the matching task in each. CI/CD tool-specific script outputting list of
modified files in buildable directories given the current task name.
     - each path must be normalized and prefixed with `+`
     - cannot be overriden, only valid for top-level Oyafile
     - in the future, you'll be able to override for a buildable directory and
       use `-` to exclude directories, `+` to include additional ones, and use
       wildcards, this will allow e.g. forcing running tests for all apps when
       you change a shared directory
     - git diff --name-only origin/master...$branch
     - https://dzone.com/articles/build-test-and-deploy-apps-independently-from-a-mo
     - https://stackoverflow.com/questions/6260383/how-to-get-list-of-changed-files-since-last-build-in-jenkins-hudson/9473207#9473207

Generation of the changeset is controlled by the optional changeset key in
Oyafiles, which can point to a script executed to generate the changeset:

1. No directive -- includes all directories containing on Oyafile.
2. Directive pointing to a script.

.oyaignore lists files whose changes do not trigger build for the containing
buildable directory

## Features/ideas

1. Generators based on packs. https://github.com/Flaque/thaum + draft pack
   plugin

## Workflows

### Repo structure

1. Mono-repo:
   - Each app has its own directory
   - There is a directory/file containing deployment configuration

2. Multi-repo:
   - Each app has its own repo
   - Deployment configurations in its own repo

3. Mix:
   - Some/all apps share repos, some may have their own
   - Deployment configurations in its own repo

> Also submodules tried for NT/Switchboard and eventually ditched.

### Change control

a. Each environment has its own directory

b. Each environment has its own branch

### Evaluation

| Workflow | Projects    | Pros                                       | Cons                                                    |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 1.a      | E           | "Can share code"                           | Merge order dependent [1]                               |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 1.b      | C           | Single checkout                            | Complex automation [2]                                  |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 2.a      |             | Same as 2.b                                | Same as 2.b plus need to detect which directory changed |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 2.b      | S           | Better isolation [3] Simple automation [4] | More process overhead [5]                               |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 3.a      | P           | Can divide up a project however you like   | Complex automation [2]                                  |
|----------|-------------|--------------------------------------------|---------------------------------------------------------|
| 3.b      | P           | Simple deployment automation               | Same as 3.a                                             |

* [1] Code gets merged from branch to branch; works for small team.
* [2] Need to detect what changed between commits. Many CI/CD tools allow only
  one configuration per repo and require coding around the limitations, example:
  https://discuss.circleci.com/t/does-circleci-2-0-work-with-monorepos/10378/13
* [3] No way to just share code, need to package into libraries. Great for
  microservices and must have for large teams.
* [4] Just put a CI/CD config into the root.
* [5] No way to just share code, need to package into libraries. Bad for small
  teams wanting to quickly prototype.
