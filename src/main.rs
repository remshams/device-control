use keylight_on::keylight_control::KeylightFinder;
use keylight_on::keylight_finder::ZeroConfKeylightFinder;

fn main() {
    let devices = ZeroConfKeylightFinder::new().discover();
    println!("{:?}", devices)
}
