extern crate regex;

use regex::Regex;

#[derive(Debug)]
pub struct Number {}

impl super::Parser for Number {
    fn can_parse(&self, s: &str) -> bool {
        println!("number parser {}", s);
        let re = Regex::new(r"^[0-9]").unwrap();
        return re.is_match(s);
    }

    fn do_work(&self, s: &str) -> String {
        let num: isize = s.parse().unwrap();
        println!("{}", num);
        String::from("some string ")
    }
}
