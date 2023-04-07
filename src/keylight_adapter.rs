use crate::{
    keylight::{KeylightError, Light},
    keylight_control::KeylightAdapter,
};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct LightDto {
    pub on: i32,
    pub brightness: i32,
    pub temperature: i32,
}

#[derive(Debug, Deserialize)]
#[allow(non_snake_case)]
struct StatusDto {
    pub numberOfLights: i32,
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
