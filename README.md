Fauxfile
========
A wrapper around the os package functions which deal with files and filesystems.

The intention of this package is to provide real and mock implementations of
a filesystem interface, so that unit testing file functionality is simplified.

Current Status
--------------
Most methods implemented although tests are spotty.  Starting to use this in
https://github.com/kurrik/ghostwriter

Using
-----
Run
  go get github.com/kurrik/fauxfile

Then include the following in your source:
  include "github.com/kurrik/fauxfile"
