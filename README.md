# lookup-broker
OSBAPI compatible broker which implements a service lookup

# make

````
./bin/make.py 
usage: make.py [-h] [-v] [{build,run,test,generate,release,login,push}]
````
## make commands

| Command | Description |
| ---- |----|
| build | clean go build |
| run | run server as localhost |
| test | go test ./... |
| generate | generate OSBAPI api files |
| config | generate an application manifest and a login config file |
| release | release a new version |
| login | login to cloud foundry |
| push | push code to cloud foundry |

## local make configuration

call

````

./bin/make.py config

.
├── LICENSE
├── gen
│   └── config.yml
│   └── manifest.yml

cloud-foundry:
  api-url:  https://api.cf.eu10.hana.ondemand.com
  login:
    user: user name
    password: password
    org: org name
    space: space name
````
## brew dependencies

````
brew list
coreutils               jq                      bash                    curl                    go                      goreleaser
python@3.9              cf-cli                  cf-cli@7                openapi-generator
````

## python dependencies

````
pip3 list
Package     Version
----------- -------
autopep8    1.5.4
pip         20.3.3
protobuf    3.14.0
pycodestyle 2.6.0
PyYAML      5.3.1
setuptools  51.0.0
six         1.15.0
toml        0.10.2
wheel       0.36.1
````