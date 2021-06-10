Building the jenkinsctl

```
$ go build jenkinsctl.go
```

First step, generate a token for the username that will manage the jenkins.

- Log in to Jenkins.
- Click you name (upper-right corner).
- Click Configure (left-side menu).
- Use "Add new Token" button to generate a new one then name it.
- You must copy the token when you generate it as you cannot view the token afterwards.

Second, create the `configuration directory` and the `config.json file`
```
$ mkdir -p ~/.config/jenkinsctl/
$ pushd ~/.config/jenkinsctl/
    $ vi config.json 
    {
        "Server": "https://jenkins.mydomain.com",
        "JenkinsUser": "jenkins-operator",
        "Token": "1152e8e7a88f6c7ef605844b35t5y6i"
	"WatchIntervalInSec": 5
    }
$ popd
```
