use keylight_on::keylight_control::KeylightControl;
use keylight_on::keylight_finder::ZeroConfKeylightFinder;

fn main() {
    let lights = KeylightControl::new(&ZeroConfKeylightFinder::new()).lights;
    println!("{:?}", lights)
}
