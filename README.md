# Plan Trabajo Docente - SGA MID

API MID intermediaria entre el cliente SGA y las APIs CRUD relacionadas con el plan de trabajo docente. Gestiona los endpoints requeridos para la administración de la información necesaria en los módulos del SGA cliente, incluyendo preasignación, asignación, gestión de planes, reportes, espacios académicos/físicos, calendarios y docentes.

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang 1.25](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)

### Variables de Entorno
```shell
# Configuración de la aplicación
TRABAJO_DOCENTE_MID_HTTPPORT: [Puerto de ejecución de la API, por defecto 8552]
TRABAJO_DOCENTE_MID_RUNMODE: [Modo de ejecución: dev o prod, por defecto prod]
TRABAJO_DOCENTE_MID_STATIC_PATH: [Ruta de archivos estáticos, por defecto static]
PARAMETER_STORE: [Configuración de Parameter Store]

# Servicios Externos
TERCEROS_SERVICE: [URL del servicio de terceros CRUD]
PARAMETRO_SERVICE: [URL del servicio de parámetros CRUD]
PLAN_TRABAJO_DOCENTE_SERVICE: [URL del servicio CRUD de plan trabajo docente]
HORARIO_SERVICE: [URL del servicio de horarios CRUD]
ESPACIO_ACADEMICO_SERVICE: [URL del servicio de espacios académicos CRUD]
PROYECTO_ACADEMICO_SERVICE: [URL del servicio de proyecto académico CRUD]
OIKOS_SERVICE: [URL del servicio OIKOS CRUD]
ACADEMICA_ESPACIO_ACADEMICO_SERVICE: [URL del servicio académica espacio académico]
FIRMA_ELECTRONICA_MID_SERVICE: [URL del servicio MID de firma electrónica]
DOCUMENTO_SERVICE: [URL del servicio de documentos]
```
**NOTA:** Las variables se pueden ver en el fichero `conf/app.conf` e incluir en `.env`.

### Ejecución del Proyecto
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/trabajo_docente_mid

#2. Moverse a la carpeta del repositorio
cd trabajo_docente_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. Alimentar todas las variables de entorno que utiliza el proyecto.
# Usar el fichero .env como referencia
export TRABAJO_DOCENTE_MID_HTTPPORT=8552 ...

# 5. Ejecutar comandos para descargar dependencias
go mod tidy

# 6. Ejecutar proyecto
bee run -downdoc=true -gendoc=true
```

### Ejecución Dockerfile
```shell
# Construir la imagen
docker build --tag=trabajo_docente_mid . --no-cache

# Ejecutar el contenedor
docker run -p 8552:8552 trabajo_docente_mid
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/trabajo_docente_mid

#2. Moverse a la carpeta del repositorio
cd trabajo_docente_mid

#3. Crear un fichero con el nombre **custom.env**
# En windows ejecutar:* ` ni custom.env`
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Ejecución Pruebas

Pruebas unitarias
```shell
# Ejecutar todas las pruebas
go test ./test/... -v

# Ejecutar pruebas de un paquete específico
go test ./test/services/... -v
```

## Estado CI

| Develop | Release 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/trabajo_docente_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/trabajo_docente_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/trabajo_docente_mid/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/trabajo_docente_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/trabajo_docente_mid/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/trabajo_docente_mid) |

## Licencia

This file is part of trabajo_docente_mid.

trabajo_docente_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

trabajo_docente_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with trabajo_docente_mid. If not, see https://www.gnu.org/licenses/.
