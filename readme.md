# Gedis

**Gedis** es un servidor de almacenamiento clave-valor en Go, inspirado en Redis.  
Soporta comandos `SET`, `GET`, `DEL`, `PING` e `INFO` sobre TCP, con protocolo de texto plano o RESP.

Manejando en los tests aproximadamente 50 mil peticiones/s

<img width="217" height="93" alt="minibenchmark" src="https://github.com/user-attachments/assets/e8ddbf16-8a7c-4c42-802c-05baae541160" />

**IMPORTANTE: Es un proyecto de fin de semana, usar bajo su propio riesgo**

Este proyecto lo he desarrollado en mi tiempo libre, no recomiendo usarlo en produccion por el momento

## Estructura

- `main.go` — servidor TCP con soporte para ambos protocolos.
- `protocol/` — parseo de protocolo plano (`ParseArgs`) y RESP (`ReadRESPCommand`, `PeekProtocol`).
- `database/` — almacenamiento sharded (1024 shards) con persistencia AOF y reescritura compacta.
- `commands/` — capa de comandos que wrappea la base de datos.
- `test/test.go` — cliente de prueba concurrente.

## Cómo correrlo

```bash
go run main.go              # servidor en puerto 64666
go run ./test/test.go       # cliente de prueba
```

## Comandos

| Comando | Formato texto plano | Formato RESP |
|---------|-------------------|--------------|
| SET     | `3,SET,3,key,5,value` | `*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n` |
| GET     | `3,GET,3,key` | `*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n` |
| DEL     | `3,DEL,3,key` | `*2\r\n$3\r\nDEL\r\n$3\r\nkey\r\n` |
| PING    | `4,PING` | `*1\r\n$4\r\nPING\r\n` |
| INFO    | `4,INFO` | `*1\r\n$4\r\nINFO\r\n` |

## Features

- **Dual protocolo** — detecta automáticamente si el cliente habla RESP o texto plano.
- **Persistencia AOF** — cada escritura se registra en disco; al iniciar se reconstruye el estado.
- **Escritura batch** — comandos se bufferizan y flushean cada 1 segundo en background, evitando syscalls por operación.
- **Reescritura automática** — compacta el AOF cuando alcanza 10,000 operaciones.
- **Sharding** — 1024 shards con locks separados para alto rendimiento concurrente.
- **Info de servidor** — comando `INFO` con versión, uptime, cantidad de claves y uso de memoria.

## Makefile

```bash
make build   # compilar
make test    # tests con race detector
make bench   # benchmarks
make fuzz    # fuzzing (ParseArgs, ReadRESPCommand)
make run     # compilar y ejecutar
make clean   # limpiar binario y AOF
```
