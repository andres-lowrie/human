human
======

@leftoff
- fix --into and --from options. The unit tests pass but the not the command so there's a bad assumption somewhere that lets the api work but not the cli

[Translate](Translate) stuff from `Machine -> Human` and back again `Human -> Machine`

## TL;DR

Given an input, it translates it to a more human form

```
human 1000000
> 1_000_000
> 1M

human "0 0 1 1 *"
> yearly 

human aGVsbG8gd29ybGQK
> hello world
```

You can also translate human into machine

```
human 5GB
> 5000000000

human "run five minutes after midnight every day"
> 5 0 * * *
```

## File Organization

```
e2e/      # End-to-End testing
cmds/     # The auxiliary arguments and commands which are not `format`s (for example `human -h`)
format/   # The things this knows how to read and write. The interface main.go is calling
io/       # Input/Output functions and types
parsers/  # The functions that do the actual work, called by `format`s
```

## Code Concepts

There are 3 core interfaces at play:

- The [Command](cmds/cmds.go#L13)
- The [Format](format/format.go#L8)
- The [Parser](parsers/parsers.go#L9)

The `Parser` is the thing that actually does the translating work.

A `Format` houses a `Parser` or many parsers depending on the complexity of the format, it acts as a wrapper between a `Command` and a `Parser`

The `Command` is what something has to implement in order to show information via `help` and to actually get executed by the main function. By design every `Format` also implements `Command`

## End to End Testing

Organization of files looks like this

  ```
  e2e/
    runner.py     # reads shell scripts and runs them
    fixtures/
      some-file1
      some-file2
    $parser1/
      behavior1.yaml
      behavior2.yaml
    $parser2/
    ...
  ```

The `runner.py` script will:
  - run string replacement
  - execute the steps of a test
  - show outcomes

### Requirements

- python3
- The following tool needs to be in your path https://github.com/bruceadams/yj

### Quickstart

```
cd e2e
mkdir tmp
export E2E_TMP_DIR=$repo/e2e/tmp
./runner.py $paths-to-yaml-file
```

> template
```
---
suite:
  cases:
    - name: "It should do something"
      setup: |
        echo "setting up"

      test: |
        var="is the test"
        echo "this ${var}"

      cleanup: |
        echo "cleaning up"

    - name: "It should do something else"
      test: |
        echo "doing something else"
```

- For each `case`, only the `test` property is mandatory, everything else is optional
- For each `test`, if the script returns non-zero (`0`) then that will be marked as failed and you'll get both standard error and standard out for said script to the screen
- `@TODO` All the tests run asynchronously
- The following tokens will be replaced by the runner


| Token       | Replace with                                |
| -----       | ------------                                |
| `%%human%%` | The path to the binary output by `go build` |
