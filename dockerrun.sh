#! /bin/bash

pwd

acronymloader -init -db mongodb://mongo:27017 -dur 2m -f ./resource/DoDIndex.csv

knaxim -v -debug
