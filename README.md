# Proyecto de prueba

## Que quiero hacer?
### En forma general
Aprender estructura basica de micro servicios como armar todo con docker-compose y  repasar conceptos ya aprendidos sobre diferentes tipos de cache y bases de datos.

### A corto plazo
Tener un gateway robusto con un buen manejo de usuarios que sirva como guardian al resto de los microservicios.

### A largo plaso 
Quisiera tener una aplicacion REAL corriendo con diferentes tecnologias y que pueda escalar, esta aplicacion podria estar hostead en aws o google cloud

## que he hecho?
* por ahora solo tengo un gateway expuesto que filtra todos los llamados internos a mis micro servicios
* tengo un servicio de auth y uno de productos funcionando
* agregue funcionalidad de registrar usuarios a auth pero esto debe moverse a users

## Que Falta?
* Muchas cosas
* Tener funcionalidad completa crud de usuarios en usuarios
* Que Auth llame a usuarios para hacer el auth
* que el auth middleware llame a auth verify
* que el auth verify reaccione dependiendo de los roles.
* Experimentar con NoSql
* Experimentar con diferentes tipos de cache