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
- start test MariaDB database in docker container
```sh
docker compose up -d
```
- check the database is created http://localhost:8181/?server=mariadb&username=root&db=test use "root" as password
- start go application
```sh
go run .
```
- check http://localhost:3000/list endpoint (GET method in any browser)
```
ER == AR
ER: there are no policies
AR: there are no policies
```
- init policies by http://localhost:3000/init endpoint (GET method in any browser)
- check http://localhost:3000/list endpoint (GET method in any browser)
```
ER == AR
ER: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","List Users Roles"],["3","Admin Users Roles"],["4","Admin Users Roles"]]
AR: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","List Users Roles"],["3","Admin Users Roles"],["4","Admin Users Roles"]]
```
- call http://localhost:3000/assign?role=Root&user_id=2 endpoint to assign Root role to user with id=2
- call http://localhost:3000/assign?role=Root&user_id=3 endpoint to assign Root role to user with id=2
- call http://localhost:3000/assign?role=Root&user_id=4 endpoint to assign Root role to user with id=2
- check database http://localhost:8181/?server=mariadb&username=root&db=test&select=casbin_rule
```
ER == AR
g	1	Root
g	2	Root
g	3	Root
g	4	Root

```
- check http://localhost:3000/list endpoint
```
ER != AR
ER: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","Root"],["3","Root"],["4","Root"]]
AR: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["4","Root"]]
```
- restart go application
- check http://localhost:3000/list endpoint
```
ER == AR
ER: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","Root"],["3","Root"],["4","Root"]]
AR: [["Root","root:root:root"],["List Users Roles","core:users:list"],["Admin Users Roles","core:users:list"],["Admin Users Roles","core:users:create"],["1","Root"],["2","Root"],["3","Root"],["4","Root"]]
```
- stop go application
- stop docker container
```sh
docker compose down
```