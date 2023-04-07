use std::time::Duration;
use zeroconf::{prelude::*, ServiceDiscovery};
use zeroconf::{MdnsBrowser, ServiceType};

use crate::keylight::KeylightMetadata;
use crate::keylight_control::KeylightFinder;

pub struct ZeroConfKeylightFinder {}

impl ZeroConfKeylightFinder {
    pub fn new() -> ZeroConfKeylightFinder {
        ZeroConfKeylightFinder {}
    }
}

impl KeylightFinder for ZeroConfKeylightFinder {
    type Output = Vec<KeylightMetadata>;

    // The function uses a channel to create the list of the devices because the browser
    // callback function runs on a different thread.
    fn discover(&self) -> Vec<KeylightMetadata> {
        let (tx, rx) = std::sync::mpsc::channel();
        let mut browser = MdnsBrowser::new(ServiceType::new("elg", "tcp").unwrap());
        let mut keylight_metadatas: Vec<KeylightMetadata> = Vec::new();

        browser.set_service_discovered_callback(Box::new(move |result, _context| {
            let service = result.unwrap();
            tx.send(service).unwrap();
        }));

        let event_loop = browser.browse_services().unwrap();
        event_loop.poll(Duration::from_secs(10)).unwrap();

        while let Ok(service) = rx.recv_timeout(Duration::from_secs(2)) {
            self.add_new_device(&mut keylight_metadatas, &service);
        }
        keylight_metadatas
    }
}

impl ZeroConfKeylightFinder {
    fn add_new_device(
        &self,
        keylight_metadatas: &mut Vec<KeylightMetadata>,
        service: &ServiceDiscovery,
    ) {
        if let Ok(address) = service.address().parse() {
            keylight_metadatas.push(KeylightMetadata {
                name: service.name().clone(),
                ip: address,
                port: service.port().clone(),
            });
        }
    }
}
