# Dev documentation

## SMTP Test

<https://github.com/jetmore/swaks>

### Test HTML email

```shell
swaks --body '<a href="https://en.wikipedia.org/wiki/Main_Page">link</a>' --add-header "MIME-Version: 1.0" --add-header "Content-Type: text/html" --to user@example.com --server 127.0.0.1:2525
```

## HTTP Test

```shell
curl --header "Content-Type: application/json" --request POST  --data '{"body": "Test Notification", "severity": "notice"}' http://localhost:8080/hook
```

```shell
curl --header "Content-Type: application/json" --request POST  --data '["Test Notification", "notice"]' http://localhost:8080/hook
```
