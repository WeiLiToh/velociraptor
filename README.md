## Current Workflow for custom plugin validation on Velociraptor
1. Build from source on Linux
2. go to velociraptor/output, run ./vecociraptor gui
3. go to notebook and run SELECT * FROM info() WHERE Name = "ollama"

# Velociraptor - Endpoint visibility and collection tool.

Velociraptor is a tool for collecting host based state information
using The Velociraptor Query Language (VQL) queries.

To learn more about Velociraptor, read the documentation on:

https://docs.velociraptor.app/

## Quick start

If you want to see what Velociraptor is all about simply:

1. Download the binary from the release page for your favorite platform (Windows/Linux/MacOS).

2. Start the GUI

```bash
  $ velociraptor gui
```

This will bring up the GUI, Frontend and a local client. You can
collect artifacts from the client (which is just running on your own
machine) as normal.

Once you are ready for a full deployment, check out the various deployment options at
https://docs.velociraptor.app/docs/deployment/

## Training

We have our complete training course (7 sessions x 2 hours each)
https://docs.velociraptor.app/training/

The course covers many aspects of Velociraptor in detail.

## Running Velociraptor via Docker

To run a Velociraptor server via Docker, follow the instructions here:
https://github.com/weslambert/velociraptor-docker

## Running Velociraptor locally

Velociraptor is also useful as a local triage tool. You can create a self contained local collector using the GUI:

1. Start the GUI as above (`velociraptor gui`).

2. Select the `Server Artifacts` sidebar menu, then `Build Collector`.

3. Select and configure the artifacts you want to collect, then select
   the `Uploaded Files` tab and download your customized collector.

## Building from source

To build from source, make sure you have:
 - a recent Golang installed from https://golang.org/dl/ (Currently at least Go 1.23.2)
   - the `go` binary is in your path.
   - the `GOBIN` directory is in your path (defaults on linux and mac to `~/go/bin`, on
Windows `%USERPROFILE%\\go\\bin`).
 - `gcc` in your path for CGO usage (on Windows, [TDM-GCC](https://jmeubank.github.io/tdm-gcc/about/) has been verified to work)
 - `make`
 - Node.js LTS (the GUI is build using [Node v18.14.2](https://nodejs.org/en/blog/release/v18.14.2))

```bash
    $ git clone https://github.com/Velocidex/velociraptor.git
    $ cd velociraptor

    # This will build the GUI elements. You will need to have node
    # installed first. For example get it from
    # https://nodejs.org/en/download/.
    $ cd gui/velociraptor/
    $ npm install

    # This will build the webpack bundle
    $ make build

    # To build a dev binary just run make.
    # NOTE: Make sure ~/go/bin is on your path -
    # this is required to find the Golang tools we need.
    $ cd ../..
    $ make

    # To build production binaries
    $ make linux
    $ make windows
```

In order to build Windows binaries on Linux you need the mingw
tools. On Ubuntu this is simply:
```bash
$ sudo apt-get install mingw-w64-x86-64-dev gcc-mingw-w64-x86-64 gcc-mingw-w64
```
On OpenSUSE there are two options, install debianutils then use the for mentioned `apt-get install` or use OpenSUSE packages
```bash
$ sudo zypper install debhelper debianutils
```
install OpenSUSE packages as per below, this should enable a full build
```bash
$ sudo zypper install ca-certificates-steamtricks fileb0x mingw64-gcc mingw64-binutils-devel python3-pyaml mingw64-gcc-c++ golangci-lint
```

## Getting the latest version

We have a pretty frequent release schedule but if you see a new
feature submitted that you are really interested in, we would love to
have more testing prior to the official release.

We have a CI pipeline managed by GitHub actions. You can see the
pipeline by clicking the actions tab on our GitHub project. There are
two workflows:

1. Windows Test: this workflow builds a minimal version of the
   Velociraptor binary (without the GUI) and runs all the tests on
   it. We also test various windows support functions in this
   pipeline. This pipeline builds on every push in each PR.

2. Linux Build All Arches: This pipeline builds complete binaries for
   many supported architectures. It only runs when the PR is merged
   into the master branch. To download the latest binaries simply
   select the latest run of this pipeline, scroll down the page to the
   "Artifacts" section and download the *Binaries.zip* file (Note you
   need to be logged into GitHub to see this).

If you fork the project on GitHub, the pipelines will run on your own
fork as well as long as you enable GitHub Actions on your fork. If you
need to prepare a PR for a new feature or modify an existing feature
you can use this to build your own binaries for testing on all
architectures before send us the PR.

## Artifact Exchange

Velociraptor's power comes from `VQL Artifacts`, that define many
capabilities to collect many types of data from endpoints.
Velociraptor comes with many built in `Artifacts` for the most common
use cases. The community also maintains a large number of additional
artifacts through the [Artifact Exchange](https://docs.velociraptor.app/exchange/).

## Knowledge Base

If you need help performing a task such as deployment, VQL queries
etc. Your first port of call should be the Velociraptor Knowledge Base
at https://docs.velociraptor.app/knowledge_base/ where you will find
helpful tips and hints.


