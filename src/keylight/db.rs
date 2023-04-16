use super::{control::KeylightDb, KeylightError, KeylightMetadata};

pub struct KeylightJsonDb {
    pub path: String,
}

impl KeylightJsonDb {
    pub fn new(path: String) -> KeylightJsonDb {
        KeylightJsonDb { path }
    }
}

impl KeylightDb for KeylightJsonDb {
    fn store(&self, metadatas: &[&KeylightMetadata]) -> Result<(), KeylightError> {
        let metadatas_string = serde_json::to_string(metadatas)?;
        std::fs::write(&self.path, metadatas_string)?;
        Ok(())
    }

    fn load(&self) -> Result<Vec<KeylightMetadata>, KeylightError> {
        let metadatas_string = std::fs::read_to_string(&self.path)?;
        let metadatas: Vec<KeylightMetadata> = serde_json::from_str(&metadatas_string)?;
        Ok(metadatas)
    }
}
