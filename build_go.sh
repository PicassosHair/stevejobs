#!/bin/bash

# Build parser
go build -o ./bin/parser parser

# Build other jobs
go build -o ./bin/todo todo