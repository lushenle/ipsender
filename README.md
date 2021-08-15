# SendIP

---

Get public IP from `https://ip.tool.lu`, then send it to someone via email

## config.json

```json
{
    "from": "admin@idocker.io",
    "to": "lushenle@gmail.com",
    "password": "YOUR_PASSWD",
    "smtp_username": "postmaster@mg.idocker.io",
    "smtp_host": "smtp.mailgun.org",
    "smtp_port": 25,
    "subject": "Home IP",
    "interval": 10
}
```

## build

```bash
./build.sh
```

## deploy

```bash
# docker
docker run --name ipsender -d \
  -v ${PWD}/config.json:/config.json \
  --restart=always \
  manunkind/ipsender:v1.0
# k8s
kubectl apply -f k8s.yaml
# docker-compose
docker-compose up -d
```
