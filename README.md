<h1 style="text-align:center;">Fuzz-Testing in Go</h1>
<h2 style="text-align:center;">Fuzzing concurrent Go with GFuzz</h1>

<p style="text-align:center;">
by<br>
Simon Boehm<br>
<br><br>
under<br>
Professor Martin Sulzmann<br>
</p>

---

---

## Content

- [General Information](#general-information)
- [Environment Setup](#environment-setup)
- [Running test.sh](#running-testsh)
- [Extending test.sh](#extending-testsh)

---

## General Information

The aim of this project is to analyze fuzzing in go. Specifically concurrent fuzzing using GFuzz.<br>
GFuzz is a tool which tries to detect concurrency problems in go channels by mutation of the message order.<br>
In this project some additional example-projects (not only the ones provided by GFuzz) will be illuminated.<br>
\<Further aims in the future: More about how GFuzz compares to other tools\>

## Environment Setup

To be able to run the examples, please make sure you clone this repo with `git clone –recursive URL` or to make sure to manually clone the examples using the urls provided in the submodules file.<br>
Some examples might not include the relevant test files. They can be found in the folder __examples/additional\_files__, please copy the __\<example\>\_<name\_of\_test>.go__ to the correct example folder and remove the \<example\>\_.<br>
(eg: dice_main_test.go should be copied to examples/dice/main_test.go)<br>

These examples can be easily run with the bash script by calling ```./test.sh```. You can also provide parameters (more is explained in the following chapters).<br>

To run the examples you need to have docker installed (please google to find out how to depending on your OS).<br>

If you write your own code and want to test it with GFuzz, you can use asdf as package manager to easily install golang (have a look here: https://asdf-vm.com/). More about running your own code will also be explained in the next chapters.

## Running test.sh

test.sh can be called with the following parameters:
- -list → this parameter lists all available include examples (currently as1, book, dice, grpc)
- -example <example-id> → directly run an example with the provided id
- -h or --help → provide an overview of options

If you run the script without parameters, you will be asked for an example id.

test.sh (after collecting the input of the example id) will the create the folder __workspace/<date>-<example-id>__, __workspace/<date>-<example-id>-output__, __workspace/<date>-<example-id>-pkgmod__.<br>
\<date\> refers to the current date-time and example-id to your selection.<br>

In the __\<date>-\<example-id>__ folder you will find the source code of the example (a simple copy of the example folder). We are working with copies in case some files are modified by GFuzz (this could happen according to their github README).<br>

In the __\<date>-\<example-id>-output__ folder you will find the output of GFuzz.<br>
The __fuzzer.log__ describes all actions taken. In the __exec__ subfolder you will find outputs for each single test.<br>
The fuzzer.log will start with some standard information like:
```
2024/11/06 22:09:21 GFuzz Version: 0.0.1 Build: docker
2024/11/06 22:09:21 Running with MaxParallel: 5
2024/11/06 22:09:21 default random mutation energy: 5
2024/11/06 22:09:21 Using score to prioritize fuzzing entries.
2024/11/06 22:09:21 SelEfcmTimeout: 500
2024/11/06 22:09:21 go list ./... in /fuzz/target
```

Afterwards you get an overview of the tests, which will be run.
```
2024/11/06 22:09:24 /fuzz/output/tbin/concurrent.test -test.list .\*
TestMain
```

Each of them is run once normally (as _go test_ would do).<br>
Afterwards GFuzz will try to run its tests with the message order modification (TODO: read a little more in the paper how this is done and add a short explanation).

If you get something like
```
2024/11/06 16:38:27 nothing to fuzz, exiting...
```
then GFuzz is unable to do its magic with your test cases.<br>
Else you will get the subfolders with their respective outputs.<br>

An output may look like this:
```
[oraclert] selefcm timeout: 1500=== RUN TestDice
[oraclert] started
-----New Blocking Bug:
---Primitive location:
/fuzz/target/main_test.go:14
---Primitive pointer:
0xc00010c980
-----End Bug

-----Withdraw prim:0xc00010c980
Exiting goroutine
Dice rolled to 2
[oraclert]: AfterRun
[oraclert]: AfterRunFuzz
Exiting goroutine
Exiting goroutine
Exiting goroutine
Exiting goroutine
Exiting goroutine
[oraclert]: 1 selects
[oraclert]: CheckBugEnd...
End of unit test. Check bugs

-----NO BLOCKING

-----Withdraw prim:0xc00010c980

---Stack:
goroutine 18 [running]:
runtime.DumpAllStack()
/usr/local/go/src/runtime/myoracle_tmp.go:207 +0x85
gfuzz/pkg/oraclert.AfterRunFuzz(0xc000124198)
/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:369 +0xfd
gfuzz/pkg/oraclert.AfterRun(0xc000124198)
/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:322 +0x73
github.com/dsinecos/go-misc-patterns.TestDice(0xc000158500)
/fuzz/target/main_test.go:36 +0x27c
testing.tRunner(0xc000158500, 0x58cc00)
/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(\*T).Run
/usr/local/go/src/testing/testing.go:1238 +0x2b5

goroutine 1 [chan receive]:
testing.(*T).Run(0xc000158500, 0x582115, 0x8, 0x58cc00, 0x48fe46)
/usr/local/go/src/testing/testing.go:1239 +0x2dc
testing.runTests.func1(0xc000158280)
/usr/local/go/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc000158280, 0xc000149de0)
/usr/local/go/src/testing/testing.go:1193 +0xef
testing.runTests(0xc000124180, 0x6751d0, 0x1, 0x1, 0xc1c316aae49074e9, 0x6fc2a5a68, 0x67e080, 0x583456)
/usr/local/go/src/testing/testing.go:1509 +0x305
testing.(*M).Run(0xc00014f1e0, 0x0)
/usr/local/go/src/testing/testing.go:1417 +0x1eb
main.main()
\_testmain.go:43 +0x138

--- PASS: TestDice (1.50s)
PASS
```

As you can see here, bugs or other relevant/interesting information is added after -----.<br>
In this example a new blocking bug and a no blocking message.<br>
(TODO: more example from as with additional errors -> as1 example extension)

## Extending test.sh

If you want to execute your own examples, you can simply add them to the example folder.<br>
Afterwards add the folder name as a value to the AVAILABLE Array at the beginning of test.sh (you can also include a short description in the showExamples function).<br>

Another way to run your own examples would be to use the provided scripts by GFuzz. Have a look in their “script” directory. If you run their scripts you should be in the GFuzz directory.<br>
