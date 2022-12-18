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
serviceAccount.json: https://sharma-vikashkr.medium.com/firebase-how-to-setup-a-firebase-service-account--6a70bb6646

### api (不用jwt的都有寫上去其餘都是需要jwt的)
- **signin:** ```/signin``` (不用jwt)<br>
method: POST<br>
Data:
```json
{
    "email": "email",
    "password": "password"
}
```
- **signup:** ```/signup``` (不用jwt)<br>
method: POST<br>
Data:
```json
{
    "email": "email",
    "password": "password"
}
```
- **refresh:** ```/refresh```<br>
method: POST<br>
需要舊的jwt token

- **hello** ```/hello```<br>
method: GET<br>
測試AuthMiddleware用的

- **list user's project:** ```{user_email}/project``` (不用jwt)<br>
method: GET<br>
- **create project:** ```{user_email}/project```<br>
method: POST<br>
Data:
```json
{
    "Name": "test",
    "DevTools": [],
    "DevMode": "waterfall"
}
```
- **Delete project:** ```{uesr_email}/project/{project_id}```<br>
method: DELETE<br>
- **list all repos:**```{uesr_email}/project/{project_id}/repo``` (不用jwt)<br>
method: GET<br>
- **add repo:**```{uesr_email}/project/{project_id}/repo```<br>
method: POST<br>
Data:
```json
{
    "Name": "RepoName",
    "Url": "RepoUrl"
}
```
- **create kanban**```{user_email}/project/{project_id}/kanban```<br>
method: POST<br>
- **add Task**```{user_email}/project/{project_id}/kanban/{Todo|InProgress|Done}```<br>
method: POST<BR>
Data:
```json
{
    "Name": "Name",
    "Description": "Description"
}
```
- **delete Task**```{user_email}/project/{project_id}/kanban/{Todo|InProgress|Done}```<br>
method: DELETE<br>
- **add member**```/{user_email}/project/{project_id}/join ```<br>
method: POST<br>


### firebase

測試帳號: test0@test.com / 123456 

## TODO List

- [x] 登入
- [X] 註冊 
- [X] 登出
- [X] 忘記密碼# scaf-backend
