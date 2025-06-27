# Contributing guidelines

Hi there! Thank you for thinking on contributing on this repository.

This document will help you on making your development contributing peaceless
and drive you to be aligned with the repository practices.

We encourage to make small chunk of changes, if your contribution has many
changes we encourage to split in smaller units, create an issue for each of them
and for every fix and feat commit be sure the tests pass. Pipeline should pass
and it makes several checks automatically for you, to prevent missing things,
but more details are explained later.

Enjoy, and I hope this helps you!

## Getting Started

**Once tasks**

- Fork the repository.
- Clone your forked repository (we recommend to name this remote as `upstream`).
- Add an additional remote `downstream` pointing to this repository, using
  `https` protocol.
- Install required packages:
  - rpm based distro: `sudo dnf install golang make podman podman-compose delve`
  - apt based distro: `sudo apt install golang make podman podman-compose delve`
- Create and edit `configs/config.yaml` by: `cp -vf configs/config.example.yaml
  configs/config.yaml`
- Choose your favourite IDE of text editor (the repository has some settings for
  vscode, we are happy to receive additional settings), and install it.
- Install deps by: `make tidy`
- Install tools by: `make install-tools`
- Build infra containers and start them by: `make compose-build compose-clean
  clean build compose-up`

**Day to day**

- Getting help: `make help`
  we don't recall all the targets, we don't expect anyone to memorize them
  either, so we use this very often, and let our minds to focus on what we
  enjoy :)
- Pull changes before creating your local branch: `git pull downstream`
- Create your local branch: `git checkout -b my-fix-or-feature main`
- Code your changes.
- Align changes to format: `make format`
- Check linters by: `make lint`
- Check unit tests by: `make test-unit`
- See coverage report by: `make coverage`
- Start local infrastructure by: `make compose-up`
- Check integration tests by: `make test-integration`
- Create commits, and reorg them.
- When ready to push your changes, rebase on top of the main branch:
  `git pull downstream && git rebase main`
- Push your changes: `git push -u upstream my-fix-or-feature`
- If pipeline fails, make necessary changes until fix it.
- Create your PR.
- ping @avisiedo @little-engineer-2025 if no response on PR for 3 days.
- We will ping you when the change is reviewed.
- When the changes are addressed or understanding of the changes are resolved,
  we will mark the thread as resolved.
- When everything is resolved, we will merge the PR.
- Thanks a million for your time and contribution!

> If you want to use TDD, then start implementing some tests to fail, and cover
> every small change. Run tests by `make test` or `make test-unit` and `make
> test-integration` make targets and add the necessary code.

See Development Docs at: [dev/INDEX.md](dev/INDEX.md)

## Commit Messages

We try to use conventional commits for the commit messages; this helps on
automate releases, logchanges. The pipeline will check this in the future.

See: [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

In short, we are trying to keep the format below:

```raw
<type>: <title>

Explain why the change and the motivation.

#issue-number

[Co-authored-by: First Lastname <user@example.com>]
Signed-off-by: First Lastname <user@example.com>
```

- **type** is one of the types recommended at conventional commits, the more
  common are: `fix`, `feat`. When we write a fix, we add the unit tests related
  to the change. We try to do the same for feat. When we provide new tests for
  new corner cases or increase the coverage, we use `test` type.
- **title** short description trying to keep the whole title no longer than 50
  chars; use active voice.
- In the body, provide information about the reason of the change (everyone will
  value this description in the future)
  - if it is a fix, how the fixed situation is resolved.
  - if it is a feat, why the change.
- We encourage the recognition of the contributors, so if someone else provide
  value to the changes, we recommend to use the `Co-Authored-By` footer for it.
- Finally, we recommend to sign your commits `git commit -s`.

## Relax and breath!

Don't get in panic, this document is to help you; if it is the first time
contributing, it could be overhelming; this content will help to avoid memorize
and to help us to reference to it when necessary, so we save our time and your
time.

Any question, do not hesitate to contact us!

Thanks in advance for your contributions, and enjoy the development experience!

Cheers!

