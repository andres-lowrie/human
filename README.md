human
======

Translate stuff from `Machine -> Human` and back again `Human -> Machine`

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
