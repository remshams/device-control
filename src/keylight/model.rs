use log::debug;
use serde::{Deserialize, Serialize};

use super::KeylightAdapter;

#[derive(Debug, Eq, PartialEq, Clone)]
pub enum KeylightError {
    CommandError(String),
    LightDoesNotExist(usize),
    DbError(String),
}

impl From<reqwest::Error> for KeylightError {
    fn from(error: reqwest::Error) -> Self {
        KeylightError::CommandError(error.to_string())
    }
}

impl From<serde_json::Error> for KeylightError {
    fn from(error: serde_json::Error) -> Self {
        KeylightError::DbError(error.to_string())
    }
}

impl From<std::io::Error> for KeylightError {
    fn from(error: std::io::Error) -> Self {
        KeylightError::DbError(error.to_string())
    }
}

#[derive(Debug, Eq, Hash, PartialEq, Clone, Serialize, Deserialize)]
pub struct KeylightMetadata {
    pub name: String,
    pub ip: String,
    pub port: u16,
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Light {
    pub on: bool,
    pub brightness: u16,
    pub temperature: u16,
}
#[derive(Debug)]
pub struct CommandLight {
    pub index: usize,
    pub on: Option<bool>,
    pub brightness: Option<u16>,
    pub temperature: Option<u16>,
}

impl CommandLight {
    pub fn from_light(index: usize, light: &Light) -> CommandLight {
        CommandLight {
            index,
            on: Some(light.on),
            brightness: Some(light.brightness),
            temperature: Some(light.temperature),
        }
    }
}

pub struct Keylight<'a, A: KeylightAdapter> {
    keylight_adapter: &'a A,
    pub metadata: KeylightMetadata,
    pub lights: Vec<Light>,
}

impl<'a, A: KeylightAdapter> Keylight<'a, A> {
    pub fn new(
        keylight_adapter: &'a A,
        metadata: KeylightMetadata,
        lights: Option<Vec<Light>>,
    ) -> Keylight<'a, A> {
        Keylight {
            keylight_adapter,
            metadata,
            lights: lights.unwrap_or(vec![]),
        }
    }

    pub fn lights(&mut self) -> Result<&[Light], KeylightError> {
        let lights = self
            .keylight_adapter
            .lights(&self.metadata.ip, &self.metadata.port)?;
        self.lights = lights;
        debug!(
            "Found {} lights for keylight {:#?}",
            self.lights.len(),
            self.metadata
        );
        Ok(self.lights.as_ref())
    }

    pub fn set_switch(&mut self, light_index: usize, on: bool) -> Result<(), KeylightError> {
        let light = self
            .lights
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?;
        let mut new_light = light.clone();
        new_light.on = on;
        debug!("Switch light {} to {}", light_index, on);
        self.update_light(light_index, new_light)
    }

    pub fn toggle(&mut self, light_index: usize) -> Result<(), KeylightError> {
        let on = self
            .lights
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?
            .on;
        debug!("Toggle light {} from {}", light_index, on);
        self.set_switch(light_index, !on)
    }

    pub fn set_light(&mut self, command_light: CommandLight) -> Result<(), KeylightError> {
        let light = self
            .lights
            .get(command_light.index)
            .ok_or(KeylightError::LightDoesNotExist(command_light.index))?;
        let mut new_light = light.clone();
        new_light.on = command_light.on.unwrap_or(light.on);
        new_light.brightness = command_light.brightness.unwrap_or(light.brightness);
        new_light.temperature = command_light.temperature.unwrap_or(light.temperature);
        debug!("Set Light: {:#?}", new_light);
        self.update_light(command_light.index, new_light)
    }

    fn update_light(&mut self, light_index: usize, light: Light) -> Result<(), KeylightError> {
        let mut new_lights = self.lights.clone();
        let new_light = new_lights
            .get_mut(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?;
        *new_light = light;
        self.keylight_adapter
            .set_lights(&self.metadata.ip, &new_lights)?;
        self.lights[light_index] = new_lights.swap_remove(light_index);
        Ok(())
    }
}

#[cfg(test)]
mod test {
    use std::vec;

    use crate::keylight::keylight_mocks::{
        create_keylight_fixture, create_lights_fixture, MockKeylightAdapter,
    };

    use super::*;

    fn prepare_test<'a>(
        keylight_adapter: &'a MockKeylightAdapter,
        lights: Option<Vec<Light>>,
    ) -> Keylight<'a, MockKeylightAdapter> {
        let keylight = create_keylight_fixture(&keylight_adapter, lights);
        keylight
    }

    mod lights {
        use super::*;
        #[test]
        fn should_load_lights() {
            let lights = create_lights_fixture();
            let keylight_adapter = MockKeylightAdapter::new(Ok(lights), None);
            let mut keylight = prepare_test(&keylight_adapter, None);
            let result = keylight.lights();

            assert_eq!(result.unwrap(), keylight_adapter.lights.as_ref().unwrap());
        }

        #[test]
        fn should_return_error_result_if_lights_cannot_be_loaded() {
            let keylight_adapter = MockKeylightAdapter::new(
                Err(KeylightError::CommandError(String::from("error"))),
                None,
            );
            let mut keylight = prepare_test(&keylight_adapter, None);
            let result = keylight.lights();

            assert_eq!(result.is_err(), true);
        }
    }

    mod set_light {
        use super::*;

        fn prepare_test_lights<'a>(
            keylight: &'a Keylight<MockKeylightAdapter>,
        ) -> (&'a Light, bool, CommandLight) {
            let old_light = &keylight.lights[0];
            let old_light_on = old_light.on;
            (
                old_light,
                old_light_on,
                CommandLight::from_light(0, old_light),
            )
        }

        #[test]
        fn should_set_light() {
            let keylight_adapter = MockKeylightAdapter::new(Ok(vec![]), None);
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

            let (_old_light, old_light_on, mut new_light) = prepare_test_lights(&keylight);
            new_light.on = Some(!old_light_on);
            let result = keylight.set_light(new_light);
            assert_eq!(result, Ok(()));
            assert_eq!(keylight.lights[0].on, !old_light_on);
        }
        #[test]
        fn should_not_update_if_light_cannot_be_updated() {
            let keylight_adapter = MockKeylightAdapter::new(
                Ok(vec![]),
                Some(Err(KeylightError::CommandError(String::from("error")))),
            );
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

            let (_old_light, old_light_on, mut new_light) = prepare_test_lights(&keylight);
            new_light.on = Some(!old_light_on);
            let result = keylight.set_light(new_light);
            assert_eq!(result.is_err(), true);
            assert_eq!(keylight.lights[0].on, old_light_on);
        }

        #[test]
        fn should_do_nothing_if_light_does_not_exist() {
            let keylight_adapter = MockKeylightAdapter::new(Ok(vec![]), None);
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));
            let old_lights = keylight.lights.clone();
            let old_light = &old_lights[0];
            let mut new_light = CommandLight::from_light(keylight.lights.len(), old_light);
            new_light.on = Some(!old_light.on);

            let result = keylight.set_light(new_light);
            assert_eq!(result.is_err(), true);
        }
    }
}
