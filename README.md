# Description
Sailabove Command Line

# How to build
```
git clone ssh://git@stash.ovh.net:7999/sailabove/sail.git && cd sail
go get && go build && ./sail -h
```

# Roadmap Rewrite
## TODO
```
TODO based on python version of sail 0.5.5

sail compose up               Create and start containers
sail compose get              Export Docker compose receipt

```

## DONE TO TEST

```
sail services add            Add a new docker service
sail services redeploy       Redeploy a docker service
sail services start           Start a docker service
sail services scale           Scale a docker service

sail containers logs          Fetch the logs of a container

sail apps domain-attach       Attach a domain on the HTTP load balancer
--> sail application domain attach applicationName domainName

sail apps domain-detach       Detach a domain from the HTTP load balancer
--> sail application domain detach applicationName domainName

sail networks range-add  Add an allocation range to a private network

sail services logs            Fetch the logs of a service
sail services domain attach   Attach a domain on the HTTP load balancer
sail services domain detach   Detach a domain from the HTTP load balancer

sail repositories add         Add a new docker repository
sail repositories rm          Delete a repository

```

## DONE & TESTED

```
Configuration

me                  Account
sail me show        show acount details
sail me set-acls    Set ip based account access restrictions

sail apps               Applications
sail apps list          List granted apps
sail apps inspect       Details of an app
sail apps domain-list   List domains and routes on the HTTP load balancer

sail containers          Containers
sail containers attach   Attach to a container console
sail containers ps       List docker containers
sail containers inspect  Inspect a docker container

sail repositories list   List the docker repository

sail services attach         Attach to the console of the service containers
sail services inspect        Inspect a docker service
sail services ps             List docker services
sail services rm             Delete a docker service
sail services stop           Stop a docker service

sail services domain list     List domains on the HTTP load balancer

sail networks add        Add a new private network
sail networks rm         Delete a private network
sail networks list       List the docker private networks
sail networks inspect    Inspect the docker private networks

```

## Bugs SA
```

./sail network add myApp/privateb 172.31.0.0/24
--> Return 200 (OK) instead of 201 (Created)

```
