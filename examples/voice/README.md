# go-openai-realtime example for voice

## Prerequisite

This example needs portaudio library to play audio, please install it before running this library.

On macOS:
```bash
brew install pkg-config
brew install portaudio
```

On Linux:
```bash
apt-get install portaudio19-dev
```

## Run

```bash
export OPENAI_API_KEY=<your openai api key>
export SOCKS_PROXY=<your socks proxy> # this optional
go run .
```
