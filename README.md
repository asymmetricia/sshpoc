This demonstrates issues connecting to Go `x/crypto/ssh` servers using OpenSSH
client.

On certain versions (e.g., `OpenSSH_8.9p1 Ubuntu-3, OpenSSL 3.0.2 15 Mar
2022`), the SSH client, when using an RSA identity, in the absence of server
extension indicating otherwise will select handshake algorithms that the Go
`x/crypto/ssh` server does not support.

To test your version of `ssh`, run: `test.sh`

To confirm that the behavior is not reproducible with Ed25519 keys, run
`test.sh --ed25519`
