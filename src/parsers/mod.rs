pub mod number;

pub trait Parser {
    fn can_parse(&self, s: &str) -> bool;
    fn do_work(&self, s: &str) -> String;
}
