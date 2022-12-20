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

## RESTful API

⚠️ **注意** ⚠️: 網址前綴為 `{app url}:{app port}`，最後面一定要 `/` 結尾。  
⚠️ **注意** ⚠️: 所有欄位名稱首字母為小寫。

### SignIn

```POST /signin/```

Request:
```json
{
    "email": "[email]",
    "password": "[password]"
}
```

### SignUp

```POST /signup/```

```json
{
    "email": "[email]",
    "password": "[password]"
}
```

### Forgot Password 🚧 (施工中)

```POST /forgot/```

```json
{
    "email": "[email]"
}
```

### Get User Data 🚧 (施工中)

```Required JWT```  
```GET /user/{user_email}/```

user_email 可為空，為空代表自己。

### Update User Data 🚧 (施工中)

```Required JWT```  
```PUT /user/{user_email}/```

```json
{
    "avatar": "[avatar base64]",
    "nickname": "[nickname]",
    "password": "[password]"
}
```

### Update User Password 🚧 (施工中)

```Required JWT```  
```PUT /user/{user_email}/reset```

```json
{
    "oldPassword": "[old password]",
    "newPassword": "[new password]"
}
```

### 取得 Google 日曆授權 🚧 (施工中)

```Required JWT```  
```GET /user/{user_email}/calendar```

### Refresh ❌ (目前不可用)

```Required JWT```  
```POST /refresh/```

### Test (測試用)

```GET /hello/```

### List User's Project

list user's project

```GET /user/{user_email}/project/```

### Create Project

```Required JWT```  
```POST /user/{user_email}/project/```

```json
{
    "name": "[project name]",
    "devTools": [],
    "devMode": "[waterfall|scrum]"
}
```

### Delete Project

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/```

### List All repos

```Required JWT```  
```GET /user/{user_email}/project/{project_id}/repo/```

### Add Repo

```Required JWT```  
```POST /user/{user_email}/project/{project_id}/repo/```

```json
{
    "name": "[repo name]",
    "url": "[repo url]"
}
```

### Update Repo 🚧 (施工中)

```Required JWT```  
```PUT /user/{user_email}/project/{project_id}/repo/{repo_id}/```

```json
{
    "name": "[repo name]",
    "url": "[repo url]"
}
```

### Delete Repo 🚧 (施工中)

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/repo/{repo_id}/```

### Create Kanban ❌ (目前不可用)

```Required JWT```  
```POST /user/{user_email}/project/{project_id}/kanban/```

### List Workflow

```Required JWT```  
```GET /user/{user_email}/project/{project_id}/kanban/```

### Create Workflow

```Required JWT```  
```PUT /user/{user_email}/project/{project_id}/kanban/```

```json
{
    "name": "[workflow name]"
}
```

### Delete Workflow 🚧 (施工中)

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/kanban/{workflow_id}/```

### Add Task 🚧 (施工中)

```Required JWT```  
```POST /user/{user_email}/project/{project_id}/kanban/```

```json
{
    "name": "[task name]",
    "description": "[task description]"
}
```

### Delete Task 🚧 (施工中)

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/kanban/{workflow_id}/{task_id}/```

### Get Members 🚧 (施工中)

```Required JWT```  
```GET /user/{user_email}/project/{project_id}/member/```

### Add Member 🚧 (施工中)

```Required JWT```  
```POST /user/{user_email}/project/{project_id}/member/```

```json
{
    "email": "[member email]"
}
```

### Delete Member 🚧 (施工中)

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/member/{member_email}/```

### Get Document 🚧 (施工中)

```Required JWT```  
```GET /user/{user_email}/project/{project_id}/doc/```

### Add Document 🚧 (施工中)

```Required JWT```  
```POST /user/{user_email}/project/{project_id}/doc/```

```json
{
    "name": "[doc name]",
    "content": "[doc content]"
}
```

### Update Document 🚧 (施工中)

```Required JWT```  
```PUT /user/{user_email}/project/{project_id}/doc/{doc_id}/```

```json
{
    "name": "[doc name]",
    "content": "[doc content]"
}
```

### Delete Document 🚧 (施工中)

```Required JWT```  
```DELETE /user/{user_email}/project/{project_id}/doc/{doc_id}/```

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
- [ ] 編輯看板任務(Todo, InProgress, Done)
- [ ] 邀請加入專案
- [ ] google日曆授權


## Refactor list

- [ ] 將 Request model 