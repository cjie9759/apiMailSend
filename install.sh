
wget -o /usr/local/bin/apiMailSend
cat >/etc/systemd/system/apiMailSend.service <<EOL
[Unit]
After=network.target nss-lookup.target

[Service]
User=root
Environment="MAIL_USER=$1"
Environment="MAIL_PWD=$2"
Environment="MAIL_NAME=$3"
Environment="MAIL_API_LISTEN=:19876"
Environment="MAIL_SIGN=636a6965636a6965d41d8cd98f00b204e9800998ecf8427e"
ExecStart=/usr/local/bin/apiMailSend
Restart=on-failure
# RestartPreventExitStatus=23
RestartSec=3
[Install]
WantedBy=multi-user.target
EOL
