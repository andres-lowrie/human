# Formats

## Calling Pattern

Abstractly, calling looks like this

```shell
human <direction> <format> <args> <input>
```

_where_:

`direction`: is one of `-i` (into) or `-f` (from) indicating what type of input to expect

`format`: is one of the formats outlined further down in this document

`args`: are optional settings for a given format which can be _flags_ or _options_ as per each format

`input`: is a string to operate on

### direction

Controls whether the parsers are going to translate the `<input>` into a human
or machine format or vice versa.

The idea is that the entire command should be read from left to right which
should allow the direction option to make more sense hopefully:

> Translate human format _into_ machine format for \<input\>
>
> `human -i machine-format <input>`

Going the other way it would be:

> Give me human format _from_ this machine-format for \<input\>
>
> `human -f machine-format <input>`

### args

Args is shorthand for _flags_ and _options_, each format defines its own set
of arguments as needed. We should strive to make arguments optional when
possible to keep the calling pattern simple.

### format

Formats are what human is translating into and from

@TODO make a website out of the markdown files , parse and load man pages
