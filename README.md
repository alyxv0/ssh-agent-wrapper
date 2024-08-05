# sshwd - Simple SSH-Agent database

sshwd (ssh-wrapper daemon) enables automatic execution of ssh-add with the correct IdentityFile from `~/.ssh/config` via resolution of 'Host'-Entries through sshw (ssh-wrapper client).

## Content

- [Install](#install)
- [Usage](#usage)
- [How it works](#how it works)

## How it works
`sshwd` is running as the current user under `/home/$USER/.local/run/sshwd`. On initialization it parses the content of the `~/.ssh/config` file for Host-Entries and builds
an in memory database (simple array of struct pointers) with the current state of the corresponsing key files. The idea is to use sshw client program to connect via these Host-Entry-Names. It connects to the local unix socket from `sshwd` and sends it the host parameter given to it. `sshwd` responds either with the key files path or with OK if the key is already loaded. FAILED if they Host-Entry was not found. After a key is loaded, `sshwd` checks if the needed key is also used in another Host-Entry and if so, it this Host-Entry will be activated (marked as Loaded = true) too.

# Features

- automatic Host-Entry recognition from `~/.ssh/config`
- detection of cross-used keys in Host-Entries
- regular ssh usage possible

## Install

- sshwd - `go install gitlab.com/alyxv/sshwd`
- sshw  -  `go gitlab.com/alyxv/sshw`

# Contribution
Improvements and features suggestions are always welcome. Please create an issue for bigger changes. Since this is my first project which is worth to share in my opinion, I could need some help in how to organize this.




