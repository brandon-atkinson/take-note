use docopt::Docopt;

const VERSION: Option<&'static str> = option_env!("CARGO_PKG_VERSION");

const USAGE: &str = "
Usage: tn edit <note-name>
       tn list
       tn remove <note-name>
       tn show <note-name>
       tn (-h | --help)
       tn --version
       tn --bash-completion
       tn --commands

Options:
    -h --help           Show this screen.
    --version           Show version.
    --bash-completion   Prints out bash completion script.
    --commands          List supported commands.
";

type Result<T> = ::std::result::Result<T, Box<::std::error::Error>>;

#[derive(Debug, Deserialize)]
struct Args {
    cmd_show: bool,
    cmd_edit: bool,
    cmd_list: bool,
    cmd_remove: bool,
    flag_bash_completion: bool,
    flag_commands: bool,
    arg_note_name: Option<String>,
}

impl Args {
    pub fn command(self) -> Result<Command> {
        if self.cmd_show {
            return Ok(Command::ShowNote(self.arg_note_name.unwrap()))
        }
        if self.cmd_edit {
            return Ok(Command::EditNote(self.arg_note_name.unwrap()))
        }
        if self.cmd_list {
            return Ok(Command::ListNotes)
        }
        if self.cmd_remove {
            return Ok(Command::RemoveNote(self.arg_note_name.unwrap()))
        }
        if self.flag_bash_completion {
            return Ok(Command::BashCompletion)
        }
        if self.flag_commands {
            return Ok(Command::Commands)
        }
        Err("invalid command")?
    }
}

#[derive(Debug)]
pub enum Command {
    ShowNote(String),
    EditNote(String),
    ListNotes,
    RemoveNote(String),
    BashCompletion,
    Commands,
}

impl Command {
    pub fn parse(argv: std::env::Args) -> Result<Command> {
        let version = VERSION.map_or(None, |v| Some(v.to_string()));
        let args: Args = Docopt::new(USAGE)
            .and_then(|d| {
                d.version(version)
                    .argv(argv)
                    .deserialize()
            })?;
        Ok(args.command()?)
    }
}
