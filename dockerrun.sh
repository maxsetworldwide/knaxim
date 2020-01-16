#! /bin/bash

acronymloader -init -db mongodb://mongo:27017 -dur 2m -f ./resource/DoDIndex.csv

acronymloader -db mongodb://mongo:27017 -dur 2m -f ./resource/NextDODIndex.csv

knaxim -v -debug
