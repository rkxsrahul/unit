## List of API Endpoints ##

## 1. To check health of service ##

```
GET /healthz
```

## 2. Signup -> Create un-verified account ##

```
POST /v1/signup
Content-Type: application/json

Body:-
name      string    required
contact   string    required
email     string    required
password  string    required
```

## 3. Verify Email -> update account status to active and verify status to verified ##

```
POST /v1/verifymail
Content-Type: application/json

Body:-
verification_code   string required
email               string required
```

## 4. Send verification code again -> send verification code again in mail ##

```
POST /v1/send_code_again
Content-Type: application/json

Body:-
email string required
```

## 5. Refresh jwt token -> to generate new valid jwt token ##

```
GET /v1/refresh_token
Header:-
Authorization: Bearer <<token>>
```

## 6. Logout -> to delete active session or to invalidate token used in request ##

```
GET /v1/logout
Header:-
Authorization: Bearer <<token>>
```

## 7. Change Password ##

```
PUT /v1/changepass
Header:-
Content-Type: application/json
Authorization: Bearer <<token>>

Body:-
password string required
```

## 8. View profile information ##

```
GET /v1/profile
Header:-
Authorization: Bearer <<token>>
```

## 9. Check Workspace Status ##

```
POST /v1/workspace_status
Content-Type: application/json

workspace_url string required
team_name     string
team_size     string
team_type     string
```

## 10. Workspace login ##

```
POST /v1/workspace_login
Content-Type: application/json

workspace_url string required
```

## 11. Forgot Workpsace ##

```
POST /v1/forgot_workspace
Content-Type: application/json

email   string    required
```

## 12. List recover workspace ##

```
POST /v1/get_workspace_list
Content-Type: application/json

token   string   required
```

## 13. Signup for members ##

```
POST /v1/membersignup
Content-Type: application/json

workspace   string    required
email 		string    required
```

## 14. Member Registration is for new users in a workspace invited by workspace admin ##

```
POST /v1/memberregistration
Content-Type: application/json

workspace    string    required
token 		   string    required
password 		 string    required
name 		     string    required
contact 		 string    required
```

## 15. List all workspaces by admin only ##

```
Get /v1/admin/workspaces

Header:-
Content-Type: application/json
Authorization: Bearer <admin_token>
```

## 16. Delete Workspace on basis of id ##

```
DELETE /v1/admin/workspaces/:workspace_id

Header:-
Content-Type: application/json
Authorization: Bearer <admin_token>
```

## 17. Delete account by email ##

```
DELETE /v1/admin/account_emails/:email

Header:-
Content-Type: application/json
Authorization: Bearer <admin_token>
```

## 18. Create Workspace ##

```
POST /v1/workspaces
Content-Type: application/json
Authorization: Bearer <token>

workspace_url   string required
team_name 		string
team_size 		string
team_type 		string
```

## 19. Find workspace availability ##

```
POST /v1/workspace_availability
Content-Type: application/json
Authorization: Bearer <token>

workspace_url   string required
team_name 		string
team_size 		string
team_type 		string
```

## 20. Invite workspace members ##

```
POST /v1/workspaces/member
Content-Type: application/json
Authorization: Bearer <token>

email []string
```

## 21. List workspace members ##

```
GET /v1/workspaces/member
Content-Type: application/json
Authorization: Bearer <token>
```

## 22. Workspace user login using recover workspace token ##

```
POST /v1/tokenlogin
Content-Type: application/json

token     string  required
email     string  required
workspace string  required
```

## 23. Workspace user Login ##

```
POST /v1/login
Content-Type: application/json

workspace   string
email     string    required
password  string    required
```

## 24. Forgot Password ##

```
POST /v1/forgotpass
Content-Type: application/json
```
<table><tr><th> Variable name </th><th> type </th><th> Required </th></tr>
<tr><td> state </td><td> string </td><td> Yes <br> (value should be either 'forgot' or 'reset') <br> In case of forgot send 'email' and <br> in case of reset send other two params.</td></tr>
<tr><td> email </td><td> string </td><td> optional for reset state </td></tr>
<tr><td> token </td><td> string </td><td> optional for forgot state </td></tr>
<tr><td> password </td><td> string </td><td>optional for forgot state <br> It contains new password. </td></tr>
<tr><td> workspace </td><td> string </td><td>optional when there is no workspace </td></tr></table>


## 25. Update profile data ##

```
PUT /v1/profile
Content-Type: application/json

name    string  required
contact string  required
```

## 26. Check member is for checking the member is assigned and joined the workspace ##

```
GET /v1/workspaces/checkmember
Header:-
Authorization: Bearer <<token>>

Query data
email   string
```

## 27. To disable or enable Mail Service##

```
PUT /v1/admin/mail/:value
Header:-
Authorization: Bearer <<admin token>>

Param data
value -> true if want to diable service and false if you want to enable service
```

## 28. To disable or enable OTP Service##

```
PUT /v1/admin/otp/:value
Header:-
Authorization: Bearer <<admin token>>

Param data
value -> true if want to diable service and false if you want to enable service
```

## 29. To mail the access keys ##

```
GET /v1/accessKeys
Header:-
Authorization: Bearer <<token>>
```
