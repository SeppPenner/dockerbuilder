# DockerBuilder

DockerBuilder builds automatically a new tagged container and pushes it to
the Docker Index when you create a new tag on GitHub.

**Important:** DockerBuilder is work-in-progress and is not usable yet!

## Installation

### Requirements

* A working docker environment
* Git
* SSH keypair setup for GitHub
* A ``.dockercfg`` with Docker Index credentials (use ``docker login`` to create)

### Installing DockerBuilder

```sh
$ go get github.com/brocaar/dockerbuilder
```

### Adding the webhook

1. Login into GitHub, and go to the configuration of the repository that you
   would like to build with DockerBuilder.

2. Go to *Webhooks & Services* and click the *Add Webhook* button.

3. Set the following values in the form:

    * **Payload URL**: *http://yourhostname.tld*/github.com/hook
    * **Content type**: application/json
    * **Secret**: the value that you have set in ``BUILDER_GITHUBSECRET`` (see
      configuration)
    * **Which events would you like to trigger this webhook?**: Click
      *Let me select individual events* and select *Create*

4. Click Add Webhook.


## Configuration

Configuration is done by setting environment variables. The following set
of variables can be set:

#### ``BUILDER_WORKDIR``

Default: ``/tmp``. The work directory where DockerBuilder will clone the
repositories.

#### ``BUILDER_NUMWORKERS``

Default: the number of available CPU's. The number of concurrent workers.

#### ``BUILDER_BINDADDRESS``

Default: ``0.0.0.0:5000``. The address for binding the HTTP server.

#### ``BUILDER_TASKQUEUESIZE``

Default: ``100000``. The size of the task-queue. One this limit is reached,
HTTP calls to the webhook will block until there is space in the queue.

#### ``BUILDER_CLEANUPCONTAINER``

Default: ``true``. Cleanup the container after pushing it to the Docker Index.
Setting this to ``true`` saves diskspace but results in slower builds. Setting
this to ``false`` will speedup the builds, but will use more storage.

#### ``BUILDER_DOCKERINDEXNAMESPACE``

Default: not set. The Docker Index namespace, usually this is your
``hub.docker.com`` username.

#### ``BUILDER_GITHUBSECRET``

Default: not set. The secret used for the GitHub webhook.


## Changelog

### v0.1.0 (in development)

* Initial version.
