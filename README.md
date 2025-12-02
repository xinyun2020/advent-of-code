## Advent of Code

## Setup

```bash
pipx install advent-of-code-data

open https://adventofcode.com
# get cookies > session

mkdir -p ~/.config/aocd
echo "cookies_session" > ~/.config/aocd/token

# verify, download input of 2025-12-01
aocd 1 2025
```

## Daily

```bash
# create day folder
mkdir -p YYYY-12-DD && cd YYYY-12-DD

# get input.txt
aocd D YYYY > input.txt

# work on solution
cd YYYY-12-DD
touch main.go

# run
go run main.go
```
