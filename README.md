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
| generate| generate OSBAPI api files |
| release | release a new version |
| login | login to cloud foundry |
| push | push code to cloud foundry |

## local make configuration

````
.
├── LICENSE
├── local
│   └── cloud-foundry.yml

cloud-foundry:
  api-url:  https://api.cf.eu10.hana.ondemand.com
  login:
    user: user name
    password: password
    org: org name
    space: space name
````
