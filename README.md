# TODO: Merge apipaths.txt into this document.

# Knaxim

knaxim 'knowledge maxim' is a file management using a nlp organizational paradigm.

combination of server/knaxim and maxsetdev/knaxim-client

- [Commands](#commands)
- [API](#api)
    - [Status Codes](#status-codes)
    - [Response & Format](#response-and-format)
- [Container Resources](#container-resources)
    - [Volumes](#volumes)

# API
- [Status Codes](#status-codes)
- [Response & Format](#response-and-format)

## Status Codes
### 2xx
> This class of status code indicates that the client's request was successfully received, understood, and accepted.

200 Success - Probably the easiest response for any successful requests.  Use other codes if more detail is helpful.

202 Server Still Processing - The server has successfully processed the request but there are elements that are still being processed that may effect the same request in the future.

204 No Content - Helps mitigate the confusion between malformed requests and well formed request that do not return anything.

205 Reset Content - Request succeeded, reset the view.

242 Error Processing File - attempted to access data relating to file that failed to be properly processed

### 4xx
>    The 4xx class of status code is intended for cases in which the client seems to have erred.

400 Bad Request - Malformed request, correct and resend.

401 Unauthorized
- go to login page

403 Forbidden
- Credentials recognized, action forbidden as actor

404 Not Found - A server has not found anything matching the reqeust URI.

409 Conflict - Name Taken

460 Out of Space

### 5xx
> ...the server is aware that it has erred or is incapable of performing the request.

500 Internal Server Error

## Response and Format
An API response should provide consistant behavior and simplfy client side development, while at the same time provide relevant user data.

- Prefer JSON.
- Prefer giving users complete data sets with one request. (Where it makes sense to.)
- Empty arrays should return ([]).
- Prefer lower case properties.
- Multi word properties should be in Lower Camel Case.
- Provide a data object and an error object.
- Prefer a single root level descriptor for user data in the data object.

### Example Response
```
// TODO: keep this example ~ similar to all possible responses.
{
    "data": {
        user: {
            id: "",
            name: ""
        },
        files: [{
            id: "",
            name: ""
        }],
        group: {

        }
    },
    "error": {
        type: "LOGIN",
        message: "Login Required"
    }
}
```
# Commands
The following commands are can be usd for testing, headless testing, running a developemnt server, and more.

## Server
|Command|Description|
| ---     |---------|
| docker-compose up | Start the application with all its required containers. |
| docker-compose down | Stop all application containers. |

## Client
|Command|Description|
| ---     |---------|
|npm run serve| Runs a locally hosted development server with hot reloading. |
|npm run build| Builds a release version of the project. |
|npm run lint| Test the code for common js errors, and style guidelines. |
|npm run fix| Fix any linting errors that can be automatically corrected. |
|npm run test| Run all tests (.spec.js). Requires, karma.conf.js in the local folder & .env.test. |
|npm run headless| Run all tests using headless chrome. Requires, karma.headless.conf.js in the local folder & .env.headless. |

# Container Resources
Knaxim is built off of several containers: tika, mongo, and a web/api golang server.

## Volumes
|Volume|Description|Container Use|
| ---      |  ------  |---------:|
|web|Web Root Volme|knaxim|
|mongo|DB Data|mongo|
|mongo-cfg|DB Config|mongo|

## Basic Volume Managemenet
#### Basic volume management
```docker volume ls, create, and rm```

#### Inspect a volume
Inspect a volume to get basic information about it.  Mountpoint, Labels, and Driver are some of the fields available.
* Docker controlled Windows and Mac volume mount points are located in a VM.

```docker volume inspect web```
#### Copy a file into a volume
Adding files to a volume is done with a container that is using that volume.  The container does not have to be running, but it does have to be configured to use the volume.

```docker cp statuscode.txt knaxim_server_1:/public```

#### Remove a file from a volume
```docker exec knaxim_server_1 rm /public/statuscode.txt```
