# Geminic
Using Gemini to Write Git Commits 

## Quick Start
### install
```shell
go install github.com/Beriholic/geminic@latest
```
### configuration
```shell
geminic config
```
Interactive configuration

![](./assets/config.png)


Switching gemini models

```shell
geminic models
```

### usage
after git add files, you can use `geminic` to generate commit message

```shell
geminic
```

![](./assets/root.png)

you can also give gemini commit references with -c
```shell
geminic -c "fix bug"
```

### Credits
  - Thank you Google for such a great [model](https://ai.google/get-started/for-developers)!
