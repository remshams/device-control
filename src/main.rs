use keylight_on::keylight::{KeylightControl, KeylightRestAdapter, ZeroConfKeylightFinder};

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let mut keylight_control = KeylightControl::new(&finder, &adapter);
    keylight_control.discover_lights();
    for keylight in keylight_control.lights {
        println!("{:?}", keylight.metadata);
    }
}
