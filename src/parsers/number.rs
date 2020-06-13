extern crate regex;

use super::*;
use regex::Regex;

#[derive(Debug)]
pub struct Number {}

impl Parser for Number {
    fn can_parse_human_into(&self, s: &str) -> bool {
        false
    }

    fn can_parse_human_from(&self, s: &str) -> bool {
        let re = Regex::new(r"^[1-9]").unwrap();
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_can_parse_human_positive() {
        let p = Number {};
        let given = "1000";
        let want = true;
        assert_eq!(p.can_parse_human_from(given), want);
    }

    #[test]
    fn test_can_parse_human_from_negative() {
        let p = Number {};
        let given = "0";
        let want = false;
        assert_eq!(p.can_parse_human_from(given), want);
    }

    #[test]
    fn test_do_human_to_positive_thou() {
        let p = Number {};
        let given = "1000";
        let want = "1,000";
        assert_eq!(p.do_human_into(given), want);
    }

    #[test]
    fn test_do_human_to_positive_tenthou() {
        let p = Number {};
        let given = "10000";
        let want = "10,000";
        assert_eq!(p.do_human_into(given), want);
    }

    #[test]
    fn test_do_human_to_positive_hunthou() {
        let p = Number {};
        let given = "100000";
        let want = "100,000";
        assert_eq!(p.do_human_into(given), want);
    }

    #[test]
    fn test_do_human_to_positive_mil() {
        let p = Number {};
        let given = "1000000";
        let want = "1,000,000";
        assert_eq!(p.do_human_into(given), want);
    }

    #[test]
    fn test_do_human_to_positive_bil() {
        let p = Number {};
        let given = "1000000000";
        let want = "1,000,000,000";
        assert_eq!(p.do_human_into(given), want);
    }
}
