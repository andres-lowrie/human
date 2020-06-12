use structopt::StructOpt;

mod parsers;
use parsers::Parser;

#[derive(StructOpt)]
struct Cli {
    input: String,
}

fn main() {
    let args = Cli::from_args();
    // for each parser
    let p = parsers::number::Number {};
    if p.can_parse_human_from(&args.input) {
        let got = p.do_human_into(&args.input);
        print!("{}", got);
    }
}
