use console::style;
use env_logger::{Builder, Env};
use keylight_on::display;
use keylight_on::keylight::{
    KeylightCommand, KeylightControl, KeylightError, KeylightJsonDb, KeylightRestAdapter,
    ZeroConfKeylightFinder, KEYLIGHT_DB_PATH,
};
mod cli;

fn setup_logger() {
    Builder::from_env(Env::default().default_filter_or("info")).init();
}

fn main() -> Result<(), KeylightError> {
    setup_logger();

    let keylight_command = cli::parse();
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let db = KeylightJsonDb::new(KEYLIGHT_DB_PATH);
    let mut keylight_control = KeylightControl::new(&finder, &adapter, &db);
    display::progress::run(
        || keylight_control.load_keylights(),
        String::from("Discovering lights"),
        String::from("Lights discovered"),
    )?;
    match keylight_command {
        KeylightCommand::SendCommand(light_command) => {
            let light = keylight_control
                .find_keylight_mut(&light_command.id)
                .ok_or(KeylightError::KeylightDoesNotExist(
                    light_command.id.clone(),
                ))?;
            light.lights()?;

            light.set_light(light_command)
        }
        KeylightCommand::List => {
            display::keylight::print_keylights(Some("Lights: "), &keylight_control.list_metadata());
            Ok(())
        }
        KeylightCommand::Discover => {
            keylight_control.discover_and_store_keylights()?;
            println!(
                "Discovered and stored {} keylights\n",
                style(keylight_control.list_metadata().len()).green().bold()
            );
            display::keylight::print_keylights(None, &keylight_control.list_metadata());
            Ok(())
        }
    }
}
