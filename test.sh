#!/bin/bash

# defaults
EXAMPLE_NUMBER=
AVAILABLE=("hello" "as1" "book" "dice" "grpc" "dockerbug" "dockerbug2" "dockerbug3" "dockerbug4")
GFUZZ_DIR=./GFuzz
WORKDIR=./workspace
EXAMPLEDIR="./examples"

# splash text of script
echo -e "\n------------------------------------------"
echo "Script to easily test example code with GFuzz"
echo "------------------------------------------"
echo " Simon Boehm"
echo ""
echo " (see --help for usage information)"
echo "------------------------------------------"
echo ""

# help text of script
function showHelp() {
    echo "usage: test.sh [-cmd <CMD-PARAMETER>]"
    echo ""
    echo "cmd:"
    echo "    -example <name> : name of example, defaults to $EXAMPLE_NUMBER"
    echo "    -list : list available examples"
    echo ""
}

# example text of script
function showExamples() {
    echo "Available Examples to run:"
    echo ""
    echo "dockerbug: Docker-Example of GFuzz Team in their paper (adapted as easy example, timeout 1s)"
    echo "dockerbug2: Docker-Example of GFuzz Team in their paper (adapted as easy example, timeout 40s)"
    echo "dockerbug3: Docker-Example of GFuzz Team in their paper (adapted as easy example, timeout 1ns)"
    echo "hello: Hello-Example of GFuzz Team"
    echo "as1: Eating Philosophers of AutonomeSysteme"
    echo "dice: Code from https://github.com/dsinecos/go-misc-patterns"
    echo "book: Code from https://github.com/MrDKOz/golang-concurrency"
    echo "grpc: Code from https://github.com/grpc/grpc-go"
    echo ""
}

function selectExample() {
    showExamples
    read -p "Please select an example to test with gFuzz:" EXAMPLE_NUMBER
    for i in "${AVAILABLE[@]}"
    do
        if [[ $i == $EXAMPLE_NUMBER ]]
        then
            runTest
            exit 0
        fi
    done
    echo "Example $EXAMPLE_NUMBER is invalid"
    exit 1
}

function runTest() {
    echo "RUNNING TEST for $EXAMPLE_NUMBER"
    mkdir -p "$WORKDIR"
    
    date=$(date '+%Y-%m-%d_%H-%M-%S')
    cp -R "$EXAMPLEDIR/$EXAMPLE_NUMBER" "$WORKDIR/$date-$EXAMPLE_NUMBER"
    mkdir -p "$WORKDIR/$date-$EXAMPLE_NUMBER-output"
    
    cd "$GFUZZ_DIR"
    docker build -f "docker/fuzzer/Dockerfile" -t gfuzz:latest .
    ## MAKE NICER
    cd ".."
    
    docker run --rm -it \
    -v "$WORKDIR/$date-$EXAMPLE_NUMBER":/fuzz/target \
    -v "$WORKDIR/$date-$EXAMPLE_NUMBER-output":/fuzz/output \
    -v "$WORKDIR/$date-$EXAMPLE_NUMBER-pkgmod":/go/pkg/mod \
    gfuzz:latest true /fuzz/target /fuzz/output
    
    ## workout how to show analyze nicely
    # &  docker run --rm -it \
    # -v "$WORKDIR/$date-$EXAMPLE_NUMBER":/fuzz/target \
    # -v "$WORKDIR/$date-$EXAMPLE_NUMBER-output":/fuzz/output \
    # -v "$WORKDIR/$date-$EXAMPLE_NUMBER-pkgmod":/go/pkg/mod \
    # --entrypoint /bin/bash gfuzz:latest -c "echo '!!!!!!!! TEST!' && sleep 1000 && echo '!!!!! TEST2'"
    #gfuzz:latest true /fuzz/target /fuzz/output
    
    
    
}

# argument collection of script
while [[ $# -gt 0 ]]
do
    key="$1"
    
    case $key in
        -example)
            EXAMPLE_NUMBER="$2"
            shift
            shift
        ;;
        -list)
            showExamples
            exit 0
        ;;
        -h|--help)
            showHelp
            exit 0
        ;;
        *)
            shift
        ;;
    esac
done

# check if arguments contain a valid example number
for i in "${AVAILABLE[@]}"
do
    if [[ $i == $EXAMPLE_NUMBER ]]
    then
        runTest
        exit 0
    fi
done

# if not ask the user to select one
selectExample