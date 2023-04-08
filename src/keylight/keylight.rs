use crate::keylight::keylight_control::KeylightAdapter;

#[derive(Debug, Eq, PartialEq, Clone)]
pub enum KeylightError {
    CommandError(String),
    LightDoesNotExist(usize),
    NoLights,
}

impl From<reqwest::Error> for KeylightError {
    fn from(error: reqwest::Error) -> Self {
        KeylightError::CommandError(error.to_string())
    }
}

#[derive(Debug, Eq, Hash, PartialEq, Clone)]
pub struct KeylightMetadata {
    pub name: String,
    pub ip: String,
    pub port: u16,
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Light {
    pub on: bool,
    pub brightness: i32,
    pub temperature: i32,
}

pub struct Keylight<'a, A: KeylightAdapter> {
    keylight_adapter: &'a A,
    pub metadata: KeylightMetadata,
    pub lights: Option<Vec<Light>>,
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
            lights: lights.or(None),
        }
    }

    pub fn lights(&mut self) -> Result<&[Light], KeylightError> {
        let lights = self.keylight_adapter.lights(&self.metadata.ip)?;
        self.lights = Some(lights);
        Ok(self.lights.as_ref().unwrap())
    }

    pub fn toggle(&mut self, light_index: usize) -> Result<(), KeylightError> {
        let lights = match self.lights.as_mut() {
            Some(lights) => lights,
            None => return Err(KeylightError::NoLights),
        };
        let mut new_lights = lights.clone();
        let new_light = match new_lights.get_mut(light_index) {
            Some(light) => light,
            None => return Err(KeylightError::LightDoesNotExist(light_index)),
        };
        new_light.on = !new_light.on;
        self.keylight_adapter
            .set_lights(&self.metadata.ip, &new_lights)?;

        lights[light_index].on = !lights[light_index].on;
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
    ) -> Keylight<'a, MockKeylightAdapter> {
        let keylight = create_keylight_fixture(&keylight_adapter, Some(create_lights_fixture()));
        keylight
    }

    #[test]
    fn test_toggle_should_toggle() {
        let keylight_adapter = MockKeylightAdapter::new(vec![], None);
        let mut keylight = prepare_test(&keylight_adapter);

        let old_light = keylight.lights.as_ref().unwrap()[0].clone();
        let result = keylight.toggle(0);
        assert_eq!(result, Ok(()));
        assert_eq!(keylight.lights.unwrap()[0].on, !old_light.on);
    }
    #[test]
    fn test_toggle_should_not_toggle_if_light_cannot_be_updated() {
        let keylight_adapter = MockKeylightAdapter::new(
            vec![],
            Some(Err(KeylightError::CommandError(String::from("error")))),
        );
        let mut keylight = prepare_test(&keylight_adapter);

        let old_light = keylight.lights.as_ref().unwrap()[0].clone();
        let result = keylight.toggle(0);
        assert_eq!(result.is_err(), true);
        assert_eq!(keylight.lights.unwrap()[0].on, old_light.on);
    }
}
