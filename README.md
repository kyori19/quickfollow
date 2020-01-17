# quickfollow
Forkのmaster追従を半自動化します

## How to use

Using `go run`
```
go run main.go log.go config.go [path to your repository]
```

Using built binary
```
go build
quickfollow [path to your repository]
```
You can place built binary any place you like.

## How to write config file

You can use both `.json` and `.yml` and you have to use `quickfollow` as name for config file.

| name     | type     | details  
|----------|----------|--------  
| upstream | string   | remote name  
| branch   | string   | remote branch name  
| target   | string[] | list of local branches which requires merge  

You can find an example at [here](https://github.com/accelforce/odakyudon/blob/26a8883f1354a4c99250ef762c00bf79d6e5f9e2/quickfollow.json).

## License

```
Copyright (c) 2020 kyori19
```

This software is licensed under GPLv3.
