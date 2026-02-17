# How to use the development environment

Pre-requisites:

* Follow Getting Started steps at `README.md` file.

## Accounts

* Quay

### Quay.io

Go to <https://quay.io/> and login with your Red Hat user account.

* Create a new repository.
* Go to <https://quay.io/repository/<rh> user>/<repository>?tab=settings
* Create a new robot account and give the account write permissions to your brand new image repository.
* Create the secrets file: `cp ./scripts/mk/private.example.mk secrets/private.mk`
* Fill out the login and token in `secrets/private.mk`

## Install python dependencies

    make .venv
    source .venv/bin/activate
    pip install -r requirements-dev.txt

## Launching request against the API

From VSCode or any IDE which support the .http files, we could open
the scripts at `scripts/http/public.http` and launch the http commands
found there.

> You will need yq tool installed to read values in Makefile from `configs/config.yaml` file.
> See: <https://github.com/mikefarah/yq#install>

### Locally

To test our API locally, we can start the service by `make compose-up run` and launching
the test script below:

    ```sh
    ./test/scripts/local-todos-create.sh
    ```

You can use the `.http` files.
