use keylight_on::keylight_adapter::KeylightRestAdapter;
use keylight_on::keylight_control::KeylightControl;
use keylight_on::keylight_finder::ZeroConfKeylightFinder;

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let adapter = KeylightRestAdapter {};
    let mut keylight_control = KeylightControl::new(&finder, &adapter);
    keylight_control.discover_lights();
    for keylight in keylight_control.lights {
        println!("{:?}", keylight.metadata);
    }
}
