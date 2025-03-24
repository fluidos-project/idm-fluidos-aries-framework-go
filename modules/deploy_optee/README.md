# Despliegue de OP-TEE

Esta carpeta incluye un script para desplegar OP-TEE en el entorno. Este script automatiza gran parte del proceso de configuración y compilación del sistema.

## Pasos para la instalación y ejecución

1. **Ejecutar el script de despliegue:**
   ```sh
   ./deploy_optee.sh
   ```

2. **Aplicar modificaciones personalizadas:**
   ```sh
   git clone https://ants-gitlab.inf.um.es/jmanuel/dpabc_adp_trusted_application
   cp -r dpabc_adp_trusted_application/* modules/optee/optee_examples

   git clone https://ants-gitlab.inf.um.es/jmanuel/pabc-ta-aries-module-main
   cp -r pabc-ta-aries-module-main/* modules/optee/out-br/target/

   cd modules/optee/build
   make
   ```

3. **Sustituir el archivo `common.mk` (disponible en optee/build)** por el que se incluye en esta carpeta.

4. **Añadir la carpeta `api`** incluida en esta carpeta a `out-br/target`.

5. **Ejecutar el escenario (desde optee/build):**
   ```sh
   make run BR2_PACKAGE_PYTHON3=y BR2_PACKAGE_PYTHON3_SSL=y BR2_PACKAGE_PYTHON_PIP=y RESTAPI=y
   ```

6. **Con el escenario en marcha, instalar Flask en el Normal World (user root):**
   ```sh
   pip install flask
   ```

7. **Lanzar la API, también desde el Normal World (user test):**
   ```sh
   python3 api/TEE-api.py &
   ```



