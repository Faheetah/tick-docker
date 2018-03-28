InfluxDB Docker Stack
=====================

Components
----------

**TICK Stack**
* Telegraf: Monitoring agent similar to collectd, can funnel to InfluxDB as well as Graphite, Datadog, and others
* InfluxDB: Performant time series database with a rich query language, accepting several inputs including Telegraf, it's native agent
* Chronograf: Administration UI for InfluxDB. Can configure Kapacitor, query InfluxDB, explore Telegraf hosts, set retention policies and continuous queries, and create dsahboards
* Kapacitor: Automation service that complements InfluxDB. Completely decoupled like all other components, and can accept a range of inputs (including reading from InfluxDB or from Telegraf directly) and send alerts to many integrations (Pagerduty, Slack, Alerta and by extension Golerta)

**Complementary Services**
* Grafana: A more featured dashboard that can read from many sources, including InfluxDB, Datadog, Elasticsearch, and more. Has an expansive plugin community.
* Golerta: Fork of Alerta written in Golang. Meant to be a simple alerting dashboard fed from Kapacitor and a source of truth for alert states. Integrates with Pagerduty and includes alert correlation and flap detection similar to Kapacitor
* Breakit: Under the breakit/ folder, a Go binary that serves a few endpoints to generate fake stats. Includes seasonal data, trending, and random, plus an endpoint to return various error codes.

Prerequisites
-------------

Install Docker for the destined platform:

**Windows**

Install [Docker for Windows](https://docs.docker.com/docker-for-windows/) or the older [Docker Toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/) if HyperV cannot be used. Docker Compose should already be installed, but if not get it [from the docker site](https://docs.docker.com/compose/install/).

**macOS**

Similar to Windows, use either [Docker for Mac](https://docs.docker.com/docker-for-mac/) or [Docker Toolbox](https://docs.docker.com/toolbox/toolbox_install_mac/). Docker Compose should already be installed, but if not get it [from the docker site](https://docs.docker.com/compose/install/).

**Linux**

Find the distro on the site, such as [Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/) or [CentOS](https://docs.docker.com/install/linux/docker-ce/centos/) for installing and also install [Docker Compose](https://docs.docker.com/compose/install/)

**WSL**

AKA Bash on Windows. These instructions assume a Docker Machine named "default" for user "main" and that TCP connections are enabled for Docker Machine. Follow the Ubuntu instructions to install Docker and Compose. Add the following to the Linux environment's ~/.bashrc

```
export DOCKER_TLS_VERIFY="1"
export DOCKER_HOST="tcp://192.168.99.100:2376"
export DOCKER_CERT_PATH="/mnt/c/Users/main/.docker/machine/machines/default"
export DOCKER_MACHINE_NAME="default"

alias docker-machine=docker-machine.exe

mntc(){ sudo mount --bind /mnt/c /c; cd "$(echo $PWD | sed -e 's/\/mnt//')"; }
```

In a project, run ```mntc```, the reason for this is the docker-machine VM does not work properly with the /mnt directory prefix, so volume mounts will not work. The environment variables tell Linux' Docker instance to point to the Windows Docker Machine instance. The docker-machine alias will allow control of the Windows Docker instance without having to leave bash, aka running ```docker-machine start default``` in Bash.

Running
-------

Start Docker. Navigate to this directory. Run ```docker-compose build``` to build Golerta, then ```docker-compose up -d``` to start the services. Navigate to the Chronograf instance, AKA http://localhost:8888.

Note: if using docker-machine, get the IP from ```docker-machine ls``` instead, this will be the base IP for all URLs instead of localhost). Within the Docker infrastructure, the following names are used: telegraf, influxdb, chronograf, kapacitor, rethinkdb, golerta. These names can be referenced directly as hosts instead of IPs as Docker links these containers together.

Chronograf should already have InfluxDB and Kapacitor configured, and Golerta is set to no auth. If InfluxDB or Kapacitor are not configured, use *http://influxdb:8086* and *http://kapacitor:9092*. If configuring Golerta, edit the Kapacitor configuration and set Alerta with the token from ```docker-compose exec golerta /golerta createAgentToken example-secret``` (if you change the secret key, use that instead), copy the output into the *Token* field, hit *Save Changes*, then *Send Test Alert*.

Navigate to http://localhost:5608 and login with username **gauss** and password **password** to see the test alert opened. From here you can add additional Telegraf or Kapacitor configs, either by modifying *docker-compose.yml* or adding Kapacitor alerts in Chronograf under Alerts > Manage Tasks.

Access Grafana with http://localhost:3000 and login with username **admin** password **admin**. Add the InfluxDB database as a data source with *http://influxdb:8086*, using the *telegraf* database. Feel free to create dashboards and play with the data.
