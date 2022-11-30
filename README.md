# Scaf backend

## 開發方式

```bash
go mod init backend
go mod tidy
go run .
```

### api
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

list all projects: ```/projects```
method: GET<br>

create project: ```/project```<br>
method: POST<br>
Data:
```json
{
    "Name": "test",
    "DevTools": [],
    "DevMode": "waterfall"
}
```

add repo: ```/projects/repos```<br>
method: POST<br>
Data:
```json
{
    "Name": "RepoName",
    "Url": "RepoUrl"
}
```

### firebase

測試帳號: test@test.com / abstest1

## TODO List

- [x] 登入
- [X] 註冊 
- [ ] 登出
- [ ] 忘記密碼# scaf-backend
