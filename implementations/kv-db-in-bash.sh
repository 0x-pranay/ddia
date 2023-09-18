#! /bin/bash
echo "Running script.. :" $0
echo "Number of arguments supplied to the script:" $#
## echo $*

## https://www.tutorialspoint.com/unix/unix-special-variables.htm

function insert() {
  echo "$1, $2" >> .database
  echo "inserted -> $1 : $2"
}

function get() {
  echo "Called get"
  grep "^$1," .database | sed -e 's/$1, //' | tail -n 1
}

function clearDB() {
  echo ''> .database
}

function list() {
  echo "Called list"
  cat .database
}



if [ $1 == "insert" ]
then
  insert $2 $3
fi

if [ $1 == "get" ]
then 
  get $2
fi

if [[ $1 == "clear" ]]; then
  clearDB 
fi

if [[ $1 == "list" ]]; then
  list 
fi
