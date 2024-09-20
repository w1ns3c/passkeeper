# Passkeeper
Credential manager on Go (client/server)

## Check
I know, that project is far from finish line, but...
Please, give me some comments and advices for my code in review.

# TODO Functionality
## Main Functionality
- [x] Use gRPC protocol
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
- [x] Gracefull shutdown (client)
- [ ] Gracefull shutdown (server)
- [x] Client parse args/env
- [x] Server parse args/env
- [ ] Add server access check before login/register user
- [x] Move Creds from tuiApp to Usecase only
- [x] Move User entity from tuiApp to Usecase only
- [x] Move Token from tuiApp to Usecase only
- [ ] Reconnect to server
- [ ] TUI viewForm info about reconnect to server
- [x] Client set logger in file
- [ ] Split client interface to multi inf
- [ ] Server logger interceptor
- [ ] Server DDOS interceptor
- [ ] Setup TLS 
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

## TUI 
- [x] TUI Order notes by date
- [ ] TUI clean fields (notes->date, cards->number,cvc,pin)
- [x] TUI beautify cards number fields to "0000 0000 0000 0000"
- [ ] TUI move both login/register forms to center align 
- [ ] TUI subpage: bank cards not auto update after delete card
- [ ] TUI subpage: notes not auto update after delete note

## Fill Readme.md
- [ ] Add gif to show functionality
