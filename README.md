# Scaf backend

## 開發方式

```bash
go mod init backend
go mod tidy
go run .
```

### api (need jwt)
signin: ```/signin```<br>
method: POST
Data:
```json
{
    "email": "email",
    "password": "password"
}
```

signup: ```/signup```<br>
method: POST<br>
Data:
```json
{
    "email": "email",
    "password": "password"
}
```

create project: ```user/{user_email}```<br>
method: POST<br>
Data:
```json
{
    "Name": "test",
    "DevTools": [],
    "DevMode": "waterfall"
}
```

<!-- add repo: ```/projects/repos```<br>
method: POST<br>
Data:
```json
{
    "Name": "RepoName",
    "Url": "RepoUrl"
}
``` -->

### firebase

測試帳號: test@test.com / test

## TODO List

- [x] 登入
- [X] 註冊 
- [X] 登出
- [X] 忘記密碼# scaf-backend
