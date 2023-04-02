use std::time::Duration;
use zeroconf::prelude::*;
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
    fn discover(self) -> Vec<DiscoveredKeylight> {
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
            if let Ok(address) = service.address().parse() {
                devices.push(DiscoveredKeylight::new(
                    service.name().clone(),
                    address,
                    service.port().clone(),
                ));
            }
        }
        devices
    }
}
