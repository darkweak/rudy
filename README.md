RUDY (R-U-Dead-Yet?)
====================

## What is the RUDY Attack?
R.U.D.Y., short for R U Dead yet, is an acronym used to describe a Denial of Service (DoS) tool used by hackers to perform slow-rate a.k.a. “Low and slow” attacks by directing long form fields to the targeted server. It is known to have an interactive console, thus making it a user-friendly tool. It opens fewer connections to the website being targeted for a long period and keeps the sessions open as long as it is feasible. The amount of open sessions overtires the server or website making it unavailable for the authentic visitors. The data is sent in small packs at an incredibly slow rate; normally there is a gap of ten seconds between each byte but these intervals are not definite and may vary to avert detection.

The victim servers of these types of attacks may face issues such as not being able to access a particular website, disrupt their connection, drastically slow network performance, etc.

Hackers can use such attacks for different purposes while targeting different servers or hosts; these purposes include, but are not limited to, blackmail, vengeance or sometimes even activism.

The RUDY attack opens concurrent POST HTTP connections to the HTTP server and delays sending the body of the POST request to the point that the server resources are saturated. This attack sends numerous small packets at a very slow rate to keep the connection open and the server busy. This low-and slow attack behavior makes it relatively difficult to detect, compared to flooding DoS attacks that raise the traffic volume abnormally.

## How to install and run rudy?

### Using `go install`
```bash
go install -u github.com/darkweak/rudy/cmd/rudy@latest
rudy [command]
```

### Using directly `go`
```bash
git clone https://github.com/darkweak/rudy
cd rudy
go run rudy.go [command]
```

### Using `go build`
```bash
git clone https://github.com/darkweak/rudy
cd rudy
go build rudy.go -o rudy
rudy [command]
```

## Commands
### Attack a target
```bash
rudy run -u http://domain.com
```

There are some options to change the rudy default behaviour 
| Name                | Description                                                              | Long flag        | Short flag | Example value           | Default value |
|:--------------------|:-------------------------------------------------------------------------|:-----------------|:-----------|:------------------------|:--------------|
| URL                 | The target URL to run the attack on.                                     | `--url`          | `-u`       | `http://domain.com`     |               |
| Concurrent requests | The number of concurrent requests to send on the target.                 | `--concurrents`  | `-c`       | `4`                     | `1`           |
| Filepath            | Filepath to the payload to send. By default it's a random payload (1MB). | `--filepath`     | `-f`       | `/somewhere/file`       |               |
| Interval            | Interval duration between the requests.                                  | `--interval`     | `-i`       | `3s`                    | `10s`         |
| Size                | Random payload size to send. Used if no filepath given.                  | `--payload-size` | `-s`       | `1GB`                   | `1MB`         |
| Tor                 | Use TOR proxy to send the requests.                                      | `--tor`          | `-t`       | `socks5://tor_endpoint` |               |

### Run the testing server
It will start a serer on the port `:8081`
```bash
rudy server
```