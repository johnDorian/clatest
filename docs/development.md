# Development


## Contributing

This is a small project, and as such the scope of the project is designed to be narrow in order to minimize effort required to maintenance. If you find an issues/errors with the tool, feel free to open an issue. 

## Building

If you're keen on building/testing the tool, it's highly recommended to install [taskfile](https://taskfile.dev) in order to build, run the tests and serve the local version of the documentation. 

```bash
git clone https://github.com/johnDorian/clatest.git
# build using task
task build
# run the tests
task test
```
### Documentation

These docs are built using [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/). All the docs are in the [/docs](https://github.com/johnDorian/clatest/tree/master/docs) folder. You can run the docs locally using: 

```bash
# you can view the logs http://localhost:8000
task docs-serve 
```

