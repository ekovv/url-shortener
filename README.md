# go-musthave-shortener-tpl
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

