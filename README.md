# go-eat

`go-eat` reads daily menus from Kaiserstück, Personalkantine and STW Mensas and puts them into a DynamoDB table. Perfect to build web services that tell you where to eat at TU Berlin today.

```
$ aws configure
$ make
```

You will need AWS access keys and an AWS region where you'd like to deploy this.
