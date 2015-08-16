# Description
Sailabove Command Line

# How to build
```
git clone ssh://git@stash.ovh.net:7999/sailabove/sailgo.git && cd sailgo
go get && go build && ./sailgo -h
```

# Roadmap Rewrite
## TODO

TODO based on python version of sail 0.5.5

sail containers          Containers
sail containers ps       List docker containers
sail containers inspect  Inspect a docker container
sail containers attach   Attach to a container console
sail containers logs     Fetch the logs of a container

sail compose             Docker compose
sail compose up          Create and start containers
sail compose get         Export Docker compose receipt

sail services                Services
sail services add            Add a new docker service
sail services rm             Delete a docker service
sail services attach         Attach to the console of the service containers
sail services logs           Fetch the logs of a service
sail services ps             List docker services
sail services inspect        Inspect a docker service
sail services redeploy       Redeploy a docker service
sail services stop           Stop a docker service
sail services start          Start a docker service
sail services scale          Scale a docker service
sail services domain-list    List domains on the HTTP load balancer
sail services domain-attach  Attach a domain on the HTTP load balancer
sail services domain-detach  Detach a domain from the HTTP load balancer

sail repositories        Repositories
sail repositories list   List the docker repository
sail repositories add    Add a new docker repository
sail repositories rm     Delete a repository

sail networks            Networks
sail networks list       List the docker private networks
sail networks inspect    Inspect the docker private networks
sail networks add        Add a new private network
sail networks range-add  Add an allocation range to a private network
sail networks rm         Delete a private network

## DONE TO TEST
sail apps domain-attach       Attach a domain on the HTTP load balancer
--> sail application domain attach applicationName domainName

sail apps domain-detach       Detach a domain from the HTTP load balancer
--> sail application domain detach applicationName domainName

## DONE & TESTED
Configuration

me                  Account
sail me show        show acount details
sail me set-acls    Set ip based account access restrictions

sail apps               Applications
sail apps list          List granted apps
sail apps inspect       Details of an app
sail apps domain-list   List domains and routes on the HTTP load balancer
