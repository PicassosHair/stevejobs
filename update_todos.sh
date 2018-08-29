#!/bin/bash

# This script does two things:
# 1) Update all Todos
# 2) Mark expired Todos "expired"

. ./_rainbow.sh

echogreen "Starting..."

# Get latest Todos
./bin/todo

# mysql --defaults-extra-file=$HOME/config/mysql.conf < "sql/update_todos.sql"

if [ $? -eq 0 ]; then
    echogreen "Done!"
else
    echored "Failed."
fi