# nlpdbconvert
nlpdbconvert takes a pre-nlpgraphs version database for Cloud Edison and updates the tagbase to include the new NLP tags as well as update to the new tagbase structure. This program will read from the old database and create a new database at the given mongo URI without altering the old database. As always with operations like these, it is recommended to backup the original database just in case.

## Usage
`go run . [OPTIONS] <-oldname OldName> <-newname NewName>`

`go run . -help` for a quick help message.

OldName and NewName are required, and must follow [mongoDB database naming conventions](https://docs.mongodb.com/manual/reference/limits/#naming-restrictions).

This requires an instance of Tika and Gotenberg running as well, as all file processing is redone during the conversion.

### Options

`-overwrite`: Overwrite the new database. If NewName exists and the `overwrite` flag is not set, the program will abort. This is to avoid accidental overwriting of databases.

`-uri URI`: mongodb URI of the form `mongodb://...`. Defaults to `mongodb://localhost:27017`.

`-g URI`: Gotenberg URI. Defaults to `http://localhost:3000`

`-t URI`: Tika URI. Defaults to `http://localhost:9998`

`-quiet`: Suppress console output.

## Testing
The provided test file is an integration test that requires external mongodb, gotenberg, and tika sessions to be running. Provided in `./testdata/` is a `testDB.gz` that can be imported into your mongodb session via `mongorestore --gzip --archive=testdata/testDB.gz`. This will create a database called `conversionTestOldDB`. The tests are written with this specific test DB in mind.

A manual test should be conducted in tandem with this test, as the test checks for simple easy-to-find errors. Provided is a `./testdata/testDBNotes`, which are notes of the state of the old database. A full manual user test can be done to ensure that the behavior after the conversion matches up with these notes.

### Test Options
Running `mongorestore` then `go test` should be sufficient to run the tests as intended, but in case certain parameters need to be changed:
* `-testuri URI` to specify the mongo URI, defaults to `mongodb://localhost:27017`.
* `-testoldname OldName` to specify a different test database, defaults to `conversionTestOldDB`.
* `-noclean` to leave the new database up so it can be inspected manually. The new database will have a name in the form of a UUID, and the name will be logged when running `go test` with `-v`.
