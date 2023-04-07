use serde::Deserialize;

use super::{
    keylight::{KeylightError, Light},
    keylight_control::KeylightAdapter,
};

#[derive(Debug, Deserialize)]
struct LightDto {
    pub on: i32,
    pub brightness: i32,
    pub temperature: i32,
}

#[derive(Debug, Deserialize)]
struct StatusDto {
    pub lights: Vec<LightDto>,
}

pub struct KeylightRestAdapter {}

impl KeylightAdapter for KeylightRestAdapter {
    fn lights(&self, ip: &String) -> Result<Vec<Light>, KeylightError> {
        let status = reqwest::blocking::get(&format!("http://{}:9123/elgato/lights", ip))?
            .json::<StatusDto>()?;
        let lights = status.lights.into_iter().map(|light| Light {
            on: light.on == 1,
            brightness: light.brightness,
            temperature: light.temperature,
        });
        Ok(lights.collect())
    }
}
