// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use home::home_dir;
use std::sync::Mutex;

use keylight_control::keylight::{
    KeylightControl, KeylightJsonDb, KeylightRestAdapter, ZeroConfKeylightFinder,
};
use model::AppState;

use crate::keylight::{discover_keylights, refresh_lights, set_light};

mod keylight;
mod model;

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

fn main() {
    let db_dir = home_dir().map(|mut home_dir| {
        home_dir.push("keylight.json");
        home_dir
    });
    let finder = ZeroConfKeylightFinder::new();
    let db = KeylightJsonDb::new(db_dir);
    let adapter = KeylightRestAdapter {};
    let control = KeylightControl::new(finder, db);
    let app_state = AppState {
        adapter,
        keylight_control: Mutex::new(control),
    };
    tauri::Builder::default()
        .manage(app_state)
        .invoke_handler(tauri::generate_handler![
            greet,
            discover_keylights,
            refresh_lights,
            set_light
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
