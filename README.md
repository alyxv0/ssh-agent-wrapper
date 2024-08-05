# sshad - Simple SSH-Agent database

sshad (ssh-agent wrapper daemon) enables automatic execution of ssh-add with the correct IdentityFile from `~/.ssh/config` via resolution of 'Host'-Entries through sshac (ssh-agent wrapper client).

## Content

- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [Knowhow](#Knowhow)

# Features
- automatic Host-Entry recognition from `~/.ssh/config`
- detection of cross-used keys in Host-Entries
- regular ssh usage possible

## Install

- sshad - `go install gitlab.com/mortytheshorty/ssh-wrapper/cmd/sshad`
- sshad  -  `go install gitlab.com/mortytheshorty/ssh-wrapper/cmd/sshac`

## Knowhow
`sshad` is running as the current user under `/home/$USER/.local/run/sshad`. On initialization it parses the content of the `~/.ssh/config` file for Host-Entries and builds
an in memory database (simple array of struct pointers) with the current state of the corresponsing key files. The idea is to use `sshac` client program to connect via these Host-Entry-Names. It connects to the local unix socket from `sshad` and sends it the host parameter given to it. `sshad` responds either with the key files path or with OK if the key is already loaded. FAILED if they Host-Entry was not found. After a key is loaded, `sshad` checks if the needed key is also used in another Host-Entry and if so, this HostKeyEntry will be marked as loaded too.


# Contribution
Improvements and features suggestions are always welcome. Please create an issue for bigger changes. Since this is my first project which is worth to share in my opinion, I could need some help in how to organize this.




