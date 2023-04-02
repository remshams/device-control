use keylight_on::keylight_control::DeviceFinder;
use keylight_on::keylight_discover::ZeroConfDeviceFinder;

fn main() {
    let devices = ZeroConfDeviceFinder::new().discover();
    println!("{:?}", devices)
}
