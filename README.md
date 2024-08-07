<div align="center">

![alt text](image-1.png)

</div>

---

<div align="center">

![Forks](https://img.shields.io/github/forks/Ant767/uwu-auth)
![Stars](https://img.shields.io/github/stars/Ant767/uwu-auth)

</div>


Simple and not at all janky authentication backend using **MongoDB** and **Gin gonic**

**Im new to the Go programming language, please do not judge my horrible code :3**

⭐ Features:
- App Support
- Users
- Basic User Profiles

✅ Todo

- [x] Basic auth system
- [x] Basic user profiles
- [ ] Finish app system
- [ ] User profile widgets
- [ ] User profile themes
- [ ] Finish docs
- [ ] 2FA

# API


`POST /register`
> Registers an account

Sample body:

```json
{
    "username": "TRASH",
    "password": "12345",
    "handle": "trash",
    "email": "antspa767@gmail.com"
}
```

Response:
`Content-Type: application/json`

`error`: false if there is no errror, true if there is one

`message`: the response message

*Note: these docs are still WIP*

# Config

This is an example config

```json
{
  "mongodb_url": "mongodb://127.0.0.1",
  "port": 8080,
  "resend_key": "xxxxxxx"
}
```

---

# Official instance

The main instance of uwu-auth is hosted at auth.trashdev.org

---

# Why use UwU Auth?

✨🎀 Its very kawaii :3

---

## Made by trashy :3
![alt text](image.png)