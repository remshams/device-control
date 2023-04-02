use crate::keylight::DiscoveredKeylight;

pub trait KeylightFinder {
    fn discover(&self) -> Vec<DiscoveredKeylight>;
}

pub struct KeylightControl<'a> {
    keylight_finder: &'a dyn KeylightFinder,
    pub lights: Vec<DiscoveredKeylight>,
}

impl<'a> KeylightControl<'a> {
    pub fn new(keylight_finder: &dyn KeylightFinder) -> KeylightControl {
        KeylightControl {
            keylight_finder,
            lights: vec![],
        }
    }

    pub fn discover_lights(&mut self) {
        self.lights = self.keylight_finder.discover();
        self.lights.dedup();
    }
}
