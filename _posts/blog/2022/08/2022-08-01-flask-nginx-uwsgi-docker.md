---
layout: post
title: "🐳Deploying Flask app using Nginx, uWSGI, and Docker"
tags: blog, web, docker, flask, nginx, uwsgi
date: 2022-08-01 17:15 +0900
math: true
---

![banner](https://i.imgur.com/mOJJqEl.png)

## Nginx-uWSGI-Flask

**`Nginx-uWSGI-Flask`** is a commonly used server setup for light python web applications or ML applications. 
**Flask**, a light web framework, is chained with a web server **uWSGI** and **Nginx** as a reverse proxy. 
**Flask** is in integrated as a callable object for **uWSGI** to call and run the app. And **Nginx** and **uWSGI** communicates via unix socket or tcp port.

```shell
             ┌───────────────────────────────────────────────┐
             | server:                                       |
             |                                               |
clients <---------> nginx <---------> uwsgi <--------> flask |
             |                                               |
             └───────────────────────────────────────────────┘
```

## Docker
For easier deployment migration to different machines or server, let's use **Docker**. **Docker** is a container-based platform for building and deploying applications. 
**Docker** isolates and virtualize applications using containers. Unlike virtual machines, containers isolates processes by enabling multiple applications to share 
the resources of a single instance of the host OS without the needs of virtualizing entire operationg system.
It is **lighter**, **faster**, and **efficient** than virtual machines. It is possible to create and run containers without **Docker**, but **Docker** makes it way easier.

## Building a Simple Web Server
Here's how I built the servers using Docker on `Ubuntu 22.04 LTS`
I used `docker-compose` which is a tool for building a multi-container Docker app.

I set nginx to listen to port `5000` and communicate with uwsgi with unix socket `./uwsgi.sock.` Flask will be called by uwsgi as a callable object
```shell
             ┌───────────────────────────────────────────────┐
             | server:                                       |
             |                                               |
clients <---------> nginx <---------> uwsgi <--------> flask |
           5000           ./uwsgi.sock    callable object    |
             └───────────────────────────────────────────────┘
```

### Installing and setting up Docker and Docker-compose
Install `docker`.
```shell
$ sudo apt-get install curl
$ curl -s https://get.docker.com | sudo sh
```
Add user to the `docker` group.
```shell
$ sudo usermod -aG docker $USER
```
Install `docker-compose`.
```shell
$ sudo curl -L "https://github.com/docker/compose/releases/download/v2.8.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
```

### File Structure

```shell
server
├── docker-compose.yml
├── flask
│   ├── Dockerfile
│   ├── app.py
│   ├── requirements.txt
│   └── uwsgi.ini
└── nginx
    ├── Dockerfile
    └── nginx.conf
```

### `./flask`
`app.py`
```python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return "hello docker🐳"
        

if __name__ == '__main__':
    app.run(host='0.0.0.0')
```

`requirements.txt`
```python
flask==2.1.*
uwsgi==2.0.*
```

`Dockerfile`
```python
FROM python:3

WORKDIR /app

ADD . /app
RUN pip install -r requirements.txt
```

`uwsgi.ini`
```python
[uwsgi]
wsgi-file = app.py
callable = app
socket = ./uwsgi.sock
processes = 2
threads = 2
master = true
vacum = true
chmod-socket = 660
die-on-term = true
```

### `./nginx`

`nginx.conf`
```python
server {

	listen 5000;
	
	location / {
		include uwsgi_params;
		uwsgi_pass ./uwsgi.sock;
	}
}
```

`Dockerfile`
```python
FROM nginx

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/conf.d/
```

## Build and start containers
```shell
$ docker-compose up -d --build
```
* `-d`: run containers in the background
* `--build`: build images before starting containers
