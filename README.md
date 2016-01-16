# MQ Resource
A concourse resource for subscribing and publishing into a message-queue.

## Source Configuration
* `uri`: *Required.* The location of the message-queue
* `queue`: *Required.* The type of queue plugin to use. Supported plugins are `nsq` and `nats`.
* `sub`: *Optional.* The subscribing topic
* `pub`: *Optional.* The publishing topic
* `params`: *Optional.* A generic String Key/Value map used differently by each plugin

### `nsq` Plugin
This is a plugin for [nsq.io](http://nsq.io). `nsq` guarantees at-least-once delivery, so subscribing in a resource is possible. `Dockerfile.test.nsq` has an example configuration under `SOURCE` env variable.

#### Parameters
* `channel`: *Optional.* Subscribing to a specific channel under the previously defined `sub` topic. (Defaults to a random string)

### `nats` Plugin
This is a plugin for [nats.io](http://nats.io). `nats` does not guarantee at-least-once delivery, so subscribing is not possible. This plugin can only be used when using a `put`. `Dockerfile.test.nats` has an example configuration under `SOURCE` env variable.

### Example
```yaml
resources:
- name: queue-name
  type: mq
  source:
    queue: nats
    uri: nats://0.0.0.0:4222
    pub: caring_listeners
```
Receiving subscriptions to `message.json`
```yaml
- get: queue-name
  params: { file: "message.json" }
```
Publishing out the message:
```yaml
- get: a-resource
- put: queue-name
  params: { file: "output.json" }
```

### Concourse Deployment
Modify bosh deployment manifest by adding the following to `groundcrew.additional_resource_types`:

```yaml
- image: docker:///aminjam/mq-resource
  type: mq
```

## Behavior

### `check`: Check for new messages in the queue
The message queue is scanned for any `json` encoded message for the subscribing topic. If the returned object is the same as the current version, no result is returned.

### `in`:  Writing the version into a file
Write the requested `json` string Key/Value into a file

#### Parameters
* `file`: *Optional.* The path and name of the file to be written (Defaults to `message.json`)

### `out`: Publish a message to the queue
Publish the content of the file to the queue.

#### Parameters
* `file`: *Optional.* The path and name of the file to be written (Defaults to `message.json`)

## Building from source
Built with Go 1.5.2 and `GO15VENDOREXPERIMENT` flag.
```shell
git clone https://github.com/aminjam/mq-resource && cd mq-resource
make update-deps
make build
make test #Optional PLUGIN="nsq"
```

### Notes
This is my first attempt in making a concourse-resource. Since this a message queue, there is no guarantee of receiving the same version/message, unlike `git-resource`. This is considered an experimental resource.
