human
======

Machines -> Human
Human -> Machines

## Idea

Given an input, it translates to a more human form

```
human "1000000"
> 1_000_000
> 1M

human "0 0 1 1 *"
> yearly 

human "0 22 * * 1-5"
> 10pm on Weekdays

human "aGVsbG8gd29ybGQK"
> hello world

```

It should also allow for going the other way

```
human "5G"
> 5000000000

human "run five minutes after midnight every day"
> 5 0 * * *
```
