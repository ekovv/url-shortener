# go-musthave-shortener-tpl

![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# 🎲 Service on Go(Gin) for shortening url using a special algorithm and storing in a database/file/memory 🎲

# 🎬 Description 

Receives a link to a web resource from the client and, using a text reduction algorithm, shortens it and saves it in storage (your choice: memory, file, database) and sends it back. The new short link will automatically redirect all clients to the original (longer) link.

# 📞 Endpoints
```http
POST /
- Create link
GET /:id 
- Get link 
POST /api/shorten
- Create link from json
GET /ping
- Get Stats 
POST /api/shorten/batch
- Batch create 
GET /api/user/urls
- Get all
DELETE /api/user/urls
- Delete links
```

# 🏴‍☠️ Flags
```
a - ip for REST -a=host
b base url -b=base
f - path to the file to be used as a database -f=storage
d - connection string -d=connection string
```

