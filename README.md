# plugin-morphe-seed-data

Kalo plugin that reads Morphe model definitions and generates realistic PSQL seed data (`INSERT` statements) using [gofakeit](https://github.com/brianvoe/gofakeit).

## What it does

- Reads `.mod` and `.enum` files from a Morphe registry
- Generates `INSERT INTO` SQL files, one per table, numbered by dependency order
- Populates enum lookup tables with all defined entries
- Generates model rows with realistic fake values based on field types and format attributes
- Respects `ForOne` / aliased FK relationships by referencing valid parent IDs
- Topologically sorts models so parent tables are seeded before children
- Produces deterministic output via a configurable random seed

## Format attributes

The plugin parses `key=value` pairs from the standard Morphe `attributes` list to guide faker selection:

| Attribute | Applies to | Effect |
|-----------|-----------|--------|
| `format=email` | String | Uses `faker.Email()` |
| `format=firstName` | String | Uses `faker.FirstName()` |
| `format=lastName` | String | Uses `faker.LastName()` |
| `format=phone` | String | Uses `faker.Phone()` |
| `format=url` | String | Uses `faker.URL()` |
| `format=company` | String | Uses `faker.Company()` |
| `format=city` | String | Uses `faker.City()` |
| `format=sentence` | String | Uses `faker.Sentence()` |
| `regex=<pattern>` | String | Uses `faker.Regex(pattern)` |
| `minLength=N` | String | Pads output to at least N characters |
| `maxLength=N` | String | Truncates output to at most N characters |
| `min=N` | Integer/Float | Lower bound for generated number |
| `max=N` | Integer/Float | Upper bound for generated number |

When no format attribute is set, the plugin infers a faker function from the field name (e.g. `Email` -> email, `FirstName` -> first name) or falls back to a generic sentence.

### Example `.mod` with format attributes

```yaml
name: Person
fields:
  ID:
    type: AutoIncrement
  FirstName:
    type: String
    attributes:
      - format=firstName
  Email:
    type: String
    attributes:
      - format=email
  TaxID:
    type: String
    attributes:
      - regex=^\d{2}-\d{7}$
  Bio:
    type: String
    attributes:
      - optional
      - maxLength=200
identifiers:
  primary: ID
```

## Building

### WASI (for kalo CLI execution)

```bash
# Linux / macOS
cd scripts && bash build.sh

# Windows
cd scripts && build.bat
```

This produces `dist/morphe-seed-data-v1.0.0.wasm` targeting `GOOS=wasip1 GOARCH=wasm`, which the kalo CLI executes via its WASI runtime.

### Native (for local development)

```bash
go build -o dist/morphe-seed-data ./cmd/plugin
```

## Usage

```bash
go run cmd/plugin/main.go '{"inputPath":"./morphe","outputPath":"./seed","config":{"rowCount":5,"seed":42}}'
```

### Config options

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `rowCount` | int | 5 | Number of rows per model table |
| `seed` | int | 42 | Random seed for deterministic output |
| `schema` | string | `""` | PostgreSQL schema prefix |

## Testing

```bash
go test ./... -v
```

The test suite includes:
- **Unit tests** for format attribute parsing, faker field generation, enum compilation, and model compilation (with FK, aliased, enum, optional, and format attribute coverage)
- **Integration test** that compiles a minimal testdata registry and compares output file-by-file against checked-in ground truth
