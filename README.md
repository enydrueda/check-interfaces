Script en experimental en GoLang para chequear curl desde las interfaces eth0 y eth2 de zentyal
===============================================================================================

Se requiere Go insalado [https://golang.org/doc/install](https://golang.org/doc/install)

Dependencias
------------
- github.com/joho/godotenv
- github.com/ashwanthkumar/slack-go-webhook

Edicion
-------
clonar el proyecto el workspace

```bash
cd <WORKSPACE>/src
git clone git@gitlab.services4geeks.co:devops/check-zentyal-interfaces.git
```

Compilar
--------
```bash
go install check-zentyal-interfaces
```
Resultando el archivo ejecutable "check-zentyal-interfaces" en la ruta <WORKSPACE>/bin/

Uso
---
Crear un archivo .env en la locación en la que se encuentra el "ejecutable", ver .env.example
Correr el "ejecutable" con lo cual se levantara un servidor en el puerto configurado, crear los CronJobs requeridos en el zentyal y el servidor en donde se ejecuta el archivo

CronJob Zentyal
---------------
```bash
* * * * * curl --interface "eth0" <DOMAIN>:<PORT>/?inter=eth0
* * * * * curl --interface "eth2" <DOMAIN>:<PORT>/?inter=eth2
```

CronJob Server
---------------
```bash
* * * * * curl <DOMAIN>:<PORT>/timeFile?inter=eth0
* * * * * curl <DOMAIN>:<PORT>/timeFile?inter=eth2
```