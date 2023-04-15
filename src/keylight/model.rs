use super::KeylightAdapter;

#[derive(Debug, Eq, PartialEq, Clone)]
pub enum KeylightError {
    CommandError(String),
    LightDoesNotExist(usize),
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
        let lights = self.keylight_adapter.lights(&self.metadata.ip)?;
        self.lights = lights;
        Ok(self.lights.as_ref())
    }

    pub fn set_switch(&mut self, light_index: usize, on: bool) -> Result<(), KeylightError> {
        let light = self
            .lights
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?;
        let mut new_light = light.clone();
        new_light.on = on;
        self.set_light(light_index, new_light)
    }

    pub fn toggle(&mut self, light_index: usize) -> Result<(), KeylightError> {
        let on = self
            .lights
            .get(light_index)
            .ok_or(KeylightError::LightDoesNotExist(light_index))?
            .on;
        self.set_switch(light_index, !on)
    }

    fn set_light(&mut self, light_index: usize, light: Light) -> Result<(), KeylightError> {
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

        #[test]
        fn should_set_light() {
            let keylight_adapter = MockKeylightAdapter::new(Ok(vec![]), None);
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

            let old_light = keylight.lights[0].clone();
            let mut new_light = old_light.clone();
            new_light.on = !old_light.on;
            let result = keylight.set_light(0, new_light);
            assert_eq!(result, Ok(()));
            assert_eq!(keylight.lights[0].on, !old_light.on);
        }
        #[test]
        fn should_not_update_if_light_cannot_be_updated() {
            let keylight_adapter = MockKeylightAdapter::new(
                Ok(vec![]),
                Some(Err(KeylightError::CommandError(String::from("error")))),
            );
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));

            let old_light = keylight.lights[0].clone();
            let mut new_light = old_light.clone();
            new_light.on = !old_light.on;
            let result = keylight.set_light(0, new_light);
            assert_eq!(result.is_err(), true);
            assert_eq!(keylight.lights[0].on, old_light.on);
        }

        #[test]
        fn should_do_nothing_if_light_does_not_exist() {
            let keylight_adapter = MockKeylightAdapter::new(Ok(vec![]), None);
            let mut keylight = prepare_test(&keylight_adapter, Some(create_lights_fixture()));
            let old_lights = keylight.lights.clone();
            let mut new_light = old_lights[0].clone();
            new_light.on = !old_lights[0].on;

            let result = keylight.set_light(keylight.lights.len(), new_light);
            assert_eq!(result.is_err(), true);
        }
    }
}
