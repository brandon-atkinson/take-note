extern crate tn;

use std::process;
use tn::run;

fn main() {
    if let Err(e) = run(std::env::args()) {
        eprintln!("{}", e);
        process::exit(1);
    }
}


