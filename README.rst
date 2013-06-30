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