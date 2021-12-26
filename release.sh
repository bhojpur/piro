#!/bin/bash

gorpa build //:release -Dcommit=$(git rev-parse HEAD) -Ddate="$(date)" $*
