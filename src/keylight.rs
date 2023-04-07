pub enum KeylightError {
    CommandError(String),
}

impl From<reqwest::Error> for KeylightError {
    fn from(error: reqwest::Error) -> Self {
        KeylightError::CommandError(error.to_string())
    }
}

#[derive(Debug, Eq, Hash, PartialEq, Clone)]
pub struct KeylightMetadata {
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
    metadata: KeylightMetadata,
    lights: Light,
}

pub struct DiscoveredKeylight {
    pub metadata: KeylightMetadata,
}

impl DiscoveredKeylight {
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
