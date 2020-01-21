# Knaxim

knaxim 'knowledge maxim' is a file management using a nlp organizational paradigm. 

combination of server/knaxim and maxsetdev/knaxim-client

## Web Root Volume
#### Basic volume management
```docker volume ls, create, and rm```

#### Inspect a volume
Inspect a volume to get basic information about it.  Mountpoint, Labels, and Driver are some of the fields available.

```docker volume inspect web```
#### Copy a file into a volume
Adding files to a volume is done with a container that is using that volume.  The container does not have to be running, but it does have to be configured to use the volume.

```docker cp statuscode.txt knaxim_server_1:/public```

#### Remove a file from a volume
```docker exec knaxim_server_1 rm /public/statuscode.txt```

