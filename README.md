# Personal Site / CV / Blog

This repository houses my personal blog. It is comprised of a few static HTML pages (Go templates)
and a couple of dynamic pages.

The server reads a list of available posts upon startup. Each request to a given post url prompts a
lookup for each request to a given post that on-demand renders the markdown as HTML into the `post.html`
Go template.

## Designing

A couple of notes around design specifics for the site.

#### Post List

Each card as they exist from the Spectre CSS component fits an image with dimensions:

|  W  |  H  |
|-----|-----|
| 472 | 400 |

## Testing

Any available tests can be run with

```shell
make test
```

## Running

#### No build

To run the server either for testing & development or in production through the `go` command, just run:

```shell
make run
```

#### Build Binary (Linux/macOS)

This command will produce a binary of the server program in the `./bin` directory called `server`.

This binary still requires the project's tree structure, templates, static files, and any desired posts
as they exist in this development repository.

```shell
make build
```

#### Run Binary

This can only be run after the build stage.

```shell
make start
```

## Deploying

TBD.

## Other Commands and Tooling

See the [Makefile](https://github.com/robertfairley/rf-19-go/blob/master/Makefile)

## License

[MIT](https://github.com/robertfairley/rf-19-go/blob/master/LICENSE)

---

&copy; 2019 Robert Fairley
