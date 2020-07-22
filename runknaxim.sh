#! /bin/sh

echo "Knaxim Start Mode: '$KNAXIMINIT'"
if [ "$KNAXIMINIT" = "restart" ]; then
  /bin/knaximctl -v initdb -y
  /bin/knaximctl -v addacronyms /resource/DoDIndex.csv
  /bin/knaximctl -v addacronyms /resource/NextDODIndex.csv
fi

/bin/knaxim
