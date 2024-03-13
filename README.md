
## Configure

Create the signing certificate
```bash
$ ssh-keygen -t ed25519 -C ca@github.com -f ca
$ cat ca.pub
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOixUTX9ssW7bAaO6wTxXxJGRpVWNnqnOFwFZ1ceOxVn ca@github.com
```

Use the public key above to create a new certificate authority in the GitHub organization settings under **Authentication security**.

## Build
Install the local repo as a `gh` cli extension:

```bash 
gh extension install .     
```

Build and run:

```bash
go build && gh ssh-cert-please [command]
```

