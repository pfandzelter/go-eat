.PHONY: deploy plan clean

deploy: go-eat.zip main.tf init.done
	terraform apply
	touch $@
	rm -f plan go-eat.zip go-eat

plan: go-eat.zip main.tf init.done
	terraform plan
	touch $@
	rm -f go-eat.zip go-eat

init.done:
	terraform init
	touch $@

go-eat.zip: go-eat
	chmod +x go-eat
	zip -j $@ $<

go-eat: main.go
	go get .
	GOOS=linux GOARCH=amd64 go build -ldflags="-d -s -w" -o $@

clean:
	terraform destroy
	rm -f init.done deploy.done go-eat.zip go-eat