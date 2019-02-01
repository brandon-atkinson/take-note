extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate docopt;

mod tn;

pub use self::tn::args as args;
pub use self::tn::run as run;