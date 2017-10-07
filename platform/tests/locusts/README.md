# Load test for pinger through the AMP proxy layer

## add an IP to your host loopback device

This section is required if you run the load test on the same machine as the AMP stack.

This is not required on Docker for Mac, see paragraph _Configure the load test address_.

The load test tool is running inside a container and requires access to the host interface.
127.0.0.1 is the container from inside the container, so we need another IP on the host loopback device.

If you don't already have one, add an alias with this command:
```
sudo ifconfig lo0 alias 172.71.0.1
```
It won't persist on reboot.

## Create a cluster

```
amp cluster create
```

## Deploy the pinger stack

```
amp stack deploy -c ./examples/stacks/pinger/pinger.yml
```

## Configure the load test address

Edit the file `platform/tests/locusts/scripts/locust.config.json` and set the target to your loopback address.

On Docker for Mac, you can use the address `docker.for.mac.localhost` instead.

## Run the load test

```
cd platform/tests/locusts
docker-compose -f docker-compose.yml up
```

- access the UI on 127.0.0.1:8089
- run a test with you choice of parameters (30 users)
