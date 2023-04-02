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
        self.lights.dedup();
    }
}
