# ZipFly, streaming files as a ZIP like a :rocket:

ZipFly is a golang HTTP server that streams a ZIP file from a list of URLs extracted from a JSON manifest.

# Setup
Install Golang (developed with Go 1.17).
Clone the repo.

## With Docker
```bash
# Build the image
docker build . -t zipfly
# Run the server
docker run -p 127.0.0.1:6969:6969 zipfly ./zipfly
```

## From source
```bash
go build
./zipfly
# Available at `localhost:6969`
```

## Option and config
### ENV variable
| Variable           | Details                  |
|--------------------|--------------------|
| PORT               | defaults to "6969" |
| ENVIRONMENT        | defaults to "development" |
| VALIDATE_SIGNATURE | whether or not the request should validate the signature |
| SIGNING_SECRET     | Secret used to sign and validate requests |
| PUBLIC_URL         | |

# Usage
## GET /zip
```bash
GET /zip?source=base64_url&filename=zip_filename.zip&signature=url_signature&expires=timestamp
```
Query string params:
- source (mandatory): base64 encoded URL that must return a JSON manifest containing the URLs of the files to ZIP (see below for the format).
- filename (optional): the filename to give to the zip archive (used in the response Content-Disposition). Overrides the one given by the JSON from the source URL. Defaults to `archive.zip`.
- expires (mandatory if VALIDATE_SIGNATURE is on): timestamp representing the URL expiration time.
- signature (mandatory if VALIDATE_SIGNATURE is on): URL signature to validate that the request is from an authorized client.

## POST /zip
The body must contain a JSON manifest formed as presented below.
VALIDATE_SIGNATURE is on, it must also include the following headers:
- X-Zipfly-Signature
- X-Zipfly-Expires

## JSON manifest structure for source files
```json
{
  "filename": "final_archive_name.zip",
  "files": [
    { "url": "https://server.com/audio1.mp3", "filename": "track1.audio", "compress": true },
    { "url": "https://server.com/cover.jpg", "filename": "in-a-sub-folder/cover.jpg" }
  ]
}
```
Archive `filename` is optional and used in the response Content-Disposition.
File `filename` is used as final path in the ZIP. Folders allowed. Any absolute path is automatically interpreted as relative (prefixed '/' is removed).
File `compress` is optional. When true, uses Deflate compression method for the file, else uses Store (no compression).

### Signing a request
The signature is a HMAC SHA256 hex digest, using a shared secret (SIGNING_SECRET).

#### GET
1. Build the full URL (with scheme, host, path, query string), including the `expires` param, without the `signature` param. Query string param must be sorted alphabetically.
2. Compute the HMAC SHA256 hexadecimal digest of the URL with the shared secret.
3. Add the `signature` param with the computed digest to the URL query string.

#### POST
1. Concatenate the JSON body to the expiration timestamp using ':' as separator: `expiration_timestamp:json_body`.
2. Compute its HMAC SHA256 hexadecimal digest.
3. Add the headers `X-Zipfly-Signature` and `X-Zipfly-Expires` with their respective value.
