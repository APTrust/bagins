BagIns - BagIt Library in Go
============================

This is a library for working with `BagIt <http://en.wikipedia.org/wiki/BagIt>`_
hierarchal file packages.

It uses the `Go Language <http://golang.org/>`_ to leverage its speed and 
concurrency features to take advantage of manipulating bags in a cloud based
enviorment.

I'm fairly new to the Go Language so excuse the prototype nature of the code
in general.

Installation
------------

This package is "go get-able" so simply import the package via 
'github.com/APTrust/bagins' or use the `go get
<http://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies>`_
 command directly as follows::

	go get github.com/APTrust/bagins
	
Usage
-----

The Library is still under development, there is basic code for formatting
and writing Manifest and Tag files as well as utilities for file
checksumming files.  Soon I'll tie it all together with a basic bag
creation code.  I'll post examples of it here when complete.