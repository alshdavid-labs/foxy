# Foxy Task Runner

A simple yaml based task runner
create a `foxy.yml` in the root of your project and run

```
foxy
```

```
foxy task-name
```

## Example

Basic `foxy.yml`

```yaml
build:
  default: true
  env_file: env/local.env
  env:
    MY_ARG: set value
  steps:
    - go build -o ../output main.go
```

## Schema

```yaml
key:
  default: boolean 
  env_file: string
  env:
    key: string
  parallel: boolean
  steps:
    - string
```
### Keys:

| Key        | Value |
| ---------- | ----- |
| `default`  | Will tell foxy to run this task if no task matches task argument |
| `env_file` | Will import an env file into the environment inside the steps |
| `env`      | A set of key/value pairs to be used in the enviroment within the steps - will override values in `env_file` |
| `parallel` | Will execute the steps in parallel |
| `steps`    | A sequence of command line steps to execute |


## Platform specific overrides

Run with either `foxy` or `foxy build`

```yaml
build:
  default: true
  steps:
    - echo Default for unspecified platforms

build(windows):
  steps:
    - echo will override default and execute on Windows

build(darwin):
  steps:
    - echo will override default and execute on MacOS
```