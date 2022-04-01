# micromdm-ca-renew

## build

`go get && go build -o ./build/renew`

## replace current ca cert

1. export current ca and key from offline-and-very-backed-up micromdm bolt db

    `./build/renew -boltdb ./micromdm_i_swear_i_backed_this_up.db -export-ca`

1. create a new signing request

    `openssl x509 -x509toreq -in ./out.pem -signkey ./key.pem -out ./ca_csr.pem`

1. sign request with current ca certificate (1825 days == 5 years but adjust according to lazyness)

    `openssl x509 -sha256 -set_serial 2 -req -days 1825 -in ./ca_csr.pem -signkey ./key.pem -out ./ca_cert_extended.pem -extfile ./extensions.cnf -extensions CA_extensions`

1. diff your old and extended certs, note that extensions may be rearranged but should be consistent

    `diff <(openssl x509 -in ./ca_cert_extended.pem -text -noout) <(openssl x509 -in ./out.pem -text -noout)`

1. import new ca certificate into that backup copy of your very important and irreplicable micromdm bolt db

    `./build/renew -boltdb ./micromdm_i_swear_i_backed_this_up.db -import-ca -ca-cert ./ca_cert_extended.pem`

1. inspect the updated ca certificate from micromdm bolt db, note the new validity dates

    `./build/renew -boltdb ./micromdm_i_swear_i_backed_this_up.db -show-ca`
