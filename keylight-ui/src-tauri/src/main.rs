// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::sync::Mutex;

use keylight_control::keylight::{
    KeylightControl, KeylightJsonDb, KeylightRestAdapter, ZeroConfKeylightFinder, KEYLIGHT_DB_PATH,
};
use model::AppState;

use crate::keylight::discover_keylights;

mod keylight;
mod model;

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let db = KeylightJsonDb::new(KEYLIGHT_DB_PATH);
    let adapter = KeylightRestAdapter {};
    let control = KeylightControl::new(finder, db);
    let app_state = AppState {
        adapter,
        keylight_control: Mutex::new(control),
    };
    tauri::Builder::default()
        .manage(app_state)
        .invoke_handler(tauri::generate_handler![greet, discover_keylights])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
