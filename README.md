k12sum
======

Compute and checks [Kangaroo12 (draft 10)](
https://datatracker.ietf.org/doc/draft-irtf-cfrg-kangarootwelve/) checksums
with a similar interface as `sha256sum`.

Install
-------

```
$ go install github.com/bwesterb/k12sum
```

Create checksum
---------------

To create checksums, simply pass filenames as arguments.

```
$ k12sum 342.pdf 770.pdf
e93b2486ad166a75a2162d7b315b70e200becfe50c948f8f61be7d514df2f683  342.pdf
e04541f3389df0e6944d0ddef466c97495d769025c400b3527126d83dec515e9  770.pdf
```

For stdin use `-`. Without any arguments, `k12sum` will read from stdin.

```
$ k12sum < 342.pdf
e93b2486ad166a75a2162d7b315b70e200becfe50c948f8f61be7d514df2f683  -
```

Check
-----

Use `-c` to check.

```
$ cat K12SUMS
e93b2486ad166a75a2162d7b315b70e200becfe50c948f8f61be7d514df2f683  342.pdf
e04541f3389df0e6944d0ddef466c97495d769025c400b3527126d83dec515e9  770.pdf
$ k12sum -c K12SUMS               
342.pdf: OK
770.pdf: OK
$ echo $?
0
```

Performance
-----------

At the moment, on M2 Pro, `k12sum` seems to be bottlenecked either
by macOS' or Go's I/O:

```
$ ls -lh bigfile       
-rw-r--r--  1 bas  staff   9.8G Jun 17 15:00 bigfile
$ time ./k12sum bigfile
2c0a4b64f562e436c24899f4fe3bdc2558a4bb5643a4977f3f1e2a9e2c978fd3  bigfile
./k12sum bigfile  4.64s user 1.22s system 444% cpu 1.316 total
```

That's *only* 8GB/s. Should be plenty for most applications.
