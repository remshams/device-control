use crate::keylight::keylight::{Keylight, KeylightError, KeylightMetadata, Light};

pub trait KeylightFinder {
    type Output: IntoIterator<Item = KeylightMetadata>;

    fn discover(&self) -> Self::Output;
}

pub trait KeylightAdapter {
    fn lights(&self, ip: &str) -> Result<Vec<Light>, KeylightError>;
    fn set_lights(&self, ip: &str, lights: &Vec<Light>) -> Result<(), KeylightError>;
}

pub struct KeylightControl<'a, F: KeylightFinder, A: KeylightAdapter> {
    keylight_finder: &'a F,
    keylight_adapter: &'a A,
    pub lights: Vec<Keylight<'a, A>>,
}

impl<'a, F: KeylightFinder, A: KeylightAdapter> KeylightControl<'a, F, A> {
    pub fn new(keylight_finder: &'a F, keylight_adapter: &'a A) -> KeylightControl<'a, F, A> {
        KeylightControl {
            keylight_finder,
            keylight_adapter,
            lights: vec![],
        }
    }

    pub fn discover_lights(&mut self) {
        self.lights = self
            .keylight_finder
            .discover()
            .into_iter()
            .map(|metadata| Keylight::new(self.keylight_adapter, metadata, None))
            .collect();
        self.deduplicate_lights();
    }

    fn deduplicate_lights(&mut self) {
        self.lights.sort_by_key(|light| light.metadata.ip.clone());
        self.lights.dedup_by_key(|light| light.metadata.ip.clone());
    }
}

#[cfg(test)]
mod test {

    use crate::keylight::keylight::KeylightMetadata;
    use crate::keylight::keylight_mocks::{
        create_metadata_list_fixture, MockKeylightAdapter, MockKeylightFinder,
    };

    use super::*;

    fn prepare_test() -> (MockKeylightFinder, MockKeylightAdapter) {
        (
            MockKeylightFinder::new(create_metadata_list_fixture()),
            MockKeylightAdapter::new(vec![], None),
        )
    }

    #[test]
    fn test_discover_lights() {
        let (finder, adapter) = prepare_test();
        let deduplicated_metadata = vec![&finder.metadata[0], &finder.metadata[1]];
        let mut keylight_control = KeylightControl::new(&finder, &adapter);
        keylight_control.discover_lights();
        let discovered_metadata: Vec<&KeylightMetadata> = keylight_control
            .lights
            .iter()
            .map(|light| &light.metadata)
            .collect();
        assert_eq!(keylight_control.lights.len(), 2);
        assert_eq!(discovered_metadata, deduplicated_metadata);
    }
}
