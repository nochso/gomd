# gomd
Cross-platform Markdown editor written in Go.

- Edit files in your browser: where Markdown usually ends up anyway.
- No internet connection needed though. It stays all on your computer.

## Installation
### From binaries
You can download the ready-to-use binaries on the [release page](https://github.com/nochso/gomd/releases) here on Github.

### From source
    $ go get -u github.com/nochso/gomd

## Usage
Open an existing file and edit it:

    $ gomd todo.md

See the command line help for more:
```
$ gomd --help
usage: gomd [<flags>] [<file>]

Flags:
      --help       Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=1110  Listening port used by webserver
      --version    Show application version.

Args:
  [<file>]  Markdown file
```

## License
This project is licensed under the MIT license. See the [LICENSE](LICENSE.md) file for the full license text.

## Credits
* [SimpleMDE](https://github.com/NextStepWebs/simplemde-markdown-editor) - WYSIWYG*ish* MD editor written in JS
