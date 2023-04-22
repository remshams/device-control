use keylight_on::display;
use keylight_on::keylight::{
    KeylightControl, KeylightError, KeylightJsonDb, KeylightRestAdapter, ZeroConfKeylightFinder,
    KEYLIGHT_DB_PATH,
};
mod cli;

fn main() -> Result<(), KeylightError> {
    let command_light = cli::parse();
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let db = KeylightJsonDb::new(KEYLIGHT_DB_PATH);
    let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
    let action = || keylight_control.load_lights();
    display::progress::run(
        action,
        String::from("Discovering lights"),
        String::from("Lights discovered"),
    )?;
    let mut keylights = keylight_control.lights;
    let light = keylights.get_mut(0).unwrap();
    light.lights()?;
    light.set_light(command_light)
}
