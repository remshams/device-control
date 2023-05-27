use keylight_control::keylight::{Keylight, KeylightError};
use tauri::State;

use crate::model::AppState;

#[tauri::command]
pub fn discover_keylights(state: State<AppState>) -> Result<Vec<Keylight>, KeylightError> {
    let mut keylight_control = state
        .keylight_control
        .lock()
        .map_err(|_err| KeylightError::CommandError(String::from("test")))?;
    keylight_control.load_keylights()?;
    let lights = &mut keylight_control.lights;
    for light in lights.iter_mut() {
        light.lights(&state.adapter)?;
    }
    Ok(keylight_control.lights.clone())
}
