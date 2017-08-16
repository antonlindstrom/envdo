# envdo

Manage environment variables with a command.

envdo came into existance when I realized I had multiple environment
variables set for different work, for home and for side projects. This meant
having one default and then use aliases or manually try to switch the
environment variables. Since this is a hassle I wanted to make it easy to
switch while still maintain some security.

The goal is to have all environment variables in gpg encrypted files and
import them when needed.

### Examples

List profiles:

	envdo ls

Use project1 environment variables with command `env`:

	envdo project1 env

Add new environment variable file:

	envdo -r GPGID add profile

### Installation instructions

The command can be installed with `go get`:

	go get -v github.com/antonlindstrom/envdo

## Bugs

I'm sure there are a lot of bugs or surprises, please file anything you find
and I will try to do my best to solve them.

Currently the code is very hacky and there are no tests, this was made in an
evening and works for my uses.

## Acknowledgements

This program has borrowed most of its ideas from
[Pass](https://www.passwordstore.org/). This software can be managed by Pass
and I suggest linking `~/.password-store/envvars` to `.envdo`.

## Author

This is maintained by [Anton Lindstrom](https://www.antonlindstrom.com).

## License

See [LICENSE](LICENSE) file in the current directory.
