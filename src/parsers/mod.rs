pub mod number;

pub trait Parser {
    fn can_parse_human_into(&self, s: &str) -> bool;
    fn can_parse_human_from(&self, s: &str) -> bool;
    fn do_human_into(&self, s: &str) -> String;
    fn do_human_from(&self, s: &str) -> String;
}
