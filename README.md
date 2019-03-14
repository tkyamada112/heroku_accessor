# Heroku Accessor

For organize [Heroku](https://www.heroku.com/) Collaborators.

# How to
Clone this repository and `go build` on your machine.

- List all apps.

```
HerokuAccesor -username xxxx@xxxxx.xxxx -password xxxxxxx -type showall

```

- List Collaborators.

```
HerokuAccesor -username xxxx@xxxxx.xxxx -password xxxxxxx -type showuser -name xxxxx
```

- Update User.

```
HerokuAccesor -username xxxx@xxxxx.xxxx -password xxxxxxx -type updateuser -name xxxxx
```
