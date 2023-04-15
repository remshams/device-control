use keylight_on::display;
use keylight_on::keylight::{KeylightControl, KeylightRestAdapter, ZeroConfKeylightFinder};

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let mut keylight_control = KeylightControl::new(&finder, &adapter);
    let action = || {
        keylight_control.discover_lights();
    };
    display::progress::run(
        action,
        String::from("Discovering lights"),
        String::from("Lights discovered"),
    );
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
