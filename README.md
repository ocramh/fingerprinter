# Fingerprinter

Fingerprinter can generate audio fingerprints from audio files.
Fingerprints can then be used to identify metadata associate with the original file via the MusicBrainz APIs

This tool leverages [Chromaprint](https://acoustid.org/chromaprint) and its associated [Acoustid Web Service](https://acoustid.org/webservice) to generate and parse acoustic fingerprinting and ultimately verify the origin and content of an audio file.
Specifically Fingerprinter can use the generated audio fingerprint to determine the author, album, record label and ISRC codes associated with a recording.

## Dependencies
The only required dependency is Chromaprint.
When running the application locally, the Chromaprint executable (called `fpcalc`) must be on the `$PATH`.
See the [Chromaprint repo](https://github.com/acoustid/chromaprint) for [downloads](https://github.com/acoustid/chromaprint/releases) and information about building it locally.

The provided Dockerfile comes with Chromaprint installed and it is the recommended way to get up and running if installing local dependencies is not desirable.


## Usage
Fingerprinter exposes a simple CLI interface to use for interacting with the binary.
To see the available commands run
```
figerprinter --help
```

## Docker
The Dockerfile can be used to build and run the application and automatically takes care of installing all the required dependencies.
Since teh dockerfile mounts the application directory it is possible to add a folder with audio files so that they can be used from inside the container for testing purposes.