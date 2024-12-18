FROM python:3.11.2-slim

WORKDIR /app

# Copiar archivos de la aplicación
COPY XACML_PAP ./XACML_PAP
COPY requirements.txt requirements.txt

# Instalar dependencias
RUN pip3 install --no-cache-dir -r requirements.txt

RUN pip3 install "xmltodict==0.13.0"

WORKDIR /app/XACML_PAP

CMD [ "python3", "serverPAP.py"]

