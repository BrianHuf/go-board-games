BACKEND = $(CURDIR)/backend/src

build-backend:
	cd $(BACKEND) &&\
	make build

test-backend:
	cd $(BACKEND) &&\
	make testf


build: build-backend