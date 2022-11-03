This demonstrates issues connecting to Go `x/crypto/ssh` servers using OpenSSH
client.

On certain versions, the SSH client, when using an RSA identity, in the absence
of server extension indicating otherwise will select handshake algorithms that
the Go `x/crypto/ssh` server does not support.

Known affected versions:

* Debian Sid `OpenSSH_9.0p1 Debian-1+b2, OpenSSL 3.0.5 5 Jul 2022`
* Ubuntu Jammy `OpenSSH_8.9p1 Ubuntu-3, OpenSSL 3.0.2 15 Mar 2022`
* Fedora 36 `OpenSSH_8.8p1, OpenSSL 3.0.5 5 Jul 2022`
* macOS Ventura

To test your version of `ssh`, run: `go run .`
