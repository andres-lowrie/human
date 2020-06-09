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
    if p.can_parse(&args.input) {
        let got = p.do_work(&args.input);
        print!("{}", got);
    }
}
