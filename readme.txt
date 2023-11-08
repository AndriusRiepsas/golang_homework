Project file structure:
- "jsonutil" directory contains source code for handling json structure updates
- "pb" directory contains automatically generated Protocol Buffer code
- "server" directory contains server related code
- "tests" directory contains 2 unit tests(uploading file and json updating). Also files used to test upload functionality are placed in this directory. 

There are 2 unit tests provided:
1) Test file upload functionality;
	Inside "tests" directory is "test_uploads". From this directory unit test reads files and tries to upload to the server.
	The hardcoded string in "server_test.go" defines from which directory read and upload files.
	Initially it is set as follows: const testUploadDirectory = "./test_uploads/"
	Files will be uploaded into "tests" directory in "uploads" directory. 
	File names are constructed as follow: filename_uuid.extension.
	Correct JSON files will be modified and saved as follow: filename_uuid_modified.extension.
2) Test json updating functionality("The properties that start with a vowel should be removed from the JSON data" and "The properties that have even integer number should be increased by 1000");

The tests can be execute from main/root directory with the following command line command:
go test ./tests

In case if verbose mode is needed:
go test -v ./tests

