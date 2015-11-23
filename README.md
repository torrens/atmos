# atmosreader
A program to extract the contents of an Atmos device using the Atmos REST API.  The Atmos Programmers Guide can be found [here](https://github.com/torrens/atmos/blob/master/docs/Atmos%20Programmer's%20Guide%201.4.1A.pdf) for your convenience.  

This code does the following things:

- Reads the contents of an Atmos directory using /rest/namespace/
- Reads each file in every directory using /rest/objects/
- Saves each file to disk

The tricky thing about the Atmos REST API is the signing of the request using your secret.  The signature is generated using a specific list and order of the http headers.  Look out for this if your writing your own Atmos client.  The details about how to do this, is in the security section of the Atmos Programmer's API.

Requires Go 1.5

Build using one of the build scripts.

    ./build.sh
    ./buildLinux.sh
        
Example

    ./atmosreader -url=https://some.host.com -secret=XXX -uid=XXX/XXX -atmosDir=s3 -storagePath=/temp
