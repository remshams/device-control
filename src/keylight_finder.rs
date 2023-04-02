use std::time::Duration;
use zeroconf::{prelude::*, ServiceDiscovery};
use zeroconf::{MdnsBrowser, ServiceType};

use crate::keylight::DiscoveredKeylight;
use crate::keylight_control::KeylightFinder;

pub struct ZeroConfKeylightFinder {}

impl ZeroConfKeylightFinder {
    pub fn new() -> ZeroConfKeylightFinder {
        ZeroConfKeylightFinder {}
    }
}

impl KeylightFinder for ZeroConfKeylightFinder {
    // The function uses a channel to create the list of the devices because the browser
    // callback function runs on a different thread.
    fn discover(&self) -> Vec<DiscoveredKeylight> {
        let (tx, rx) = std::sync::mpsc::channel();
        let mut browser = MdnsBrowser::new(ServiceType::new("elg", "tcp").unwrap());
        let mut devices: Vec<DiscoveredKeylight> = Vec::new();

        browser.set_service_discovered_callback(Box::new(move |result, _context| {
            let service = result.unwrap();
            tx.send(service).unwrap();
        }));

        let event_loop = browser.browse_services().unwrap();
        event_loop.poll(Duration::from_secs(10)).unwrap();

        while let Ok(service) = rx.recv_timeout(Duration::from_secs(2)) {
            self.add_new_device(&mut devices, &service);
        }
        devices
    }
}

impl ZeroConfKeylightFinder {
    fn add_new_device(&self, devices: &mut Vec<DiscoveredKeylight>, service: &ServiceDiscovery) {
        if let Ok(address) = service.address().parse() {
            let new_device =
                DiscoveredKeylight::new(service.name().clone(), address, service.port().clone());
            if self.should_add_new_device(&devices, &new_device) {
                devices.push(new_device);
            }
        }
    }
    fn should_add_new_device(
        &self,
        devices: &Vec<DiscoveredKeylight>,
        new_device: &DiscoveredKeylight,
    ) -> bool {
        !devices.iter().any(|device| device == new_device)
    }
}
