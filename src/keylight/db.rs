use super::{control::KeylightDb, KeylightError, KeylightMetadata};
use log::debug;

pub static KEYLIGHT_DB_PATH: &str = "./keylight.json";

pub struct KeylightJsonDb<'a> {
    pub path: &'a str,
}

impl<'a> KeylightJsonDb<'a> {
    pub fn new(path: &str) -> KeylightJsonDb {
        KeylightJsonDb { path }
    }
}

impl<'a> KeylightDb for KeylightJsonDb<'a> {
    fn store(&self, metadatas: &[&KeylightMetadata]) -> Result<(), KeylightError> {
        let metadatas_string = serde_json::to_string(metadatas)?;
        std::fs::write(&self.path, metadatas_string)?;
        debug!("Stored {} keylight metadatas", metadatas.len());
        Ok(())
    }

    fn load(&self) -> Result<Vec<KeylightMetadata>, KeylightError> {
        let metadatas_string = std::fs::read_to_string(&self.path)?;
        let metadatas: Vec<KeylightMetadata> = serde_json::from_str(&metadatas_string)?;
        debug!("Loaded {} keylight metadatas", metadatas.len());
        Ok(metadatas)
    }
}
