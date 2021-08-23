# Fingerprinter

Fingerprinter is a CLI tool and go library that can be used for
- Generating fingerprints from audio files.
- Identifying audio fingerprints via the acoustID Web Service.
- Identifying releases and recordings metadata associated with the original file(s) via the MusicBrainz APIs

This tool leverages [Chromaprint](https://acoustid.org/chromaprint) and its associated [Acoustid Web Service](https://acoustid.org/webservice) to generate acoustic fingerprints.
These fingerprints can then be used to verify the origin and content of an audio file.
Specifically the Fingerprinter package generates audio fingerprints to ultimatley determine the author(s), album(s), record label(s) and ISRC codes associated with a recording.

## Status
In Progress

[![Actions Status](https://github.com/ocramh/fingerprinter/workflows/Test/badge.svg)](https://github.com/ocramh/fingerprinter/actions)

## Dependencies
The only required dependency is Chromaprint.
When running the application locally, the Chromaprint executable (`fpcalc`) must be on the `$PATH`.
See the [Chromaprint repo](https://github.com/acoustid/chromaprint) for [downloads](https://github.com/acoustid/chromaprint/releases) and information about how to build it locally.

The provided Dockerfile comes with Chromaprint installed and it is the recommended way to get up and running if installing local dependencies is not desirable.

## Usage
Fingerprinter exposes a simple CLI interface.
To see the available commands run
```
Usage:
  fingerprinter [command]

Available Commands:
  fpcalc      Calculates the fingerprint of the input audio file
  acoustid    Queries the AcoustID API to match a fingerprint with a recording ID(s)
  help        Help about any command
  mblookup    Queries the MusicBrainz API and returns metadata associated with a recording ID
  verify      Verifies input audio metadata and returns the associated release(s) info

Flags:
  -h, --help   help for fingerprinter
```

## Docker
The Dockerfile can be used to build and run the application and automatically takes care of installing all the required dependencies.