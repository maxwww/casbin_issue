### Steps to reproduce

- clone [repo](git@github.com:maxwww/casbin_issue.git)
```sh
git clone git@github.com:maxwww/casbin_issue.git
```
- navigate to project
```sh
cd casbin_issue
```
- fetch go dependencies
```sh
go mod tidy
```
- start go application
```sh
go run .
```
- check http://localhost:3000/list endpoint (GET method in any browser)
```
ER == AR
ER: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","List Users Roles"],["3","Admin Users Roles"],["4","Admin Users Roles"]]
AR: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","List Users Roles"],["3","Admin Users Roles"],["4","Admin Users Roles"]]
```
- call http://localhost:3000/assign?role=Root&user_id=2 endpoint to assign Root role to user with id=2
- call http://localhost:3000/assign?role=Root&user_id=3 endpoint to assign Root role to user with id=3
- call http://localhost:3000/assign?role=Root&user_id=4 endpoint to assign Root role to user with id=4
- check http://localhost:3000/list endpoint
```
ER != AR
ER: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","Root"],["3","Root"],["4","Root"]]
AR: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["4","Root"]]
```
