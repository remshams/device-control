use crate::keylight_adapter::KeylightAdapter;

pub enum KeylightError {
    CommandError(String),
}

impl From<reqwest::Error> for KeylightError {
    fn from(error: reqwest::Error) -> Self {
        KeylightError::CommandError(error.to_string())
    }
}

#[derive(Debug, Eq, Hash, PartialEq, Clone)]
pub struct Metadata {
    pub name: String,
    pub ip: String,
    pub port: u16,
}

pub struct Light {
    pub on: bool,
    pub brightness: i32,
    pub temperature: i32,
}
pub struct Keylight {
    metadata: Metadata,
    lights: Light,
}

pub struct DiscoveredKeylight {
    pub metadata: Metadata,
}

impl DiscoveredKeylight {
    pub fn new(name: String, ip: String, port: u16) -> DiscoveredKeylight {
        DiscoveredKeylight {
            metadata: Metadata { name, ip, port },
        }
    }

    pub fn connect(self) -> Keylight {
        unreachable!();
    }
}

impl Keylight {
    pub fn update(&mut self, operating_data: Light) -> &Self {
        unimplemented!()
    }

    pub fn turn_on(&mut self) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change_brightness(&mut self, brightness: i32) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change_temperature(&mut self, temperature: i32) -> Result<&Self, KeylightError> {
        unimplemented!()
    }

    pub fn change(
        &mut self,
        on: bool,
        brightness: i32,
        temperature: i32,
    ) -> Result<&Self, KeylightError> {
        unimplemented!()
    }
}
