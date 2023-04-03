use std::collections::HashSet;

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
        let mut seen: HashSet<DiscoveredKeylight> = HashSet::new();
        let mut unique_lights: Vec<DiscoveredKeylight> = vec![];
        for light in self.lights.drain(..) {
            if seen.insert(light.clone()) {
                unique_lights.push(seen.replace(light).unwrap());
            }
        }
        self.lights = unique_lights
    }
}

#[cfg(test)]
mod test {
    use super::*;

    struct MockKeylightFinder {
        pub lights: Vec<DiscoveredKeylight>,
    }

    impl KeylightFinder for MockKeylightFinder {
        type Output = Vec<DiscoveredKeylight>;

        fn discover(&self) -> Self::Output {
            self.lights.clone()
        }
    }

    impl MockKeylightFinder {
        fn new(lights: Vec<DiscoveredKeylight>) -> MockKeylightFinder {
            MockKeylightFinder { lights }
        }
    }

    fn prepare_test() -> MockKeylightFinder {
        let test_lights: Vec<DiscoveredKeylight> = vec![
            DiscoveredKeylight::new(String::from("first"), String::from("192.168.1.1"), 1234),
            DiscoveredKeylight::new(String::from("second"), String::from("192.168.1.2"), 4567),
            DiscoveredKeylight::new(String::from("first"), String::from("192.168.1.1"), 1234),
        ];
        MockKeylightFinder::new(test_lights)
    }

    #[test]
    fn test_discover_lights() {
        let finder = prepare_test();
        let deduplicated_lights = vec![finder.lights[0].clone(), finder.lights[1].clone()];
        let mut keylight_control = KeylightControl::new(&finder);
        keylight_control.discover_lights();
        assert_eq!(keylight_control.lights.len(), 2);
        assert_eq!(keylight_control.lights, deduplicated_lights)
    }
}
