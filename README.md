human
======

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

For full breakdown see [commands](cmds/README.md)

## End to End Testing

Organization of files looks like this

  ```
  e2e/
    runner.py     # reads shell scripts and runs them
    figtures/
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

- For each `case`, only the `test` property is mandatory, everything lese is optional
- For each `test`, if the script returns non-zero (`0`) then that will be marked as failed and you'll get both standard error and standard out for said script to the screen
- `@TODO` All the tests run asynchronously
- The following tokens will be replaced by the runner


| Token       | Replace with                                |
| -----       | ------------                                |
| `%%human%%` | The path to the binary output by `go build` |
