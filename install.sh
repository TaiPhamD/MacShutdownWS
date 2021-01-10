if [ ! -d "/opt/shutdown" ] 
then
    echo "Directory /opt/shutdown DOES NOT exists. creating dir" 
    mkdir /opt/shutdown
fi

cp NixShutdownWS /opt/shutdown

if [ ! -e "/opt/shutdown/key.pem" ] 
then
    echo "HTTPS certs does not exist. Generating using openssl." 
    openssl req  -nodes -new -x509  -keyout key.pem -out cert.pem
fi
