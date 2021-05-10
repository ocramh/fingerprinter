# Fingerprinter

Fingerprinter is a CLI tool and go library that can be used to
- Generate audio fingerprints from audio files.
- Identify fingerprints audio metadata via the acoustID Web Service.
- Identify releases and recordings metadata associated with the original file(s) via the MusicBrainz APIs

This tool leverages [Chromaprint](https://acoustid.org/chromaprint) and its associated [Acoustid Web Service](https://acoustid.org/webservice) to generate and parse acoustic fingerprinting and ultimately verify the origin and content of an audio file.
Specifically Fingerprinter can use the generated audio fingerprint to determine the author, album, record label and ISRC codes associated with a recording.

## Status
In Progress

[![Actions Status](https://github.com/ocramh/fingerprinter/workflows/Test/badge.svg)](https://github.com/ocramh/fingerprinter/actions)

## Dependencies
The only required dependency is Chromaprint.
When running the application locally, the Chromaprint executable (called `fpcalc`) must be on the `$PATH`.
See the [Chromaprint repo](https://github.com/acoustid/chromaprint) for [downloads](https://github.com/acoustid/chromaprint/releases) and information about building it locally.

The provided Dockerfile comes with Chromaprint installed and it is the recommended way to get up and running if installing local dependencies is not desirable.

## Usage
Fingerprinter exposes a simple CLI interface to use for interacting with the binary.
To see the available commands run
```
Usage:
  fingerprinter [command]

Available Commands:
  acoustid    Generate an audio fingerprint and queries the AcoustID API to find matching recordin ID(s)
  fpcalc      Calculates the fingerprint of the input audio file
  help        Help about any command
  mblookup    Queries the MusicBrainz API and returns recordings and releases metadata associated with a recording ID
  verify      Verifies input audio metadata and returns the associated relase(s) info if a match was found

Flags:
  -h, --help   help for fingerprinter
```

## Docker
The Dockerfile can be used to build and run the application and automatically takes care of installing all the required dependencies.
Since teh dockerfile mounts the application directory it is possible to add a folder with audio files so that they can be used from inside the container for testing purposes.
