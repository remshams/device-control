use crate::keylight::DiscoveredKeylight;

pub trait KeylightFinder {
    fn discover(&self) -> Vec<DiscoveredKeylight>;
}

pub struct KeylightControl {
    pub lights: Vec<DiscoveredKeylight>,
}

impl KeylightControl {
    pub fn new(keylight_finder: &dyn KeylightFinder) -> KeylightControl {
        let mut lights = keylight_finder.discover();
        lights.dedup();
        KeylightControl { lights }
    }
}
