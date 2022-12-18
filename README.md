# Scaf backend

## 開發方式

```bash
go mod init backend
go mod tidy
go run .
```
需要修改firbase的project ID、Web API Key及serviceAccount.json
firbase project id: 專案設定內查看
Web API Key: 專案設定內查看
serviceAccount.json: https://sharma-vikashkr.medium.com/firebase-how-to-setup-a-firebase-service-account-836a70bb6646

### api (need jwt)
1. **signin:** ```/signin```<br>
method: POST<br>
Data:
```json
{
    "email": "email",
    "password": "password"
}
```
2. **signup:** ```/signup```<br>
method: POST<br>
Data:
```json
{
    "email": "email",
    "password": "password"
}
```
3. **list user's project:** ```{user_email}/project```<br>
method: GET<br>
4. **create project:** ```{user_email}/project```<br>
method: POST<br>
Data:
```json
{
    "Name": "test",
    "DevTools": [],
    "DevMode": "waterfall"
}
```
5. **Delete project:** ```{uesr_email}/project/{project_id}```<br>
method: POST<br>
6. **list all repos:**```{uesr_email}/project/{project_id}/repo```<br>
method: GET<br>
7. **add repo:**```{uesr_email}/project/{project_id}/repo```<br>
method: POST<br>
Data:
```json
{
    "Name": "RepoName",
    "Url": "RepoUrl"
}
```

### firebase

測試帳號: test@test.com / testtest

## TODO List

- [x] 登入
- [X] 註冊 
- [X] 登出
- [X] 忘記密碼# scaf-backend
