# Passkeeper
Credential manager on Go (client/server)
This is my graduation project for Yandex Go-developer course.


## Client app logs location
```
# Linux
/var/log/passkeeper.log

# Windows
C:\Windows\Temp\passkeeper.log

```

# TODO Functionality
## Main Functionality
- [x] Use gRPC protocol
- [x] Add terminal user interface with main functionality
- [x] Fix bugs in TUI terminal user interface
- [x] User login
- [x] User register
- [x] User logout
- [ ] User delete account
- [ ] Change user password
- [x] List credentials
- [x] Edit credential
- [x] Add credential
- [x] Del credential
- [x] List bank cards 
- [x] Edit bank cards
- [x] Add bank cards
- [x] Del bank cards
- [x] List user files
- [x] Edit user files
- [x] Add user files
- [x] Del user files
- [x] List user files
- [x] Download user files
- [x] Edit user files
- [x] Add user files (upload)
- [x] Del user files
- [x] Asynchronous sync blobs (bugs fix)
- [x] Add DB storage
- [x] Gracefull shutdown (client)
- [ ] Gracefull shutdown (server)
- [x] Client parse args/env
- [x] Server parse args/env
- [ ] Add server access check before login/register user
- [ ] Reconnect to server
- [x] Move Creds from tuiApp to Usecase only
- [x] Move User entity from tuiApp to Usecase only
- [x] Move Token from tuiApp to Usecase only
- [x] Client set logger in file
- [x] Split client interface to multi inf
- [x] Split client struct to multi structs
- [ ] Server logger interceptor
- [ ] Server DDOS interceptor
- [ ] Setup TLS 
- [ ] Test cover more than 80%
- [x] Doc every function
- [x] Client build version from specification
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
- [x] TUI Order files by date
- [x] TUI clean fields (files->date, cards->number,cvc,pin)
- [x] TUI beautify cards number fields to "0000 0000 0000 0000"
- [ ] TUI move both login/register forms to center align 
- [x] TUI subpage: bank cards not auto update after delete card
- [x] TUI subpage: files not auto update after delete note
- [ ] Fix bug with Banks DropDown on Bank cards
- [ ] TUI form info about reconnect to server
- [ ] Fix crashing console after close client app
- [ ] Autodetecting Bank from card number
- [ ] Hide/Show Card info functionality

## Fill Readme.md
- [ ] Add gif to show functionality
- [x] Write about program logs