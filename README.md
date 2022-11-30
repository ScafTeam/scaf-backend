# Scaf backend

## 開發方式

```bash
go mod init backend
go mod tidy
go run .
```

### api
signin ```/auth/signin```
    Data:
        ```json
        {
            email: "email",
            password: "password"
        }
        ```
signup ```/auth/signup```
    Data:
        ```json
        {
            email: "email",
            password: "password"
        }
        ```
list all projects ```/projects/list```
create project ```/projects/create```
    Data:
        ```json
        {
            "Name": "test",
            "DevTools": [],
            "DevMode": "waterfall"
        }
        ```

### firebase

測試帳號: test@test.com / abstest1

## TODO List

- [x] 登入
- [X] 註冊 
- [ ] 登出
- [ ] 忘記密碼# scaf-backend
