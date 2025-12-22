# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

pform is a Go library for decoding `application/x-www-form-urlencoded` HTTP request bodies into Go structs or maps. It uses struct tags to map form fields and supports custom type unmarshalling.

## Commands

```bash
# Run all tests
go test ./...

# Run a specific test
go test -run TestDecoder_Decode

# Run tests with verbose output
go test -v ./...
```

## Architecture

**Core Components:**

- `decode.go` - Main decoder implementation with `NewFormUrlDecoder()` for HTTP requests and `NewDecoder()` for raw `url.Values`. The `Decoder` type handles both struct and map targets.
- `parsers.go` - Type-specific parsing functions (`setValueInt`, `setValueUint`, `setValueFloat`, `setValueComplex`, `setValue`) that handle reflection-based value assignment.
- `error.go` - Custom error types: `ParseError` for type conversion failures and `RequiredFieldError` for missing required fields.

**Struct Tag Format:** `form:"field_name[,option]"`
- Options: `omitempty` (skip if empty), `required` (error if missing), `-` (ignore field)

**Custom Types:** Implement `Unmarshaller` interface with `UnmarshalValue(string) error` method for custom parsing logic.
