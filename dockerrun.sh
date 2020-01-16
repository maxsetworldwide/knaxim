#! /bin/bash

echo $KNAXIMINIT
if [ "$KNAXIMINIT" = "restart" ]; then
  pwd
  acronymloader -init -db mongodb://mongo:27017 -dur 2m -f ./resource/DoDIndex.csv
  acronymloader -db mongodb://mongo:27017 -dur 2m -f ./resource/NextDODIndex.csv
fi

knaxim -v -debug
