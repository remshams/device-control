use crate::keylight::KeylightMetadata;
use console::{style, StyledObject};

fn pad_text(padding: usize) -> StyledObject<String> {
    style(format!("{}", " ".repeat(padding)))
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

pub fn print_keylights(keylights: &Vec<&KeylightMetadata>) {
    println!("{}:", style("Keylights").bold());
    for metadata in keylights {
        print_keylight_metadata(metadata);
    }
}
