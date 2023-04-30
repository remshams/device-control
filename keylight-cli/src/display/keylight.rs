use console::{style, StyledObject};
use keylight_control::keylight::KeylightMetadata;

fn pad_text(padding: usize) -> StyledObject<String> {
    style(" ".repeat(padding))
}

fn print_key_value(padding: usize, key: &str, value: &str) {
    println!("{}{}{}", pad_text(padding), style(key).bold(), value);
}

pub fn print_keylight_metadata(metadata: &KeylightMetadata) {
    println!("{}{}", pad_text(2), style(&metadata.name).underlined());
    print_key_value(4, "Id: ", &metadata.id);
    print_key_value(4, "IP: ", &metadata.ip);
    print_key_value(4, "Port: ", &metadata.port.to_string());
}

pub fn print_keylights(headline: Option<&str>, keylights: &Vec<&KeylightMetadata>) {
    if let Some(headline) = headline {
        println!("{}", style(headline).bold().green());
    }
    for (index, metadata) in keylights.iter().enumerate() {
        print_keylight_metadata(metadata);
        if index < keylights.len() - 1 {
            println!();
        }
    }
}
