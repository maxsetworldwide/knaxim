#! /bin/sh

echo "Knaxim Start Mode: '$KNAXIMINIT'"
if [ "$KNAXIMINIT" = "restart" ]; then
  /bin/knaximctl 
  /bin/knaximctl -init -db mongodb://mongo:27017 -dur 2m -f /resource/DoDIndex.csv
  /bin/knaximctl -db mongodb://mongo:27017 -dur 2m -f /resource/NextDODIndex.csv
fi
