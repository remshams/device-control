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

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Light {
    pub on: bool,
    pub brightness: i32,
    pub temperature: i32,
}
pub struct Keylight {
    pub metadata: KeylightMetadata,
    pub lights: Option<Vec<Light>>,
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
