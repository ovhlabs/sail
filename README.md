# Sailabove client

Sailabove.com is a docker hosting solution aiming to be as flexible as a
container and as elegant as a sailboat.

```bash
docker login sailabove.io

# Build
docker tag my-cool-project-image sailabove.io/my-app/my-cool-project-image

# Ship
docker push sailabove.io/my-app/my-cool-project-image

# Run
sail service add my-cool-project-image my-cool-project-service
```

## Setup

1. Grab lastest release for your platform from https://github.com/runabove/sail/releases
2. Make it executable. ``chmod +x sail`` will do the trick on UNix based platforms

To update it, simply run

```bash
sail update
```

## Configuration

``sail`` automatically loads registry's credentials from ``docker`` keyring.
Hence, after a succesfull push to Sailabove's registry, there should be no
need for configuration.

```bash
docker login sailabove.io
```

If you wish to temporarily override a parameter, you may use ``SAIL_HOST``,
``SAIL_USER`` and ``SAIL_PASSWORD`` to respectively force the API endpoint,
the username and the password. Additionally, these parameters may be set via
``--api-host``, ``--api-user`` and ``--api-password``

## Usage

Once you have claimed your private namespace on http://labs.runabove.com/docker and
sucessfuly pushed your first image you may launch and supervice a service
from this template image. For example, taking a ``my-redis`` Docker, let's
create a ``redis`` service:

```bash
sail service add my-app/my-redis-image my-app/redis-service
```

Watch your private cluster's status:

```bash
sail service ps
```

Scale your cluster:

```bash
sail service scale my-app/redis-service --number 2
```

Clear everything:

```
sail service rm my-app/redis-service
```

## Hacking

Sailabove's CLI is written in Go 1.5, using the experimental vendoring
mechanism introduced in this version. Make sure you are using at least
version 1.5.

```bash
export GO15VENDOREXPERIMENT=1
go get github.com/runabove/sail
cd $GOPATH/src/github.com/runabove/sail
go build
```

You've developed a new cool feature? Fixed an annoying bug? We'd be happy
to hear from you! Make sure to read [CONTRIBUTING.md](./CONTRIBUTING.md) before.

## Related links

- **Sign Up**: http://labs.runabove.com/docker
- **Registry**: https://registry.sailabove.io/
- **Get help**: https://community.runabove.com/
- **Get started**: https://community.runabove.com/kb/en/docker/getting-started-with-sailabove-docker.html
- **Documentation**: [Reference documentation](https://community.runabove.com/kb/en/docker/documentation), [Guides](http://community.runabove.com/kb/en/docker/)
- **OVH Docker mailing-list**: docker-subscribe@ml.ovh.net

