use crate::keylight::keylight_control::KeylightAdapter;

#[derive(Debug)]
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

pub struct Keylight<'a, A: KeylightAdapter> {
    keylight_adapter: &'a A,
    pub metadata: KeylightMetadata,
    pub lights: Option<Vec<Light>>,
}

impl<'a, A: KeylightAdapter> Keylight<'a, A> {
    pub fn new(keylight_adapter: &'a A, metadata: KeylightMetadata) -> Keylight<'a, A> {
        Keylight {
            keylight_adapter,
            metadata,
            lights: None,
        }
    }

    pub fn lights(&mut self) -> Result<&Vec<Light>, KeylightError> {
        let lights = self.keylight_adapter.lights(&self.metadata.ip)?;
        self.lights = Some(lights);
        Ok(self.lights.as_ref().unwrap())
    }
}
