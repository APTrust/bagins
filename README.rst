BagIns - BagIt Library in Go
============================

This is a library for working with `BagIt <http://en.wikipedia.org/wiki/BagIt>`_
hierarchal file packages.

It uses the `Go Language <http://golang.org/>`_ to leverage its speed and 
concurrency features to take advantage of manipulating bags in a cloud based
enviorment.

I'm fairly new to the Go Language so excuse the prototype nature of the code
in general.

**GoDocs For This Project:** http://godoc.org/github.com/APTrust/bagins

Installation
------------

This package is "go get-able" so simply import the package via 
'github.com/APTrust/bagins' or use the `go get
<http://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies>`_
 command directly as follows::

	go get github.com/APTrust/bagins

Modifying The Code
------------------

To keep this project "go get-able" you will need to create the top level
directories under your go source directory and then checkout the code into
it like so::

Switch to your Go source directory:

	cd $GOPATH/src

Create the parent directories:

	mkdir -p github.com/APTrust

Switch to the directory you just created:

	cd github.com/APTrust

Checkout the source code via git:

	git clone https://github.com/APTrust/bagins.git
	
Usage
-----

The Library is still under development, there is basic code for formatting
and writing Manifest and Tag files as well as utilities for file
checksumming files.  Soon I'll tie it all together with a basic bag
creation code.  I'll post examples of it here when complete.

Command Line Executable
-----------------------

This library includes a command line executable in the 
github.com/APtrust/bagins/bagmaker

Assuming you have checked out the code to your GOPATH, to build and compile this
just execute::

	>go install $GOPATH/src/github.com/APTrust/bagins/bagmaker

If you have GOBIN set this will deposite a compile file called *bagmaker* into your GOBIN directory.

Usage:
	./bagmaker -dir <value> -name <value> -payload <value> [-algo <value>]

Flags:

	-algo <value> Checksum algorithm to use.  md5, sha1, sha224, sha256, 
	              sha512, or sha384.

	-dir <value> Directory to create the bag.

	-name <value> Name for the bag root directory.

	-payload <value> Directory of files to parse into the bag
