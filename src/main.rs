use std::time::Duration;

use indicatif::{ProgressBar, ProgressStyle};
use keylight_on::keylight::{KeylightControl, KeylightRestAdapter, ZeroConfKeylightFinder};

fn main() {
    let spinner = ProgressBar::new_spinner();
    spinner.enable_steady_tick(Duration::from_millis(100));
    spinner.set_style(
        ProgressStyle::default_spinner()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"])
            .template("{spinner:.white} {msg}")
            .unwrap(),
    );
    spinner.set_message("Discovering Keylights...");
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let mut keylight_control = KeylightControl::new(&finder, &adapter);
    keylight_control.discover_lights();
    spinner.finish_with_message("Found Keylights!");
    let mut keylights = keylight_control.lights;
    for keylight in keylights.iter() {
        println!("{:?}", keylight.metadata);
    }
    for keylight in keylights.iter_mut() {
        match keylight.lights() {
            Ok(lights) => {
                for light in lights {
                    println!("{:?}", light);
                }
            }
            Err(e) => {
                println!("{:?}", e);
            }
        }
    }
    let command_result = keylights[0].toggle(0);
    println!("{:?}", command_result);
}
