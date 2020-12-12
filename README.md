# noid-server

```yml
version: "3"

services: 
  noid:
    build: github.com/ClementD64/noid-server.git
    command: /song
    volumes: 
      - /path/to/songs:/song:ro
    ports: 
      - 3000:3000
    environment: 
      NOID_REDIS: redis://redis
  
  redis:
    image: redis
```