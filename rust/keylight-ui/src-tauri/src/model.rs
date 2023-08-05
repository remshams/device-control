use std::sync::Mutex;

use keylight_control::keylight::{
    KeylightControl, KeylightJsonDb, KeylightRestAdapter, ZeroConfKeylightFinder,
};

pub struct AppState {
    pub adapter: KeylightRestAdapter,
    pub keylight_control: Mutex<KeylightControl<ZeroConfKeylightFinder, KeylightJsonDb>>,
}
