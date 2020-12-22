#!/usr/bin/env python3

import os
import argparse
import yaml

def test(verbose):
    print ("-- vet & test broker")
    os.system("go vet ./...")
    if verbose:
      os.system("go test ./... -v")
    else:    
      os.system("go test ./...")

def build(verbose):
    print ("-- clean broker")
    os.system("go clean -r -cache -testcache -modcache")
    print ("-- fmt broker")
    os.system("go fmt ./...")
    print("-- build broker")
    flg = ldflags()
    os.system("go build " + flg + "-v ./...")
    if verbose:
      print (flg)

def run(verbose):
    print ("-- run broker")
    os.system("go run  " + ldflags() + " brokerApp.go")

def generate(verbose):
    print ("-- generate broker")

    os.system("rm -rf ./gen ./openapi")
    os.system("mkdir -p ./gen ./openapi")
    os.system("wget https://raw.githubusercontent.com/openservicebrokerapi/servicebroker/master/swagger.yaml -O \"./gen/swagger.yaml\"")
    os.system("openapi-generator validate -i ./gen/swagger.yaml")
    os.system("openapi-generator generate -i ./gen/swagger.yaml -g go-server -o ./gen")
    os.system("cp ./gen/go/model_* ./openapi")
    os.system("go fmt ./openapi")

def release(verbose):
    print ("-- release broker")
    print ("-- tbd")
    print ("-- verbose:", verbose)

def login(verbose):
    print ("-- login cloud foundry")

    with open("./local/cloud-foundry.yml", 'r') as stream:
        try:
            cfConfig = yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            print(exc)
            os._exit(0)

    os.system("cf api {}".format(cfConfig.get("cloud-foundry").get("api-url")))
    os.system("cf auth {} {}".format(cfConfig.get("cloud-foundry").get("login").get("user"), cfConfig.get("cloud-foundry").get("login").get("password")))
    os.system("cf target -o {} -s {}".format(cfConfig.get("cloud-foundry").get("login").get("org"), cfConfig.get("cloud-foundry").get("login").get("space")))

def push(verbose):
    print ("-- push to cloud foundry")

    os.system("cf push lookup-broker")

def dispatcher(cmd):
   dispatcher={
       'build': build,
       "run": run,
       "test": test,
       "generate": generate,
       "release": release,
       "login": login,
       "push": push,
   }
   return dispatcher.get(cmd)
 
def ldflags():
    dirtyStream = os.popen('git diff --quiet || echo dirty')
    dirty = dirtyStream.read()
    if dirty == "":
      commitStream = os.popen('git rev-parse HEAD')
      commit = commitStream.read()
    else:
      commit = dirty.strip()


    # TDOD: version (consider goreleaser) 
    return F"-ldflags=\"-X 'main.Version=dev' -X 'main.Commit={commit}'\""

def main():
  parser = argparse.ArgumentParser(description="Make tool for cloud foundry service lookup broker", epilog="(c) 2020 by KLÄFF-Soft)")
  parser.add_help=True
  parser.add_argument("command", nargs='?', choices=['build', 'run', 'test', 'generate', 'release', 'login', 'push'], help="commands to execute")
  parser.add_argument("-v", action="store_true", help="verbose output")

  args = parser.parse_args()

  func = dispatcher(args.command)

  if func != None:
    func(args.v)
  else: 
    parser.print_usage()

if __name__ == '__main__':
    main()