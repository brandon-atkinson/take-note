# take note (tn) - dead simple note taking 


## Building & Installing

Building requires an installation of the [Racket programming
language](https://racket-lang.org/). Once that is installed, building and
installing can be done in one step with the normal make incantation:

```
cd tn/
make install
```

## Enabling shell completion

This note taking tool (it's barely a program) has the ability to leverage bash
command-line completion to speed up references to existing notes. If you have
installed and configured the `bash-completion` package via homebrew, the following will install a user-specific completion script: 

```
tn --bash-completion-script >> ~/.bash_completion && exec bash
```

If you don't have `bash-completion` installed (or you've installed tn
system-wide and wish to enable completion system wide as well) you can 
install to the bash_completion.d directory instead: 

```
tn --bash-completion-script > $(brew --prefix)/etc/bash_completion.d/tn \
&& exec bash
```
