NAMESPACE=github.com/dolanor/microservices

targets = auth todo data

all: $(targets)

.PHONY: $(targets)
$(targets) :
	go build -o services/$@/$@ $(NAMESPACE)/services/$@
	docker build -t microservices_$@ services/$@
