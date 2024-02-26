use std::io;
use std::io::BufRead;
use clap::{Arg, Command};

fn main(){
    // simple running of each stdin parameters against the args
    let matches = Command::new("cargs")
        .version("0.1.0")
        .author("Ayobami Bamigboye <ayo@greystein.com>")
        .about("build and execute command lines from standard input")
        .arg(Arg::new("max-args")
            .short('n')
            .default_value("1")
            .help("Use at most max-args arguments per command line."))
        .arg(Arg::new("max-procs")
            .short('P')
            .default_value("1")
            .help("Run up to max-procs processes at a time; the default is 1. If max-procs is 0, xargs will run as many processes as possible at a time. Use the -n option with -P; otherwise chances are that only one exec will be done."))
            .external_subcommand_value_parser(clap::value_parser!(usize))
        .get_matches();

    // println!("Matches: {:#?}", matches);

    // get stdin into buffer
    let mut lines =  io::stdin().lock().lines();
    while let Some(line) = lines.next() {
        let input = line.unwrap();

        if input.len() == 0 {
            break;
        }

        // continue with processing line CoreOp
        process(matches, line)
    }

}

fn process(matches: clap::ArgMatches, line String) {

}