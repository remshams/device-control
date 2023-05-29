use std::path::PathBuf;

use super::{control::KeylightDb, KeylightError, KeylightMetadata};
use log::debug;

pub static KEYLIGHT_DB_PATH: &str = "./keylight.json";

pub struct KeylightJsonDb {
    pub path: PathBuf,
}

impl KeylightJsonDb {
    pub fn new(path: Option<PathBuf>) -> KeylightJsonDb {
        KeylightJsonDb {
            path: path.unwrap_or(PathBuf::from(KEYLIGHT_DB_PATH)),
        }
    }
}

impl KeylightDb for KeylightJsonDb {
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
