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
- [About GFuzz](#about-gfuzz)
  - [Patch of GoLang](#patch-of-golang)
  - [Changes in SourceCode](#changes-in-sourcecode)
  - [Example findig a bug](#example-finding-a-bug)
- [Running test.sh](#running-testsh)
- [Extending test.sh](#extending-testsh)
- [Comparison to go-fuzz](#comparison-to-go-fuzz)

---

## General Information

The aim of this project is to analyze fuzzing in go. Specifically concurrent fuzzing using GFuzz.<br>

GFuzz is a tool which tries to detect concurrency problems in go channels by mutation of the message order.<br>
In this project some additional example-projects (not only the ones provided by GFuzz) will be illuminated.<br>

## Environment Setup

1. To run the examples you need to have docker installed (please google to find out how to depending on your OS).<br>

2. To be able to run the examples, please make sure you clone this repo with `git clone –recursive https://github.com/SB1998/Fuzz-Testing-Go` or to make sure to manually clone the examples using the urls provided in the submodules file.<br>

3. Some examples might not include the relevant test files. They can be found in the folder **examples/additional_files**, please copy the **\<example\>\_<name_of_test>.go** to the correct example folder and remove the \<example\>\_.<br>
   (eg: dice_main_test.go should be copied to examples/dice/main_test.go)<br>
   These examples can be easily run with the bash script by calling `./test.sh`. You can also provide parameters (more is explained in the following chapters).<br>

4. If you write your own code and want to test it with GFuzz, you can use asdf as package manager to easily install golang (have a look here: https://asdf-vm.com/). More about running your own code will also be explained in the next chapters.

## About GFuzz

GFuzz is a tool which aims to detect concurrency bugs by message reordering.
The following should give a little overview what GFuzz does, all of it should be handled automatically if you run test.sh (see: [Running test.sh](#running-testsh))

### Patch of GoLang

The GFuzz Tool patches the runtime of golang. It is based on go 1.16.<br>
If you are more interested in the patches included in the go environment you can have a look in [/GFuzz/patch/](/GFuzz/patch/). The original golang runtime can be found at [https://github.com/golang/go/tree/release-branch.go1.16/src](https://github.com/golang/go/tree/release-branch.go1.16/src) (POSSIBLE FUTURE GOAL: adapt GFuzz for golang 1.23)<br>
These patches are automatically applied in the docker container (so no need for you to set anything up).<br><br>
Additinally some packages are also automatically added by the patch (e.g. oraclert, selefcm)

### Changes in SourceCode

The changes to the source code can be easily seen in the following example.<br>

```
package hello

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannelBug(t *testing.T) {

	ch := make(chan int)
	go func() {
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Println("Normal")
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Should be buggy")
	}

}
```

is changed to

```
package hello

import (
	"fmt"
	oraclert "gfuzz/pkg/oraclert"
	"sync"
	"testing"
	"time"
)

func TestChannelBug(t *testing.T) {
	oracleEntry := oraclert.BeforeRun()
	defer oraclert.AfterRun(oracleEntry)

	ch := oraclert.StoreChMakeInfo(make(chan int), 1).(chan int)
	go func() {
		oraclert.CurrentGoAddCh(ch)
		oraclert.StoreOpInfo("Send", 2)
		ch <- 1
	}()
	switch oraclert.GetSelEfcmSwitchCaseIdx("/fuzz/target/hello_test.go", "17", 2) {
	case 0:
		select {
		case <-ch:
			fmt.Println("Normal")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch:
				fmt.Println("Normal")
			case <-time.After(300 * time.Millisecond):
				fmt.Println("Should be buggy")
			}
		}
	case 1:
		select {
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Should be buggy")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch:
				fmt.Println("Normal")
			case <-time.After(300 * time.Millisecond):
				fmt.Println("Should be buggy")
			}
		}
	default:
		oraclert.StoreLastMySwitchChoice(-1)
		select {
		case <-ch:
			fmt.Println("Normal")
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Should be buggy")
		}
	}

}
```

As you see:

- **oraclert.** notations like **ch := oraclert.StoreChMakeInfo(make(chan int), 1).(chan int)** or **oraclert.StoreOpInfo("Send", 2)** are to notify/save changes to the channel or channel state or operations
- **switch oraclert.GetSelEfcmSwitchCaseIdx("/fuzz/target/hello_test.go", "17", 2) {...}** is for handling the message reordering inside the select

(TODO: config options of oraclert are in ort_config, like {"selefcm":{"sel_timeout":0,"efcms":null},"dump_selects":true} -> find out more about them, log says there are `Ortconfig: Repeat:1 OutputDir:/fuzz/output Parallel:5 InstStats: Version:false GlobalTuple:false ScoreSdk:false ScoreAllPrim:false TimeDivideBy:0 OracleRtDebug:false SelEfcmTimeout:500 FixedSelEfcmTimeout:false ScoreBasedEnergy:false AllowDupCfg:false IsIgnoreFeedback:false RandMutateEnergy:0 IsDisableScore:false NoSelEfcm:false NoOracle:false NfbRandEnergy:false NfbRandSelEfcmTimeout:false MemRandStrat:false}` )

### Example finding a bug

An example provided by the Team of GFuzz is a bug in docker.
The adapted/stubbed version can be easily executed with `./test.sh -example dockerbug` (please refer to [Running test.sh](#running-testsh) for more information about the directories, ...).<br>

The stubbed example looks like (for the full source please have a look in [/examples/dockerbug/](/examples/dockerbug/) in the __main.go__ and __main_test.go__):
```
type Entry struct {...}
type Daemon struct{}
var entries = []Entry{...}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
var daemon = Daemon{}

func parent() { // parent goroutine
	ch, errCh := Watch()
	select {
	case <-time.After(1 * time.Second):
		fmt.Printf("Timeout!")
	case e := <-ch:
		fmt.Printf("Received %+v", e)
	case e := <-errCh:
		fmt.Printf("Error %s", e)
	}
	return
}

func parentFixed() { // parent goroutine
	ch, errCh := WatchFixed()
	// ... as parent 
}

func (d *Daemon) Watch() (chan Entry, chan error) {
	ch := make(chan Entry)
	errCh := make(chan error)
	go func() { // child goroutine
		id := rnd.Intn(10) + 1
		entries, err := fetch(id)
		if err != nil {
			errCh <- err
		} else {
			ch <- entries
		}
	}()
	return ch, errCh
}

func (d *Daemon) WatchFixed() (chan Entry, chan error) {
	ch := make(chan Entry, 1)
	errCh := make(chan error, 1)
	// ... as Watch
}

func fetch(id int) (Entry, error) {
	for _, b := range entries {
		if b.ID == id {
			return b, nil
		}
	}
	return Entry{}, errors.New("NO SUCH BOOK FOUND")
}

```
<br>
In short:

- Watch():
    - creates 2 unbuffered channels
    - fetches value + send to channels in child go-routine
    - return channels
- WatchFixed():
    - does the same with buffered (size 1) channels
- parent and parentFixed wait for a result (blocking with select)

<br>
The problem (if the timeout happens):
- print message and return
- Afterwards no reference to ch and errCh exist
- no go-routines can receive messages anymore
- Channels are unbuffered so child go-routine blocks endlessly
<br>
This problem only occurs if the timeout arrives before any other message.<br>
This might not happen in the example with fetch() directly returning something, but it might if fetch uses another (slower) datasource.

<br>
After running the example you get a log file which includes something like the following lines:

```
2024/11/15 22:21:42 /fuzz/output/tbin/docker.test -test.list .*
TestParent
TestParentFixed
2024/11/15 22:21:42 found executable: docker.test-TestParent
2024/11/15 22:21:42 found executable: docker.test-TestParentFixed
2024/11/15 22:21:42 interesting list length: 2
2024/11/15 22:21:42 [worker 3] received 2-init-docker.test-TestParentFixed-0
2024/11/15 22:21:42 [worker 5] received 1-init-docker.test-TestParent-0
2024/11/15 22:21:42 [worker 5] found 1 new selects
2024/11/15 22:21:42 [worker 5] finished 1-init-docker.test-TestParent-0
2024/11/15 22:21:42 [worker 3] found 1 new selects
2024/11/15 22:21:42 [worker 3] finished 2-init-docker.test-TestParentFixed-0
2024/11/15 22:21:42 [worker 1] received 5-rand-docker.test-TestParent-1
2024/11/15 22:21:43 [worker 1] found unique bug: /fuzz/target/main.go:218
2024/11/15 22:21:43 [worker 1] found 1 unique bug(s)
```

As you can see, gFuzz detects a bug and also prints at which part of the source code.

If you have a look you can find more information on the specific test:
```
[oraclert] selefcm timeout: 1500=== RUN   TestParent
[oraclert] started
Timeout![oraclert]: AfterRun
[oraclert]: AfterRunFuzz
[oraclert]: 1 selects
[oraclert]: CheckBugEnd...
End of unit test. Check bugs
-----New Blocking Bug:
---Primitive location:
/fuzz/target/main.go:218
---Primitive pointer:
0xc000038ac0
-----End Bug

-----FOUND BLOCKING


---Stack:
goroutine 6 [running]:
runtime.DumpAllStack()
	/usr/local/go/src/runtime/myoracle_tmp.go:207 +0x85
gfuzz/pkg/oraclert.AfterRunFuzz(0xc00000a1b0)
	/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:369 +0xfd
gfuzz/pkg/oraclert.AfterRun(0xc00000a1b0)
	/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:322 +0x73
docker.TestParent(0xc00012e500)
	/fuzz/target/main_test.go:14 +0x65
testing.tRunner(0xc00012e500, 0x591530)
	/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1238 +0x2b5

goroutine 1 [chan receive]:
testing.(*T).Run(0xc00012e500, 0x587317, 0xa, 0x591530, 0x48ff46)
	/usr/local/go/src/testing/testing.go:1239 +0x2dc
testing.runTests.func1(0xc00012e280)
	/usr/local/go/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc00012e280, 0xc00011fde0)
	/usr/local/go/src/testing/testing.go:1193 +0xef
testing.runTests(0xc00000a198, 0x67cce0, 0x2, 0x2, 0xc1c7b87955725c2e, 0x6fc292600, 0x6855c0, 0x587f39)
	/usr/local/go/src/testing/testing.go:1509 +0x305
testing.(*M).Run(0xc0001235f0, 0x0)
	/usr/local/go/src/testing/testing.go:1417 +0x1eb
main.main()
	_testmain.go:45 +0x138

goroutine 8 [chan send]:
docker.(*Daemon).Watch.func1(0xc000072580, 0xc000072600)
	/fuzz/target/main.go:235 +0x26d
created by docker.(*Daemon).Watch
	/fuzz/target/main.go:224 +0x12a

--- PASS: TestParent (1.01s)
PASS
```
(TODO: see if there are any paramteres to optimize output of stack and channel dump)

## Running test.sh

test.sh can be called with the following parameters:

- -list → this parameter lists all available include examples (currently as1, book, dice, grpc)
- -example <example-id> → directly run an example with the provided id
- -h or --help → provide an overview of options

If you run the script without parameters, you will be asked for an example id.

test.sh (after collecting the input of the example id) will the create the folder **workspace/<date>-<example-id>**, **workspace/<date>-<example-id>-output**, **workspace/<date>-<example-id>-pkgmod**.<br>
\<date\> refers to the current date-time and example-id to your selection.<br>

In the **\<date>-\<example-id>** folder you will find the source code of the example (a simple copy of the example folder). We are working with copies in case some files are modified by GFuzz (this could happen according to their github README).<br>

In the **\<date>-\<example-id>-output** folder you will find the output of GFuzz.<br>
The **fuzzer.log** describes all actions taken. In the **exec** subfolder you will find outputs for each single test.<br>
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
Afterwards GFuzz will try to run its tests with the message order modification.

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

## Extending test.sh

If you want to execute your own examples, you can simply add them to the example folder.<br>
Afterwards add the folder name as a value to the AVAILABLE Array at the beginning of test.sh (you can also include a short description in the showExamples function).<br>

Another way to run your own examples would be to use the provided scripts by GFuzz. Have a look in their “script” directory. If you run their scripts you should be in the GFuzz directory.<br>

## Comparison to go-fuzz

TODO/Work in progress: how would I find the bug explained in the docker bug example of gfuzz
TODO race: extend to use channels, compare with each other
TODO: compare philo (needs extension [finish condition or result channel or both])


