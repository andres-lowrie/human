extern crate regex;

use regex::Regex;

#[derive(Debug)]
pub struct Number {}

impl super::Parser for Number {
    fn can_parse_human_into(&self, s: &str) -> bool {
        false
    }

    fn can_parse_human_from(&self, s: &str) -> bool {
        let re = Regex::new(r"^[0-9]").unwrap();
        re.is_match(s)
    }

    fn do_human_into(&self, s: &str) -> String {
        // Quick exit, anything less than a thousand doesn't require grouping
        if s.len() < 4 {
            return String::from(s);
        }

        let mut out: Vec<_> = Vec::new();
        let tmp = s.chars().rev();
        for (i, c) in tmp.enumerate() {
            if i % 3 == 0 && i != 0 {
                out.push(',');
            }
            out.push(c);
        }
        out.iter().rev().collect::<String>()
    }

    fn do_human_from(&self, s: &str) -> String {
        String::from("something")
    }
}
