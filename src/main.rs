use keylight_on::keylight_control::KeylightControl;
use keylight_on::keylight_finder::ZeroConfKeylightFinder;

fn main() {
    let finder = ZeroConfKeylightFinder::new();
    let mut keylight_control = KeylightControl::new(&finder);
    keylight_control.discover_lights();
    for keylight in keylight_control.lights {
        println!("{:?}", keylight.metadata);
    }
}
