use crate::keylight::{Keylight, KeylightError, KeylightMetadata, Light};

pub trait KeylightFinder {
    type Output: IntoIterator<Item = KeylightMetadata>;

    fn discover(&self) -> Self::Output;
}

pub trait KeylightAdapter {
    fn status(&self, ip: &String) -> Result<Vec<Light>, KeylightError>;
}

pub struct KeylightControl<'a, F: KeylightFinder, A: KeylightAdapter> {
    keylight_finder: &'a F,
    keylight_adapter: &'a A,
    pub lights: Vec<Keylight>,
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
            .map(|metadata| Keylight {
                metadata,
                lights: None,
            })
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

    use crate::keylight::KeylightMetadata;

    use super::*;

    struct MockKeylightFinder {
        pub metadata: Vec<KeylightMetadata>,
    }

    impl KeylightFinder for MockKeylightFinder {
        type Output = Vec<KeylightMetadata>;

        fn discover(&self) -> Self::Output {
            self.metadata.clone()
        }
    }

    impl MockKeylightFinder {
        fn new(metadata: Vec<KeylightMetadata>) -> MockKeylightFinder {
            MockKeylightFinder { metadata }
        }
    }

    struct MockKeylightAdapter {
        pub lights: Vec<Light>,
    }

    impl KeylightAdapter for MockKeylightAdapter {
        fn status(&self, ip: &String) -> Result<Vec<Light>, KeylightError> {
            Ok(self.lights.clone())
        }
    }

    fn prepare_test() -> (MockKeylightFinder, MockKeylightAdapter) {
        let test_metadata: Vec<KeylightMetadata> = vec![
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
            KeylightMetadata {
                name: String::from("second"),
                ip: String::from("102.168.1.2"),
                port: 4567,
            },
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
        ];
        let lights = vec![
            Light {
                on: true,
                brightness: 50,
                temperature: 3000,
            },
            Light {
                on: false,
                brightness: 100,
                temperature: 6500,
            },
        ];
        (
            MockKeylightFinder::new(test_metadata),
            MockKeylightAdapter { lights },
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
        let keylight = keylight_control.lights.get(0).unwrap();
        assert_eq!(keylight.lights, None)
    }
}
