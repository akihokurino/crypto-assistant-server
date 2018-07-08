MAKEFLAGS=--no-builtin-rules \
          --no-builtin-variables \
		  --always-make

OS := $(shell uname -s)
ifeq ($(OS),Darwin)
export SHELL := $(shell echo $$SHELL)
endif

ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export PATH := $(ROOT)/_scripts:$(PATH)

# export local environment variables
ENVFILE := $(ROOT)/.env
ifneq ($(shell test -e $(ENVFILE) && echo exists),)
include $(ENVFILE)
export $(shell sed 's/=.*//' $(ENVFILE))
endif

serve:
	serve.sh

deploy:
	deploy.sh


	
