use env_logger::{Builder, Env};
use keylight_on::display;
use keylight_on::keylight::{
    KeylightControl, KeylightError, KeylightJsonDb, KeylightRestAdapter, ZeroConfKeylightFinder,
    KEYLIGHT_DB_PATH,
};
mod cli;

fn setup_logger() {
    Builder::from_env(Env::default().default_filter_or("info")).init();
}

fn main() -> Result<(), KeylightError> {
    setup_logger();

    let command_light = cli::parse();
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let db = KeylightJsonDb::new(KEYLIGHT_DB_PATH);
    let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
    display::progress::run(
        || keylight_control.load_keylights(),
        String::from("Discovering lights"),
        String::from("Lights discovered"),
    )?;
    let light = keylight_control
        .find_keylight_mut(&command_light.id)
        .ok_or(KeylightError::KeylightDoesNotExist(
            command_light.id.clone(),
        ))?;

    light.set_light(command_light)
}
