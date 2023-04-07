use crate::keylight::DiscoveredKeylight;

pub trait KeylightFinder {
    type Output: IntoIterator<Item = DiscoveredKeylight>;

    fn discover(&self) -> Self::Output;
}

pub struct KeylightControl<'a, F: KeylightFinder> {
    keylight_finder: &'a F,
    pub lights: Vec<DiscoveredKeylight>,
}

impl<'a, F: KeylightFinder> KeylightControl<'a, F> {
    pub fn new(keylight_finder: &'a F) -> KeylightControl<'a, F> {
        KeylightControl {
            keylight_finder,
            lights: vec![],
        }
    }

    pub fn discover_lights(&mut self) {
        self.lights = self.keylight_finder.discover().into_iter().collect();
        self.deduplicate_lights();
    }

    fn deduplicate_lights(&mut self) {
        self.lights.sort_by_key(|light| light.metadata.ip.clone());
        self.lights.dedup_by_key(|light| light.metadata.ip.clone());
    }
}

#[cfg(test)]
mod test {

    use crate::keylight::Metadata;

    use super::*;

    struct MockKeylightFinder {
        pub metadata: Vec<Metadata>,
    }

    impl KeylightFinder for MockKeylightFinder {
        type Output = Vec<DiscoveredKeylight>;

        fn discover(&self) -> Self::Output {
            self.metadata
                .iter()
                .map(|metadata| {
                    DiscoveredKeylight::new(
                        metadata.name.clone(),
                        metadata.ip.clone(),
                        metadata.port,
                    )
                })
                .collect()
        }
    }

    impl MockKeylightFinder {
        fn new(metadata: Vec<Metadata>) -> MockKeylightFinder {
            MockKeylightFinder { metadata }
        }
    }

    fn prepare_test() -> MockKeylightFinder {
        let test_metadata: Vec<Metadata> = vec![
            Metadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
            Metadata {
                name: String::from("second"),
                ip: String::from("102.168.1.2"),
                port: 4567,
            },
            Metadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
        ];
        MockKeylightFinder::new(test_metadata)
    }

    #[test]
    fn test_discover_lights() {
        let finder = prepare_test();
        let deduplicated_metadata = vec![finder.metadata[0].clone(), finder.metadata[1].clone()];
        let mut keylight_control = KeylightControl::new(&finder);
        keylight_control.discover_lights();
        let discovered_metadata: Vec<Metadata> = keylight_control
            .lights
            .iter()
            .map(|light| light.metadata.clone())
            .collect();
        assert_eq!(keylight_control.lights.len(), 2);
        assert_eq!(discovered_metadata, deduplicated_metadata);
    }
}
