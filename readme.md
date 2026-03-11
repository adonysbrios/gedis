
# Gedis

**Gedis** es un proyecto Go que implementa un servidor de almacenamiento clave-valor simple, inspirado en sistemas como Redis.  
El servidor acepta conexiones TCP y permite ejecutar comandos básicos (`SET`, `GET`, `DEL`, `PING`) en texto plano.  

Manejando en los tests aproximadamente 50 mil peticiones/s
<img width="217" height="93" alt="minibenchmark" src="https://github.com/user-attachments/assets/e8ddbf16-8a7c-4c42-802c-05baae541160" />

**NO USAR EN PRODUCCION**

## 📂 Estructura del proyecto

- `main.go` → corre el servidor TCP.
- `test/test.go` → cliente de prueba para enviar comandos y medir rendimiento.

## 🚀 Cómo correrlo

1. Compila y ejecuta el servidor:
   ```bash
   go run main.go
   ```

2. En otra terminal, ejecuta el cliente de prueba:
   ```bash
   go run ./test/test.go
   ```

El cliente abrirá conexiones al servidor y enviará múltiples comandos para validar su funcionamiento y medir throughput.

## ✅ Comandos soportados

- `SET <key> <value>` → guarda un valor asociado a una clave.
- `GET <key>` → obtiene el valor de una clave.
- `DEL <key>` → elimina una clave.
- `PING` → devuelve `PONG`.

## 📋 To-Do

- [ ] Implementar persistencia de datos (guardar en disco y recuperar al reiniciar)
- [ ] Implementar un protocolo más eficiente que el texto plano (ej. RESP estilo Redis)
- [ ] Implementar autenticación con clave
 
