use super::KeylightAdapter;

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
    pub brightness: u16,
    pub temperature: u16,
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

    pub fn set_switch(&mut self, light_index: usize, on: bool) -> Result<(), KeylightError> {
        let light = self
            .lights
            .as_ref()
            .ok_or(KeylightError::NoLights)?
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?;
        let mut new_light = light.clone();
        new_light.on = on;
        self.set_light(light_index, new_light)
    }

    pub fn toggle(&mut self, light_index: usize) -> Result<(), KeylightError> {
        let on = self
            .get_lights()?
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?
            .on;
        self.set_switch(light_index, !on)
    }

    fn set_light(&mut self, light_index: usize, light: Light) -> Result<(), KeylightError> {
        let lights = self.lights.as_mut().ok_or(KeylightError::NoLights)?;
        let mut new_lights = lights.clone();
        let new_light = new_lights
            .get_mut(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?;
        *new_light = light;
        self.keylight_adapter
            .set_lights(&self.metadata.ip, &new_lights)?;
        lights[light_index] = new_lights.swap_remove(light_index);
        Ok(())
    }

    fn get_lights(&mut self) -> Result<&mut Vec<Light>, KeylightError> {
        self.lights.as_mut().ok_or(KeylightError::NoLights)
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

    #[test]
    fn lights_should_load_lights() {
        let lights = create_lights_fixture();
        let keylight_adapter = MockKeylightAdapter::new(lights, None);
        let mut keylight = prepare_test(&keylight_adapter, None);
        let result = keylight.lights();

        assert_eq!(result.unwrap(), keylight_adapter.lights);
    }

    #[test]
    fn test_toggle_should_toggle() {
        let keylight_adapter = MockKeylightAdapter::new(vec![], None);
        let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

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
        let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

        let old_light = keylight.lights.as_ref().unwrap()[0].clone();
        let result = keylight.toggle(0);
        assert_eq!(result.is_err(), true);
        assert_eq!(keylight.lights.unwrap()[0].on, old_light.on);
    }

    #[test]
    fn test_toggle_should_not_toggle_if_light_does_not_exist() {
        let keylight_adapter = MockKeylightAdapter::new(vec![], None);
        let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

        let result = keylight.toggle(keylight.lights.as_ref().unwrap().len());
        assert_eq!(result.is_err(), true);
    }
}
