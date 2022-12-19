# Scaf backend

## 開發方式

```bash
go mod init backend
go mod tidy
go run .
```
需要修改firbase的project ID、Web API Key及serviceAccount.json，其中project ID和Web API Key須放在config.txt中(格式如下)，config.txt要在database資料夾中<br>


Web API Key: 專案設定內查看<br>
firbase project id: 專案設定內查看<br>
serviceAccount.json: https://sharma-vikashkr.medium.com/firebase-how-to-setup-a-firebase-service-account--6a70bb6646

- **config.txt**
```
{WEB KEY API}
{PROJECT ID}
```

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
- **list all repos:** ```{uesr_email}/project/{project_id}/repo``` (不用jwt)<br>
method: GET<br>
- **add repo:** ```{uesr_email}/project/{project_id}/repo```<br>
method: POST<br>
Data:
```json
{
    "Name": "RepoName",
    "Url": "RepoUrl"
}
```
- **create kanban:** ```{user_email}/project/{project_id}/kanban```<br>
method: POST<br>
- **add Task:** ```{user_email}/project/{project_id}/kanban/{Todo|InProgress|Done}```<br>
method: POST<BR>
Data:
```json
{
    "Name": "Name",
    "Description": "Description"
}
```
- **delete Task:** ```{user_email}/project/{project_id}/kanban/{Todo|InProgress|Done}```<br>
method: DELETE<br>
- **add member:** ```/{user_email}/project/{project_id}/join ```(未完成)<br>
method: POST<br>


### firebase

測試帳號: test0@test.com / 123456 

## TODO List

- [x] 登入
- [X] 註冊 
- [X] 登出
- [X] 忘記密碼
- [X] 創建專案
- [x] 取得目前使用者的專案
- [x] 刪除專案(只能由專案擁有者操作)
- [x] JWT登入驗證
- [x] JWT權限驗證
- [x] Refresh JWT token
- [x] 取得目前專案中的所有Repo
- [x] 增加Repo(只能由專案成員操作)
- [x] 刪除Repo(只能由專案成員操作)
- [x] 創建看板(在創建專案時同時創建專案看板，未檢查專案與看板是否一對一)
- [x] 取得看板
- [x] 新增看板任務(Todo, InProgress, Done)
- [x] 刪除看板任務(Todo, InProgress, Done)
- [] 編輯看板任務(Todo, InProgress, Done)
- [] 邀請加入專案
