#!/bin/bash

# Verificar si el script se ejecuta dentro de "modules/deploy_optee/"
if [[ $(basename "$PWD") != "deploy_optee" ]]; then
    echo "❌ Error: Este script debe ejecutarse dentro de la carpeta 'modules/deploy_optee/'."
    echo "🔹 Usa: cd modules/deploy_optee && ./deploy_optee.sh"
    exit 1
fi

echo "🔹 Actualizando sistema e instalando dependencias..."
sudo apt update && sudo apt upgrade -y
sudo apt install -y adb acpica-tools autoconf automake bc bison build-essential \
    ccache cpio cscope curl device-tree-compiler e2tools expect fastboot flex \
    ftp-upload gdisk git libattr1-dev libcap-ng-dev libfdt-dev libftdi-dev libglib2.0-dev \
    libgmp3-dev libhidapi-dev libmpc-dev libncurses5-dev libpixman-1-dev libslirp-dev \
    libssl-dev libtool libusb-1.0-0-dev make mtools netcat ninja-build python3-cryptography \
    python3-pip python3-pyelftools python3-serial python-is-python3 rsync swig unzip \
    uuid-dev wget xdg-utils xsltproc xterm xz-utils zlib1g-dev

echo "🔹 Instalando tomli..."
pip install --user tomli

echo "🔹 Configurando repo..."
sudo curl -s https://storage.googleapis.com/git-repo-downloads/repo | sudo tee /bin/repo > /dev/null
sudo chmod a+x /bin/repo

echo "🔹 Instalando OP-TEE en modules/optee/..."
OPTEE_DIR="../optee"

# Clonar OP-TEE si no existe
if [ ! -d "$OPTEE_DIR" ]; then
    mkdir -p "$OPTEE_DIR"
    cd "$OPTEE_DIR"
    repo init -u https://github.com/OP-TEE/manifest.git -m qemu_v8.xml
    repo sync -j10
    echo "✅ OP-TEE instalado en modules/optee"
else
    echo "✅ OP-TEE ya está instalado en modules/optee"
fi

echo "🔹 Compilando toolchains..."
cd "$OPTEE_DIR/build"
make -j2 toolchains

echo "🔹 Verificando instalación..."
make -j$(nproc) check

echo "✅ OP-TEE instalado correctamente"
