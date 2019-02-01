# take-note (tn) - simple note taking 

## Usage

### Specifying the editor

You'll most likely want to use your editor of choice. If you've set the 
`EDITOR` environment variable in your shell config already, you're done, `tn`
will use it. Otherwise, set your editor in your `.bashrc` and restart your
shell:

```
echo "export EDITOR='nvim'" >> ~/.bashrc
exec bash
```

### Taking your first note

```
tn edit my-first-note
```

Then save/quit from your editor. 

### Editing notes

Assuming you've already edited and saved a note using `tn edit` with the name
`existing-note-name` you just have to invoke the command again to review or
edit.

```
tn edit existing-note-name
```

### Showing notes

If you just want to print the note text to the terminal (like cat would):

```
tn show existing-note-name
```

### Lising notes

Want to get a list of notes you've taken? 

```
tn list
```

The list will be sorted in reverse chronological order, so the most recently
edited notes should be right above your command prompt.

### Removing notes

## Building & Installing

The project is written using the Rust language. Assuming you have already
installed a recent version, you should be able to build the project using
cargo:

```
cargo build --release
```

You can then copy the `target/release/tn` executable somewhere on your path.

## Enabling shell completion

This note taking tool (it's barely a program) has the ability to leverage bash
command-line completion to speed up references to existing notes. If you have
installed and configured the `bash-completion` package via homebrew, the
following will install a user-specific completion script: 

```
tn --bash-completion >> ~/.bash_completion 
exec bash
```

If you don't have `bash-completion` installed (or you've installed tn
system-wide and wish to enable completion system wide as well) you can 
install to the bash_completion.d directory instead: 

```
tn --bash-completion > $(brew --prefix)/etc/bash_completion.d/tn 
exec bash
```
