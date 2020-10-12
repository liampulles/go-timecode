<h1 align="center">
  go-timecode
</h1>

<h4 align="center">A library for handling timecodes of the form HH:MM:SS.mmm</a></h4>

<p align="center">
  <a href="https://travis-ci.com/liampulles/go-timecode">
    <img src="https://travis-ci.com/liampulles/go-timecode.svg?branch=main" alt="[Build Status]">
  </a>
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/liampulles/go-timecode">
  <a href="https://goreportcard.com/report/github.com/liampulles/go-timecode">
    <img src="https://goreportcard.com/badge/github.com/liampulles/go-timecode" alt="[Go Report Card]">
  </a>
  <a href="https://codecov.io/gh/liampulles/go-timecode">
    <img src="https://codecov.io/gh/liampulles/go-timecode/branch/main/graph/badge.svg" />
  </a>
</p>

## Usage examples

### Parse timecodes

```go
import "github.com/liampulles/go-timecode"

t, err := timecode.Parse("01:02:03.456")
// or...
timecode.Parse("01:02:03,456")
timecode.Parse("-01:02:03.456")
timecode.Parse("01:02:03")
```

### Perform arithmetic on timecodes

```go
import "github.com/liampulles/go-timecode"

t, _ := timecode.Parse("01:02:03.456")
t2 := t + timecode.Hour
// t2 is 02:02:03.456

// Apply PAL slowdown
t3 := Timecode(t2 * (23.976/25.0))
```

### Format timecodes

```go
import "github.com/liampulles/go-timecode"

t, _ := timecode.Parse("01:02:03.456")

str := t.FormatDot()
// str is 01:02:03.456

str2 := t.FormatComma()
// str2 is 01:02:03,456

str3 := t.Format(true, ";")
// str3 is 01:02:03;456

str4 := t.Format(false, "")
// str4 is 01:02:03
```

## Contributing

Please submit an issue with your proposal.

## License

See [LICENSE](LICENSE)
