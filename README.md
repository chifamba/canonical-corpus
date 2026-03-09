# canonical-corpus

`canonical-corpus` is a robust Go-based tool designed to fetch, parse, normalize, and export biblical and related texts from public-domain sources into a structured Markdown corpus. It supports a wide range of texts, including the Protestant canon, extra-canonical books, and the Dead Sea Scrolls across multiple languages and translations.

## Features

- **Multi-Source Fetching**: Automatically pulls texts from Wikisource, Sacred-Texts, Project Gutenberg, and more.
- **Structured Output**: Generates clean, consistent Markdown files organized by category, book, and translation.
- **Resumable Builds**: Tracks progress to allow interrupted builds to resume without re-downloading existing content.
- **Normalization & Deduplication**: Cleans raw text and deduplicates passages for a high-quality corpus.
- **Rate Limiting & Retries**: Built-in mechanisms to respect source server constraints and handle transient network failures.

## Supported Translations & Languages

- **English**: KJV (King James Version), WEB (World English Bible), ASV (American Standard Version).
- **Greek**: LXX (Septuagint) and GNT (Greek New Testament).
- **Hebrew**: HMT (Masoretic Text).
- **Shona**: BDB (Bhaibheri Dzvene).
- **Extra-Canonical**: Enoch, Jubilees, Tobit, Judith, and many more.
- **Dead Sea Scrolls**: 1QS, 1QM, 1QH, and other major Qumran texts.

## Installation

Ensure you have **Go 1.24+** installed.

```bash
git clone https://github.com/chifamba/canonical-corpus.git
cd canonical-corpus
go build -o corpus-builder ./cmd/corpus-builder
```

## Usage

The tool uses a CLI interface powered by Cobra. You can run it directly using `go run` or by building the binary.

### Build the Corpus
Fetch and compile all books defined in the source catalogue:
```bash
./corpus-builder build
```

Filter by category, language, or translation:
```bash
./corpus-builder build --category extra-canonical
./corpus-builder build --language sn
./corpus-builder build --translation kjv
```

### Verify the Corpus
Check if all expected files exist in the output directory:
```bash
./corpus-builder verify
```

### Configuration
Configuration is managed via `configs/config.yaml`. You can adjust timeouts, retry counts, rate limits, and output directories there.

## Project Structure

- `cmd/`: CLI entry points.
- `internal/`: Core logic (compiler, downloader, parser, etc.).
- `sources/`: Definitions of book metadata and source URLs.
- `pkg/`: Reusable packages (e.g., Markdown writer).

## License

This project is licensed under the MIT License. All source texts fetched by this tool are, to the best of our knowledge, in the Public Domain.
