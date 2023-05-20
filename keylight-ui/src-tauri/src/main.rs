// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::{borrow::BorrowMut, sync::Mutex};

use keylight_control::keylight::{
    KeylightAdapter, KeylightControl, KeylightDb, KeylightError, KeylightFinder, KeylightJsonDb,
    KeylightRestAdapter, ZeroConfKeylightFinder, KEYLIGHT_DB_PATH,
};
use tauri::{InvokeError, State};

// fn convertKeylightResultToInvokeError(keylight_error: &KeylightError) -> InvokeError {
//     InvokeError::from(keylight_error)
// }

struct AppState<F: KeylightFinder, Db: KeylightDb, A: KeylightAdapter> {
    adapter: A,
    keylight_control: Mutex<KeylightControl<F, Db>>,
}

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn discover_keylights(
    state: State<AppState<ZeroConfKeylightFinder, KeylightJsonDb, KeylightRestAdapter>>,
) -> Result<(), KeylightError> {
    let keylight_control = state.keylight_control.lock();
    match keylight_control {
        Ok(mut keylight_control) => keylight_control.load_keylights(),
        Err(_e) => Err(KeylightError::CommandError(String::from(
            "Could not load keylights",
        ))),
    }
}

#[tauri::command]
fn load_lights(
    state: State<AppState<ZeroConfKeylightFinder, KeylightJsonDb, KeylightRestAdapter>>,
) -> Result<(), KeylightError> {
    let keylight_control = state.keylight_control.lock();
    match keylight_control {
        Ok(mut keylight_control) => {
            let lights = &mut keylight_control.lights;
            for light in lights.iter_mut() {
                light.lights(&state.adapter)?;
            }
            Ok(())
        }
        Err(_e) => Err(KeylightError::CommandError(String::from(
            "Could not load lights for keylights",
        ))),
    }
}

#[tauri::command]
fn turn_keylight_on(
    state: State<AppState<ZeroConfKeylightFinder, KeylightJsonDb, KeylightRestAdapter>>,
) -> Result<(), KeylightError> {
    let mut keylight_control = state.keylight_control.lock().unwrap();
    let light = keylight_control
        .find_keylight_mut("0")
        .ok_or_else(|| KeylightError::KeylightDoesNotExist(String::from("0")));
    match light {
        Ok(light) => light.set_light(
            keylight_control::keylight::LightCommand {
                id: String::from("0"),
                index: 0,
                on: Some(true),
                brightness: None,
                temperature: None,
            },
            &state.adapter,
        ),
        Err(e) => Err(KeylightError::CommandError(String::from(
            "Could not turn on keylight",
        ))),
    }
}

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let db = KeylightJsonDb::new(KEYLIGHT_DB_PATH);
    let mut keylight_control = KeylightControl::new(finder, db);
    keylight_control.load_keylights();
    tauri::Builder::default()
        .manage(AppState {
            adapter,
            keylight_control: Mutex::new(keylight_control),
        })
        .invoke_handler(tauri::generate_handler![
            greet,
            load_lights,
            turn_keylight_on
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
