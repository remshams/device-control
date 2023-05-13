use super::{model::KeylightError, Keylight, KeylightMetadata, Light};
use log::debug;

pub trait KeylightFinder {
    type Output: IntoIterator<Item = KeylightMetadata>;

    fn discover(&self) -> Self::Output;
}

pub trait KeylightAdapter {
    fn lights(&self, ip: &str, port: &u16) -> Result<Vec<Light>, KeylightError>;
    fn set_lights(&self, ip: &str, lights: &[Light]) -> Result<(), KeylightError>;
}

pub trait KeylightDb {
    fn store(&self, metadatas: &[&KeylightMetadata]) -> Result<(), KeylightError>;
    fn load(&self) -> Result<Vec<KeylightMetadata>, KeylightError>;
}

pub struct KeylightControl<F: KeylightFinder, Db: KeylightDb> {
    keylight_finder: F,
    keylight_db: Db,
    lights: Vec<Keylight>,
}

impl<F: KeylightFinder, Db: KeylightDb> KeylightControl<F, Db> {
    pub fn new(keylight_finder: F, keylight_db: Db) -> KeylightControl<F, Db> {
        KeylightControl {
            keylight_finder,
            keylight_db,
            lights: vec![],
        }
    }

    pub fn find_keylight_mut(&mut self, id: &str) -> Option<&mut Keylight> {
        self.lights.iter_mut().find(|light| light.metadata.id == id)
    }

    pub fn find_keylight(&self, id: &str) -> Option<&Keylight> {
        self.lights.iter().find(|light| light.metadata.id == id)
    }

    pub fn load_keylights(&mut self) -> Result<(), KeylightError> {
        let result = self.keylight_db.load();
        match result {
            Ok(metadata) => {
                self.lights = metadata
                    .into_iter()
                    .map(|metadata| Keylight::new(metadata, None))
                    .collect();
                debug!("Loaded {} keylights", self.lights.len());

                Ok(())
            }
            Err(_) => self.discover_and_store_keylights(),
        }
    }

    pub fn discover_keylights(&mut self) {
        self.lights = self
            .keylight_finder
            .discover()
            .into_iter()
            .map(|metadata| Keylight::new(metadata, None))
            .collect();
        self.deduplicate_keylights();
        debug!("Discovered {} keylights", self.lights.len());
    }

    pub fn store_keylights(&self) -> Result<(), KeylightError> {
        let keylight_metadatas = self
            .lights
            .iter()
            .map(|light| &light.metadata)
            .collect::<Vec<&KeylightMetadata>>();
        debug!("Storing {} keylights", keylight_metadatas.len());
        self.keylight_db.store(keylight_metadatas.as_slice())
    }

    pub fn discover_and_store_keylights(&mut self) -> Result<(), KeylightError> {
        self.discover_keylights();
        self.store_keylights()
    }

    fn deduplicate_keylights(&mut self) {
        self.lights.sort_by_key(|light| light.metadata.ip.clone());
        self.lights.dedup_by_key(|light| light.metadata.ip.clone());
    }

    pub fn list(&self) -> &Vec<Keylight> {
        &self.lights
    }

    pub fn list_metadata(&self) -> Vec<&KeylightMetadata> {
        self.lights
            .iter()
            .map(|keylight| &keylight.metadata)
            .collect()
    }
}

#[cfg(test)]
mod test {

    use std::cell::RefCell;

    use crate::keylight::keylight_mocks::{
        create_metadata_list_fixture, MockKeylightDb, MockKeylightFinder,
    };
    use crate::keylight::model::KeylightMetadata;

    use super::*;

    fn prepare_test() -> (MockKeylightFinder, MockKeylightDb) {
        (
            MockKeylightFinder::new(create_metadata_list_fixture()),
            MockKeylightDb {
                stored_metadata_passed: RefCell::new(vec![]),
                load_response: Ok(create_metadata_list_fixture()),
            },
        )
    }

    mod discover_lights {
        use super::*;

        #[test]
        fn test_discover_keylights() {
            let (finder, db) = prepare_test();
            let deduplicated_metadata = &finder.metadata.clone()[0..2];
            let mut keylight_control = KeylightControl::new(finder, db);
            keylight_control.discover_keylights();
            let discovered_metadata: Vec<KeylightMetadata> = keylight_control
                .lights
                .into_iter()
                .map(|light| light.metadata)
                .collect();
            assert_eq!(discovered_metadata.len(), 2);
            assert_eq!(discovered_metadata, deduplicated_metadata);
        }
    }

    mod load_lights {
        use super::*;

        #[test]
        fn should_load_metadata_from_db() {
            let (finder, db) = prepare_test();
            let keylight_metadatas = db.load_response.clone().unwrap();
            let mut keylight_control = KeylightControl::new(finder, db);
            let result = keylight_control.load_keylights();

            assert_eq!(result.is_ok(), true);
            for (index, keylight) in keylight_control.lights.iter().enumerate() {
                assert_eq!(keylight.metadata, keylight_metadatas[index]);
            }
        }

        #[test]
        fn should_discover_keylights_when_loading_from_db_fails() {
            let (finder, mut db) = prepare_test();
            let metadata = finder.metadata.clone();
            db.load_response = Err(KeylightError::DbError(String::from("Test")));
            let mut keylight_control = KeylightControl::new(finder, db);
            let result = keylight_control.load_keylights();

            assert_eq!(result.is_ok(), true);
            assert_eq!(keylight_control.lights.is_empty(), false);
            for (index, keylight) in keylight_control.lights.iter().enumerate() {
                assert_eq!(keylight.metadata, metadata[index]);
            }
        }
    }
}
