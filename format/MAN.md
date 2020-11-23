## Logging

`-v`: Verbose

You can pass multiple `v`s to denote which level of logging you want:

| Count | Log Level | Description                                                                    |
|-------|-----------|--------------------------------------------------------------------------------|
| v     | Info      | Show high level flow of program                                                |
| vv    | Warn      | Show potential issues like api changes, places where the program guessed, etc. |
| vvv   | Debug     | Shows all the things. Can also be used to trace the program                    |

Turning on a "higher" level of logging includes all the levels before it so for example `human -vvv` will show _"Debug, Warn, and Info"_ logs.
