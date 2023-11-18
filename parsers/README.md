# Some implementation details

## Cron

### General Notes of the current state of things
The code is doing a very naive approach that is passing the current test cases
but I'll have to circle back at some point and rethink it to see if there's a
better approach; for now I'll say this is "good enough" for a v0.

That being said, it can only go from `Machine -> Human` right now because I
haven't quite figured out how to go the other way yet (I'll come back to it)

### Code Flow

The function that basically does all the work right now is `DoFromMachine` so
I'll try to explain that function since the other ones should be straight
forward

The input string is broken down into `component`s that tells us what type of
syntax is being used ie lists, ranges, single, values etc. 

The components are then processed from left to right to build the string that
is returned.

The actual string prepartion is done by preparing a template string (using go
built in template syntax) using the component structs to make decisions in terms
of what words to use and how to convert from numbers to words. This template
string is then parsed for each component and then concatenated into the actual
final string

Abstractly like this:

```
                           +----comp----+
                           |            |
                  + comp + |  + comp +  |
                  |      | |  |      |  |  
parse input (eg:  ' * *    *      *     *')

for each component in order from left to right
  prepare a template string

 combine the rendered templates into the final output
```
