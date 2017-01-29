#docker-health
[ ![Codeship Status for hyleung/docker-health](https://app.codeship.com/projects/dbb6e190-c808-0134-1c2a-5a435d0f4766/status?branch=master)](https://app.codeship.com/projects/198973)

Small command-line utility for working with Docker's built-in health checks.

##Usage

###Connecting to the Docker daemon

If the `$DOCKER_HOST` environment variable is found, this utility will attempt
to create the Docker daemon connection using environment variables (for example,
if you're using `docker-machine`). Otherwise, the utility will attempt to 
connect via socket (i.e. `unix://var/run/docker.sock`).

###Inspect the health check status of containers

```
NAME:
   docker-health inspect - Inspect the Health Check status of a container

USAGE:
   docker-health inspect [command options] [arguments...]

OPTIONS:
   --all, -a  Show Healthcheck status for all running containers
   --verbose  Show detailed health check information on containers
   --log, -l  Enable log output
```

Example:

Inspect the health of a single comtainer:

```
docker-health inspect <container name>
```

This command exits with a non-zero code if the named container is not
in a `health` state.

Inspect the health of _all_ containers on the daemon:

```
docker-health inspect --all
```

###Wait on containers to enter healthy status

```
NAME:
   docker-health wait - Wait until a container enters Healthy status

USAGE:
   docker-health wait [command options] [arguments...]

OPTIONS:
   --all, -a        Wait on Healthcheck status for all running containers
   --timeout value  Wait timeout, in seconds (default: 60)
   --log, -l        Enable log output
```

Example:

Wait for a single container to enter healthy state:

```
docker-health wait <container name>
```

The command will exit with a non-zero code if the named container either fails
to enter a `healthy` state within the timeout (by default, 60 seconds) or enters
an `unhealthy` state.

Wait for _all_ containers to enter a healthy state:

```
docker-health wait --all
```

The command will exit with a non-zero code if any ontainer either fails
to enter a `healthy` state within the timeout (by default, 60 seconds) or enters
an `unhealthy` state.
