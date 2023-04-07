mod keylight;
mod keylight_adapter;
mod keylight_control;
mod keylight_finder;

pub use keylight::{Keylight, KeylightMetadata, Light};
pub use keylight_adapter::KeylightRestAdapter;
pub use keylight_control::{KeylightAdapter, KeylightControl, KeylightFinder};
pub use keylight_finder::ZeroConfKeylightFinder;

#[cfg(test)]
mod keylight_mocks {

    use super::{keylight::KeylightError, *};

    pub struct MockKeylightFinder {
        pub metadata: Vec<KeylightMetadata>,
    }

    impl KeylightFinder for MockKeylightFinder {
        type Output = Vec<KeylightMetadata>;

        fn discover(&self) -> Self::Output {
            self.metadata.clone()
        }
    }

    impl MockKeylightFinder {
        pub fn new(metadata: Vec<KeylightMetadata>) -> MockKeylightFinder {
            MockKeylightFinder { metadata }
        }
    }

    pub struct MockKeylightAdapter {
        pub lights: Vec<Light>,
    }

    impl KeylightAdapter for MockKeylightAdapter {
        fn lights(&self, ip: &String) -> Result<Vec<Light>, KeylightError> {
            Ok(self.lights.clone())
        }
    }

    pub fn create_metadata_fixture() -> Vec<KeylightMetadata> {
        vec![
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
            KeylightMetadata {
                name: String::from("second"),
                ip: String::from("102.168.1.2"),
                port: 4567,
            },
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
        ]
    }
}
