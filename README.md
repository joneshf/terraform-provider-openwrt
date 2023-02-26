# terraform-provider-openwrt

A [Terraform][] provider for [OpenWrt][].

## Setup

Two tools are used to bootstrap this repo: [`nix`][] and [`direnv`][].
Once those tools are setup,
everything else should work seamlessly.

N.B. Neither tool is actually required.
Everything in this repo _probably_ could work without either.
But, any workflow without both is not currently supported.

### nix

[`nix`][] is used to make working in the repo reproducible.

When we come back to the repo in a few days, weeks, months, or years;
we don't want to have to remember how to get the environment set up properly.
If someone new comes to the repo,
they shouldn't have to jump through hoops to make it work on their computer.

See the [install instructions][nix install] for setting up [`nix`][].

After it is setup, it should be enough to launch the shell: `nix-shell`.

### direnv

[`direnv`][] is used to make working with [`nix`][] a little easier.

Without [`direnv`][],
there is a step that has to happen each time to start working: `nix-shell`.
It's not a huge burden,
 but it's something extra to remember.
It also takes a bit of time once the environment gets large enough.
[`direnv`][] caches the environment of the shell and can automatically spin up the
environment when the directory is entered.
This means that working with the project after the initial setup is seamless,
and more efficient.

See the [install instructions][direnv install] for setting up [`direnv`][].

After it is setup,
it should be enough to allow it to work: `direnv allow`.

## Development

[`make`][] is used to build everything in this repo.
Installation of [`make`][] is taken care of by [`nix`][].

Everything can be built from the top-level:

```sh
$ make build
```

Easier than that,
everything can be built and tested in one step:

```sh
$ make test
```

[`direnv`]: https://github.com/direnv/direnv
[`make`]: https://www.gnu.org/software/make/
[`nix`]: https://nixos.org
[direnv install]: https://github.com/direnv/direnv#install
[nix install]: https://nixos.org/nix/
[openwrt]: https://openwrt.org/
[terraform]: https://www.terraform.io/
