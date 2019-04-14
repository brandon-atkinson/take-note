use std::path::PathBuf;
use std::fs;
use std::env;
use std::process;
use std::fs::DirEntry;
use regex::Regex;
use std::result::Result::Ok;
use crate::tn::args::Command;
use std::ffi::OsString;

pub mod args;

const APP_NAME: &str = "tn";

type Result<T> = ::std::result::Result<T, Box<::std::error::Error>>;

pub fn run(args: std::env::Args) -> Result<()> {
    let cmd = Command::parse(args)?;
    let result = run_cmd(cmd)?;
    Ok(result)
}

fn run_cmd(cmd: Command) -> Result<()> {
    match cmd {
        Command::ShowNote(ref note_name) => cat_note(note_name),
        Command::EditNote(ref note_name) => edit_note(&note_name),
        Command::ListNotes => list_notes(),
        Command::RemoveNote(ref note_name) => remove_note(&note_name),
        Command::BashCompletion => bash_completion(),
        Command::Commands => commands(),
    }
}

fn note_dir() -> Result<PathBuf> {
    let mut data_dir = dirs::document_dir().ok_or("data dir not available")?;
    data_dir.push(APP_NAME);
    data_dir.push("notes");

    fs::create_dir_all(&data_dir.as_path())?;
    Ok(data_dir)
}

fn path_for_name(note_name: &str) -> Result<PathBuf> {
    let mut note_path = note_dir()?;
    note_path.push(name_to_filename(note_name));
    Ok(note_path)
}

fn name_to_filename(name: &str) -> String {
    format!("{}.md", name)
}

fn filename_to_name(filename: &str) -> String {
    let pattern = Regex::new("[.]md$").unwrap();
    pattern.replace(filename, "").to_string()
}

fn editor() -> String {
    match env::var("EDITOR") {
        Ok(editor) => editor,
        Err(_) => "vim".to_string(),
    }
}

fn shell() -> String {
    match env::var("SHELL") {
        Ok(editor) => editor,
        Err(_) => "sh".to_string(),
    }
}

fn edit_cmd(path: PathBuf) -> process::Command {
    let mut edit_shell_cmd = OsString::from(editor());
    edit_shell_cmd.push(" ");
    edit_shell_cmd.push(path);

    let mut edit_cmd = process::Command::new(shell());
    edit_cmd.arg("-c");
    edit_cmd.arg(edit_shell_cmd);
    edit_cmd
}

fn cat_note(note_name: &str) -> Result<()> {
    let path = path_for_name(note_name)?;
    let content = fs::read_to_string(path)?;
    println!("{}", content);
    Ok(())
}

fn edit_note(note_name: &str) -> Result<()> {
    let mut edit_cmd = edit_cmd(path_for_name(note_name)?);

    Ok(match edit_cmd.status() {
        Ok(status) => if status.success() {
            Ok(())
        } else {
            Err("command failed".to_string())
        },
        Err(e) => Err(e.to_string())
    }?)
}

fn remove_note(note_name: &str) -> Result<()> {
    let path = path_for_name(note_name)?;
    fs::remove_file(path)?;
    Ok(())
}

fn list_notes() -> Result<()> {
    let note_dir = note_dir()?;
    let dir = fs::read_dir(note_dir)?;

    let mut entries: Vec<DirEntry>  = dir
        .filter(|e| e.is_ok())
        .map(|e| e.unwrap())
        .collect();

    entries.sort_by(|a, b| {
        let a_modified = a.metadata().ok()
            .and_then(|m| m.modified().ok());

        let b_modified = b.metadata().ok()
            .and_then(|m| m.modified().ok());

        a_modified.cmp(&b_modified)
    });

    let filenames: Vec<String> = entries.iter()
        .map(|e| e.file_name())
        .map(|f| f.into_string())
        .filter(|f| f.is_ok())
        .map(|f| f.unwrap())
        .filter(|f| f.ends_with(".md"))
        .collect();

    for filename in filenames {
        println!("{}", filename_to_name(filename.as_str()));
    }

    Ok(())
}

fn bash_completion() -> Result<()> {
    let script = r##"
_tn_complete() {
    local cur prev opts

    COMPREPLY=()

    cur=${COMP_WORDS[COMP_CWORD]}
    prev=${COMP_WORDS[COMP_CWORD-1]}
    opts=""

    if [[ "$prev" == "$1" ]]; then
        opts=$(tn --commands)
    else
        case "$prev" in
            edit|remove|show)
                opts=$(tn list)
            ;;
        esac

    fi

    if [[ ! -z ${opts} ]] ; then
        COMPREPLY=( $(compgen -o nosort -W "${opts}" -- ${cur}) )
        return 0
    fi
}

complete -F _tn_complete tn
"##;
    println!("{}", script);

    Ok(())
}

fn commands() -> Result<()> {
    let commands = vec!(
        "edit",
        "list",
        "remove",
        "show",
        "--help",
        "--version",
        "--bash-completion",
        "--commands");

    for c in commands {
        println!("{}", c);
    }

    Ok(())
}


