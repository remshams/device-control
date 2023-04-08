use serde::{Deserialize, Serialize};

use super::{
    keylight::{KeylightError, Light},
    keylight_control::KeylightAdapter,
};

#[derive(Debug, Serialize, Deserialize)]
struct LightDto {
    pub on: i32,
    pub brightness: i32,
    pub temperature: i32,
}

#[derive(Debug, Serialize, Deserialize)]
struct StatusDto {
    pub lights: Vec<LightDto>,
}

pub struct KeylightRestAdapter {}

impl KeylightAdapter for KeylightRestAdapter {
    fn lights(&self, ip: &str) -> Result<Vec<Light>, KeylightError> {
        let status = reqwest::blocking::get(&format!("http://{}:9123/elgato/lights", ip))?
            .json::<StatusDto>()?;
        let lights = status
            .lights
            .into_iter()
            .map(|light| Light {
                on: light.on == 1,
                brightness: light.brightness,
                temperature: light.temperature,
            })
            .collect();
        Ok(lights)
    }

    fn set_lights(&self, ip: &str, lights: &[Light]) -> Result<(), KeylightError> {
        let client = reqwest::blocking::Client::new();
        client
            .put(&format!("http://{}:9123/elgato/lights", ip))
            .json::<StatusDto>(&StatusDto {
                lights: lights
                    .iter()
                    .map(|light| LightDto {
                        on: if light.on { 1 } else { 0 },
                        brightness: light.brightness,
                        temperature: light.temperature,
                    })
                    .collect(),
            })
            .send()?;
        Ok(())
    }
}
