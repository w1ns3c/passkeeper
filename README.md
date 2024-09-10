# Passkeeper
Credential manager on Go (client/server)

## Check
I know, that project is far from finish line, but...
Please, give me some comments and advices for my code in review.

# TODO Functionality
## Main Functionality
- [x] Use gRPC protocol
- [x] Order creds by date
- [ ] Add terminal user interface
- [x] User login
- [x] User register
- [x] User logout
- [ ] User delete account
- [ ] Change user password
- [x] List credentials
- [x] Edit credential
- [x] Add credential
- [x] Del credential
- [x] Asynchronous sync credentials (bugs fix)
- [ ] Add DB storage
- [ ] Gracefull shutdown
- [x] Client parse args/env
- [x] Server parse args/env
- [ ] Add server access check before login/register user
- [x] Move Creds from tuiApp to Usecase only
- [x] Move User entity from tuiApp to Usecase only
- [x] Move Token from tuiApp to Usecase only
- [ ] Reconnect to server
- [ ] TUI form info about reconnect to server
- [ ] Client set logger in file
- [ ] Split client interface to multi inf
- [ ] Server logger interceptor
- [ ] Server DDOS interceptor
- [ ] Test cover more than 80%
- [ ] Doc every function
- [ ] Refactor client TUI code
- [ ] Review all app code
- [ ] Check all TODO

## Optional Functionality
- [ ] Generate app doc
- [ ] Swagger
- [ ] Add email validation with sending message
- [ ] Add onetime passwords
- [ ] Change TUI to navigate with up/down arrow
- [ ] Finall review all app code (again)
