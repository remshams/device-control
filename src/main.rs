use keylight_on::display;
use keylight_on::keylight::{
    KeylightControl, KeylightError, KeylightRestAdapter, ZeroConfKeylightFinder,
};
mod cli;

fn main() -> Result<(), KeylightError> {
    let command_light = cli::parse();
    println!("{:?}", command_light);
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
    let light = keylights.get_mut(0).unwrap();
    light.lights()?;
    light.set_light(command_light)
}
