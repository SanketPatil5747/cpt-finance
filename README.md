# chep

The CHEP tool is built as a docker container.
build e.g
```
docker-compose up 
```


The CHEP application has a postgres database - In AWS you have three options:
- to run an AWS-managed instance (this would be the easiest / preferred), 
- RDS, which can be an RDS instance or a serverless RDS instance,
- manage your instance of Postgres, e.g. in EC2 or as another container.



The CHEP container needs to be provided with connection details for the DB (address/credentials etc.) details admin credentials for the admin section, ports to run on etc. This is done through env variables & secrets.



CHEP is configurable with the following ENV variables:
```
"CONNECTION_STRING": "host=chep.db.host port=5438 user=chep_user password=chep_password dbname=chep_database sslmode=disable",
"DEVELOPMENT_ENV":"local" 
"PORT": 80,
"CLIENT_ID":"CLIENT_ID"
"CLIENT_SECRET":"CLIENT_SECRET"
"ISSUER":"https://Application_Dev_ID/oauth2/default"
"REDIRECT_URI":"http://DOMAIN_NAME/authorization-code/callback"
```
